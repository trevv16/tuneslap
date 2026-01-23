package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// TestApp wraps a Fiber app for testing
type TestApp struct {
	App *fiber.App
}

// SetupTestApp creates a new Fiber app with test configuration
func SetupTestApp() *TestApp {
	// Set up test environment variables
	os.Setenv("JWT_SECRET", "test-secret-key")

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

	// Add minimal middleware for testing
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	return &TestApp{App: app}
}

// SetupTestAppWithLogger creates a test app with logging enabled
func SetupTestAppWithLogger() *TestApp {
	testApp := SetupTestApp()
	testApp.App.Use(logger.New())
	return testApp
}

// AddRoute adds a route to the test app
func (ta *TestApp) AddRoute(method, path string, handler fiber.Handler) {
	switch method {
	case "GET":
		ta.App.Get(path, handler)
	case "POST":
		ta.App.Post(path, handler)
	case "PATCH":
		ta.App.Patch(path, handler)
	case "DELETE":
		ta.App.Delete(path, handler)
	case "PUT":
		ta.App.Put(path, handler)
	}
}

// AddGroup adds a group to the test app
func (ta *TestApp) AddGroup(prefix string, handlers ...fiber.Handler) fiber.Router {
	return ta.App.Group(prefix, handlers...)
}

// TestRequest creates a test HTTP request
func (ta *TestApp) TestRequest(method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return ta.App.Test(req)
}

// TestRequestWithAuth creates a test HTTP request with JWT authentication
func (ta *TestApp) TestRequestWithAuth(method, path, token string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	return ta.App.Test(req)
}

// CreateTestRequest creates a raw HTTP test request
func CreateTestRequest(method, path string, body interface{}) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// CreateAuthenticatedRequest creates an HTTP request with JWT token
func CreateAuthenticatedRequest(method, path, token string, body interface{}) (*http.Request, error) {
	req, err := CreateTestRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	return req, nil
}

// GenerateTestToken generates a JWT token for testing
// Note: This uses a simple implementation to avoid import cycles
func GenerateTestToken(userID string) (string, error) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	// Import services at runtime is not possible, so we provide a simple token for testing
	// For actual JWT generation, use services.GenerateJWT directly in handler tests
	return "test-token-" + userID, nil
}

// ParseResponseBody parses the response body into a target struct
func ParseResponseBody(resp *http.Response, target interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.Unmarshal(body, target)
}

// GetResponseBody reads and returns the response body as string
func GetResponseBody(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return string(body), nil
}

// AssertResponseStatus checks if the response has the expected status code
func AssertResponseStatus(resp *http.Response, expectedStatus int) bool {
	return resp.StatusCode == expectedStatus
}
