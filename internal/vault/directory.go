package vault

import (
	"archive/tar"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

const (
	EncryptedDirExt = ".dxrk.enc"
	BackupExt       = ".tar.gz.enc"
	MaxFileSize     = 10 * 1024 * 1024 * 1024 // 10GB max
	DirManifestName = "manifest.json"
)

// DirectoryManifest contains metadata about an encrypted directory
type DirectoryManifest struct {
	Version   int         `json:"version"`
	Algorithm string      `json:"algorithm"`
	CreatedAt time.Time   `json:"created_at"`
	SourceDir string      `json:"source_dir"`
	FileCount int         `json:"file_count"`
	TotalSize int64       `json:"total_size"`
	Checksum  string      `json:"checksum"`
	Files     []FileEntry `json:"files"`
}

// FileEntry represents a single encrypted file in the directory
type FileEntry struct {
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	Encrypted string    `json:"encrypted_path"`
	ModTime   time.Time `json:"mod_time"`
	Checksum  string    `json:"checksum"`
}

// EncryptDirectory encrypts an entire directory using AES-256-GCM
// Creates a tar.gz archive first, then encrypts the archive
func (v *Vault) EncryptDirectory(srcDir, destPath string) error {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.key == nil {
		return fmt.Errorf("vault: not initialized")
	}

	// Validate source directory
	srcInfo, err := os.Stat(srcDir)
	if err != nil {
		return fmt.Errorf("vault: source directory not found: %w", err)
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("vault: source is not a directory")
	}

	// Create temporary tar.gz file
	tmpFile, err := os.CreateTemp("", "dxrk-vault-*.tar.gz")
	if err != nil {
		return fmt.Errorf("vault: failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Compress directory to tar.gz
	if err = compressDirectory(srcDir, tmpFile); err != nil {
		tmpFile.Close()
		return fmt.Errorf("vault: compression failed: %w", err)
	}
	tmpFile.Close()

	// Read compressed data
	compressedData, err := os.ReadFile(tmpPath)
	if err != nil {
		return fmt.Errorf("vault: failed to read compressed data: %w", err)
	}

	// Calculate checksum of compressed data
	checksum := sha256.Sum256(compressedData)
	checksumStr := base64.StdEncoding.EncodeToString(checksum[:])

	// Encrypt the compressed data
	ciphertext, err := v.encryptBytes(compressedData)
	if err != nil {
		return fmt.Errorf("vault: encryption failed: %w", err)
	}

	// Create destination directory
	destDir := filepath.Dir(destPath)
	if err = os.MkdirAll(destDir, 0700); err != nil {
		return fmt.Errorf("vault: failed to create destination: %w", err)
	}

	// Write encrypted file
	if err = os.WriteFile(destPath, ciphertext, 0600); err != nil {
		return fmt.Errorf("vault: failed to write encrypted file: %w", err)
	}

	// Create and write manifest
	manifest, err := createManifest(srcDir, checksumStr)
	if err != nil {
		return fmt.Errorf("vault: failed to create manifest: %w", err)
	}

	manifestPath := destPath + ".manifest"
	if err := writeManifest(manifestPath, manifest); err != nil {
		return fmt.Errorf("vault: failed to write manifest: %w", err)
	}

	return nil
}

// DecryptDirectory decrypts an encrypted directory archive
func (v *Vault) DecryptDirectory(encryptedPath, destDir string) error {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.key == nil {
		return fmt.Errorf("vault: not initialized")
	}

	// Read encrypted file
	ciphertext, err := os.ReadFile(encryptedPath)
	if err != nil {
		return fmt.Errorf("vault: failed to read encrypted file: %w", err)
	}

	// Decrypt the data
	compressedData, err := v.decryptBytes(ciphertext)
	if err != nil {
		return fmt.Errorf("vault: decryption failed: %w", err)
	}

	// Create destination directory
	if err = os.MkdirAll(destDir, 0700); err != nil {
		return fmt.Errorf("vault: failed to create destination: %w", err)
	}

	// Create temporary file for decompression
	tmpFile, err := os.CreateTemp("", "dxrk-vault-dec-*.tar.gz")
	if err != nil {
		return fmt.Errorf("vault: failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Write compressed data to temp file
	if _, err := tmpFile.Write(compressedData); err != nil {
		tmpFile.Close()
		return fmt.Errorf("vault: failed to write temp file: %w", err)
	}
	tmpFile.Close()

	// Decompress and extract
	if err := extractArchive(tmpPath, destDir); err != nil {
		return fmt.Errorf("vault: extraction failed: %w", err)
	}

	return nil
}

// EncryptDirectoryFiles encrypts each file in a directory individually
// This allows partial access and better security
func (v *Vault) EncryptDirectoryFiles(srcDir, destDir string) error {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.key == nil {
		return fmt.Errorf("vault: not initialized")
	}

	// Create destination directory
	if err := os.MkdirAll(destDir, 0700); err != nil {
		return fmt.Errorf("vault: failed to create destination: %w", err)
	}

	var manifest DirectoryManifest
	manifest.Version = version
	manifest.Algorithm = "AES-256-GCM"
	manifest.CreatedAt = time.Now()
	manifest.SourceDir = filepath.Base(srcDir)

	// Walk through source directory
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Skip .vault directory
		if strings.Contains(path, ".vault") {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("vault: failed to read %s: %w", path, err)
		}

		if int64(len(content)) > MaxFileSize {
			return fmt.Errorf("vault: file %s exceeds max size", path)
		}

		// Encrypt file
		encrypted, err := v.encryptBytes(content)
		if err != nil {
			return fmt.Errorf("vault: failed to encrypt %s: %w", path, err)
		}

		// Create relative path structure in destination
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return fmt.Errorf("vault: failed to get relative path: %w", err)
		}

		destPath := filepath.Join(destDir, relPath+".enc")
		destDirPath := filepath.Dir(destPath)
		if err := os.MkdirAll(destDirPath, 0700); err != nil {
			return fmt.Errorf("vault: failed to create dir: %w", err)
		}

		// Write encrypted file
		if err := os.WriteFile(destPath, encrypted, 0600); err != nil {
			return fmt.Errorf("vault: failed to write encrypted: %w", err)
		}

		// Calculate checksum
		checksum := sha256.Sum256(content)

		// Add to manifest
		manifest.Files = append(manifest.Files, FileEntry{
			Path:      relPath,
			Size:      info.Size(),
			Encrypted: relPath + ".enc",
			ModTime:   info.ModTime(),
			Checksum:  base64.StdEncoding.EncodeToString(checksum[:]),
		})
		manifest.FileCount++
		manifest.TotalSize += info.Size()

		return nil
	})

	if err != nil {
		return fmt.Errorf("vault: directory walk failed: %w", err)
	}

	// Write manifest
	manifestPath := filepath.Join(destDir, DirManifestName)
	if err := writeManifest(manifestPath, &manifest); err != nil {
		return fmt.Errorf("vault: failed to write manifest: %w", err)
	}

	return nil
}

// DecryptDirectoryFiles decrypts individually encrypted files back to directory
func (v *Vault) DecryptDirectoryFiles(encryptedDir, destDir string) error {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.key == nil {
		return fmt.Errorf("vault: not initialized")
	}

	// Read manifest
	manifestPath := filepath.Join(encryptedDir, DirManifestName)
	manifest, err := readManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("vault: failed to read manifest: %w", err)
	}

	// Create destination directory
	if err := os.MkdirAll(destDir, 0700); err != nil {
		return fmt.Errorf("vault: failed to create destination: %w", err)
	}

	// Decrypt each file
	for _, file := range manifest.Files {
		encPath := filepath.Join(encryptedDir, file.Encrypted)
		destPath := filepath.Join(destDir, file.Path)

		// Create destination directory
		destDirPath := filepath.Dir(destPath)
		if err := os.MkdirAll(destDirPath, 0700); err != nil {
			return fmt.Errorf("vault: failed to create dir: %w", err)
		}

		// Read encrypted file
		ciphertext, err := os.ReadFile(encPath)
		if err != nil {
			return fmt.Errorf("vault: failed to read %s: %w", encPath, err)
		}

		// Decrypt
		plaintext, err := v.decryptBytes(ciphertext)
		if err != nil {
			return fmt.Errorf("vault: failed to decrypt %s: %w", encPath, err)
		}

		// Verify checksum
		checksum := sha256.Sum256(plaintext)
		if base64.StdEncoding.EncodeToString(checksum[:]) != file.Checksum {
			return fmt.Errorf("vault: checksum mismatch for %s", file.Path)
		}

		// Write decrypted file
		if err := os.WriteFile(destPath, plaintext, 0600); err != nil {
			return fmt.Errorf("vault: failed to write %s: %w", destPath, err)
		}
	}

	return nil
}

// CreateEncryptedBackup creates an encrypted tar.gz backup
func (v *Vault) CreateEncryptedBackup(srcDir, backupPath, password string) error {
	v.mu.RLock()
	defer v.mu.RUnlock()

	// Create temporary tar.gz file
	tmpFile, err := os.CreateTemp("", "dxrk-backup-*.tar.gz")
	if err != nil {
		return fmt.Errorf("vault: failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Compress directory
	if err = compressDirectory(srcDir, tmpFile); err != nil {
		tmpFile.Close()
		return fmt.Errorf("vault: compression failed: %w", err)
	}
	tmpFile.Close()

	// Read compressed data
	compressedData, err := os.ReadFile(tmpPath)
	if err != nil {
		return fmt.Errorf("vault: failed to read compressed data: %w", err)
	}

	// Derive backup key from password
	salt := make([]byte, SaltLen)
	if _, err = io.ReadFull(rand.Reader, salt); err != nil {
		return fmt.Errorf("vault: failed to generate salt: %w", err)
	}
	backupKey := deriveBackupKey(password, salt)

	// Encrypt with backup key
	ciphertext, err := encryptWithKey(backupKey, compressedData)
	if err != nil {
		return fmt.Errorf("vault: encryption failed: %w", err)
	}

	// Prepend salt to ciphertext
	finalData := append(salt, ciphertext...)

	// Create backup directory
	backupDir := filepath.Dir(backupPath)
	if err = os.MkdirAll(backupDir, 0700); err != nil {
		return fmt.Errorf("vault: failed to create backup dir: %w", err)
	}

	// Write backup file
	if err = os.WriteFile(backupPath, finalData, 0600); err != nil {
		return fmt.Errorf("vault: failed to write backup: %w", err)
	}

	return nil
}

// RestoreEncryptedBackup restores from an encrypted backup
func (v *Vault) RestoreEncryptedBackup(backupPath, destDir, password string) error {
	v.mu.RLock()
	defer v.mu.RUnlock()

	// Read backup file
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("vault: failed to read backup: %w", err)
	}

	if len(data) < SaltLen {
		return fmt.Errorf("vault: backup file corrupted")
	}

	// Extract salt and ciphertext
	salt := data[:SaltLen]
	ciphertext := data[SaltLen:]

	// Derive backup key
	backupKey := deriveBackupKey(password, salt)

	// Decrypt
	compressedData, err := decryptWithKey(backupKey, ciphertext)
	if err != nil {
		return fmt.Errorf("vault: decryption failed (wrong password?): %w", err)
	}

	// Create destination directory
	if err = os.MkdirAll(destDir, 0700); err != nil {
		return fmt.Errorf("vault: failed to create destination: %w", err)
	}

	// Create temporary file for decompression
	tmpFile, err := os.CreateTemp("", "dxrk-restore-*.tar.gz")
	if err != nil {
		return fmt.Errorf("vault: failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Write compressed data
	if _, err = tmpFile.Write(compressedData); err != nil {
		tmpFile.Close()
		return fmt.Errorf("vault: failed to write temp file: %w", err)
	}
	tmpFile.Close()

	// Extract archive
	if err = extractArchive(tmpPath, destDir); err != nil {
		return fmt.Errorf("vault: extraction failed: %w", err)
	}

	return nil
}

// Helper functions

func (v *Vault) encryptBytes(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(v.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, versionAAD), nil
}

func (v *Vault) decryptBytes(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(v.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, versionAAD)
}

func deriveBackupKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)
}

func compressDirectory(srcDir string, writer io.Writer) error {
	gzWriter := gzip.NewWriter(writer)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .vault and .git directories
		if info.IsDir() && (info.Name() == ".vault" || info.Name() == ".git" || info.Name() == ".dxrk-vault") {
			return filepath.SkipDir
		}

		// Skip encrypted files
		if !info.IsDir() && filepath.Ext(path) == ".enc" {
			return nil
		}

		// Create header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return fmt.Errorf("vault: failed to create tar header: %w", err)
		}

		// Update header name to relative path
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return fmt.Errorf("vault: failed to get relative path: %w", err)
		}
		header.Name = filepath.ToSlash(relPath)

		// Write header
		if err = tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("vault: failed to write tar header: %w", err)
		}

		// Skip directories (no content to write)
		if info.IsDir() {
			return nil
		}

		// Write file content
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("vault: failed to open file %s: %w", path, err)
		}
		defer file.Close()

		if _, err := io.Copy(tarWriter, file); err != nil {
			return fmt.Errorf("vault: failed to write file %s: %w", path, err)
		}

		return nil
	})
}

