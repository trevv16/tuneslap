package handlers

import (
	"tuneslap/repositories"
	"tuneslap/validation"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BaseHandler provides common CRUD operations and error handling
type BaseHandler[T any, CreateRequest any, UpdateRequest any] struct {
	repo      *repositories.Repository[T]
	validator interface {
		Validate(data interface{}) validation.ValidationResult
	}
}

// NewBaseHandler creates a new base handler
func NewBaseHandler[T any, CreateRequest any, UpdateRequest any](
	repo *repositories.Repository[T],
	validator interface {
		Validate(data interface{}) validation.ValidationResult
	},
) *BaseHandler[T, CreateRequest, UpdateRequest] {
	return &BaseHandler[T, CreateRequest, UpdateRequest]{
		repo:      repo,
		validator: validator,
	}
}

// HandleGetAll handles GET all requests with pagination and optional projection
func (h *BaseHandler[T, CreateRequest, UpdateRequest]) HandleGetAll(
	c *fiber.Ctx,
	filterFunc func(authorId primitive.ObjectID) bson.M,
	responseMapper func(T) interface{},
	projection bson.M,
	dataFieldName string,
) error {
	// Get and validate authorId
	authorId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}

	// Validate pagination parameters
	pagination, err := validatePaginationParams(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid pagination parameters", err)
	}

	// Create find options with pagination and projection
	findOptions := options.Find().
		SetSkip(int64(pagination.Skip)).
		SetLimit(int64(pagination.Limit)).
		SetSort(bson.M{"createdAt": -1})

	if projection != nil {
		findOptions = findOptions.SetProjection(projection)
	}

	// Get results with pagination
	results, err := h.repo.FindAll(filterFunc(authorId), findOptions)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Error retrieving data", err)
	}

	// Get total count
	totalCount, err := h.repo.Count(filterFunc(authorId))
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Error counting data", err)
	}

	// Map results if mapper provided
	var data interface{}
	if responseMapper != nil {
		mappedResults := make([]interface{}, len(results))
		for i, result := range results {
			mappedResults[i] = responseMapper(result)
		}
		data = mappedResults
	} else {
		data = results
	}

	// Create pagination metadata
	// Safety check: ensure Limit is not zero to avoid division by zero
	currentPage := 1
	if pagination.Limit > 0 {
		currentPage = (pagination.Skip / pagination.Limit) + 1
	}
	paginationMeta := PaginationMeta{
		CurrentPage:  currentPage,
		PageSize:     pagination.Limit,
		TotalResults: totalCount,
	}

	return SendPaginatedResponse(
		c,
		fiber.StatusOK,
		"Data retrieved successfully",
		data,
		paginationMeta,
		dataFieldName,
	)
}

// HandleGetById handles GET by ID requests with optional projection
func (h *BaseHandler[T, CreateRequest, UpdateRequest]) HandleGetById(
	c *fiber.Ctx,
	idParam string,
	authorId primitive.ObjectID,
	responseMapper func(T) interface{},
	projection bson.M,
) error {
	// Parse ID
	id, err := GetValidObjectId(c, idParam)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	// Create find one options with projection
	var findOptions *options.FindOneOptions
	if projection != nil {
		findOptions = options.FindOne().SetProjection(projection)
	}

	// Get result with author filter and projection
	var result T
	if findOptions != nil {
		result, err = h.repo.FindOne(bson.M{
			"_id":      id,
			"authorId": authorId,
		}, findOptions)
	} else {
		result, err = h.repo.FindOne(bson.M{
			"_id":      id,
			"authorId": authorId,
		})
	}

	if err != nil {
		return SendErrorResponse(c, fiber.StatusNotFound, "Resource not found", err)
	}

	// Map result if mapper provided
	var data interface{}
	if responseMapper != nil {
		data = responseMapper(result)
	} else {
		data = result
	}

	// Return data directly (not wrapped) to match OpenAPI spec for GET by ID endpoints
	return c.Status(fiber.StatusOK).JSON(data)
}

