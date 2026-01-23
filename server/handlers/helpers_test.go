package handlers

import (
	"errors"
	"net/http/httptest"
	"testing"
	"tuneslap/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestValidatePaginationParams(t *testing.T) {
	tests := []struct {
		name        string
		queryParams map[string]string
		expectError bool
		expectSkip  int
		expectLimit int
	}{
		{
			name:        "default values",
			queryParams: map[string]string{},
			expectError: false,
			expectSkip:  0,
			expectLimit: 25, // DEFAULT_PAGE_SIZE
		},
		{
			name: "custom page and limit",
			queryParams: map[string]string{
				"page":  "2",
				"limit": "10",
			},
			expectError: false,
			expectSkip:  10, // (2-1) * 10
			expectLimit: 10,
		},
		{
			name: "page 1",
			queryParams: map[string]string{
				"page":  "1",
				"limit": "25",
			},
			expectError: false,
			expectSkip:  0,
			expectLimit: 25,
		},
		{
			name: "limit exceeds max returns error",
			queryParams: map[string]string{
				"limit": "200",
			},
			expectError: true, // Implementation returns error for limit > MAX_PAGE_SIZE
		},
		{
			name: "zero limit",
			queryParams: map[string]string{
				"limit": "0",
			},
			expectError: true,
		},
		{
			name: "valid skip",
			queryParams: map[string]string{
				"skip":  "20",
				"limit": "10",
			},
			expectError: false,
			expectSkip:  20,
			expectLimit: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			var result *PaginationParams
			var resultErr error

			app.Get("/test", func(c *fiber.Ctx) error {
				result, resultErr = validatePaginationParams(c)
				return c.SendStatus(fiber.StatusOK)
			})

			// Build URL with query params
			url := "/test"
			if len(tt.queryParams) > 0 {
				url += "?"
				first := true
				for k, v := range tt.queryParams {
					if !first {
						url += "&"
					}
					url += k + "=" + v
					first = false
				}
			}

			req := httptest.NewRequest("GET", url, nil)
			_, err := app.Test(req)
			assert.NoError(t, err)

			if tt.expectError {
				assert.Error(t, resultErr)
			} else {
				assert.NoError(t, resultErr)
				if result != nil {
					assert.Equal(t, tt.expectSkip, result.Skip)
					assert.Equal(t, tt.expectLimit, result.Limit)
				}
			}
		})
	}
}

func TestGetAuthorId(t *testing.T) {
	validID := primitive.NewObjectID()

	tests := []struct {
		name        string
		localValue  interface{}
		expectError bool
		expectedID  primitive.ObjectID
	}{
		{
			name:        "valid ObjectID string",
			localValue:  validID.Hex(),
			expectError: false,
			expectedID:  validID,
		},
		{
			name:        "missing authorId",
			localValue:  nil,
			expectError: true,
		},
		{
			name:        "invalid ObjectID",
			localValue:  "invalid",
			expectError: true,
		},
		{
			name:        "empty string",
			localValue:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			var result primitive.ObjectID
			var resultErr error

			app.Get("/test", func(c *fiber.Ctx) error {
				if tt.localValue != nil {
					c.Locals("authorId", tt.localValue)
				}
				result, resultErr = GetAuthorId(c)
				return c.SendStatus(fiber.StatusOK)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			_, err := app.Test(req)
			assert.NoError(t, err)

			if tt.expectError {
				assert.Error(t, resultErr)
			} else {
				assert.NoError(t, resultErr)
				assert.Equal(t, tt.expectedID, result)
			}
		})
	}
}

func TestGetValidObjectId(t *testing.T) {
	validID := primitive.NewObjectID()

	tests := []struct {
		name        string
		paramValue  string
		expectError bool
		expectedID  primitive.ObjectID
	}{
		{
			name:        "valid ObjectID",
			paramValue:  validID.Hex(),
			expectError: false,
			expectedID:  validID,
		},
		{
			name:        "invalid ObjectID",
			paramValue:  "invalid",
			expectError: true,
		},
		{
			name:        "too short string",
			paramValue:  "abc123",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			var result primitive.ObjectID
			var resultErr error

			app.Get("/test/:id", func(c *fiber.Ctx) error {
				result, resultErr = GetValidObjectId(c, "id")
				return c.SendStatus(fiber.StatusOK)
			})

			req := httptest.NewRequest("GET", "/test/"+tt.paramValue, nil)
			_, err := app.Test(req)
			assert.NoError(t, err)

			if tt.expectError {
				assert.Error(t, resultErr)
			} else {
				assert.NoError(t, resultErr)
				assert.Equal(t, tt.expectedID, result)
			}
		})
	}
}

