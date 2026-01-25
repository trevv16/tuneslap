package handlers

import (
	"errors"
	"tuneslap/validation"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// use a single instance of Validate, it caches struct info
var validate = validator.New()

const (
	DEFAULT_PAGE_SIZE = 25
	MAX_PAGE_SIZE     = 100
)

type PaginationParams struct {
	Skip  int
	Limit int
}

type PaginationMeta struct {
	CurrentPage  int   `json:"currentPage"`
	PageSize     int   `json:"pageSize"`
	TotalResults int64 `json:"totalResults"`
}

func validatePaginationParams(c *fiber.Ctx) (*PaginationParams, error) {
	page := c.QueryInt("page", 0)
	skip := c.QueryInt("skip", -1)
	limit := c.QueryInt("limit", DEFAULT_PAGE_SIZE)

	// Validate limit if provided
	if limit <= 0 || limit > MAX_PAGE_SIZE {
		return nil, errors.New("limit must be between 1 and 100")
	}

	// Cannot use both page and skip
	if page > 0 && skip != -1 {
		return nil, errors.New("cannot use both page and skip parameters")
	}

	// Calculate skip value
	calculatedSkip := 0
	if page > 0 {
		calculatedSkip = (page - 1) * limit
	} else if skip != -1 {
		if skip < 0 {
			return nil, errors.New("skip cannot be negative")
		}
		calculatedSkip = skip
	}

	return &PaginationParams{
		Skip:  calculatedSkip,
		Limit: limit,
	}, nil
}

func GetAuthorId(c *fiber.Ctx) (primitive.ObjectID, error) {
	authorId := c.Locals("authorId")
	if authorId == nil {
		return primitive.NilObjectID, errors.New("authorId not found")
	}

	authorObjectId, err := primitive.ObjectIDFromHex(authorId.(string))
	if err != nil {
		return primitive.NilObjectID, errors.New("invalid authorId")
	}

	return authorObjectId, nil
}

func GetValidObjectId(c *fiber.Ctx, idKey string) (primitive.ObjectID, error) {
	value := c.Params(idKey)
	return primitive.ObjectIDFromHex(value)
}

// SendErrorResponse sends a standardized error response
func SendErrorResponse(c *fiber.Ctx, status int, message string, err error) error {
	var errData interface{}
	if err != nil {
		errData = err.Error()
	}
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"message": message,
		"data":    errData,
	})
}

// SendSuccessResponse sends a standardized success response
func SendSuccessResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// SendPaginatedResponse sends a standardized paginated response
// Wraps data array in an object with the field name and pagination (matches OpenAPI spec)
func SendPaginatedResponse(c *fiber.Ctx, status int, message string, data interface{}, pagination interface{}, dataFieldName string) error {
	if dataFieldName == "" {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "dataFieldName is required", errors.New("dataFieldName cannot be empty"))
	}
	return c.Status(status).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data": fiber.Map{
			dataFieldName: data,
			"pagination":  pagination,
		},
	})
}

// SendValidationErrorResponse sends a standardized validation error response
func SendValidationErrorResponse(c *fiber.Ctx, result validation.ValidationResult) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"message": "Validation failed",
		"data":    result.Errors,
	})
}
