package core

import (
	"fmt"

	"github.com/clh021/lhkeymanager/utils"
)

// StoreAPIKeyForTest is a helper function for tests that bypasses the key validation
// It's only used in tests and should not be used in production code
func StoreAPIKeyForTest(apiKey, envName, encryptionKey, envFilePath string) (string, error) {
	// Check for empty key to simulate encryption failure in tests
	if encryptionKey == "" {
		return "", fmt.Errorf("encryption failed: empty key")
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

// LoadAPIKeysForTest is a helper function for tests that bypasses the key validation
// It's only used in tests and should not be used in production code
func LoadAPIKeysForTest(encryptionKey, envFilePath string) (map[string]string, error) {
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
		if value[:11] == "enc:AES256:" {
			encData := value[11:]

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
