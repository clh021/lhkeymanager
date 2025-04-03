package main

import (
	"os"
	"testing"
)

// TestMain is a simple test to ensure the main package compiles
func TestMain(t *testing.T) {
	// This is just a placeholder test to ensure the main package compiles
	// We can't easily test the interactive parts of the main function
	// Those are better tested through integration tests or manual testing
	
	// Test that we can create and remove a file
	tempFile, err := os.CreateTemp("", "keymanager_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	
	// Clean up
	err = os.Remove(tempPath)
	if err != nil {
		t.Fatalf("Failed to remove temp file: %v", err)
	}
}

// TestClearString tests the clearString function
func TestClearString(t *testing.T) {
	// Test cases
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "Empty string",
			input: "",
		},
		{
			name:  "Short string",
			input: "test",
		},
		{
			name:  "Long string",
			input: "this is a long string that should be cleared from memory",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Make a copy of the input
			s := tc.input
			
			// Clear the string
			clearString(&s)
			
			// Verify that the string is empty
			if s != "" {
				t.Errorf("Expected empty string, got %q", s)
			}
		})
	}

	// Test with nil pointer
	var nilStr *string
	clearString(nilStr) // This should not panic
}

// TestSecureDeleteFile tests the secureDeleteFile function
func TestSecureDeleteFile(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "keymanager_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	
	// Verify that the file exists
	_, err = os.Stat(tempPath)
	if err != nil {
		t.Fatalf("Failed to stat temp file: %v", err)
	}
	
	// Delete the file
	secureDeleteFile(tempPath)
	
	// Verify that the file no longer exists
	_, err = os.Stat(tempPath)
	if !os.IsNotExist(err) {
		t.Errorf("Expected file to be deleted, but it still exists")
	}
}
