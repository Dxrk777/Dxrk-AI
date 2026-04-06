package vault_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/Dxrk777/Dxrk/internal/vault"
)

// TestVaultNubeBrainEncryptDecrypt tests the full Vault encryption round-trip
// simulating how NubeBrain stores and retrieves encrypted memories.
func TestVaultNubeBrainEncryptDecrypt(t *testing.T) {
	dir := t.TempDir()
	v, err := vault.New(dir, "NubeBrain!Test2026")
	if err != nil {
		t.Fatalf("vault.New failed: %v", err)
	}

	// Simulate NubeBrain storing encrypted memory snapshots
	memories := []struct {
		title   string
		content string
		typ     string
	}{
		{"user-pref", "theme: dark, language: en", "preference"},
		{"conversation", "user asked about Go generics", "chat"},
		{"knowledge", "GMP scheduler: G=goutine, M=OS thread, P=processor", "learned"},
	}

	var encryptedPayloads [][]byte
	for _, mem := range memories {
		// NubeBrain marshals to JSON then encrypts
		payload := []byte(mem.title + "|" + mem.content + "|" + mem.typ)
		encrypted, encErr := v.Encrypt(payload)
		if encErr != nil {
			t.Fatalf("Encrypt failed for '%s': %v", mem.title, encErr)
		}

		// Encrypted data must differ from plaintext
		if bytes.Equal(encrypted, payload) {
			t.Errorf("Encrypted payload for '%s' should differ from plaintext", mem.title)
		}

		encryptedPayloads = append(encryptedPayloads, encrypted)
	}

	// Decrypt and verify round-trip
	for i, mem := range memories {
		decrypted, decErr := v.Decrypt(encryptedPayloads[i])
		if decErr != nil {
			t.Fatalf("Decrypt failed for '%s': %v", mem.title, decErr)
		}

		expected := []byte(mem.title + "|" + mem.content + "|" + mem.typ)
		if !bytes.Equal(decrypted, expected) {
			t.Errorf("Round-trip mismatch for '%s': got %q, want %q", mem.title, decrypted, expected)
		}
	}
}

