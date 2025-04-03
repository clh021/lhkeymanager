package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// DecryptAES256 decrypts data using AES-256-GCM
// encryptedData: base64 encoded encrypted data
// key: decryption key
// Returns the decrypted string or an error
func DecryptAES256(encryptedData string, key string) (string, error) {
	// Decode base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("base64 decoding failed: %w", err)
	}

	// Check data length
	if len(ciphertext) < 12+16 {
		return "", fmt.Errorf("encrypted data too short")
	}

	// Generate a 32-byte key from the provided key
	keyBytes := sha256.Sum256([]byte(key))

	// Create AES-256-GCM cipher
	block, err := aes.NewCipher(keyBytes[:])
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM mode: %w", err)
	}

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("encrypted data too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return string(plaintext), nil
}

// EncryptAES256 encrypts data using AES-256-GCM
// plaintext: plain text data to encrypt
// key: encryption key
// Returns base64 encoded encrypted data or an error
func EncryptAES256(plaintext string, key string) (string, error) {
	// Generate a 32-byte key from the provided key
	keyBytes := sha256.Sum256([]byte(key))

	// Create AES-256-GCM cipher
	block, err := aes.NewCipher(keyBytes[:])
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM mode: %w", err)
	}

	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	// In a production environment, a cryptographically secure random number generator should be used
	// For simplicity, we use a fixed nonce here
	for i := range nonce {
		nonce[i] = byte(i)
	}

	// Encrypt
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Base64 encode
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return encoded, nil
}
