package database

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartMongoDB(t *testing.T) {
	tests := []struct {
		name        string
		mongoURI    string
		database    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "successful connection with valid environment",
			mongoURI:    "mongodb://localhost:27017",
			database:    "testdb",
			expectError: false,
		},
		{
			name:        "missing MONGODB_URI",
			mongoURI:    "",
			database:    "testdb",
			expectError: true,
			errorMsg:    "you must set your 'MONGODB_URI' environmental variable",
		},
		{
			name:        "missing DATABASE",
			mongoURI:    "mongodb://localhost:27017",
			database:    "",
			expectError: true,
			errorMsg:    "you must set your 'DATABASE' environmental variable",
		},
		{
			name:        "both environment variables missing",
			mongoURI:    "",
			database:    "",
			expectError: true,
			errorMsg:    "you must set your 'MONGODB_URI' environmental variable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original environment variables
			originalMongoURI := os.Getenv("MONGODB_URI")
			originalDatabase := os.Getenv("DATABASE")
			defer func() {
				os.Setenv("MONGODB_URI", originalMongoURI)
				os.Setenv("DATABASE", originalDatabase)
			}()

			// Set test environment variables
			if tt.mongoURI != "" {
				os.Setenv("MONGODB_URI", tt.mongoURI)
			} else {
				os.Unsetenv("MONGODB_URI")
			}

			if tt.database != "" {
				os.Setenv("DATABASE", tt.database)
			} else {
				os.Unsetenv("DATABASE")
			}

			// Reset global variables for clean test
			mongoClient = nil
			dbName = ""

			// Test StartMongoDB
			err := StartMongoDB()

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Nil(t, mongoClient)
			} else {
				// Note: This will only work if MongoDB is actually running
				// In a real test environment, you might want to mock this
				if err == nil {
					assert.NotNil(t, mongoClient)
					assert.Equal(t, tt.database, dbName)
				} else {
					// If MongoDB is not running, that's expected in test environment
					t.Logf("MongoDB connection failed (expected if MongoDB not running): %v", err)
				}
			}
		})
	}
}

func TestStartMongoDB_InvalidURI(t *testing.T) {
	// Save original environment variables
	originalMongoURI := os.Getenv("MONGODB_URI")
	originalDatabase := os.Getenv("DATABASE")
	defer func() {
		os.Setenv("MONGODB_URI", originalMongoURI)
		os.Setenv("DATABASE", originalDatabase)
	}()

	// Set invalid MongoDB URI
	os.Setenv("MONGODB_URI", "invalid://uri")
	os.Setenv("DATABASE", "testdb")

	// Reset global variables
	mongoClient = nil
	dbName = ""

	// This should return an error due to invalid URI
	err := StartMongoDB()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect to MongoDB")
}

func TestCloseMongoDB(t *testing.T) {
	tests := []struct {
		name        string
		setupClient bool
		expectPanic bool
	}{
		{
			name:        "close with nil client",
			setupClient: false,
			expectPanic: true,
		},
		{
			name:        "close with valid client",
			setupClient: true,
			expectPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupClient {
				// Set up a valid client first
				os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
				os.Setenv("DATABASE", "testdb")

				err := StartMongoDB()
				if err != nil {
					t.Skipf("Skipping test - MongoDB not available: %v", err)
				}
			} else {
				mongoClient = nil
			}

			if tt.expectPanic {
				assert.Panics(t, func() {
					CloseMongoDB()
				})
			} else {
				assert.NotPanics(t, func() {
					CloseMongoDB()
				})
			}
		})
	}
}

func TestGetCollection(t *testing.T) {
	tests := []struct {
		name        string
		collection  string
		setupClient bool
		expectPanic bool
	}{
		{
			name:        "get collection with nil client",
			collection:  "users",
			setupClient: false,
			expectPanic: true,
		},
		{
			name:        "get collection with valid client",
			collection:  "users",
			setupClient: true,
			expectPanic: false,
		},
		{
			name:        "get collection with empty name",
			collection:  "",
			setupClient: true,
			expectPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupClient {
				// Set up a valid client first
				os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
				os.Setenv("DATABASE", "testdb")

				err := StartMongoDB()
				if err != nil {
					t.Skipf("Skipping test - MongoDB not available: %v", err)
				}
			} else {
				mongoClient = nil
			}

			if tt.expectPanic {
				assert.Panics(t, func() {
					GetCollection(tt.collection)
				})
			} else {
				assert.NotPanics(t, func() {
					collection := GetCollection(tt.collection)
					assert.NotNil(t, collection)
					assert.Equal(t, tt.collection, collection.Name())
				})
			}
		})
	}
}

