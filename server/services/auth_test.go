package services

import (
	"io"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		jwtSecret   string
		expectError bool
	}{
		{
			name:        "successful JWT generation with valid user ID",
			userID:      "507f1f77bcf86cd799439011",
			jwtSecret:   "test-secret-key",
			expectError: false,
		},
		{
			name:        "JWT generation with empty user ID",
			userID:      "",
			jwtSecret:   "test-secret-key",
			expectError: false, // Should not error, just create token with empty user ID
		},
		{
			name:        "JWT generation with empty secret",
			userID:      "507f1f77bcf86cd799439011",
			jwtSecret:   "",
			expectError: false, // Should not error, just use empty secret
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original JWT secret
			originalSecret := os.Getenv("JWT_SECRET")
			defer os.Setenv("JWT_SECRET", originalSecret)

			// Set test JWT secret
			os.Setenv("JWT_SECRET", tt.jwtSecret)

			// Generate JWT
			token, err := GenerateJWT(tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// Verify token can be parsed
				parsedToken, err := ParseJWT(token)
				assert.NoError(t, err)
				assert.NotNil(t, parsedToken)

				// Verify claims
				claims, ok := parsedToken.Claims.(jwt.MapClaims)
				assert.True(t, ok)
				assert.Equal(t, tt.userID, claims["userId"])

				// Verify expiration is set
				exp, ok := claims["exp"].(float64)
				assert.True(t, ok)
				assert.Greater(t, exp, float64(time.Now().Unix()))
			}
		})
	}
}

func TestParseJWT(t *testing.T) {
	// Set up test environment
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	tests := []struct {
		name        string
		token       string
		expectError bool
	}{
		{
			name:        "parse valid JWT",
			token:       "", // Will be generated
			expectError: false,
		},
		{
			name:        "parse invalid JWT",
			token:       "invalid.token.here",
			expectError: true,
		},
		{
			name:        "parse empty token",
			token:       "",
			expectError: true,
		},
		{
			name:        "parse malformed token",
			token:       "not.a.valid.jwt",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var token string
			if tt.name == "parse valid JWT" {
				// Generate a valid token for this test
				generatedToken, err := GenerateJWT("507f1f77bcf86cd799439011")
				assert.NoError(t, err)
				token = generatedToken
			} else {
				token = tt.token
			}

			parsedToken, err := ParseJWT(token)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, parsedToken)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, parsedToken)
				assert.True(t, parsedToken.Valid)
			}
		})
	}
}

func TestJWTProtected(t *testing.T) {
	// Set up test environment
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// Create Fiber app for testing
	app := fiber.New()

	// Add protected route
	protected := app.Group("/protected", JWTProtected())
	protected.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "protected"})
	})

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid JWT token",
			authHeader:     "", // Will be set with valid token
			expectedStatus: fiber.StatusOK,
			expectedBody:   `{"message":"protected"}`,
		},
		{
			name:           "missing Authorization header",
			authHeader:     "",
			expectedStatus: fiber.StatusUnauthorized,
			expectedBody:   "Missing Authorization header",
		},
		{
			name:           "invalid Authorization header format",
			authHeader:     "InvalidFormat token123",
			expectedStatus: fiber.StatusUnauthorized,
			expectedBody:   "Invalid token",
		},
		{
			name:           "invalid JWT token",
			authHeader:     "Bearer invalid.token.here",
			expectedStatus: fiber.StatusUnauthorized,
			expectedBody:   "Invalid token",
		},
		{
			name:           "expired JWT token",
			authHeader:     "", // Will be set with expired token
			expectedStatus: fiber.StatusUnauthorized,
			expectedBody:   "Token expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var authHeader string

			switch tt.name {
			case "valid JWT token":
				// Generate valid token
				token, err := GenerateJWT("507f1f77bcf86cd799439011")
				assert.NoError(t, err)
				authHeader = "Bearer " + token
			case "expired JWT token":
				// Generate expired token
				claims := jwt.MapClaims{
					"userId": "507f1f77bcf86cd799439011",
					"exp":    time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, err := token.SignedString([]byte("test-secret-key"))
				assert.NoError(t, err)
				authHeader = "Bearer " + tokenString
			default:
				authHeader = tt.authHeader
			}

			// Create request
			req := httptest.NewRequest("GET", "/protected/", nil)
			if authHeader != "" {
				req.Header.Set("Authorization", authHeader)
			}

			// Test request
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Read response body
			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Contains(t, string(body), tt.expectedBody)
		})
	}
}

func TestJWTProtected_Context(t *testing.T) {
	// Test that user ID is properly set in context
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	app := fiber.New()

	// Add protected route that checks context
	protected := app.Group("/protected", JWTProtected())
	protected.Get("/", func(c *fiber.Ctx) error {
		authorId := c.Locals("authorId")
		if authorId == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "authorId not found"})
		}
		return c.JSON(fiber.Map{"authorId": authorId})
	})

	// Generate valid token
	userID := "507f1f77bcf86cd799439011"
	token, err := GenerateJWT(userID)
	assert.NoError(t, err)

	// Test request
	req := httptest.NewRequest("GET", "/protected/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), userID)
}

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectError bool
	}{
		{
			name:        "hash valid password",
			password:    "mySecurePassword123",
			expectError: false,
		},
		{
			name:        "hash empty password",
			password:    "",
			expectError: false, // bcrypt allows empty passwords
		},
		{
			name:        "hash long password",
			password:    "veryLongPasswordWithManyCharacters123456789012345678901234567890",
			expectError: false,
		},
		{
			name:        "hash password with special characters",
			password:    "p@ssw0rd!@#$%^&*()",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
				assert.NotEqual(t, tt.password, hash) // Hash should be different from original

				// Verify hash can be verified
				err = CheckPasswordHash(tt.password, hash)
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	// Set up test password and hash
	password := "testPassword123"
	hash, err := HashPassword(password)
	assert.NoError(t, err)

	tests := []struct {
		name        string
		password    string
		hash        string
		expectError bool
	}{
		{
			name:        "correct password and hash",
			password:    password,
			hash:        hash,
			expectError: false,
		},
		{
			name:        "incorrect password",
			password:    "wrongPassword",
			hash:        hash,
			expectError: true,
		},
		{
			name:        "empty password",
			password:    "",
			hash:        hash,
			expectError: true,
		},
		{
			name:        "empty hash",
			password:    password,
			hash:        "",
			expectError: true,
		},
		{
			name:        "invalid hash format",
			password:    password,
			hash:        "invalid-hash-format",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestJWTExpiration(t *testing.T) {
	// Test JWT expiration handling
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// Generate token with short expiration
	claims := jwt.MapClaims{
		"userId": "507f1f77bcf86cd799439011",
		"exp":    time.Now().Add(1 * time.Millisecond).Unix(), // Expires in 1ms
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("test-secret-key"))
	assert.NoError(t, err)

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	// Try to parse expired token
	parsedToken, err := ParseJWT(tokenString)
	assert.Error(t, err)
	assert.Nil(t, parsedToken)
}

// Benchmark tests
func BenchmarkGenerateJWT(b *testing.B) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	userID := "507f1f77bcf86cd799439011"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateJWT(userID)
	}
}

func BenchmarkHashPassword(b *testing.B) {
	password := "testPassword123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashPassword(password)
	}
}

func BenchmarkCheckPasswordHash(b *testing.B) {
	password := "testPassword123"
	hash, _ := HashPassword(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckPasswordHash(password, hash)
	}
}
