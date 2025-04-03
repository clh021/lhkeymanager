package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// SaveToEnvFile saves a key-value pair to the .env file
// name: environment variable name
// value: environment variable value
// envFilePath: path to the .env file
// Returns an error if the operation fails
func SaveToEnvFile(name, value, envFilePath string) error {
	// Check if .env file exists
	var file *os.File
	var err error

	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		// File doesn't exist, create a new one
		file, err = os.Create(envFilePath)
		if err != nil {
			return err
		}
	} else {
		// File exists, open in append mode
		file, err = os.OpenFile(envFilePath, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	// Set file permissions to 600 (readable and writable only by the owner)
	err = os.Chmod(envFilePath, 0600)
	if err != nil {
		return err
	}

	// Write the environment variable
	_, err = fmt.Fprintf(file, "%s=%s\n", name, value)
	return err
}

// ReadEnvFile reads and parses the .env file
// envFilePath: path to the .env file
// Returns a map of environment variable names to values and an error if the operation fails
func ReadEnvFile(envFilePath string) (map[string]string, error) {
	// Check if .env file exists
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf(".env file does not exist")
	}

	// Open the .env file
	file, err := os.Open(envFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read and parse the file
	envVars := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse the environment variable
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Skip malformed lines
		}

		name := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		envVars[name] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envVars, nil
}

// CleanEnvVarName removes the suffix from an environment variable name
// name: environment variable name with suffix
// Returns the cleaned name
func CleanEnvVarName(name string) string {
	// If the name has a suffix pattern like "_PROD", "_DEV", etc.
	// We only want to remove the last part if it's a suffix, not any underscore
	parts := strings.Split(name, "_")
	if len(parts) > 1 && len(parts[len(parts)-1]) <= 5 {
		// Join all parts except the last one
		return strings.Join(parts[:len(parts)-1], "_")
	}
	return name
}
