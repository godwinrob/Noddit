package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
	"crypto/sha1"
)

const (
	// Iterations for PBKDF2 (matching Java's 100,000)
	iterations = 100000
	// Salt length in bytes
	saltLength = 128
	// Key length in bytes (128 bits = 16 bytes)
	keyLength = 16
)

// HashPassword creates a PBKDF2 hash of the password with a random salt
// Returns the base64-encoded hash and salt
func HashPassword(password string) (hash string, salt string, err error) {
	// Generate random salt
	saltBytes := make([]byte, saltLength)
	if _, err := rand.Read(saltBytes); err != nil {
		return "", "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Hash password using PBKDF2WithHmacSHA1 (matching Java implementation)
	hashBytes := pbkdf2.Key([]byte(password), saltBytes, iterations, keyLength, sha1.New)

	// Encode to base64
	hashStr := base64.StdEncoding.EncodeToString(hashBytes)
	saltStr := base64.StdEncoding.EncodeToString(saltBytes)

	return hashStr, saltStr, nil
}

// VerifyPassword verifies a password against a hash and salt
func VerifyPassword(password, hash, salt string) (bool, error) {
	// Decode salt from base64
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return false, fmt.Errorf("invalid salt encoding: %w", err)
	}

	// Decode expected hash from base64
	expectedHash, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return false, fmt.Errorf("invalid hash encoding: %w", err)
	}

	// Hash the provided password with the same salt
	actualHash := pbkdf2.Key([]byte(password), saltBytes, iterations, keyLength, sha1.New)

	// Constant-time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare(expectedHash, actualHash) == 1, nil
}
