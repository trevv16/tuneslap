package testutils

import (
	"context"
	"os"
	"time"
	"tuneslap/database"
)

const (
	TestRedisURL = "localhost:6379"
)

// SetupTestRedis initializes a test Redis connection
// Returns a cleanup function to be called after tests
func SetupTestRedis() (func(), error) {
	// Save original environment
	originalURL := os.Getenv("REDIS_URL")
	originalUsername := os.Getenv("REDIS_USERNAME")
	originalPassword := os.Getenv("REDIS_PASSWORD")

	// Set test environment
	os.Setenv("REDIS_URL", TestRedisURL)
	os.Setenv("REDIS_USERNAME", "")
	os.Setenv("REDIS_PASSWORD", "")

	// Initialize Redis
	database.InitRedis()

	// Check if Redis is available
	client := database.GetRedisClient()
	if client == nil {
		// Restore environment on error
		os.Setenv("REDIS_URL", originalURL)
		os.Setenv("REDIS_USERNAME", originalUsername)
		os.Setenv("REDIS_PASSWORD", originalPassword)
		return nil, ErrRedisUnavailable
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Ping(ctx).Err()
	if err != nil {
		// Restore environment on error
		os.Setenv("REDIS_URL", originalURL)
		os.Setenv("REDIS_USERNAME", originalUsername)
		os.Setenv("REDIS_PASSWORD", originalPassword)
		return nil, err
	}

	// Return cleanup function
	cleanup := func() {
		ClearTestCache()
		database.CloseRedis()
		// Restore original environment
		os.Setenv("REDIS_URL", originalURL)
		os.Setenv("REDIS_USERNAME", originalUsername)
		os.Setenv("REDIS_PASSWORD", originalPassword)
	}

	return cleanup, nil
}

// SetupTestRedisWithSkip sets up Redis and skips the test if unavailable
func SetupTestRedisWithSkip(t TestingT) func() {
	cleanup, err := SetupTestRedis()
	if err != nil {
		t.Skipf("Skipping test - Redis not available: %v", err)
		return func() {}
	}
	return cleanup
}

// ClearTestCache clears all test cache keys
func ClearTestCache() error {
	client := database.GetRedisClient()
	if client == nil {
		return ErrRedisUnavailable
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	// Flush all keys in the test database
	return client.FlushDB(ctx).Err()
}

// ClearCachePattern clears cache keys matching a pattern
func ClearCachePattern(pattern string) error {
	client := database.GetRedisClient()
	if client == nil {
		return ErrRedisUnavailable
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	// Use SCAN to find all keys matching the pattern
	var keys []string
	iter := client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	// Delete all found keys
	if len(keys) > 0 {
		return client.Del(ctx, keys...).Err()
	}

	return nil
}

// SetTestCache sets a test cache value
func SetTestCache(key, value string, ttl time.Duration) error {
	return database.SetCache(key, value, ttl)
}

// GetTestCache gets a test cache value
func GetTestCache(key string) (string, error) {
	return database.GetCache(key)
}

// DeleteTestCache deletes a test cache key
func DeleteTestCache(key string) error {
	return database.DeleteCache(key)
}