// TestVaultKeyRotation tests that data remains accessible after rotating the password.
func TestVaultKeyRotation(t *testing.T) {
	dir := t.TempDir()
	oldPassword := "Original!Pass2024"
	newPassword := "Rotated!Pass2025"

	// Create vault with old password and encrypt some data
	v, err := vault.New(dir, oldPassword)
	if err != nil {
		t.Fatalf("vault.New failed: %v", err)
	}

	testData := []byte("critical secret: API key sk-abc123xyz")
	encryptedBefore, err := v.Encrypt(testData)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Save encrypted file to disk (simulating NubeBrain state)
	encPath := filepath.Join(dir, ".vault", "secret.dat.enc")
	os.MkdirAll(filepath.Dir(encPath), 0700)
	if err := os.WriteFile(encPath, encryptedBefore, 0600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Rotate the key
	if err := v.KeyRotation(oldPassword, newPassword); err != nil {
		t.Fatalf("KeyRotation failed: %v", err)
	}

	// Old in-memory ciphertext should no longer decrypt (key changed)
	_, err = v.Decrypt(encryptedBefore)
	if err == nil {
		t.Error("Old ciphertext should NOT decrypt with new key after rotation")
	}

	// New encryption with rotated key should work
	newPlaintext := []byte("new data after rotation")
	encryptedAfter, err := v.Encrypt(newPlaintext)
	if err != nil {
		t.Fatalf("Encrypt after rotation failed: %v", err)
	}

	decryptedAfter, err := v.Decrypt(encryptedAfter)
	if err != nil {
		t.Fatalf("Decrypt after rotation failed: %v", err)
	}
	if !bytes.Equal(decryptedAfter, newPlaintext) {
		t.Errorf("Mismatch after rotation: got %q, want %q", decryptedAfter, newPlaintext)
	}

	// Verify wrong old password now fails to open vault
	dir2 := t.TempDir()
	// Copy the rotated key file to a new dir to test re-opening
	keyData, _ := os.ReadFile(filepath.Join(dir, ".vault", "master.key"))
	os.MkdirAll(filepath.Join(dir2, ".vault"), 0700)
	os.WriteFile(filepath.Join(dir2, ".vault", "master.key"), keyData, 0600)

	_, err = vault.New(dir2, oldPassword)
	if err == nil {
		t.Error("Old password should fail after key rotation")
	}

	v2, err := vault.New(dir2, newPassword)
	if err != nil {
		t.Fatalf("New password should work after rotation: %v", err)
	}
	_ = v2
}

// TestVaultChangePasswordReencrypt tests ChangePassword with encrypted files on disk.
func TestVaultChangePasswordReencrypt(t *testing.T) {
	dir := t.TempDir()
	oldPass := "OldPass!2024abc"
	newPass := "NewPass!2025xyz"

	v, err := vault.New(dir, oldPass)
	if err != nil {
		t.Fatalf("vault.New failed: %v", err)
	}

	// Create multiple encrypted files in the .vault directory
	testFiles := map[string][]byte{
		"config.enc":  []byte(`{"theme":"dark","lang":"en"}`),
		"memory.enc":  []byte("remember: user likes blue"),
		"session.enc": []byte("session-data-abc-123"),
	}

	vaultDir := filepath.Join(dir, ".vault")
	for name, plaintext := range testFiles {
		encrypted, encErr := v.Encrypt(plaintext)
		if encErr != nil {
			t.Fatalf("Encrypt failed for %s: %v", name, encErr)
		}
		if err := os.WriteFile(filepath.Join(vaultDir, name), encrypted, 0600); err != nil {
			t.Fatalf("WriteFile failed for %s: %v", name, err)
		}
	}

	// Change password (triggers re-encryption of all .enc files)
	if err := v.ChangePassword(oldPass, newPass); err != nil {
		t.Fatalf("ChangePassword failed: %v", err)
	}

	// All files should now be decryptable with the new key
	for name, expectedPlaintext := range testFiles {
		encrypted, readErr := os.ReadFile(filepath.Join(vaultDir, name))
		if readErr != nil {
			t.Fatalf("ReadFile failed for %s: %v", name, readErr)
		}
		decrypted, decErr := v.Decrypt(encrypted)
		if decErr != nil {
			t.Fatalf("Decrypt failed for %s after re-encryption: %v", name, decErr)
		}
		if !bytes.Equal(decrypted, expectedPlaintext) {
			t.Errorf("Re-encrypted %s mismatch: got %q, want %q", name, decrypted, expectedPlaintext)
		}
	}

	// Old password should no longer work for new vault instance
	dir2 := t.TempDir()
	keyData, _ := os.ReadFile(filepath.Join(dir, ".vault", "master.key"))
	os.MkdirAll(filepath.Join(dir2, ".vault"), 0700)
	os.WriteFile(filepath.Join(dir2, ".vault", "master.key"), keyData, 0600)

	_, err = vault.New(dir2, oldPass)
	if err == nil {
		t.Error("Old password should fail after ChangePassword")
	}

	// New password should work
	v2, err := vault.New(dir2, newPass)
	if err != nil {
		t.Fatalf("New password should work: %v", err)
	}

	// Verify new vault can encrypt/decrypt
	testMsg := []byte("post-change-password-data")
	enc, _ := v2.Encrypt(testMsg)
	dec, _ := v2.Decrypt(enc)
	if !bytes.Equal(dec, testMsg) {
		t.Error("New vault should encrypt/decrypt correctly")
	}
}

// TestVaultFileEncryptionIntegration tests EncryptFile/DecryptFile round-trip
// simulating how the system encrypts configuration files.
func TestVaultFileEncryptionIntegration(t *testing.T) {
	dir := t.TempDir()
	v, err := vault.New(dir, "FileEncrypt!Test26")
	if err != nil {
		t.Fatalf("vault.New failed: %v", err)
	}

	// Create a plaintext config file
	configContent := []byte(`{
  "api_key": "sk-secret-key-12345",
  "endpoint": "https://api.example.com",
  "model": "gpt-4"
}`)
	configPath := filepath.Join(dir, "config.json")
	if err := os.WriteFile(configPath, configContent, 0600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Encrypt the file
	if err := v.EncryptFile(configPath); err != nil {
		t.Fatalf("EncryptFile failed: %v", err)
	}

	encPath := configPath + ".enc"
	if _, err := os.Stat(encPath); os.IsNotExist(err) {
		t.Fatal("Encrypted file should exist")
	}

	// Encrypted file should NOT contain plaintext
	encBytes, _ := os.ReadFile(encPath)
	if bytes.Contains(encBytes, []byte("api_key")) {
		t.Error("Encrypted file should not contain plaintext 'api_key'")
	}

	// Decrypt and verify
	decrypted, err := v.DecryptFile(encPath)
	if err != nil {
		t.Fatalf("DecryptFile failed: %v", err)
	}
	if !bytes.Equal(decrypted, configContent) {
		t.Errorf("Decrypted content mismatch:\n  got:  %q\n  want: %q", decrypted, configContent)
	}
}

// TestVaultConcurrentEncryptDecrypt verifies thread safety under concurrent access.
func TestVaultConcurrentEncryptDecrypt(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrency test in short mode")
	}

	dir := t.TempDir()
	v, err := vault.New(dir, "Concurrent!Test26")
	if err != nil {
		t.Fatalf("vault.New failed: %v", err)
	}

	const goroutines = 50
	done := make(chan error, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			plaintext := []byte("concurrent data " + string(rune(id)))
			encrypted, encErr := v.Encrypt(plaintext)
			if encErr != nil {
				done <- encErr
				return
			}
			decrypted, decErr := v.Decrypt(encrypted)
			if decErr != nil {
				done <- decErr
				return
			}
			if !bytes.Equal(decrypted, plaintext) {
				done <- bytes.ErrTooLarge // any error
				return
			}
			done <- nil
		}(i)
	}

	for i := 0; i < goroutines; i++ {
		if err := <-done; err != nil {
			t.Errorf("Goroutine %d failed: %v", i, err)
		}
	}
}

