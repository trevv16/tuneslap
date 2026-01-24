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

// validateAndParseParentId validates auth and parses parent ID
func (h *ArrayHandler[T, CreateRequest, UpdateRequest]) validateAndParseParentId(
	c *fiber.Ctx,
	parentIdParam string,
) (primitive.ObjectID, error) {
	if _, err := GetAuthorId(c); err != nil {
		return primitive.NilObjectID, SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}
	parentId, err := GetValidObjectId(c, parentIdParam)
	if err != nil {
		return primitive.NilObjectID, SendErrorResponse(c, fiber.StatusBadRequest, "Invalid parent ID", err)
	}
	return parentId, nil
}

// parseAndValidateRequest parses body and validates request
func (h *ArrayHandler[T, CreateRequest, UpdateRequest]) parseAndValidateRequest(
	c *fiber.Ctx,
	request interface{},
) error {
	if err := c.BodyParser(request); err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}
	validationResult := h.validator.Validate(request)
	if !validationResult.IsValid {
		return SendValidationErrorResponse(c, validationResult)
	}
	return nil
}

// mapResponse maps result using responseMapper if provided
func mapResponse(data interface{}, responseMapper func(interface{}) interface{}) interface{} {
	if responseMapper != nil {
		return responseMapper(data)
	}
	return data
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
	parentId, err := h.validateAndParseParentId(c, parentIdParam)
	if err != nil {
		return err
	}

	request := new(CreateRequest)
	if err := h.parseAndValidateRequest(c, request); err != nil {
		return err
	}

	// Call pre-create hook if provided
	if len(preCreateHook) > 0 && preCreateHook[0] != nil {
		if err := preCreateHook[0](c, parentId, request); err != nil {
			return err
		}
	}

	element := createFunc(*request)
	if _, err = h.arrayRepo.CreateInArray(parentId, arrayField, element); err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create array element", err)
	}

	return SendSuccessResponse(c, fiber.StatusCreated, "Array element created successfully", mapResponse(element, responseMapper))
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
	parentId, err := h.validateAndParseParentId(c, parentIdParam)
	if err != nil {
		return err
	}

	elementId, err := GetValidObjectId(c, elementIdParam)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid element ID", err)
	}

	request := new(UpdateRequest)
	if err := h.parseAndValidateRequest(c, request); err != nil {
		return err
	}

	if err = h.arrayRepo.UpdateInArray(parentId, arrayField, elementId, updateFunc(*request)); err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to update array element", err)
	}

	updatedElement, err := h.arrayRepo.GetArrayElement(parentId, arrayField, elementId)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve updated element", err)
	}

	return SendSuccessResponse(c, fiber.StatusOK, "Array element updated successfully", mapResponse(updatedElement, responseMapper))
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
