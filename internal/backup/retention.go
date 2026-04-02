package backup

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RetentionPolicy defines how many backups to keep.
const DefaultRetentionCount = 5

// RetentionManager handles backup retention, deduplication, and compression.
type RetentionManager struct {
	retentionCount int
}

// NewRetentionManager creates a new retention manager with default settings.
func NewRetentionManager() *RetentionManager {
	return &RetentionManager{
		retentionCount: DefaultRetentionCount,
	}
}

// ComputeChecksum computes a SHA-256 checksum of all files in the manifest entries.
// Returns empty string if no files exist or all files don't exist.
func (m *RetentionManager) ComputeChecksum(entries []ManifestEntry) (string, error) {
	hasher := sha256.New()
	foundAny := false

	for _, entry := range entries {
		if !entry.Existed || entry.SnapshotPath == "" {
			continue
		}
		foundAny = true

		// Hash the original path to distinguish same content in different locations
		hasher.Write([]byte(entry.OriginalPath))
		hasher.Write([]byte{0}) // separator

		file, err := os.Open(entry.SnapshotPath)
		if err != nil {
			continue // skip files that can't be read
		}

		if _, err := io.Copy(hasher, file); err != nil {
			file.Close()
			continue
		}
		file.Close()
		hasher.Write([]byte{0}) // separator between files
	}

	if !foundAny {
		return "", nil
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// ShouldSkipBackup checks if the new backup can be skipped due to deduplication.
// Returns true if the checksum matches the previous backup's checksum.
func (m *RetentionManager) ShouldSkipBackup(currentChecksum, previousChecksum string) bool {
	if currentChecksum == "" || previousChecksum == "" {
		return false
	}
	return currentChecksum == previousChecksum
}

// PruneOldBackups removes backups beyond the retention count, respecting pinned backups.
// Backups are sorted newest-first; pinned backups are never pruned.
func (m *RetentionManager) PruneOldBackups(backups []Manifest, backupDir string) ([]Manifest, error) {
	if len(backups) <= m.retentionCount {
		return backups, nil
	}

	// Separate pinned and unpinned backups
	var pinned, unpinned []Manifest
	for _, b := range backups {
		if b.Pinned {
			pinned = append(pinned, b)
		} else {
			unpinned = append(unpinned, b)
		}
	}

	// Sort unpinned by date (newest first)
	sort.Slice(unpinned, func(i, j int) bool {
		return unpinned[i].CreatedAt.After(unpinned[j].CreatedAt)
	})

	// Keep retentionCount most recent unpinned backups (plus all pinned)
	maxUnpinned := m.retentionCount
	if len(unpinned) > maxUnpinned {
		toDelete := unpinned[maxUnpinned:]
		unpinned = unpinned[:maxUnpinned]

		// Delete old backups
		for _, b := range toDelete {
			if err := DeleteBackup(b); err != nil {
				// Log but continue - don't fail the whole operation
				fmt.Fprintf(os.Stderr, "WARNING: failed to delete old backup %s: %v\n", b.ID, err)
			}
		}
	}

	// Merge back and sort by date
	result := append(unpinned, pinned...)
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	return result, nil
}

// CompressBackup creates a tar.gz archive of the backup directory.
// The manifest is NOT included in the archive - it's kept separate.
func (m *RetentionManager) CompressBackup(manifest Manifest) (string, error) {
	if manifest.RootDir == "" {
		return "", fmt.Errorf("backup has no root directory")
	}

	filesDir := filepath.Join(manifest.RootDir, "files")
	if _, err := os.Stat(filesDir); os.IsNotExist(err) {
		return "", nil // No files to compress
	}

	archivePath := manifest.RootDir + ".tar.gz"
	archive, err := os.Create(archivePath)
	if err != nil {
		return "", fmt.Errorf("create archive: %w", err)
	}
	defer archive.Close()

	gzw := gzip.NewWriter(archive)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// Walk the files directory and add each file to the archive
	err = filepath.Walk(filesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get relative path from filesDir
		relPath, err := filepath.Rel(filesDir, path)
		if err != nil {
			return err
		}

		// Security: prevent path traversal
		if strings.Contains(relPath, "..") {
			return fmt.Errorf("path traversal attempt detected: %s", relPath)
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = relPath

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if _, err := io.Copy(tw, file); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("walk files: %w", err)
	}

	// Mark manifest as compressed
	manifest.Compressed = true
	manifestPath := filepath.Join(manifest.RootDir, ManifestFilename)
	if err := WriteManifest(manifestPath, manifest); err != nil {
		return "", fmt.Errorf("update manifest: %w", err)
	}

	return archivePath, nil
}

// ExtractCompressedBackup extracts a tar.gz archive to the backup directory.
func (m *RetentionManager) ExtractCompressedBackup(archivePath string, targetDir string) error {
	archive, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("open archive: %w", err)
	}
	defer archive.Close()

	gzr, err := gzip.NewReader(archive)
	if err != nil {
		return fmt.Errorf("create gzip reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	filesDir := filepath.Join(targetDir, "files")

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read tar: %w", err)
		}

		// Security: prevent path traversal
		target := filepath.Join(filesDir, header.Name)
		if !strings.HasPrefix(target, filesDir) {
			return fmt.Errorf("path traversal attempt detected: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return fmt.Errorf("create dir: %w", err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return fmt.Errorf("create parent dir: %w", err)
			}
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("create file: %w", err)
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return fmt.Errorf("copy file: %w", err)
			}
			outFile.Close()
		}
	}

	return nil
}

// PinBackup marks a backup as pinned (protected from auto-pruning).
func PinBackup(manifest *Manifest) {
	manifest.Pinned = true
	manifestPath := filepath.Join(manifest.RootDir, ManifestFilename)
	if err := WriteManifest(manifestPath, *manifest); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: failed to pin backup: %v\n", err)
	}
}

// UnpinBackup removes the pinned status from a backup.
func UnpinBackup(manifest *Manifest) {
	manifest.Pinned = false
	manifestPath := filepath.Join(manifest.RootDir, ManifestFilename)
	if err := WriteManifest(manifestPath, *manifest); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: failed to unpin backup: %v\n", err)
	}
}
