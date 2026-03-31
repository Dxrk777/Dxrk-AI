package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPrivateVault(t *testing.T) {
	tmpDir := t.TempDir()
	pv := NewPrivateVault(tmpDir)

	t.Run("Unlock first time", func(t *testing.T) {
		if err := pv.Unlock("DXRK-0mega-K1ng-777x-9f3a"); err != nil {
			t.Fatalf("unlock failed: %v", err)
		}
		if !pv.IsUnlocked() {
			t.Error("vault should be unlocked")
		}
	})

	t.Run("Encrypt and Decrypt", func(t *testing.T) {
		plaintext := []byte("super secret data that only the owner can read")
		encrypted, err := pv.Encrypt(plaintext)
		if err != nil {
			t.Fatalf("encrypt failed: %v", err)
		}
		if encrypted == string(plaintext) {
			t.Error("encrypted should differ from plaintext")
		}

		decrypted, err := pv.Decrypt(encrypted)
		if err != nil {
			t.Fatalf("decrypt failed: %v", err)
		}
		if string(decrypted) != string(plaintext) {
			t.Error("decrypted should match plaintext")
		}
	})

	t.Run("Put and Get", func(t *testing.T) {
		data := map[string]string{"key": "value", "secret": "42"}
		if err := pv.Put("project", "test-entry", data); err != nil {
			t.Fatalf("put failed: %v", err)
		}

		var result map[string]string
		if err := pv.Get("project", "test-entry", &result); err != nil {
			t.Fatalf("get failed: %v", err)
		}
		if result["key"] != "value" {
			t.Error("retrieved data should match")
		}
	})

	t.Run("Wrong key fails", func(t *testing.T) {
		pv2 := NewPrivateVault(tmpDir)
		err := pv2.Unlock("wrong-password")
		if err == nil {
			t.Error("should fail with wrong key")
		}
	})

	t.Run("Lock clears key", func(t *testing.T) {
		pv.Lock()
		if pv.IsUnlocked() {
			t.Error("vault should be locked")
		}

		_, err := pv.Encrypt([]byte("test"))
		if err == nil {
			t.Error("encrypt should fail when locked")
		}
	})

	t.Run("Key file exists", func(t *testing.T) {
		hashPath := filepath.Join(tmpDir, ".dxrk-vault", ".keyhash")
		if _, err := os.Stat(hashPath); os.IsNotExist(err) {
			t.Error("key hash file should exist after first unlock")
		}
	})
}
