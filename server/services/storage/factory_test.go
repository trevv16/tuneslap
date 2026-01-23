package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserUploadsStorage_MissingBucket(t *testing.T) {
	// Clear bucket env var
	originalBucket := os.Getenv("USER_UPLOADS_BUCKET")
	defer os.Setenv("USER_UPLOADS_BUCKET", originalBucket)
	os.Setenv("USER_UPLOADS_BUCKET", "")

	// Reset global state
	userUploadsStorage = nil

	_, err := GetUserUploadsStorage()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "USER_UPLOADS_BUCKET")
}

func TestGetMediaStorage_MissingBucket(t *testing.T) {
	// Clear bucket env var
	originalBucket := os.Getenv("MEDIA_BUCKET")
	defer os.Setenv("MEDIA_BUCKET", originalBucket)
	os.Setenv("MEDIA_BUCKET", "")

	// Reset global state
	mediaStorage = nil

	_, err := GetMediaStorage()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MEDIA_BUCKET")
}

func TestCreateStorage_S3_MissingCredentials(t *testing.T) {
	// Set up S3 provider without credentials
	originalProvider := os.Getenv("STORAGE_PROVIDER")
	originalAccessKey := os.Getenv("S3_ACCESS_KEY")
	originalSecretKey := os.Getenv("S3_SECRET_KEY")

	defer func() {
		os.Setenv("STORAGE_PROVIDER", originalProvider)
		os.Setenv("S3_ACCESS_KEY", originalAccessKey)
		os.Setenv("S3_SECRET_KEY", originalSecretKey)
	}()

	os.Setenv("STORAGE_PROVIDER", "s3")
	os.Setenv("S3_ACCESS_KEY", "")
	os.Setenv("S3_SECRET_KEY", "")

	// Reset global state
	userUploadsStorage = nil
	os.Setenv("USER_UPLOADS_BUCKET", "test-bucket")

	_, err := GetUserUploadsStorage()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "S3_ACCESS_KEY")
}

func TestCreateStorage_InvalidProvider(t *testing.T) {
	// Set up invalid provider
	originalProvider := os.Getenv("STORAGE_PROVIDER")
	originalBucket := os.Getenv("USER_UPLOADS_BUCKET")

	defer func() {
		os.Setenv("STORAGE_PROVIDER", originalProvider)
		os.Setenv("USER_UPLOADS_BUCKET", originalBucket)
	}()

	os.Setenv("STORAGE_PROVIDER", "invalid_provider")
	os.Setenv("USER_UPLOADS_BUCKET", "test-bucket")

	// Reset global state
	userUploadsStorage = nil

	_, err := GetUserUploadsStorage()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported storage provider")
}

func TestCreateStorage_DefaultsToS3(t *testing.T) {
	// Test that provider defaults to S3 when not set
	originalProvider := os.Getenv("STORAGE_PROVIDER")
	originalBucket := os.Getenv("USER_UPLOADS_BUCKET")
	originalAccessKey := os.Getenv("S3_ACCESS_KEY")
	originalSecretKey := os.Getenv("S3_SECRET_KEY")

	defer func() {
		os.Setenv("STORAGE_PROVIDER", originalProvider)
		os.Setenv("USER_UPLOADS_BUCKET", originalBucket)
		os.Setenv("S3_ACCESS_KEY", originalAccessKey)
		os.Setenv("S3_SECRET_KEY", originalSecretKey)
	}()

	os.Setenv("STORAGE_PROVIDER", "") // Empty = defaults to S3
	os.Setenv("USER_UPLOADS_BUCKET", "test-bucket")
	os.Setenv("S3_ACCESS_KEY", "") // Missing credentials should error

	// Reset global state
	userUploadsStorage = nil

	_, err := GetUserUploadsStorage()
	// Should fail because S3 credentials are missing, not because of invalid provider
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "S3_ACCESS_KEY")
}

func TestGetUserUploadsStorage_CachesInstance(t *testing.T) {
	originalProvider := os.Getenv("STORAGE_PROVIDER")
	originalBucket := os.Getenv("USER_UPLOADS_BUCKET")
	originalAccessKey := os.Getenv("S3_ACCESS_KEY")
	originalSecretKey := os.Getenv("S3_SECRET_KEY")
	originalEndpoint := os.Getenv("S3_ENDPOINT")

	defer func() {
		os.Setenv("STORAGE_PROVIDER", originalProvider)
		os.Setenv("USER_UPLOADS_BUCKET", originalBucket)
		os.Setenv("S3_ACCESS_KEY", originalAccessKey)
		os.Setenv("S3_SECRET_KEY", originalSecretKey)
		os.Setenv("S3_ENDPOINT", originalEndpoint)
		userUploadsStorage = nil
	}()

	os.Setenv("STORAGE_PROVIDER", "s3")
	os.Setenv("USER_UPLOADS_BUCKET", "test-bucket")
	os.Setenv("S3_ACCESS_KEY", "test-key")
	os.Setenv("S3_SECRET_KEY", "test-secret")
	os.Setenv("S3_ENDPOINT", "http://localhost:9000")

	// Reset global state
	userUploadsStorage = nil

	// First call creates storage
	storage1, err := GetUserUploadsStorage()
	assert.NoError(t, err)
	assert.NotNil(t, storage1)

	// Second call returns cached instance
	storage2, err := GetUserUploadsStorage()
	assert.NoError(t, err)
	assert.Same(t, storage1, storage2)
}

func TestGetMediaStorage_CachesInstance(t *testing.T) {
	originalProvider := os.Getenv("STORAGE_PROVIDER")
	originalBucket := os.Getenv("MEDIA_BUCKET")
	originalAccessKey := os.Getenv("S3_ACCESS_KEY")
	originalSecretKey := os.Getenv("S3_SECRET_KEY")
	originalEndpoint := os.Getenv("S3_ENDPOINT")

	defer func() {
		os.Setenv("STORAGE_PROVIDER", originalProvider)
		os.Setenv("MEDIA_BUCKET", originalBucket)
		os.Setenv("S3_ACCESS_KEY", originalAccessKey)
		os.Setenv("S3_SECRET_KEY", originalSecretKey)
		os.Setenv("S3_ENDPOINT", originalEndpoint)
		mediaStorage = nil
	}()

	os.Setenv("STORAGE_PROVIDER", "s3")
	os.Setenv("MEDIA_BUCKET", "test-media-bucket")
	os.Setenv("S3_ACCESS_KEY", "test-key")
	os.Setenv("S3_SECRET_KEY", "test-secret")
	os.Setenv("S3_ENDPOINT", "http://localhost:9000")

	// Reset global state
	mediaStorage = nil

	// First call creates storage
	storage1, err := GetMediaStorage()
	assert.NoError(t, err)
	assert.NotNil(t, storage1)

	// Second call returns cached instance
	storage2, err := GetMediaStorage()
	assert.NoError(t, err)
	assert.Same(t, storage1, storage2)
}

// Benchmark tests
func BenchmarkGetMediaKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetMediaKey("abc123def456", "audio", "test-file.mp3")
	}
}
