package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

// PrivateVault enforces single-user access with AES-256-GCM encryption.
// Only the person who knows the master key can decrypt the data.
type PrivateVault struct {
	mu       sync.RWMutex
	key      []byte
	dataDir  string
	unlocked bool
}

// VaultEntry is an encrypted record.
type VaultEntry struct {
	ID        string    `json:"id"`
	Namespace string    `json:"namespace"`
	Data      string    `json:"data"` // base64 encrypted
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewPrivateVault creates a vault that requires a master key to unlock.
func NewPrivateVault(dataDir string) *PrivateVault {
	return &PrivateVault{
		dataDir: filepath.Join(dataDir, ".dxrk-vault"),
	}
}

// deriveKey derives a 256-bit key from the master password using PBKDF2-SHA256.
func deriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)
}

// Unlock unlocks the vault with the master key. Fails if wrong.
func (pv *PrivateVault) Unlock(masterKey string) error {
	pv.mu.Lock()
	defer pv.mu.Unlock()

	// Verify key against stored hash
	hashPath := filepath.Join(pv.dataDir, ".keyhash")
	if _, err := os.Stat(hashPath); os.IsNotExist(err) {
		// First time — generate salt, derive key, store salt+hash
		os.MkdirAll(pv.dataDir, 0700)
		salt := make([]byte, 16)
		if _, err := rand.Read(salt); err != nil {
			return fmt.Errorf("vault: salt generation failed: %w", err)
		}
		key := deriveKey(masterKey, salt)
		mac := hmac.New(sha256.New, key)
		mac.Write(salt)
		verification := mac.Sum(nil)
		// Store: salt (16) + HMAC (32)
		data := append(salt, verification...)
		if err := os.WriteFile(hashPath, data, 0600); err != nil {
			return fmt.Errorf("vault: failed to write key hash: %w", err)
		}
		pv.key = key
		pv.unlocked = true
		return nil
	}

	// Read stored salt + verification
	storedData, err := os.ReadFile(hashPath)
	if err != nil {
		return fmt.Errorf("vault: cannot read key hash: %w", err)
	}
	if len(storedData) < 48 {
		return fmt.Errorf("vault: corrupted key hash file")
	}

	salt := storedData[:16]
	storedMAC := storedData[16:48]

	key := deriveKey(masterKey, salt)

	mac := hmac.New(sha256.New, key)
	mac.Write(salt)
	computedMAC := mac.Sum(nil)

	if subtle.ConstantTimeCompare(storedMAC, computedMAC) != 1 {
		return fmt.Errorf("vault: ACCESS DENIED — wrong master key")
	}

	pv.key = key
	pv.unlocked = true
	return nil
}

// Lock locks the vault, clearing the key from memory.
func (pv *PrivateVault) Lock() {
	pv.mu.Lock()
	defer pv.mu.Unlock()
	// Overwrite key in memory
	for i := range pv.key {
		pv.key[i] = 0
	}
	pv.key = nil
	pv.unlocked = false
}

// IsUnlocked returns true if the vault is unlocked.
func (pv *PrivateVault) IsUnlocked() bool {
	pv.mu.RLock()
	defer pv.mu.RUnlock()
	return pv.unlocked
}

// Encrypt encrypts data using AES-256-GCM.
func (pv *PrivateVault) Encrypt(plaintext []byte) (string, error) {
	pv.mu.RLock()
	defer pv.mu.RUnlock()

	if !pv.unlocked {
		return "", fmt.Errorf("vault: locked — unlock first")
	}

	block, err := aes.NewCipher(pv.key)
	if err != nil {
		return "", fmt.Errorf("vault: cipher creation: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("vault: GCM creation: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("vault: nonce generation: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts AES-256-GCM encrypted data.
func (pv *PrivateVault) Decrypt(encrypted string) ([]byte, error) {
	pv.mu.RLock()
	defer pv.mu.RUnlock()

	if !pv.unlocked {
		return nil, fmt.Errorf("vault: locked — unlock first")
	}

	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, fmt.Errorf("vault: base64 decode: %w", err)
	}

	block, err := aes.NewCipher(pv.key)
	if err != nil {
		return nil, fmt.Errorf("vault: cipher creation: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("vault: GCM creation: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("vault: ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("vault: decryption failed (wrong key?): %w", err)
	}

	return plaintext, nil
}

// Put stores an encrypted entry in the vault.
func (pv *PrivateVault) Put(namespace, id string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("vault: marshaling data: %w", err)
	}

	encrypted, err := pv.Encrypt(jsonData)
	if err != nil {
		return err
	}

	entry := VaultEntry{
		ID:        id,
		Namespace: namespace,
		Data:      encrypted,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	nsDir := filepath.Join(pv.dataDir, namespace)
	os.MkdirAll(nsDir, 0700)

	entryPath := filepath.Join(nsDir, id+".enc")
	entryJSON, _ := json.Marshal(entry)
	return os.WriteFile(entryPath, entryJSON, 0600)
}

// Get retrieves and decrypts an entry from the vault.
func (pv *PrivateVault) Get(namespace, id string, result interface{}) error {
	entryPath := filepath.Join(pv.dataDir, namespace, id+".enc")
	entryJSON, err := os.ReadFile(entryPath)
	if err != nil {
		return fmt.Errorf("vault: entry not found: %s/%s", namespace, id)
	}

	var entry VaultEntry
	if err := json.Unmarshal(entryJSON, &entry); err != nil {
		return fmt.Errorf("vault: parsing entry: %w", err)
	}

	plaintext, err := pv.Decrypt(entry.Data)
	if err != nil {
		return err
	}

	return json.Unmarshal(plaintext, result)
}

// ListNamespaces returns all namespace directories.
func (pv *PrivateVault) ListNamespaces() ([]string, error) {
	entries, err := os.ReadDir(pv.dataDir)
	if err != nil {
		return nil, err
	}

	var namespaces []string
	for _, e := range entries {
		if e.IsDir() && e.Name()[0] != '.' {
			namespaces = append(namespaces, e.Name())
		}
	}
	return namespaces, nil
}

// ListEntries returns all entry IDs in a namespace.
func (pv *PrivateVault) ListEntries(namespace string) ([]string, error) {
	nsDir := filepath.Join(pv.dataDir, namespace)
	entries, err := os.ReadDir(nsDir)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".enc" {
			ids = append(ids, e.Name()[:len(e.Name())-4])
		}
	}
	return ids, nil
}

// Delete removes an entry from the vault.
func (pv *PrivateVault) Delete(namespace, id string) error {
	entryPath := filepath.Join(pv.dataDir, namespace, id+".enc")
	return os.Remove(entryPath)
}