func TestSendErrorResponse(t *testing.T) {
	app := fiber.New()
	testErr := errors.New("test error details")

	app.Get("/test", func(c *fiber.Ctx) error {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Test error", testErr)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSendSuccessResponse(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return SendSuccessResponse(c, fiber.StatusOK, "Test success", map[string]string{"key": "value"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestSendPaginatedResponse(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		data := []map[string]string{{"key": "value"}}
		pagination := PaginationMeta{
			CurrentPage:  1,
			PageSize:     10,
			TotalResults: 1,
		}
		return SendPaginatedResponse(c, fiber.StatusOK, "Test", data, pagination, "items")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestSendValidationErrorResponse(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		result := validation.ValidationResult{
			IsValid: false,
			Errors: []validation.ValidationError{
				{Field: "name", Message: "name is required"},
			},
		}
		return SendValidationErrorResponse(c, result)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCreateFindOptions(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		opts, err := CreateFindOptions(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		if opts != nil {
			return c.SendStatus(fiber.StatusOK)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	})

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
	}{
		{
			name:           "default options",
			queryParams:    "",
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "with page and limit",
			queryParams:    "?page=1&limit=10",
			expectedStatus: fiber.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test"+tt.queryParams, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestCreateProjection(t *testing.T) {
	tests := []struct {
		name        string
		fields      []string
		expectCount int
	}{
		{
			name:        "no fields",
			fields:      []string{},
			expectCount: 0,
		},
		{
			name:        "single field",
			fields:      []string{"name"},
			expectCount: 1,
		},
		{
			name:        "multiple fields",
			fields:      []string{"name", "email", "createdAt"},
			expectCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projection := CreateProjection(tt.fields)
			assert.Len(t, projection, tt.expectCount)
			for _, field := range tt.fields {
				assert.Contains(t, projection, field)
				assert.Equal(t, 1, projection[field])
			}
		})
	}
}

func TestCreatePaginationMeta(t *testing.T) {
	tests := []struct {
		name         string
		skip         int64
		limit        int64
		totalCount   int64
		expectedPage int
		expectedSize int
	}{
		{
			name:         "first page",
			skip:         0,
			limit:        10,
			totalCount:   100,
			expectedPage: 1,
			expectedSize: 10,
		},
		{
			name:         "second page",
			skip:         10,
			limit:        10,
			totalCount:   100,
			expectedPage: 2,
			expectedSize: 10,
		},
		{
			name:         "large skip",
			skip:         50,
			limit:        25,
			totalCount:   200,
			expectedPage: 3,
			expectedSize: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := CreatePaginationMeta(tt.skip, tt.limit, tt.totalCount)
			assert.Equal(t, tt.expectedPage, meta.CurrentPage)
			assert.Equal(t, tt.expectedSize, meta.PageSize)
			assert.Equal(t, int64(tt.totalCount), meta.TotalResults)
		})
	}
}

// Benchmark tests
func BenchmarkValidatePaginationParams(b *testing.B) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		validatePaginationParams(c)
		return c.SendStatus(fiber.StatusOK)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/test?page=2&limit=10", nil)
		app.Test(req)
	}
}

func BenchmarkGetAuthorId(b *testing.B) {
	app := fiber.New()
	validID := primitive.NewObjectID()

	app.Get("/test", func(c *fiber.Ctx) error {
		c.Locals("authorId", validID.Hex())
		GetAuthorId(c)
		return c.SendStatus(fiber.StatusOK)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		app.Test(req)
	}
}
