package database

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitRedis(t *testing.T) {
	tests := []struct {
		name        string
		redisURL    string
		username    string
		password    string
		expectError bool
	}{
		{
			name:        "successful initialization with valid environment",
			redisURL:    "localhost:6379",
			username:    "",
			password:    "",
			expectError: false,
		},
		{
			name:        "successful initialization with authentication",
			redisURL:    "localhost:6379",
			username:    "testuser",
			password:    "testpass",
			expectError: false,
		},
		{
			name:        "initialization with empty REDIS_URL",
			redisURL:    "",
			username:    "",
			password:    "",
			expectError: false, // Redis client can be created with empty URL
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original environment variables
			originalRedisURL := os.Getenv("REDIS_URL")
			originalUsername := os.Getenv("REDIS_USERNAME")
			originalPassword := os.Getenv("REDIS_PASSWORD")
			defer func() {
				os.Setenv("REDIS_URL", originalRedisURL)
				os.Setenv("REDIS_USERNAME", originalUsername)
				os.Setenv("REDIS_PASSWORD", originalPassword)
			}()

			// Set test environment variables
			if tt.redisURL != "" {
				os.Setenv("REDIS_URL", tt.redisURL)
			} else {
				os.Unsetenv("REDIS_URL")
			}

			if tt.username != "" {
				os.Setenv("REDIS_USERNAME", tt.username)
			} else {
				os.Unsetenv("REDIS_USERNAME")
			}

			if tt.password != "" {
				os.Setenv("REDIS_PASSWORD", tt.password)
			} else {
				os.Unsetenv("REDIS_PASSWORD")
			}

			// Reset global variable for clean test
			RedisClient = nil

			// Test InitRedis
			assert.NotPanics(t, func() {
				InitRedis()
			})

			// Verify client is created
			assert.NotNil(t, RedisClient)

			// Test that client can be retrieved
			client := GetRedisClient()
			assert.Equal(t, RedisClient, client)
		})
	}
}

func TestGetRedisClient(t *testing.T) {
	tests := []struct {
		name        string
		setupClient bool
		expectNil   bool
	}{
		{
			name:        "get client when not initialized",
			setupClient: false,
			expectNil:   true,
		},
		{
			name:        "get client when initialized",
			setupClient: true,
			expectNil:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupClient {
				// Set up Redis client
				os.Setenv("REDIS_URL", "localhost:6379")
				InitRedis()
			} else {
				RedisClient = nil
			}

			client := GetRedisClient()

			if tt.expectNil {
				assert.Nil(t, client)
			} else {
				assert.NotNil(t, client)
				assert.Equal(t, RedisClient, client)
			}
		})
	}
}

func TestCloseRedis(t *testing.T) {
	tests := []struct {
		name        string
		setupClient bool
	}{
		{
			name:        "close with nil client",
			setupClient: false,
		},
		{
			name:        "close with valid client",
			setupClient: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupClient {
				// Set up Redis client
				os.Setenv("REDIS_URL", "localhost:6379")
				InitRedis()
			} else {
				RedisClient = nil
			}

			// CloseRedis should not panic
			assert.NotPanics(t, func() {
				CloseRedis()
			})
		})
	}
}

