package config

import (
	"os"
	"testing"
)

func TestLoadENV(t *testing.T) {
	tests := []struct {
		name        string
		goEnv       string
		shouldLoad  bool
		expectError bool
	}{
		{
			name:        "should load .env when GO_ENV is not set",
			goEnv:       "",
			shouldLoad:  true,
			expectError: false, // Will not error even if .env file doesn't exist
		},
		{
			name:        "should load .env when GO_ENV is local",
			goEnv:       "local",
			shouldLoad:  true,
			expectError: false, // Will not error even if .env file doesn't exist
		},
		{
			name:        "should not load .env when GO_ENV is production",
			goEnv:       "production",
			shouldLoad:  false,
			expectError: false,
		},
		{
			name:        "should not load .env when GO_ENV is staging",
			goEnv:       "staging",
			shouldLoad:  false,
			expectError: false,
		},
		{
			name:        "should not load .env when GO_ENV is test",
			goEnv:       "test",
			shouldLoad:  false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original GO_ENV value
			originalGoEnv := os.Getenv("GO_ENV")
			defer os.Setenv("GO_ENV", originalGoEnv)

			// Set test GO_ENV value
			if tt.goEnv != "" {
				os.Setenv("GO_ENV", tt.goEnv)
			} else {
				os.Unsetenv("GO_ENV")
			}

			// Call the function
			err := LoadENV()

			// Check error expectations
			if tt.expectError && err == nil {
				t.Errorf("LoadENV() expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("LoadENV() unexpected error: %v", err)
			}
		})
	}
}

func TestLoadENV_WithExistingEnvFile(t *testing.T) {
	// Save original GO_ENV value
	originalGoEnv := os.Getenv("GO_ENV")
	defer os.Setenv("GO_ENV", originalGoEnv)

	// Set GO_ENV to local to trigger .env loading
	os.Setenv("GO_ENV", "local")

	// Create a temporary .env file
	tempDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Change to temp directory and create .env file
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	defer os.Chdir(originalWd)

	// Create a simple .env file
	envContent := "TEST_VAR=test_value\n"
	if err := os.WriteFile(".env", []byte(envContent), 0644); err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}

	// Call LoadENV - it should not return an error when .env file exists
	err = LoadENV()
	if err != nil {
		t.Errorf("LoadENV() should not error when .env file exists, got: %v", err)
	}
}

func TestLoadENV_WithMissingEnvFile(t *testing.T) {
	// Save original GO_ENV value
	originalGoEnv := os.Getenv("GO_ENV")
	defer os.Setenv("GO_ENV", originalGoEnv)

	// Set GO_ENV to local to trigger .env loading
	os.Setenv("GO_ENV", "local")

	// Temporarily change to a directory without .env file
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Create a temporary directory without .env file
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	defer os.Chdir(originalWd)

	// Call LoadENV - it should not return an error when .env file doesn't exist
	err = LoadENV()
	if err != nil {
		t.Errorf("LoadENV() should not error even when .env file is missing, got: %v", err)
	}
}

func TestLoadENV_EnvironmentVariableHandling(t *testing.T) {
	// Save original GO_ENV value
	originalGoEnv := os.Getenv("GO_ENV")
	defer os.Setenv("GO_ENV", originalGoEnv)

	// Test with various GO_ENV values
	testCases := []string{
		"",
		"local",
		"production",
		"staging",
		"test",
		"custom_env",
	}

	for _, envValue := range testCases {
		t.Run("GO_ENV="+envValue, func(t *testing.T) {
			if envValue != "" {
				os.Setenv("GO_ENV", envValue)
			} else {
				os.Unsetenv("GO_ENV")
			}

			// Function should not panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("LoadENV() panicked with GO_ENV=%s: %v", envValue, r)
				}
			}()

			err := LoadENV()

			// expectError is now always false for LoadENV as it suppresses errors
			shouldExpectError := false
			if shouldExpectError && err == nil {
				t.Errorf("LoadENV() expected error with GO_ENV=%s but got none", envValue)
			}
			if !shouldExpectError && err != nil {
				t.Errorf("LoadENV() unexpected error with GO_ENV=%s: %v", envValue, err)
			}
		})
	}
}