// HandleCreate handles CREATE requests
func (h *BaseHandler[T, CreateRequest, UpdateRequest]) HandleCreate(
	c *fiber.Ctx,
	authorId primitive.ObjectID,
	createFunc func(CreateRequest, primitive.ObjectID) (T, error),
	responseMapper func(T) interface{},
) error {
	// Parse request body
	request := new(CreateRequest)
	if err := c.BodyParser(request); err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	validationResult := h.validator.Validate(request)
	if !validationResult.IsValid {
		return SendValidationErrorResponse(c, validationResult)
	}

	// Create resource
	result, err := createFunc(*request, authorId)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create resource", err)
	}

	// Map result if mapper provided
	var data interface{}
	if responseMapper != nil {
		data = responseMapper(result)
	} else {
		data = result
	}

	return SendSuccessResponse(c, fiber.StatusCreated, "Resource created successfully", data)
}

// HandleUpdate handles UPDATE requests
func (h *BaseHandler[T, CreateRequest, UpdateRequest]) HandleUpdate(
	c *fiber.Ctx,
	idParam string,
	authorId primitive.ObjectID,
	updateFunc func(primitive.ObjectID, primitive.ObjectID, UpdateRequest) (T, error),
	responseMapper func(T) interface{},
) error {
	// Parse ID
	id, err := GetValidObjectId(c, idParam)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	// Parse request body
	request := new(UpdateRequest)
	if err := c.BodyParser(request); err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	validationResult := h.validator.Validate(request)
	if !validationResult.IsValid {
		return SendValidationErrorResponse(c, validationResult)
	}

	// Update resource
	result, err := updateFunc(id, authorId, *request)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to update resource", err)
	}

	// Map result if mapper provided
	var data interface{}
	if responseMapper != nil {
		data = responseMapper(result)
	} else {
		data = result
	}

	return SendSuccessResponse(c, fiber.StatusOK, "Resource updated successfully", data)
}

// HandleDelete handles DELETE requests
func (h *BaseHandler[T, CreateRequest, UpdateRequest]) HandleDelete(
	c *fiber.Ctx,
	idParam string,
	authorId primitive.ObjectID,
) error {
	// Parse ID
	id, err := GetValidObjectId(c, idParam)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	// First verify the resource exists and belongs to the author
	_, err = h.repo.FindOne(bson.M{
		"_id":      id,
		"authorId": authorId,
	})
	if err != nil {
		return SendErrorResponse(c, fiber.StatusNotFound, "Resource not found", err)
	}

	// Delete resource
	err = h.repo.Delete(id)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete resource", err)
	}

	return SendSuccessResponse(c, fiber.StatusOK, "Resource deleted successfully", nil)
}

// HandleGetPublic handles GET public resources (no authorId filtering)
func (h *BaseHandler[T, CreateRequest, UpdateRequest]) HandleGetPublic(
	c *fiber.Ctx,
	filter bson.M,
	responseMapper func(T) interface{},
	dataFieldName string,
) error {
	// Create find options with pagination
	findOptions, err := CreateFindOptions(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid pagination parameters", err)
	}

	// Get results with pagination
	results, err := h.repo.FindAll(filter, findOptions)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Error retrieving data", err)
	}

	// Get total count
	totalCount, err := h.repo.Count(filter)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Error counting data", err)
	}

	// Map results if mapper provided
	var data interface{}
	if responseMapper != nil {
		mappedResults := make([]interface{}, len(results))
		for i, result := range results {
			mappedResults[i] = responseMapper(result)
		}
		data = mappedResults
	} else {
		data = results
	}

	// Create pagination metadata using helper function
	paginationMeta := CreatePaginationMeta(
		*findOptions.Skip,
		*findOptions.Limit,
		totalCount,
	)

	return SendPaginatedResponse(
		c,
		fiber.StatusOK,
		"Data retrieved successfully",
		data,
		paginationMeta,
		dataFieldName,
	)
}

// HandleGetCurrentUser handles GET current user by ID
func (h *BaseHandler[T, CreateRequest, UpdateRequest]) HandleGetCurrentUser(
	c *fiber.Ctx,
	responseMapper func(T) interface{},
) error {
	// Get and validate the authorId
	authorId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}

	// Get result by ID (not authorId)
	result, err := h.repo.FindByID(authorId)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusNotFound, "Resource not found", err)
	}

	// Map result if mapper provided
	var data interface{}
	if responseMapper != nil {
		data = responseMapper(result)
	} else {
		data = result
	}

	// Return data directly (not wrapped) to match OpenAPI spec for GET by ID endpoints
	return c.Status(fiber.StatusOK).JSON(data)
}

