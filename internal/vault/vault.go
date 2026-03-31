// Package vault implementa el sistema de encriptación de Dxrk v1.0
//
// El Vault proporciona almacenamiento seguro para datos sensibles:
// - Encriptación AES-256-GCM para datos
// - Derivación de claves con PBKDF2 (100,000 iteraciones)
// - Verificación de integridad con HMAC-SHA256
// - Almacenamiento seguro de claves maestras
//
// Seguridad:
// - Claves derivadas con PBKDF2 (100,000 iteraciones)
// - Nonce aleatorio para cada encriptación
// - Verificación de contraseña en tiempo constante
// - Limpieza segura de memoria (best effort)
//
// Uso:
// - New(): Crea o carga un vault existente
// - Encrypt(): Encripta datos con la clave maestra
// - Decrypt(): Desencripta datos con la clave maestra
// - ChangePassword(): Cambia la contraseña del vault
package vault

import (
	"crypto/aes"      // Para encriptación AES
	"crypto/cipher"   // Para modos de encriptación (GCM)
	"crypto/hmac"     // Para verificación de integridad
	"crypto/rand"     // Para generación de números aleatorios
	"crypto/sha256"   // Para hashing SHA-256
	"encoding/base64" // Para codificación de claves
	"errors"          // Para errores personalizados
	"fmt"             // Para formateo de errores
	"io"              // Para lectura de datos aleatorios
	"os"              // Para operaciones de sistema de archivos
	"path/filepath"   // Para manipulación de rutas
	"sync"            // Para sincronización de concurrencia

	"golang.org/x/crypto/pbkdf2" // Para derivación de claves
)

// ErrWrongPassword se retorna cuando la contraseña proporcionada no coincide
// con la clave del vault almacenada. Usar errors.Is() para distinguir
// de otros errores del vault.
var ErrWrongPassword = errors.New("vault: wrong password")

const (
	KeyLen    = 32     // AES-256 - longitud de la clave en bytes
	SaltLen   = 16     // Longitud del salt en bytes
	IterCount = 100000 // Iteraciones de PBKDF2 para derivación de claves
	version   = 1      // Versión del formato del vault
)

// versionAAD es el Additional Associated Data para AES-GCM
// Se usa para vincular la versión del formato con los datos encriptados
var versionAAD = []byte{version}

// Vault es el sistema de encriptación de Dxrk
//
// Proporciona almacenamiento seguro para datos sensibles usando:
// - AES-256-GCM para encriptación autenticada
// - PBKDF2 para derivación de claves
// - HMAC-SHA256 para verificación de integridad
//
// El vault es thread-safe y puede ser usado concurrentemente.
type Vault struct {
	mu      sync.RWMutex // Mutex para operaciones concurrentes
	key     []byte       // Clave maestra derivada (32 bytes para AES-256)
	keyPath string       // Ruta al archivo de clave maestra
	dataDir string       // Directorio de datos del vault
}

// New crea o carga un vault existente
//
// Si el archivo de clave maestra existe, carga y verifica la contraseña.
// Si no existe, crea un nuevo vault con la contraseña proporcionada.
//
// La contraseña debe tener al menos 8 caracteres y cumplir con la política
// de contraseñas (ver password.go).
//
// Retorna error si:
// - La contraseña no cumple la política
// - Falla la derivación de la clave
// - Falla la verificación de la contraseña existente
func New(dataDir, password string) (*Vault, error) {
	v := &Vault{
		dataDir: dataDir,
		keyPath: filepath.Join(dataDir, ".vault", "master.key"),
	}

	if err := os.MkdirAll(filepath.Dir(v.keyPath), 0700); err != nil {
		return nil, fmt.Errorf("vault: failed to create dir: %w", err)
	}

	if _, err := os.Stat(v.keyPath); err == nil {
		if err := v.loadKey(password); err != nil {
			return nil, fmt.Errorf("vault: %w", err)
		}
	} else {
		if err := v.initKey(password); err != nil {
			return nil, fmt.Errorf("vault: init failed: %w", err)
		}
	}

	return v, nil
}

func (v *Vault) initKey(password string) error {
	if err := ValidatePassword(password); err != nil {
		return err
	}

	salt := make([]byte, SaltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return err
	}

	key := pbkdf2.Key([]byte(password), salt, IterCount, KeyLen, sha256.New)

	mac := hmac.New(sha256.New, key)
	mac.Write(salt)
	verification := mac.Sum(nil)

	data := append(salt, verification...)

	if err := os.WriteFile(v.keyPath, data, 0600); err != nil {
		return err
	}

	v.key = key
	return nil
}

func (v *Vault) loadKey(password string) error {
	data, err := os.ReadFile(v.keyPath)
	if err != nil {
		return err
	}
	if len(data) < SaltLen+32 {
		return fmt.Errorf("corrupted key file")
	}

	salt := data[:SaltLen]
	storedMAC := data[SaltLen : SaltLen+32]

	key := pbkdf2.Key([]byte(password), salt, IterCount, KeyLen, sha256.New)

	mac := hmac.New(sha256.New, key)
	mac.Write(salt)
	computedMAC := mac.Sum(nil)

	if !hmac.Equal(storedMAC, computedMAC) {
		return ErrWrongPassword
	}

	v.key = key
	return nil
}

func (v *Vault) Encrypt(plaintext []byte) ([]byte, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

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

	return gcm.Seal(nonce, nonce, plaintext, versionAAD), nil
}

