package vault_test

import (
	"os"
	"testing"

	"github.com/Dxrk777/Dxrk-Hex/internal/vault"
)

const testPassword = "TestPass!2345"

func TestNewVault(t *testing.T) {
	dir := t.TempDir()
	v, err := vault.New(dir, testPassword)
	if err != nil {
		t.Fatalf("New vault failed: %v", err)
	}
	if v == nil {
		t.Fatal("Vault should not be nil")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	dir := t.TempDir()
	v, _ := vault.New(dir, testPassword)

	plaintext := []byte("Hello, Dxrk! This is secret data.")
	encrypted, err := v.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if string(encrypted) == string(plaintext) {
		t.Error("Encrypted should differ from plaintext")
	}

	decrypted, err := v.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("Decrypted mismatch: got '%s', want '%s'", decrypted, plaintext)
	}
}

func TestEncryptDecryptString(t *testing.T) {
	dir := t.TempDir()
	v, _ := vault.New(dir, testPassword)

	original := "my secret string 12345"
	encoded, err := v.EncryptString(original)
	if err != nil {
		t.Fatalf("EncryptString failed: %v", err)
	}

	decoded, err := v.DecryptString(encoded)
	if err != nil {
		t.Fatalf("DecryptString failed: %v", err)
	}

	if decoded != original {
		t.Errorf("Mismatch: got '%s', want '%s'", decoded, original)
	}
}

func TestEncryptDecryptFile(t *testing.T) {
	dir := t.TempDir()
	v, _ := vault.New(dir, testPassword)

	// Create a file
	content := []byte("file content to encrypt")
	filePath := dir + "/test.txt"
	os.WriteFile(filePath, content, 0644)

	// Encrypt
	if err := v.EncryptFile(filePath); err != nil {
		t.Fatalf("EncryptFile failed: %v", err)
	}

	// Encrypted file should exist
	if _, err := os.Stat(filePath + ".enc"); os.IsNotExist(err) {
		t.Error("Encrypted file should exist")
	}

	// Decrypt
	decrypted, err := v.DecryptFile(filePath + ".enc")
	if err != nil {
		t.Fatalf("DecryptFile failed: %v", err)
	}

	if string(decrypted) != string(content) {
		t.Error("Decrypted content should match original")
	}
}

func TestWrongPassword(t *testing.T) {
	dir := t.TempDir()
	_, _ = vault.New(dir, testPassword)

	_, err := vault.New(dir, "WrongPass!9999")
	if err == nil {
		t.Error("Should fail with wrong password")
	}
}

func TestSecureDelete(t *testing.T) {
	dir := t.TempDir()
	v, _ := vault.New(dir, testPassword)

	filePath := dir + "/to-delete.txt"
	os.WriteFile(filePath, []byte("sensitive data"), 0600)

	if err := v.SecureDelete(filePath); err != nil {
		t.Fatalf("SecureDelete failed: %v", err)
	}

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("File should be deleted")
	}
}

func TestStatus(t *testing.T) {
	dir := t.TempDir()
	v, _ := vault.New(dir, testPassword)

	status := v.Status()
	if status["algorithm"] != "AES-256-GCM" {
		t.Errorf("Expected AES-256-GCM, got %v", status["algorithm"])
	}
	if status["initialized"] != true {
		t.Error("Should be initialized")
	}
}

func TestEmptyPlaintext(t *testing.T) {
	dir := t.TempDir()
	v, _ := vault.New(dir, testPassword)

	encrypted, err := v.Encrypt([]byte(""))
	if err != nil {
		t.Fatalf("Encrypt empty failed: %v", err)
	}

	decrypted, err := v.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt empty failed: %v", err)
	}

	if len(decrypted) != 0 {
		t.Error("Decrypted empty should be empty")
	}
}

func TestLargeData(t *testing.T) {
	dir := t.TempDir()
	v, _ := vault.New(dir, testPassword)

	// 1MB of data
	large := make([]byte, 1024*1024)
	for i := range large {
		large[i] = byte(i % 256)
	}

	encrypted, err := v.Encrypt(large)
	if err != nil {
		t.Fatalf("Encrypt large failed: %v", err)
	}

	decrypted, err := v.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt large failed: %v", err)
	}

	if len(decrypted) != len(large) {
		t.Errorf("Size mismatch: %d vs %d", len(decrypted), len(large))
	}
}

func TestWeakPasswordRejected(t *testing.T) {
	weakPasswords := []string{
		"",               // empty
		"short",          // too short
		"alllowercase1!", // no uppercase
		"ALLUPPERCASE1!", // no lowercase
		"NoDigitsHere!!", // no digits
		"NoSpecial1234",  // no special
	}

	for _, pw := range weakPasswords {
		dir := t.TempDir()
		_, err := vault.New(dir, pw)
		if err == nil {
			t.Errorf("Weak password %q should be rejected", pw)
		}
	}
}

func TestPasswordStrength(t *testing.T) {
	tests := []struct {
		password string
		minScore int
	}{
		{"a", 0},
		{testPassword, 50},
		{"MyStr0ng!Pass#Word2026", 80},
	}

	for _, tt := range tests {
		score := vault.PasswordStrength(tt.password)
		if score < tt.minScore {
			t.Errorf("Password %q: strength %d, want at least %d", tt.password, score, tt.minScore)
		}
	}
}

func TestGeneratePassword(t *testing.T) {
	pw, err := vault.GeneratePassword(16)
	if err != nil {
		t.Fatalf("GeneratePassword failed: %v", err)
	}
	if err := vault.ValidatePassword(pw); err != nil {
		t.Errorf("Generated password should be valid: %v", err)
	}
}
