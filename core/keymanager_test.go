package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateKey(t *testing.T) {
	// Test cases
	testCases := []struct {
		name     string
		key      string
		expected bool
	}{
		{
			name:     "Valid key",
			key:      "test-key-1234!u",
			expected: true,
		},
		{
			name:     "Key too short",
			key:      "short-key!u",
			expected: false,
		},
		{
			name:     "Key without suffix",
			key:      "test-key-12345!",
			expected: false,
		},
		{
			name:     "Key without special char",
			key:      "test-key-12345u",
			expected: false,
		},
		{
			name:     "Empty key",
			key:      "",
			expected: false,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 使用简化的规则进行测试
			// 测试中使用的规则：最小长度14，必须以u结尾，必须包含至少1个!字符
			// 打印调试信息
			t.Logf("Testing key: %s, contains '!': %v, ends with 'u': %v, length: %d",
				tc.key, strings.Contains(tc.key, "!"), strings.HasSuffix(tc.key, "u"), len(tc.key))
			result := ValidateKeyWithRules(tc.key, 14, "", "u", "!", 1)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestStoreAndLoadAPIKey(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "keymanager_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test .env file path
	envFilePath := filepath.Join(tempDir, ".env")

	// Test cases
	testCases := []struct {
		name          string
		apiKey        string
		envName       string
		encryptionKey string
		shouldSucceed bool
	}{
		{
			name:          "Valid API key",
			apiKey:        "sk-1234567890abcdef",
			envName:       "API_KEY_TEST",
			encryptionKey: "test-key-1234!u",
			shouldSucceed: true,
		},
		{
			name:          "Empty API key",
			apiKey:        "",
			envName:       "EMPTY_KEY_TEST",
			encryptionKey: "test-key-1234!u",
			shouldSucceed: true,
		},
		{
			name:          "Invalid encryption key",
			apiKey:        "sk-1234567890abcdef",
			envName:       "INVALID_KEY_TEST",
			encryptionKey: "",  // Empty key will cause encryption to fail
			shouldSucceed: false,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Store the API key using test helper
			encValue, err := StoreAPIKeyForTest(tc.apiKey, tc.envName, tc.encryptionKey, envFilePath)
			if tc.shouldSucceed {
				if err != nil {
					t.Fatalf("StoreAPIKeyForTest failed: %v", err)
				}
				if encValue == "" {
					t.Errorf("Expected non-empty encrypted value")
				}
			} else {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			// Load the API keys using test helper
			decryptedVars, err := LoadAPIKeysForTest(tc.encryptionKey, envFilePath)
			if err != nil {
				t.Fatalf("LoadAPIKeys failed: %v", err)
			}

			// Verify that the API key was loaded correctly
			cleanName := tc.envName
			if len(cleanName) > 0 && cleanName[len(cleanName)-5:] == "_TEST" {
				cleanName = cleanName[:len(cleanName)-5]
			}
			if decryptedVars[cleanName] != tc.apiKey {
				t.Errorf("Expected %s=%s, got %s=%s", cleanName, tc.apiKey, cleanName, decryptedVars[cleanName])
			}
		})
	}
}

func TestLoadAPIKeys_InvalidKey(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "keymanager_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test .env file path
	envFilePath := filepath.Join(tempDir, ".env")

	// Store an API key using test helper
	apiKey := "sk-1234567890abcdef"
	envName := "API_KEY_TEST"
	encryptionKey := "test-key-1234!u"

	_, err = StoreAPIKeyForTest(apiKey, envName, encryptionKey, envFilePath)
	if err != nil {
		t.Fatalf("StoreAPIKeyForTest failed: %v", err)
	}

	// Try to load with an invalid key using test helper
	_, err = LoadAPIKeysForTest("wrong-key-1234!u", envFilePath)
	if err == nil {
		t.Errorf("Expected error when loading with wrong key")
	}
}

func TestLoadAPIKeys_NonExistentFile(t *testing.T) {
	// Try to load from a non-existent file
	_, err := LoadAPIKeys("test-key-12345u", "/non/existent/path/.env")
	if err == nil {
		t.Errorf("Expected error when loading from non-existent file")
	}
}
