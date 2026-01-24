package handlers

import (
	"tuneslap/validation"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Generic validator interface
type Validator interface {
	Validate(data interface{}) validation.ValidationResult
}

// Generic repository interface
type Repository[T any] interface {
	FindAll(filter bson.M, opts ...*options.FindOptions) ([]T, error)
	FindOne(filter bson.M) (T, error)
	FindByID(id primitive.ObjectID) (T, error)
	Create(document T) (T, error)
	Update(id primitive.ObjectID, update bson.M) (T, error)
	Delete(id primitive.ObjectID) error
	Count(filter bson.M) (int64, error)
}

// CreateFindOptions creates find options with pagination from query parameters
func CreateFindOptions(c *fiber.Ctx) (*options.FindOptions, error) {
	pagination, err := validatePaginationParams(c)
	if err != nil {
		return nil, err
	}

	findOptions := options.Find().
		SetSkip(int64(pagination.Skip)).
		SetLimit(int64(pagination.Limit)).
		SetSort(bson.M{"createdAt": -1})

	return findOptions, nil
}

// CreateProjection creates a MongoDB projection from field names
func CreateProjection(fields []string) bson.M {
	projection := bson.M{}
	for _, field := range fields {
		projection[field] = 1
	}
	return projection
}

// CreateFindOneOptions creates find one options with optional projection
func CreateFindOneOptions(projection bson.M) *options.FindOneOptions {
	if projection == nil {
		return nil
	}
	return options.FindOne().SetProjection(projection)
}

// CreatePaginationMeta creates pagination metadata from skip, limit, and total count
func CreatePaginationMeta(skip, limit, totalCount int64) PaginationMeta {
	currentPage := int((skip / limit) + 1)
	return PaginationMeta{
		CurrentPage:  currentPage,
		PageSize:     int(limit),
		TotalResults: totalCount,
	}
}