func TestSetCache(t *testing.T) {
	tests := []struct {
		name        string
		setupClient bool
		key         string
		value       interface{}
		ttl         time.Duration
		expectError bool
	}{
		{
			name:        "set cache with valid client",
			setupClient: true,
			key:         "test_key",
			value:       "test_value",
			ttl:         1 * time.Minute,
			expectError: false,
		},
		{
			name:        "set cache with nil client",
			setupClient: false,
			key:         "test_key",
			value:       "test_value",
			ttl:         1 * time.Minute,
			expectError: true,
		},
		{
			name:        "set cache with empty key",
			setupClient: true,
			key:         "",
			value:       "test_value",
			ttl:         1 * time.Minute,
			expectError: false, // Redis allows empty keys
		},
		{
			name:        "set cache with complex value",
			setupClient: true,
			key:         "complex_key",
			value:       map[string]interface{}{"nested": "value"},
			ttl:         1 * time.Minute,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupClient {
				// Set up Redis client
				os.Setenv("REDIS_URL", "localhost:6379")
				InitRedis()
			} else {
				RedisClient = nil
			}

			err := SetCache(tt.key, tt.value, tt.ttl)

			if tt.expectError {
				assert.Error(t, err)
				if RedisClient == nil {
					assert.EqualError(t, err, "Redis client is not initialized")
				}
			} else {
				// Note: This will only work if Redis is actually running
				// In a test environment, connection errors are expected
				if err != nil {
					t.Logf("SetCache failed (expected if Redis not running): %v", err)
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestGetCache(t *testing.T) {
	tests := []struct {
		name        string
		setupClient bool
		key         string
		expectError bool
	}{
		{
			name:        "get cache with valid client",
			setupClient: true,
			key:         "test_key",
			expectError: false,
		},
		{
			name:        "get cache with nil client",
			setupClient: false,
			key:         "test_key",
			expectError: true,
		},
		{
			name:        "get cache with empty key",
			setupClient: true,
			key:         "",
			expectError: false, // Redis allows empty keys
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupClient {
				// Set up Redis client
				os.Setenv("REDIS_URL", "localhost:6379")
				InitRedis()
			} else {
				RedisClient = nil
			}

			value, err := GetCache(tt.key)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, value)
				if RedisClient == nil {
					assert.EqualError(t, err, "Redis client is not initialized")
				}
			} else {
				// Note: This will only work if Redis is actually running
				// In a test environment, connection errors are expected
				if err != nil {
					t.Logf("GetCache failed (expected if Redis not running): %v", err)
				} else {
					// If no error, value should be retrieved (might be empty if key doesn't exist)
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestDeleteCache(t *testing.T) {
	tests := []struct {
		name        string
		setupClient bool
		key         string
		expectError bool
	}{
		{
			name:        "delete cache with valid client",
			setupClient: true,
			key:         "test_key",
			expectError: false,
		},
		{
			name:        "delete cache with nil client",
			setupClient: false,
			key:         "test_key",
			expectError: true,
		},
		{
			name:        "delete cache with empty key",
			setupClient: true,
			key:         "",
			expectError: false, // Redis allows empty keys
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupClient {
				// Set up Redis client
				os.Setenv("REDIS_URL", "localhost:6379")
				InitRedis()
			} else {
				RedisClient = nil
			}

			err := DeleteCache(tt.key)

			if tt.expectError {
				assert.Error(t, err)
				if RedisClient == nil {
					assert.EqualError(t, err, "Redis client is not initialized")
				}
			} else {
				// Note: This will only work if Redis is actually running
				// In a test environment, connection errors are expected
				if err != nil {
					t.Logf("DeleteCache failed (expected if Redis not running): %v", err)
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestRedisIntegration(t *testing.T) {
	// Integration test that tests the full Redis flow
	// This test requires Redis to be running

	// Save original environment
	originalRedisURL := os.Getenv("REDIS_URL")
	originalUsername := os.Getenv("REDIS_USERNAME")
	originalPassword := os.Getenv("REDIS_PASSWORD")
	defer func() {
		os.Setenv("REDIS_URL", originalRedisURL)
		os.Setenv("REDIS_USERNAME", originalUsername)
		os.Setenv("REDIS_PASSWORD", originalPassword)
	}()

	// Set up test environment
	os.Setenv("REDIS_URL", "localhost:6379")
	os.Setenv("REDIS_USERNAME", "")
	os.Setenv("REDIS_PASSWORD", "")

	// Reset global variable
	RedisClient = nil

	// Test full flow
	InitRedis()
	assert.NotNil(t, RedisClient)

	// Test getting client
	client := GetRedisClient()
	assert.Equal(t, RedisClient, client)

	// Test cache operations
	testKey := "integration_test_key"
	testValue := "integration_test_value"
	ttl := 1 * time.Minute

	// Set cache
	err := SetCache(testKey, testValue, ttl)
	if err != nil {
		t.Skipf("Skipping integration test - Redis not available: %v", err)
	}

	// Get cache
	retrievedValue, err := GetCache(testKey)
	assert.NoError(t, err)
	assert.Equal(t, testValue, retrievedValue)

	// Delete cache
	err = DeleteCache(testKey)
	assert.NoError(t, err)

	// Verify deletion
	_, err = GetCache(testKey)
	// Should return error or empty value since key was deleted
	if err == nil {
		// If no error, value should be empty
		assert.Empty(t, retrievedValue)
	}

	// Test closing
	assert.NotPanics(t, func() {
		CloseRedis()
	})
}

func TestRedisCacheOperations(t *testing.T) {
	// Test various cache operations with different data types
	// This test requires Redis to be running

	// Set up Redis client
	os.Setenv("REDIS_URL", "localhost:6379")
	InitRedis()
	defer CloseRedis()

	// Test different data types
	testCases := []struct {
		name  string
		key   string
		value interface{}
	}{
		{
			name:  "string value",
			key:   "string_key",
			value: "hello world",
		},
		{
			name:  "number value",
			key:   "number_key",
			value: 42,
		},
		{
			name:  "boolean value",
			key:   "bool_key",
			value: true,
		},
		{
			name:  "json-like string",
			key:   "json_key",
			value: `{"name": "test", "value": 123}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set cache
			err := SetCache(tc.key, tc.value, 1*time.Minute)
			if err != nil {
				t.Skipf("Skipping test - Redis not available: %v", err)
			}

			// Get cache
			retrievedValue, err := GetCache(tc.key)
			assert.NoError(t, err)
			assert.NotEmpty(t, retrievedValue)

			// Clean up
			err = DeleteCache(tc.key)
			assert.NoError(t, err)
		})
	}
}

func TestRedisTTL(t *testing.T) {
	// Test TTL functionality
	// This test requires Redis to be running

	// Set up Redis client
	os.Setenv("REDIS_URL", "localhost:6379")
	InitRedis()
	defer CloseRedis()

	testKey := "ttl_test_key"
	testValue := "ttl_test_value"
	shortTTL := 100 * time.Millisecond

	// Set cache with short TTL
	err := SetCache(testKey, testValue, shortTTL)
	if err != nil {
		t.Skipf("Skipping TTL test - Redis not available: %v", err)
	}

	// Immediately get the value
	retrievedValue, err := GetCache(testKey)
	assert.NoError(t, err)
	assert.Equal(t, testValue, retrievedValue)

	// Wait for TTL to expire
	time.Sleep(shortTTL + 50*time.Millisecond)

	// Try to get the value again - should be expired
	_, err = GetCache(testKey)
	// Should return error or empty value since key expired
	if err == nil {
		// If no error, value should be empty
		assert.Empty(t, retrievedValue)
	}
}

// Benchmark tests
func BenchmarkInitRedis(b *testing.B) {
	// Set up test environment
	os.Setenv("REDIS_URL", "localhost:6379")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		InitRedis()
		RedisClient = nil // Reset for next iteration
	}
}

func BenchmarkSetCache(b *testing.B) {
	// Set up Redis client for benchmark
	os.Setenv("REDIS_URL", "localhost:6379")
	InitRedis()
	defer CloseRedis()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("benchmark_key_%d", i)
		SetCache(key, "benchmark_value", 1*time.Minute)
	}
}

func BenchmarkGetCache(b *testing.B) {
	// Set up Redis client for benchmark
	os.Setenv("REDIS_URL", "localhost:6379")
	InitRedis()
	defer CloseRedis()

	// Pre-populate cache
	key := "benchmark_get_key"
	SetCache(key, "benchmark_get_value", 1*time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetCache(key)
	}
}