func TestGetContext(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
	}{
		{
			name:    "context with 1 second timeout",
			timeout: 1 * time.Second,
		},
		{
			name:    "context with 10 seconds timeout",
			timeout: 10 * time.Second,
		},
		{
			name:    "context with 0 timeout",
			timeout: 0,
		},
		{
			name:    "context with negative timeout",
			timeout: -1 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := GetContext(tt.timeout)

			// Verify context is not nil
			assert.NotNil(t, ctx)

			// Verify cancel function is not nil
			assert.NotNil(t, cancel)

			// Test that context can be used
			select {
			case <-ctx.Done():
				// Context was cancelled or timed out
				if tt.timeout > 0 {
					// This is expected for positive timeouts
					t.Logf("Context timed out as expected after %v", tt.timeout)
				}
			default:
				// Context is still active
				if tt.timeout <= 0 {
					// This is expected for zero or negative timeouts
					t.Logf("Context is active as expected for timeout %v", tt.timeout)
				}
			}

			// Always call cancel to clean up
			cancel()

			// Verify context is cancelled after calling cancel
			select {
			case <-ctx.Done():
				// Expected
			default:
				t.Error("Context should be cancelled after calling cancel function")
			}
		})
	}
}

func TestGetContext_WithTimeout(t *testing.T) {
	// Test that context actually times out
	timeout := 100 * time.Millisecond
	ctx, cancel := GetContext(timeout)
	defer cancel()

	start := time.Now()
	<-ctx.Done()
	duration := time.Since(start)

	// Verify timeout occurred within reasonable bounds
	assert.True(t, duration >= timeout, "Context should timeout after at least %v, but timed out after %v", timeout, duration)
	assert.True(t, duration <= timeout+50*time.Millisecond, "Context should timeout within reasonable bounds, but took %v", duration)
}

func TestDatabaseIntegration(t *testing.T) {
	// Integration test that tests the full flow
	// This test requires MongoDB to be running

	// Save original environment
	originalMongoURI := os.Getenv("MONGODB_URI")
	originalDatabase := os.Getenv("DATABASE")
	defer func() {
		os.Setenv("MONGODB_URI", originalMongoURI)
		os.Setenv("DATABASE", originalDatabase)
	}()

	// Set up test environment
	os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	os.Setenv("DATABASE", "testdb")

	// Reset global variables
	mongoClient = nil
	dbName = ""

	// Test full flow
	err := StartMongoDB()
	if err != nil {
		t.Skipf("Skipping integration test - MongoDB not available: %v", err)
	}

	// Verify client is set up
	assert.NotNil(t, mongoClient)
	assert.Equal(t, "testdb", dbName)

	// Test getting a collection
	collection := GetCollection("test_collection")
	assert.NotNil(t, collection)
	assert.Equal(t, "test_collection", collection.Name())

	// Test getting context
	ctx, cancel := GetContext(5 * time.Second)
	assert.NotNil(t, ctx)
	assert.NotNil(t, cancel)
	cancel()

	// Test closing
	assert.NotPanics(t, func() {
		CloseMongoDB()
	})
}

// Benchmark tests
func BenchmarkGetContext(b *testing.B) {
	timeout := 1 * time.Second

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx, cancel := GetContext(timeout)
		cancel()
		_ = ctx
	}
}

func BenchmarkGetCollection(b *testing.B) {
	// Set up MongoDB connection for benchmark
	os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	os.Setenv("DATABASE", "testdb")

	err := StartMongoDB()
	if err != nil {
		b.Skipf("Skipping benchmark - MongoDB not available: %v", err)
	}
	defer CloseMongoDB()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collection := GetCollection("benchmark_collection")
		_ = collection
	}
}
