package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveToEnvFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "env_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test .env file path
	envFilePath := filepath.Join(tempDir, ".env")

	// Test cases
	testCases := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "Simple variable",
			key:   "TEST_VAR",
			value: "test_value",
		},
		{
			name:  "Variable with special characters",
			key:   "API_KEY",
			value: "sk-1234!@#$%^&*()",
		},
		{
			name:  "Encrypted variable",
			key:   "ENCRYPTED_VAR",
			value: "enc:AES256:ABCDEFG",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Save to .env file
			err := SaveToEnvFile(tc.key, tc.value, envFilePath)
			if err != nil {
				t.Fatalf("SaveToEnvFile failed: %v", err)
			}

			// Read the .env file
			envVars, err := ReadEnvFile(envFilePath)
			if err != nil {
				t.Fatalf("ReadEnvFile failed: %v", err)
			}

			// Verify that the variable was saved correctly
			if envVars[tc.key] != tc.value {
				t.Errorf("Expected %s=%s, got %s=%s", tc.key, tc.value, tc.key, envVars[tc.key])
			}
		})
	}

	// Test file permissions
	info, err := os.Stat(envFilePath)
	if err != nil {
		t.Fatalf("Failed to stat .env file: %v", err)
	}

	// Check that the file permissions are 0600 (readable and writable only by the owner)
	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected file permissions 0600, got %o", info.Mode().Perm())
	}
}

func TestReadEnvFile_NonExistent(t *testing.T) {
	// Try to read a non-existent .env file
	_, err := ReadEnvFile("/non/existent/path/.env")
	if err == nil {
		t.Errorf("Expected error when reading non-existent .env file")
	}
}

func TestCleanEnvVarName(t *testing.T) {
	// Test cases
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No suffix",
			input:    "API",
			expected: "API",
		},
		{
			name:     "With suffix",
			input:    "API_KEY_PROD",
			expected: "API_KEY",
		},
		{
			name:     "Multiple underscores",
			input:    "OPENAI_API_KEY_PROD",
			expected: "OPENAI_API_KEY",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CleanEnvVarName(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}
