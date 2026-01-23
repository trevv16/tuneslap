//go:build integration
// +build integration

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"tuneslap/database"
	"tuneslap/models"
	"tuneslap/router"
	"tuneslap/services"
	"tuneslap/testutils"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// setupIntegrationTestApp creates a test app with full routes
func setupIntegrationTestApp(t *testing.T) (*fiber.App, func()) {
	// Set up environment
	os.Setenv("JWT_SECRET", "test-secret-key-for-integration-tests")

	// Setup databases
	mongoCleanup := testutils.SetupTestMongoDBWithSkip(t)
	redisCleanup := testutils.SetupTestRedisWithSkip(t)

	// Create app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	// Setup routes
	router.SetupRoutes(app)

	cleanup := func() {
		redisCleanup()
		mongoCleanup()
	}

	return app, cleanup
}

// Helper to make authenticated request
func makeAuthenticatedRequest(app *fiber.App, method, path, token string, body interface{}) (*http.Response, error) {
	var bodyReader *bytes.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonBody)
	} else {
		bodyReader = bytes.NewReader([]byte{})
	}

	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return app.Test(req, -1) // -1 = no timeout
}

func TestHealthEndpoint(t *testing.T) {
	app, cleanup := setupIntegrationTestApp(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuthFlow(t *testing.T) {
	app, cleanup := setupIntegrationTestApp(t)
	defer cleanup()

	t.Run("signup creates new user", func(t *testing.T) {
		body := map[string]string{
			"name":     "Test User",
			"email":    "signup@test.com",
			"password": "password123",
		}

		resp, err := makeAuthenticatedRequest(app, "POST", "/api/v1/auth/signup", "", body)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.True(t, result["success"].(bool))
		assert.NotEmpty(t, result["data"])
	})

	t.Run("signup fails with duplicate email", func(t *testing.T) {
		// First signup
		body := map[string]string{
			"name":     "Test User",
			"email":    "duplicate@test.com",
			"password": "password123",
		}
		_, _ = makeAuthenticatedRequest(app, "POST", "/api/v1/auth/signup", "", body)

		// Second signup with same email
		resp, err := makeAuthenticatedRequest(app, "POST", "/api/v1/auth/signup", "", body)
		require.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("signin with valid credentials", func(t *testing.T) {
		// Create user first
		signupBody := map[string]string{
			"name":     "Signin Test",
			"email":    "signin@test.com",
			"password": "password123",
		}
		_, _ = makeAuthenticatedRequest(app, "POST", "/api/v1/auth/signup", "", signupBody)

		// Sign in
		signinBody := map[string]string{
			"email":    "signin@test.com",
			"password": "password123",
		}
		resp, err := makeAuthenticatedRequest(app, "POST", "/api/v1/auth/signin", "", signinBody)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.True(t, result["success"].(bool))
		// Check token is returned
		data := result["data"].(map[string]interface{})
		assert.NotEmpty(t, data["token"])
	})

	t.Run("signin with invalid credentials", func(t *testing.T) {
		signinBody := map[string]string{
			"email":    "nonexistent@test.com",
			"password": "wrongpassword",
		}
		resp, err := makeAuthenticatedRequest(app, "POST", "/api/v1/auth/signin", "", signinBody)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestBoardsFlow(t *testing.T) {
	app, cleanup := setupIntegrationTestApp(t)
	defer cleanup()

	// Create user and get token
	signupBody := map[string]string{
		"name":     "Board Test User",
		"email":    "boardtest@test.com",
		"password": "password123",
	}
	_, _ = makeAuthenticatedRequest(app, "POST", "/api/v1/auth/signup", "", signupBody)

	signinBody := map[string]string{
		"email":    "boardtest@test.com",
		"password": "password123",
	}
	signinResp, _ := makeAuthenticatedRequest(app, "POST", "/api/v1/auth/signin", "", signinBody)
	var signinResult map[string]interface{}
	json.NewDecoder(signinResp.Body).Decode(&signinResult)
	token := signinResult["data"].(map[string]interface{})["token"].(string)

	t.Run("create board", func(t *testing.T) {
		body := map[string]string{
			"name":        "Test Board",
			"description": "A test board",
			"layout":      "grid",
		}
		resp, err := makeAuthenticatedRequest(app, "POST", "/api/v1/boards", token, body)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.True(t, result["success"].(bool))
	})

	t.Run("get boards", func(t *testing.T) {
		resp, err := makeAuthenticatedRequest(app, "GET", "/api/v1/boards", token, nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.True(t, result["success"].(bool))
	})

	t.Run("unauthorized access", func(t *testing.T) {
		resp, err := makeAuthenticatedRequest(app, "GET", "/api/v1/boards", "", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestUserProfileFlow(t *testing.T) {
	app, cleanup := setupIntegrationTestApp(t)
	defer cleanup()

	// Create user and get token
	signupBody := map[string]string{
		"name":     "Profile Test User",
		"email":    "profiletest@test.com",
		"password": "password123",
	}
	_, _ = makeAuthenticatedRequest(app, "POST", "/api/v1/auth/signup", "", signupBody)

	signinBody := map[string]string{
		"email":    "profiletest@test.com",
		"password": "password123",
	}
	signinResp, _ := makeAuthenticatedRequest(app, "POST", "/api/v1/auth/signin", "", signinBody)
	var signinResult map[string]interface{}
	json.NewDecoder(signinResp.Body).Decode(&signinResult)
	token := signinResult["data"].(map[string]interface{})["token"].(string)

	t.Run("get current user", func(t *testing.T) {
		resp, err := makeAuthenticatedRequest(app, "GET", "/api/v1/users/me", token, nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.True(t, result["success"].(bool))
	})

	t.Run("update current user", func(t *testing.T) {
		body := map[string]string{
			"name": "Updated Name",
		}
		resp, err := makeAuthenticatedRequest(app, "PATCH", "/api/v1/users/me", token, body)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