func extractArchive(archivePath, destDir string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("vault: failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("vault: tar read error: %w", err)
		}

		target := filepath.Join(destDir, header.Name)

		// Security check: prevent path traversal
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(destDir)) {
			return fmt.Errorf("vault: invalid file path in archive: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("vault: failed to create directory: %w", err)
			}
		case tar.TypeReg:
			// Create parent directory
			if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
				return fmt.Errorf("vault: failed to create parent dir: %w", err)
			}

			// Create file
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("vault: failed to create file: %w", err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("vault: failed to extract file: %w", err)
			}
			outFile.Close()
		}
	}

	return nil
}

func createManifest(srcDir string, checksum string) (*DirectoryManifest, error) {
	manifest := &DirectoryManifest{
		Version:   version,
		Algorithm: "AES-256-GCM",
		CreatedAt: time.Now(),
		SourceDir: filepath.Base(srcDir),
		Checksum:  checksum,
	}

	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel(srcDir, path)
			manifest.Files = append(manifest.Files, FileEntry{
				Path:    relPath,
				Size:    info.Size(),
				ModTime: info.ModTime(),
			})
			manifest.FileCount++
			manifest.TotalSize += info.Size()
		}
		return nil
	})

	return manifest, err
}

func writeManifest(path string, manifest *DirectoryManifest) error {
	data := fmt.Sprintf(`{"version":%d,"algorithm":"%s","created_at":"%s","source_dir":"%s","file_count":%d,"total_size":%d}`,
		manifest.Version,
		manifest.Algorithm,
		manifest.CreatedAt.Format(time.RFC3339),
		manifest.SourceDir,
		manifest.FileCount,
		manifest.TotalSize,
	)
	return os.WriteFile(path, []byte(data), 0600)
}

func readManifest(path string) (*DirectoryManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	manifest := &DirectoryManifest{}
	// Simple JSON parsing - in production use encoding/json
	_, _ = fmt.Sscanf(string(data),
		`{"version":%d,"algorithm":"%s","created_at":"%s","source_dir":"%s","file_count":%d,"total_size":%d}`,
		&manifest.Version,
		&manifest.Algorithm,
		&manifest.CreatedAt,
		&manifest.SourceDir,
		&manifest.FileCount,
		&manifest.TotalSize,
	)

	return manifest, nil
}