// HandleGetAllByCustomField handles GET all with custom field mapping
func (h *BaseHandler[T, CreateRequest, UpdateRequest]) HandleGetAllByCustomField(
	c *fiber.Ctx,
	fieldName string,
	responseMapper func(T) interface{},
	dataFieldName string,
) error {
	// Get and validate authorId
	authorId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}

	// Create find options with pagination
	findOptions, err := CreateFindOptions(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid pagination parameters", err)
	}

	// Create filter with custom field name
	filter := bson.M{fieldName: authorId}

	// Get results with pagination
	results, err := h.repo.FindAll(filter, findOptions)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Error retrieving data", err)
	}

	// Get total count
	totalCount, err := h.repo.Count(filter)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Error counting data", err)
	}

	// Map results if mapper provided
	var data interface{}
	if responseMapper != nil {
		mappedResults := make([]interface{}, len(results))
		for i, result := range results {
			mappedResults[i] = responseMapper(result)
		}
		data = mappedResults
	} else {
		data = results
	}

	// Create pagination metadata using helper function
	paginationMeta := CreatePaginationMeta(
		*findOptions.Skip,
		*findOptions.Limit,
		totalCount,
	)

	return SendPaginatedResponse(
		c,
		fiber.StatusOK,
		"Data retrieved successfully",
		data,
		paginationMeta,
		dataFieldName,
	)
}

// HandleGetByIdByCustomField handles GET by ID with custom field mapping
func (h *BaseHandler[T, CreateRequest, UpdateRequest]) HandleGetByIdByCustomField(
	c *fiber.Ctx,
	idParam string,
	fieldName string,
	responseMapper func(T) interface{},
) error {
	// Get and validate authorId
	authorId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}

	// Parse ID
	id, err := GetValidObjectId(c, idParam)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	// Get result with custom field filter
	result, err := h.repo.FindOne(bson.M{
		"_id":     id,
		fieldName: authorId,
	})
	if err != nil {
		return SendErrorResponse(c, fiber.StatusNotFound, "Resource not found", err)
	}

	// Map result if mapper provided
	var data interface{}
	if responseMapper != nil {
		data = responseMapper(result)
	} else {
		data = result
	}

	return SendSuccessResponse(c, fiber.StatusOK, "Resource found", data)
}

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

// Helper functions for public handlers

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

// ArrayHandler provides common operations for embedded arrays
type ArrayHandler[T any, CreateRequest any, UpdateRequest any] struct {
	arrayRepo *repositories.ArrayRepository[T]
	validator interface {
		Validate(data interface{}) validation.ValidationResult
	}
}

// NewArrayHandler creates a new array handler
func NewArrayHandler[T any, CreateRequest any, UpdateRequest any](
	arrayRepo *repositories.ArrayRepository[T],
	validator interface {
		Validate(data interface{}) validation.ValidationResult
	},
) *ArrayHandler[T, CreateRequest, UpdateRequest] {
	return &ArrayHandler[T, CreateRequest, UpdateRequest]{
		arrayRepo: arrayRepo,
		validator: validator,
	}
}

// HandleCreateInArray handles CREATE operations for embedded arrays
func (h *ArrayHandler[T, CreateRequest, UpdateRequest]) HandleCreateInArray(
	c *fiber.Ctx,
	parentIdParam string,
	arrayField string,
	createFunc func(CreateRequest) T,
	responseMapper func(interface{}) interface{},
	preCreateHook ...func(*fiber.Ctx, primitive.ObjectID, *CreateRequest) error,
) error {
	// Get and validate authorId for authentication
	_, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}

	// Parse parent ID
	parentId, err := GetValidObjectId(c, parentIdParam)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid parent ID", err)
	}

	// Parse request body
	request := new(CreateRequest)
	if err := c.BodyParser(request); err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	validationResult := h.validator.Validate(request)
	if !validationResult.IsValid {
		return SendValidationErrorResponse(c, validationResult)
	}

	// Call pre-create hook if provided
	if len(preCreateHook) > 0 && preCreateHook[0] != nil {
		if err := preCreateHook[0](c, parentId, request); err != nil {
			return err // Hook should return appropriate error response
		}
	}

	// Create the array element
	element := createFunc(*request)

	// Add to array
	_, err = h.arrayRepo.CreateInArray(parentId, arrayField, element)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create array element", err)
	}

	// Map result if mapper provided
	var data interface{}
	if responseMapper != nil {
		data = responseMapper(element)
	} else {
		data = element
	}

	return SendSuccessResponse(c, fiber.StatusCreated, "Array element created successfully", data)
}

