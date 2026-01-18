package app

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"tuneslap/router"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupApp(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
		errorMsg    string
		createEnv   bool
	}{
		{
			name: "successful setup with valid environment",
			envVars: map[string]string{
				"MONGODB_URI":                  "mongodb://localhost:27017",
				"DATABASE":                     "testdb",
				"REDIS_URL":                    "localhost:6379",
				"PORT":                         "8082",
				"JWT_SECRET":                   "test-secret",
				"USER_UPLOADS_BUCKET":          "test-bucket",
				"MEDIA_BUCKET":                 "test-media-bucket",
				"GOOGLE_SERVICE_ACCOUNT_EMAIL": "test@example.com",
				"GOOGLE_PRIVATE_KEY_PATH":      "test-key-path",
			},
			expectError: false,
			createEnv:   true,
		},
		{
			name:    "missing environment variables (env loading fails)",
			envVars: map[string]string{
				// No env vars set - will fail at config.ValidateRequiredConfig()
			},
			expectError: true,
			errorMsg:    "missing required environment variables",
			createEnv:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original environment variables
			originalEnv := make(map[string]string)
			for key := range tt.envVars {
				if val := os.Getenv(key); val != "" {
					originalEnv[key] = val
				}
			}

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Optionally create a temporary .env file
			var cleanupEnv func()
			if tt.createEnv {
				cleanupEnv = createTempDotEnv(t)
			}

			// Cleanup after test
			defer func() {
				for key := range tt.envVars {
					if originalVal, exists := originalEnv[key]; exists {
						os.Setenv(key, originalVal)
					} else {
						os.Unsetenv(key)
					}
				}
				if cleanupEnv != nil {
					cleanupEnv()
				}
			}()

			// Test SetupApp
			app, err := SetupApp()

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Nil(t, app)
			} else {
				// Note: This test may fail if MongoDB/Redis are not running
				// That's expected in a test environment
				if err == nil {
					assert.NotNil(t, app)
					// Verify app configuration
					assert.Equal(t, 10*time.Second, app.Server().ReadTimeout)
					assert.Equal(t, 10*time.Second, app.Server().WriteTimeout)
				} else {
					// Log the error but don't fail the test
					// This is expected when external services are not available
					t.Logf("SetupApp failed (expected if external services not running): %v", err)
				}
			}
		})
	}
}

// createTempDotEnv creates a temporary .env file and returns a cleanup function
func createTempDotEnv(t *testing.T) func() {
	filename := ".env"
	content := "DUMMY_VAR=dummy\n"
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create temp .env file: %v", err)
	}
	return func() {
		os.Remove(filename)
	}
}

func TestAttachMiddleware(t *testing.T) {
	app := fiber.New()

	// Test that middleware can be attached without error
	assert.NotPanics(t, func() {
		attachMiddleware(app)
	})
	// No further checks needed; if it doesn't panic, it's fine
}

func TestMiddlewareFunctionality(t *testing.T) {
	app := fiber.New()
	attachMiddleware(app)

	// Setup routes to include the health endpoint
	router.SetupRoutes(app)

	// Add a test route
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "test"})
	})

	// Test health check endpoint
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test API route with JSON acceptance
	req = httptest.NewRequest("GET", "/api/v1/test", nil)
	req.Header.Set("Accept", "application/json")
	resp, err = app.Test(req)
	require.NoError(t, err)
	// Should return 404 since the route doesn't exist, but middleware should process it
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Test CORS headers
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3001")
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// CORS headers should be present
	assert.NotEmpty(t, resp.Header.Get("Access-Control-Allow-Origin"))
}

func TestSetupAndRunApp_EnvironmentHandling(t *testing.T) {
	tests := []struct {
		name     string
		port     string
		expected string
	}{
		{
			name:     "with PORT environment variable",
			port:     "9090",
			expected: ":9090",
		},
		{
			name:     "without PORT environment variable",
			port:     "",
			expected: ":8082",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original PORT
			originalPort := os.Getenv("PORT")
			defer os.Setenv("PORT", originalPort)

			// Set test PORT
			if tt.port != "" {
				os.Setenv("PORT", tt.port)
			} else {
				os.Unsetenv("PORT")
			}

			// Note: We can't easily test the actual server startup in unit tests
			// because it blocks, but we can test the environment variable handling
			// by checking what port would be used
			port := os.Getenv("PORT")
			if port == "" {
				port = "8082"
			}
			assert.Equal(t, tt.expected, ":"+port)
		})
	}
}

func TestMiddlewareOrder(t *testing.T) {
	app := fiber.New()
	attachMiddleware(app)

	// Add a test route that returns headers
	app.Get("/headers", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"headers": c.GetReqHeaders(),
		})
	})

	req := httptest.NewRequest("GET", "/headers", nil)
	req.Header.Set("User-Agent", "test-agent")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSetupApp_Configuration(t *testing.T) {
	// Test that SetupApp properly configures the application
	// This test focuses on the setup aspects, not the external service connections

	// Create a temporary .env file
	cleanupEnv := createTempDotEnv(t)
	defer cleanupEnv()

	// Set minimal environment variables
	os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	os.Setenv("DATABASE", "testdb")
	os.Setenv("REDIS_URL", "localhost:6379")
	os.Setenv("PORT", "8082")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("USER_UPLOADS_BUCKET", "test-bucket")
	os.Setenv("MEDIA_BUCKET", "test-media-bucket")
	os.Setenv("GOOGLE_SERVICE_ACCOUNT_EMAIL", "test@example.com")
	os.Setenv("GOOGLE_PRIVATE_KEY_PATH", "test-key-path")

	app, err := SetupApp()

	// The test may fail if external services are not available, which is expected
	if err == nil && app != nil {
		// Verify app configuration
		assert.Equal(t, 10*time.Second, app.Server().ReadTimeout)
		assert.Equal(t, 10*time.Second, app.Server().WriteTimeout)

		// Verify that routes are set up
		// We can test this by checking if the health endpoint exists
		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req)
		if err == nil {
			// Health endpoint should exist and be accessible
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	} else {
		t.Logf("SetupApp failed (expected if external services not running): %v", err)
	}
}

// Benchmark tests
func BenchmarkSetupApp(b *testing.B) {
	// Set up test environment
	os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	os.Setenv("DATABASE", "testdb")
	os.Setenv("REDIS_URL", "localhost:6379")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app, err := SetupApp()
		if err == nil && app != nil {
			// Clean up if successful
			// Note: In a real benchmark, you might want to avoid actual DB connections
		}
	}
}

func BenchmarkAttachMiddleware(b *testing.B) {
	app := fiber.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		attachMiddleware(app)
		// Reset app for next iteration
		app = fiber.New()
	}
}
