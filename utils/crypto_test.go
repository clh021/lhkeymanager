package utils

import (
	"testing"
)

func TestEncryptDecryptAES256(t *testing.T) {
	// Test cases
	testCases := []struct {
		name      string
		plaintext string
		key       string
	}{
		{
			name:      "Simple text",
			plaintext: "hello world",
			key:       "test-key-12345u",
		},
		{
			name:      "Empty text",
			plaintext: "",
			key:       "test-key-12345u",
		},
		{
			name:      "Special characters",
			plaintext: "!@#$%^&*()_+{}|:<>?~`-=[]\\;',./",
			key:       "test-key-12345u",
		},
		{
			name:      "Long text",
			plaintext: "This is a very long text that should be encrypted and then decrypted back to the original text without any issues. Let's make it even longer to ensure that the encryption and decryption functions can handle large amounts of data.",
			key:       "test-key-12345u",
		},
		{
			name:      "API key format",
			plaintext: "sk-1234567890abcdefghijklmnopqrstuvwxyz",
			key:       "test-key-12345u",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Encrypt
			encrypted, err := EncryptAES256(tc.plaintext, tc.key)
			if err != nil {
				t.Fatalf("EncryptAES256 failed: %v", err)
			}

			// Decrypt
			decrypted, err := DecryptAES256(encrypted, tc.key)
			if err != nil {
				t.Fatalf("DecryptAES256 failed: %v", err)
			}

			// Verify
			if decrypted != tc.plaintext {
				t.Errorf("Expected %q, got %q", tc.plaintext, decrypted)
			}
		})
	}
}

func TestDecryptAES256_InvalidInput(t *testing.T) {
	// Test cases for invalid input
	testCases := []struct {
		name          string
		encryptedData string
		key           string
		expectError   bool
	}{
		{
			name:          "Invalid base64",
			encryptedData: "not-base64-data",
			key:           "test-key-12345u",
			expectError:   true,
		},
		{
			name:          "Data too short",
			encryptedData: "AAEC", // Valid base64 but too short
			key:           "test-key-12345u",
			expectError:   true,
		},
		{
			name:          "Wrong key",
			encryptedData: "AAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8=", // Some valid base64 data
			key:           "wrong-key-12345u",
			expectError:   true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DecryptAES256(tc.encryptedData, tc.key)
			if (err != nil) != tc.expectError {
				t.Errorf("Expected error: %v, got: %v", tc.expectError, err != nil)
			}
		})
	}
}

func TestEncryptAES256_DifferentKeys(t *testing.T) {
	plaintext := "test-data"
	key1 := "test-key-12345u"
	key2 := "different-key-12345u"

	// Encrypt with key1
	encrypted1, err := EncryptAES256(plaintext, key1)
	if err != nil {
		t.Fatalf("EncryptAES256 failed: %v", err)
	}

	// Encrypt with key2
	encrypted2, err := EncryptAES256(plaintext, key2)
	if err != nil {
		t.Fatalf("EncryptAES256 failed: %v", err)
	}

	// Verify that the encrypted data is different
	if encrypted1 == encrypted2 {
		t.Errorf("Expected different encrypted data for different keys")
	}

	// Try to decrypt with wrong key
	_, err = DecryptAES256(encrypted1, key2)
	if err == nil {
		t.Errorf("Expected error when decrypting with wrong key")
	}
}

func TestEncryptAES256_Deterministic(t *testing.T) {
	// Since we're using a fixed nonce, encryption should be deterministic
	plaintext := "test-data"
	key := "test-key-12345u"

	// Encrypt twice with the same key
	encrypted1, err := EncryptAES256(plaintext, key)
	if err != nil {
		t.Fatalf("EncryptAES256 failed: %v", err)
	}

	encrypted2, err := EncryptAES256(plaintext, key)
	if err != nil {
		t.Fatalf("EncryptAES256 failed: %v", err)
	}

	// Verify that the encrypted data is the same
	if encrypted1 != encrypted2 {
		t.Errorf("Expected same encrypted data for same plaintext and key")
	}
}
