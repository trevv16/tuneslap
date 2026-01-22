package handlers

import (
	"time"
	api "tuneslap/generated"
	"tuneslap/models"
	"tuneslap/repositories"
	"tuneslap/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// KeyHandler provides key-specific operations for embedded keys
type KeyHandler struct {
	*ArrayHandler[models.Key, api.CreateKeyRequest, api.UpdateKeyRequest]
	keyRepo *repositories.KeyRepository
}

// NewKeyHandler creates a new key handler
func NewKeyHandler() *KeyHandler {
	keyRepo := repositories.NewKeyRepository()

	return &KeyHandler{
		ArrayHandler: NewArrayHandler[models.Key, api.CreateKeyRequest, api.UpdateKeyRequest](
			keyRepo.ArrayRepository,
			keyRepo.GetValidator(),
		),
		keyRepo: keyRepo,
	}
}

// KeyResponseMapper maps key model to response format
func (h *KeyHandler) KeyResponseMapper(key models.Key) interface{} {
	return utils.ToKeyResponse(key)
}

// createKeyFromRequest creates a key from the request data
func (h *KeyHandler) createKeyFromRequest(request api.CreateKeyRequest) models.Key {
	// Convert AudioMediaId from string to ObjectID
	var audioMediaId primitive.ObjectID
	if request.AudioMediaId != "" {
		audioMediaId, _ = primitive.ObjectIDFromHex(request.AudioMediaId)
	}

	// Convert ImageMediaId from *string to ObjectID
	var imageMediaId primitive.ObjectID
	if request.ImageMediaId != nil && *request.ImageMediaId != "" {
		imageMediaId, _ = primitive.ObjectIDFromHex(*request.ImageMediaId)
	}

	description := ""
	if request.Description != nil {
		description = *request.Description
	}

	return models.Key{
		ID:           primitive.NewObjectID(),
		BoardId:      primitive.ObjectID{}, // Will be set by repository
		Name:         request.Name,
		Description:  description,
		AudioMediaId: audioMediaId,
		ImageMediaId: imageMediaId,
		HotKey:       request.HotKey,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// updateKeyFromRequest creates an update document from the request data
func (h *KeyHandler) updateKeyFromRequest(request api.UpdateKeyRequest) bson.M {
	update := bson.M{
		"updatedAt": time.Now(),
	}

	if request.Name != nil {
		update["name"] = *request.Name
	}
	if request.Description != nil {
		update["description"] = *request.Description
	}
	if request.HotKey != nil {
		update["hotKey"] = *request.HotKey
	}
	if request.AudioMediaId != nil && *request.AudioMediaId != "" {
		audioMediaId, _ := primitive.ObjectIDFromHex(*request.AudioMediaId)
		update["audioMediaId"] = audioMediaId
	}
	if request.ImageMediaId != nil && *request.ImageMediaId != "" {
		imageMediaId, _ := primitive.ObjectIDFromHex(*request.ImageMediaId)
		update["imageMediaId"] = imageMediaId
	}

	return update
}

// HandleCreateKey handles CREATE key using array handler
// Requires user to be author, or collaborator with "owner" or "editor" role
func (h *KeyHandler) HandleCreateKey(c *fiber.Ctx) error {
	userId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	// Parse boardId
	boardId, err := GetValidObjectId(c, "boardId")
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid board ID", err)
	}

	// Check board access
	board, err := CheckBoardAccess(boardId, userId)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusNotFound, "Board not found", err)
	}

	// Check key management permission
	if !CanManageKeys(board, userId) {
		return SendErrorResponse(c, fiber.StatusForbidden, "You do not have permission to manage keys on this board", nil)
	}

	return h.HandleCreateInArray(
		c,
		"boardId",
		"keys",
		func(req api.CreateKeyRequest) models.Key {
			key := h.createKeyFromRequest(req)
			key.BoardId = boardId
			return key
		},
		func(element interface{}) interface{} {
			if key, ok := element.(models.Key); ok {
				return h.KeyResponseMapper(key)
			}
			return element
		},
	)
}

// HandleUpdateKey handles UPDATE key using array handler
// Requires user to be author, or collaborator with "owner" or "editor" role
func (h *KeyHandler) HandleUpdateKey(c *fiber.Ctx) error {
	userId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	// Parse boardId
	boardId, err := GetValidObjectId(c, "boardId")
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid board ID", err)
	}

	// Check board access
	board, err := CheckBoardAccess(boardId, userId)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusNotFound, "Board not found", err)
	}

	// Check key management permission
	if !CanManageKeys(board, userId) {
		return SendErrorResponse(c, fiber.StatusForbidden, "You do not have permission to manage keys on this board", nil)
	}

	return h.HandleUpdateInArray(
		c,
		"boardId",
		"keyId",
		"keys",
		h.updateKeyFromRequest,
		func(element interface{}) interface{} {
			if key, ok := element.(models.Key); ok {
				return h.KeyResponseMapper(key)
			}
			return element
		},
	)
}

// HandleDeleteKey handles DELETE key using array handler
// Requires user to be author, or collaborator with "owner" or "editor" role
func (h *KeyHandler) HandleDeleteKey(c *fiber.Ctx) error {
	userId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	// Parse boardId
	boardId, err := GetValidObjectId(c, "boardId")
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid board ID", err)
	}

	// Check board access
	board, err := CheckBoardAccess(boardId, userId)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusNotFound, "Board not found", err)
	}

	// Check key management permission
	if !CanManageKeys(board, userId) {
		return SendErrorResponse(c, fiber.StatusForbidden, "You do not have permission to manage keys on this board", nil)
	}

	return h.HandleDeleteFromArray(
		c,
		"boardId", // parentIdParam - this should be the board ID
		"keyId",   // elementIdParam
		"keys",    // arrayField
	)
}

// HandleGetAllKeys handles GET all keys from a board using array handler
// Requires user to have access to the board (author or collaborator)
func (h *KeyHandler) HandleGetAllKeys(c *fiber.Ctx) error {
	userId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	// Parse boardId
	boardId, err := GetValidObjectId(c, "boardId")
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid board ID", err)
	}

	// Check board access
	_, err = CheckBoardAccess(boardId, userId)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusNotFound, "Board not found", err)
	}

	return h.HandleGetAllFromArray(
		c,
		"boardId", // parentIdParam - this should be the board ID
		"keys",    // arrayField
		func(element interface{}) interface{} {
			if key, ok := element.(models.Key); ok {
				return h.KeyResponseMapper(key)
			}
			return element
		},
	)
}

// HandleGetKeyById handles GET key by ID from a board using array handler
// Requires user to have access to the board (author or collaborator)
func (h *KeyHandler) HandleGetKeyById(c *fiber.Ctx) error {
	userId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	// Parse boardId
	boardId, err := GetValidObjectId(c, "boardId")
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid board ID", err)
	}

	// Check board access
	_, err = CheckBoardAccess(boardId, userId)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusNotFound, "Board not found", err)
	}

	return h.HandleGetByIdFromArray(
		c,
		"boardId", // parentIdParam - this should be the board ID
		"keyId",   // elementIdParam
		"keys",    // arrayField
		func(element interface{}) interface{} {
			if key, ok := element.(models.Key); ok {
				return h.KeyResponseMapper(key)
			}
			return element
		},
	)
}