// TestVaultWrongPasswordPersistence verifies that reopening a vault with the
// wrong password consistently fails.
func TestVaultWrongPasswordPersistence(t *testing.T) {
	dir := t.TempDir()
	correctPassword := "Correct!Horse2026"

	// Create vault with correct password
	v1, err := vault.New(dir, correctPassword)
	if err != nil {
		t.Fatalf("vault.New failed: %v", err)
	}

	// Encrypt some data
	plaintext := []byte("sensitive payload")
	encrypted, _ := v1.Encrypt(plaintext)

	// Try opening with various wrong passwords
	wrongPasswords := []string{
		"wrong-password",
		"Correct-Horse-Battery-Staple", // case sensitive
		"correct-horse-battery-stapl",  // one char off
		"",                             // empty
		correctPassword + " ",          // trailing space
	}

	for _, wrongPass := range wrongPasswords {
		_, err := vault.New(dir, wrongPass)
		if err == nil {
			t.Errorf("Should fail with wrong password %q", wrongPass)
		}
	}

	// Reopen with correct password should work
	v2, err := vault.New(dir, correctPassword)
	if err != nil {
		t.Fatalf("Reopen with correct password failed: %v", err)
	}

	decrypted, err := v2.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt after reopen failed: %v", err)
	}
	if !bytes.Equal(decrypted, plaintext) {
		t.Error("Decrypted data should match original after reopen")
	}
}
