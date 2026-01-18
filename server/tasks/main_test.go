package tasks

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitClient(t *testing.T) {
	tests := []struct {
		name        string
		redisURL    string
		redisPass   string
		expectError bool
	}{
		{
			name:        "successful initialization with valid environment",
			redisURL:    "localhost:6379",
			redisPass:   "",
			expectError: false,
		},
		{
			name:        "successful initialization with authentication",
			redisURL:    "localhost:6379",
			redisPass:   "testpass",
			expectError: false,
		},
		{
			name:        "initialization with empty REDIS_URL",
			redisURL:    "",
			redisPass:   "",
			expectError: false, // Asynq client can be created with empty URL
		},
		{
			name:        "initialization with invalid REDIS_URL",
			redisURL:    "invalid://url",
			redisPass:   "",
			expectError: false, // Asynq client creation doesn't validate connection immediately
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original environment variables
			originalRedisURL := os.Getenv("REDIS_URL")
			originalRedisPass := os.Getenv("REDIS_PASSWORD")
			defer func() {
				os.Setenv("REDIS_URL", originalRedisURL)
				os.Setenv("REDIS_PASSWORD", originalRedisPass)
			}()

			// Set test environment variables
			if tt.redisURL != "" {
				os.Setenv("REDIS_URL", tt.redisURL)
			} else {
				os.Unsetenv("REDIS_URL")
			}

			if tt.redisPass != "" {
				os.Setenv("REDIS_PASSWORD", tt.redisPass)
			} else {
				os.Unsetenv("REDIS_PASSWORD")
			}

			// Reset global variable for clean test
			Client = nil

			// Test InitClient
			assert.NotPanics(t, func() {
				InitClient()
			})

			// Verify client is created
			assert.NotNil(t, Client)

			// Test that client can be retrieved
			client := GetClient()
			assert.Equal(t, Client, client)
		})
	}
}

func TestGetClient(t *testing.T) {
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
				// Set up Asynq client
				os.Setenv("REDIS_URL", "localhost:6379")
				InitClient()
			} else {
				Client = nil
			}

			client := GetClient()

			if tt.expectNil {
				assert.Nil(t, client)
			} else {
				assert.NotNil(t, client)
				assert.Equal(t, Client, client)
			}
		})
	}
}

func TestCloseClient(t *testing.T) {
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
				// Set up Asynq client
				os.Setenv("REDIS_URL", "localhost:6379")
				InitClient()
			} else {
				Client = nil
			}

			// CloseClient should not panic
			assert.NotPanics(t, func() {
				CloseClient()
			})
		})
	}
}

func TestClientLifecycle(t *testing.T) {
	// Test the full client lifecycle
	// Save original environment
	originalRedisURL := os.Getenv("REDIS_URL")
	originalRedisPass := os.Getenv("REDIS_PASSWORD")
	defer func() {
		os.Setenv("REDIS_URL", originalRedisURL)
		os.Setenv("REDIS_PASSWORD", originalRedisPass)
	}()

	// Set up test environment
	os.Setenv("REDIS_URL", "localhost:6379")
	os.Setenv("REDIS_PASSWORD", "")

	// Reset global variable
	Client = nil

	// Test full lifecycle
	InitClient()
	assert.NotNil(t, Client)

	// Test getting client
	client := GetClient()
	assert.Equal(t, Client, client)

	// Test closing
	assert.NotPanics(t, func() {
		CloseClient()
	})
}

func TestInitClient_EnvironmentHandling(t *testing.T) {
	tests := []struct {
		name      string
		redisURL  string
		redisPass string
	}{
		{
			name:      "with REDIS_URL environment variable",
			redisURL:  "redis://localhost:6379",
			redisPass: "",
		},
		{
			name:      "without REDIS_URL environment variable",
			redisURL:  "",
			redisPass: "",
		},
		{
			name:      "with REDIS_PASSWORD environment variable",
			redisURL:  "localhost:6379",
			redisPass: "secretpass",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original environment
			originalRedisURL := os.Getenv("REDIS_URL")
			originalRedisPass := os.Getenv("REDIS_PASSWORD")
			defer func() {
				os.Setenv("REDIS_URL", originalRedisURL)
				os.Setenv("REDIS_PASSWORD", originalRedisPass)
			}()

			// Set test environment
			if tt.redisURL != "" {
				os.Setenv("REDIS_URL", tt.redisURL)
			} else {
				os.Unsetenv("REDIS_URL")
			}

			if tt.redisPass != "" {
				os.Setenv("REDIS_PASSWORD", tt.redisPass)
			} else {
				os.Unsetenv("REDIS_PASSWORD")
			}

			// Reset global variable
			Client = nil

			// Test initialization
			assert.NotPanics(t, func() {
				InitClient()
			})

			// Verify client is created
			assert.NotNil(t, Client)
		})
	}
}

func TestClientReinitialization(t *testing.T) {
	// Test that client can be reinitialized
	// Save original environment
	originalRedisURL := os.Getenv("REDIS_URL")
	defer os.Setenv("REDIS_URL", originalRedisURL)

	// Set test environment
	os.Setenv("REDIS_URL", "localhost:6379")

	// Reset global variable
	Client = nil

	// First initialization
	InitClient()
	firstClient := Client
	assert.NotNil(t, firstClient)

	// Second initialization (should create new client)
	InitClient()
	secondClient := Client
	assert.NotNil(t, secondClient)

	// Clients should be different instances
	assert.NotEqual(t, firstClient, secondClient)
}

// Benchmark tests
func BenchmarkInitClient(b *testing.B) {
	// Set up test environment
	os.Setenv("REDIS_URL", "localhost:6379")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		InitClient()
		Client = nil // Reset for next iteration
	}
}

func BenchmarkGetClient(b *testing.B) {
	// Set up Asynq client for benchmark
	os.Setenv("REDIS_URL", "localhost:6379")
	InitClient()
	defer CloseClient()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetClient()
	}
}