// HandleUpdateInArray handles UPDATE operations for embedded arrays
func (h *ArrayHandler[T, CreateRequest, UpdateRequest]) HandleUpdateInArray(
	c *fiber.Ctx,
	parentIdParam string,
	elementIdParam string,
	arrayField string,
	updateFunc func(UpdateRequest) bson.M,
	responseMapper func(interface{}) interface{},
) error {
	// Get and validate authorId for authentication
	_, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}

	// Parse parent ID
	parentId, err := GetValidObjectId(c, parentIdParam)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid parent ID", err)
	}

	// Parse element ID
	elementId, err := GetValidObjectId(c, elementIdParam)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid element ID", err)
	}

	// Parse request body
	request := new(UpdateRequest)
	if err := c.BodyParser(request); err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	validationResult := h.validator.Validate(request)
	if !validationResult.IsValid {
		return SendValidationErrorResponse(c, validationResult)
	}

	// Create update document
	update := updateFunc(*request)

	// Update array element
	err = h.arrayRepo.UpdateInArray(parentId, arrayField, elementId, update)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to update array element", err)
	}

	// Fetch the updated element
	updatedElement, err := h.arrayRepo.GetArrayElement(parentId, arrayField, elementId)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve updated element", err)
	}

	// Map result if mapper provided
	var data interface{}
	if responseMapper != nil {
		data = responseMapper(updatedElement)
	} else {
		data = updatedElement
	}

	return SendSuccessResponse(c, fiber.StatusOK, "Array element updated successfully", data)
}

// HandleDeleteFromArray handles DELETE operations for embedded arrays
func (h *ArrayHandler[T, CreateRequest, UpdateRequest]) HandleDeleteFromArray(
	c *fiber.Ctx,
	parentIdParam string,
	elementIdParam string,
	arrayField string,
) error {
	// Get and validate authorId for authentication
	_, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}

	// Parse parent ID
	parentId, err := GetValidObjectId(c, parentIdParam)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid parent ID", err)
	}

	// Parse element ID
	elementId, err := GetValidObjectId(c, elementIdParam)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid element ID", err)
	}

	// Delete array element
	err = h.arrayRepo.DeleteFromArray(parentId, arrayField, elementId)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete array element", err)
	}

	return SendSuccessResponse(c, fiber.StatusOK, "Array element deleted successfully", nil)
}

// HandleGetAllFromArray handles GET all operations for embedded arrays
func (h *ArrayHandler[T, CreateRequest, UpdateRequest]) HandleGetAllFromArray(
	c *fiber.Ctx,
	parentIdParam string,
	arrayField string,
	responseMapper func(interface{}) interface{},
) error {
	// Get and validate authorId for authentication
	_, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}

	// Parse parent ID
	parentId, err := GetValidObjectId(c, parentIdParam)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid parent ID", err)
	}

	// Get all elements from array
	elements, err := h.arrayRepo.GetAllArrayElements(parentId, arrayField)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve array elements", err)
	}

	// Map results if mapper provided
	var data interface{}
	if responseMapper != nil {
		mappedResults := make([]interface{}, len(elements))
		for i, element := range elements {
			mappedResults[i] = responseMapper(element)
		}
		data = mappedResults
	} else {
		data = elements
	}

	return SendSuccessResponse(c, fiber.StatusOK, "Array elements retrieved successfully", data)
}

// HandleGetByIdFromArray handles GET by ID operations for embedded arrays
func (h *ArrayHandler[T, CreateRequest, UpdateRequest]) HandleGetByIdFromArray(
	c *fiber.Ctx,
	parentIdParam string,
	elementIdParam string,
	arrayField string,
	responseMapper func(interface{}) interface{},
) error {
	// Get and validate authorId for authentication
	_, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}

	// Parse parent ID
	parentId, err := GetValidObjectId(c, parentIdParam)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid parent ID", err)
	}

	// Parse element ID
	elementId, err := GetValidObjectId(c, elementIdParam)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid element ID", err)
	}

	// Get element from array
	element, err := h.arrayRepo.GetArrayElement(parentId, arrayField, elementId)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusNotFound, "Array element not found", err)
	}

	// Map result if mapper provided
	var data interface{}
	if responseMapper != nil {
		data = responseMapper(element)
	} else {
		data = element
	}

	// Return wrapped response for array element GET by ID (collaborator/key endpoints expect wrapped format)
	return SendSuccessResponse(c, fiber.StatusOK, "Array element retrieved successfully", data)
}
