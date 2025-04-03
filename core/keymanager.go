package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/clh021/lhkeymanager/utils"
)

// Security rules for encryption keys
// These values can be overridden at build time using -ldflags

// MinKeyLength is the minimum length required for encryption keys
var MinKeyLength = "16"

// KeyPrefix is the required prefix for encryption keys (empty means no prefix required)
var KeyPrefix = "lh-"

// KeySuffix is the required suffix for encryption keys (empty means no suffix required)
var KeySuffix = "u"

// RequiredChars are characters that must be present in the encryption key (empty means no specific chars required)
var RequiredChars = "!@#$%^&*"

// MinSpecialChars is the minimum number of special characters required in the encryption key
var MinSpecialChars = "2"

// KeyContain is a string that must be contained in the encryption key (empty means no specific string required)
var KeyContain = "key"

// ValidateKey validates the encryption key
// key: encryption key to validate
// Returns true if the key is valid, false otherwise
func ValidateKey(key string) bool {
	// Convert string values to integers
	minKeyLength, _ := strconv.Atoi(MinKeyLength)
	minSpecialChars, _ := strconv.Atoi(MinSpecialChars)
	return ValidateKeyWithRules(key, minKeyLength, KeyPrefix, KeySuffix, RequiredChars, minSpecialChars, KeyContain)
}

// ValidateKeyWithRules validates the encryption key with custom rules
// This function is mainly used for testing
// key: encryption key to validate
// minLength: minimum length required for the key
// prefix: required prefix for the key (empty means no prefix required)
// suffix: required suffix for the key (empty means no suffix required)
// requiredChars: characters that must be present in the key (empty means no specific chars required)
// minSpecialChars: minimum number of special characters required
// contain: a string that must be contained in the key (empty means no specific string required)
// Returns true if the key is valid, false otherwise
func ValidateKeyWithRules(key string, minLength int, prefix, suffix, requiredChars string, minSpecialChars int, contain string) bool {
	// Check minimum length
	if len(key) < minLength {
		return false
	}

	// Check prefix if required
	if prefix != "" && !strings.HasPrefix(key, prefix) {
		return false
	}

	// Check suffix if required
	if suffix != "" && !strings.HasSuffix(key, suffix) {
		return false
	}

	// Check required characters if specified
	if requiredChars != "" {
		specialCharCount := 0
		for _, char := range requiredChars {
			if strings.ContainsRune(key, char) {
				specialCharCount++
			}
		}
		if specialCharCount < minSpecialChars {
			return false
		}
	} else {
		// If no required characters are specified, don't check minSpecialChars
		// This allows users to disable special character requirements
	}

	// Check if the key contains the required string
	if contain != "" && !strings.Contains(key, contain) {
		return false
	}

	return true
}

// StoreAPIKey encrypts and stores an API key in the .env file
// apiKey: the API key to encrypt and store
// envName: the environment variable name
// encryptionKey: the key to use for encryption
// envFilePath: path to the .env file
// Returns the encrypted value and an error if the operation fails
func StoreAPIKey(apiKey, envName, encryptionKey, envFilePath string) (string, error) {
	// Validate the encryption key
	if !ValidateKey(encryptionKey) {
		return "", fmt.Errorf("invalid encryption key")
	}

	// Encrypt the API key
	encrypted, err := utils.EncryptAES256(apiKey, encryptionKey)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %w", err)
	}

	// Format the encrypted value
	encValue := fmt.Sprintf("enc:AES256:%s", encrypted)

	// Save to .env file
	err = utils.SaveToEnvFile(envName, encValue, envFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to save to .env file: %w", err)
	}

	return encValue, nil
}

// LoadAPIKeys loads and decrypts API keys from the .env file
// encryptionKey: the key to use for decryption
// envFilePath: path to the .env file
// Returns a map of environment variable names to decrypted values and an error if the operation fails
func LoadAPIKeys(encryptionKey, envFilePath string) (map[string]string, error) {
	// Validate the encryption key
	if !ValidateKey(encryptionKey) {
		return nil, fmt.Errorf("invalid encryption key")
	}

	// Read the .env file
	envVars, err := utils.ReadEnvFile(envFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read .env file: %w", err)
	}

	// Decrypt the encrypted values
	decryptedVars := make(map[string]string)
	decryptionSuccess := false

	for name, value := range envVars {
		// Check if the value is encrypted
		if strings.HasPrefix(value, "enc:AES256:") {
			encData := strings.TrimPrefix(value, "enc:AES256:")

			// Decrypt
			decrypted, err := utils.DecryptAES256(encData, encryptionKey)
			if err != nil {
				// Skip this variable if decryption fails
				continue
			}

			decryptionSuccess = true

			// Clean the environment variable name
			cleanName := utils.CleanEnvVarName(name)
			decryptedVars[cleanName] = decrypted
		} else {
			// Non-encrypted value
			decryptedVars[name] = value
		}
	}

	// If no variables were successfully decrypted, return an error
	if !decryptionSuccess {
		return nil, fmt.Errorf("no variables were successfully decrypted")
	}

	return decryptedVars, nil
}