func (v *Vault) Decrypt(ciphertext []byte) ([]byte, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

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

func (v *Vault) EncryptString(plaintext string) (string, error) {
	encrypted, err := v.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (v *Vault) DecryptString(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	decrypted, err := v.Decrypt(data)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

func (v *Vault) EncryptFile(path string) error {
	plaintext, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	ciphertext, err := v.Encrypt(plaintext)
	if err != nil {
		return err
	}
	return os.WriteFile(path+".enc", ciphertext, 0600)
}

func (v *Vault) DecryptFile(path string) ([]byte, error) {
	ciphertext, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return v.Decrypt(ciphertext)
}

func (v *Vault) SecureDelete(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		data := make([]byte, info.Size())
		if _, err := rand.Read(data); err != nil {
			return fmt.Errorf("secure delete: rand.Read failed: %w", err)
		}
		if err := os.WriteFile(path, data, 0600); err != nil {
			return fmt.Errorf("secure delete: overwrite pass %d failed: %w", i+1, err)
		}
	}

	return os.Remove(path)
}

func (v *Vault) WipeAll() error {
	v.mu.Lock()
	defer v.mu.Unlock()

	vaultDir := filepath.Join(v.dataDir, ".vault")
	var errs []error
	err := filepath.Walk(vaultDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if err := v.SecureDelete(path); err != nil {
			errs = append(errs, fmt.Errorf("wipe %s: %w", path, err))
		}
		return nil
	})
	if err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return fmt.Errorf("wipe errors: %v", errs)
	}
	return nil
}

func (v *Vault) reencryptAll(oldKey, newKey []byte) error {
	dataDir := filepath.Join(v.dataDir, ".vault")
	return filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(path) != ".enc" {
			return nil
		}

		ciphertext, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reencrypt read %s: %w", path, err)
		}

		plaintext, err := decryptWithKey(oldKey, ciphertext)
		if err != nil {
			return fmt.Errorf("reencrypt decrypt %s: %w", path, err)
		}

		newCiphertext, err := encryptWithKey(newKey, plaintext)
		if err != nil {
			return fmt.Errorf("reencrypt encrypt %s: %w", path, err)
		}

		if err := os.WriteFile(path, newCiphertext, 0600); err != nil {
			return fmt.Errorf("reencrypt write %s: %w", path, err)
		}

		return nil
	})
}

func encryptWithKey(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
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
	return gcm.Seal(nonce, nonce, plaintext, versionAAD), nil
}

func decryptWithKey(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
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

func (v *Vault) ChangePassword(oldPassword, newPassword string) error {
	if err := ValidatePassword(newPassword); err != nil {
		return fmt.Errorf("new password: %w", err)
	}

	if err := v.loadKey(oldPassword); err != nil {
		return fmt.Errorf("wrong current password")
	}

	oldKey := make([]byte, len(v.key))
	copy(oldKey, v.key)

	salt := make([]byte, SaltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return err
	}

	newKey := pbkdf2.Key([]byte(newPassword), salt, IterCount, KeyLen, sha256.New)

	if err := v.reencryptAll(oldKey, newKey); err != nil {
		return fmt.Errorf("re-encryption failed: %w", err)
	}

	mac := hmac.New(sha256.New, newKey)
	mac.Write(salt)
	verification := mac.Sum(nil)
	data := append(salt, verification...)

	if err := os.WriteFile(v.keyPath, data, 0600); err != nil {
		return err
	}

	v.key = newKey
	return nil
}

func (v *Vault) KeyRotation(oldPassword, newPassword string) error {
	if err := v.loadKey(oldPassword); err != nil {
		return fmt.Errorf("wrong current password")
	}

	oldKey := make([]byte, len(v.key))
	copy(oldKey, v.key)

	salt := make([]byte, SaltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return err
	}

	newKey := pbkdf2.Key([]byte(newPassword), salt, IterCount, KeyLen, sha256.New)

	dataDir := filepath.Join(v.dataDir, ".vault")
	var errs []error
	err := filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(path) != ".enc" {
			return nil
		}

		ciphertext, readErr := os.ReadFile(path)
		if readErr != nil {
			errs = append(errs, fmt.Errorf("read %s: %w", path, readErr))
			return nil
		}

		plaintext, decErr := decryptWithKey(oldKey, ciphertext)
		if decErr != nil {
			errs = append(errs, fmt.Errorf("decrypt %s: %w", path, decErr))
			return nil
		}

		newCiphertext, encErr := encryptWithKey(newKey, plaintext)
		if encErr != nil {
			errs = append(errs, fmt.Errorf("encrypt %s: %w", path, encErr))
			return nil
		}

		if writeErr := os.WriteFile(path, newCiphertext, 0600); writeErr != nil {
			errs = append(errs, fmt.Errorf("write %s: %w", path, writeErr))
			return nil
		}

		return nil
	})
	if err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return fmt.Errorf("key rotation failed: %v", errs)
	}

	mac := hmac.New(sha256.New, newKey)
	mac.Write(salt)
	verification := mac.Sum(nil)
	keyData := append(salt, verification...)

	if err := os.WriteFile(v.keyPath, keyData, 0600); err != nil {
		return fmt.Errorf("key rotation: failed to update key file: %w", err)
	}

	v.key = newKey
	return nil
}

func (v *Vault) Status() map[string]interface{} {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return map[string]interface{}{
		"initialized": v.key != nil,
		"key_path":    v.keyPath,
		"algorithm":   "AES-256-GCM",
		"kdf":         "PBKDF2-SHA256",
		"iterations":  IterCount,
	}
}
