package handlers

import (
	"tuneslap/repositories"
	"tuneslap/validation"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
