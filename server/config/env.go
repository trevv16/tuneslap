package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadENV loads environment variables from .env file
// It follows the precedence: existing env vars > .env file
// This allows Docker/System to override .env values which is crucial for container orchestration
func LoadENV() error {
	// Always try to load .env file, ignoring error if it doesn't exist
	_ = godotenv.Load()
	return nil
}

// ValidateRequiredConfig checks if critical environment variables are set
func ValidateRequiredConfig() error {
	required := []string{
		"PORT",
		"JWT_SECRET",
		"DATABASE",
		"MONGODB_URI",
		"REDIS_URL",
		// Add GCS buckets to validation as they are in the user's list
		"USER_UPLOADS_BUCKET",
		"MEDIA_BUCKET",
		"GOOGLE_SERVICE_ACCOUNT_EMAIL",
		"GOOGLE_PRIVATE_KEY_PATH",
	}

	missing := []string{}
	for _, key := range required {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missing)
	}

	return nil
}

// GetMaxStorageBytes returns the maximum storage limit in bytes from MAX_STORAGE_BYTES environment variable.
// Returns -1 if not set or invalid, indicating unlimited storage.
func GetMaxStorageBytes() int64 {
	value := os.Getenv("MAX_STORAGE_BYTES")
	if value == "" {
		return -1 // -1 indicates unlimited
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed < 0 {
		return -1 // Invalid value, default to unlimited
	}

	return parsed
}

// Demo mode constants
const (
	DemoMaxFileSize   int64 = 10 * 1024 * 1024 // 10MB
	DemoMaxMediaCount int   = 5
)

// IsDemoMode returns true if the application is running in demo mode
func IsDemoMode() bool {
	return os.Getenv("DEMO_MODE") == "true"
}
