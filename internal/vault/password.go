package vault

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"unicode"
)

const (
	MinPasswordLength = 12
	MaxPasswordLength = 128
)

var (
	ErrPasswordTooShort = fmt.Errorf("vault: password must be at least %d characters", MinPasswordLength)
	ErrPasswordTooLong  = fmt.Errorf("vault: password must be at most %d characters", MaxPasswordLength)
	ErrPasswordTooWeak  = fmt.Errorf("vault: password must contain uppercase, lowercase, digit, and special character")
)

// ValidatePassword checks that a password meets the security policy.
func ValidatePassword(password string) error {
	length := len([]byte(password))

	if length < MinPasswordLength {
		return ErrPasswordTooShort
	}
	if length > MaxPasswordLength {
		return ErrPasswordTooLong
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return ErrPasswordTooWeak
	}

	return nil
}

// PasswordStrength returns a score from 0 to 100 indicating password strength.
func PasswordStrength(password string) int {
	score := 0
	length := len([]byte(password))

	// Length scoring (up to 40 points)
	switch {
	case length >= 20:
		score += 40
	case length >= 16:
		score += 30
	case length >= 12:
		score += 20
	case length >= 8:
		score += 10
	default:
		score += length
	}

	// Character diversity (up to 40 points)
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	uniqueChars := make(map[rune]bool)
	for _, ch := range password {
		uniqueChars[ch] = true
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}
	if hasUpper {
		score += 10
	}
	if hasLower {
		score += 10
	}
	if hasDigit {
		score += 10
	}
	if hasSpecial {
		score += 10
	}

	// Unique character ratio (up to 20 points)
	if length > 0 {
		ratio := float64(len(uniqueChars)) / float64(length)
		score += int(ratio * 20)
	}

	if score > 100 {
		score = 100
	}
	return score
}

// GeneratePassword creates a cryptographically random password that meets the policy.
func GeneratePassword(length int) (string, error) {
	if length < MinPasswordLength {
		length = MinPasswordLength
	}
	if length > MaxPasswordLength {
		length = MaxPasswordLength
	}

	// Ensure we have all required character classes
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?")
	password := make([]rune, length)

	// Force at least one of each class
	specialChars := []rune("!@#$%^&*()-_=+?")
	forceChars := []rune{
		unicode.ToUpper(rune(randomByte())),
		unicode.ToLower(rune(randomByte())),
		rune('0' + randomByte()%10),
		specialChars[randomByte()%byte(len(specialChars))],
	}

	for i := range forceChars {
		password[i] = forceChars[i%len(forceChars)]
	}

	// Fill rest randomly
	for i := len(forceChars); i < length; i++ {
		idx, err := randomInt(len(chars))
		if err != nil {
			return "", err
		}
		password[i] = chars[idx]
	}

	// Shuffle using Fisher-Yates
	for i := length - 1; i > 0; i-- {
		j, err := randomInt(i + 1)
		if err != nil {
			return "", err
		}
		password[i], password[j] = password[j], password[i]
	}

	return string(password), nil
}

// GeneratePasswordBase64 creates a base64-encoded random password.
func GeneratePasswordBase64(length int) (string, error) {
	if length < MinPasswordLength {
		length = MinPasswordLength
	}
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

func randomByte() byte {
	b := make([]byte, 1)
	rand.Read(b)
	return b[0]
}

func randomInt(max int) (int, error) {
	if max <= 0 {
		return 0, fmt.Errorf("max must be positive")
	}
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return 0, err
	}
	val := int(b[0])<<24 | int(b[1])<<16 | int(b[2])<<8 | int(b[3])
	if val < 0 {
		val = -val
	}
	return val % max, nil
}
