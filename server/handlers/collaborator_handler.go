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

// CollaboratorHandler provides collaborator-specific operations for embedded collaborators
type CollaboratorHandler struct {
	*ArrayHandler[models.Collaborator, api.CreateCollaboratorRequest, api.UpdateCollaboratorRequest]
	collaboratorRepo *repositories.CollaboratorRepository
	userRepo         *repositories.UserRepository
	boardRepo        *repositories.BoardRepository
}

// NewCollaboratorHandler creates a new collaborator handler
func NewCollaboratorHandler() *CollaboratorHandler {
	collaboratorRepo := repositories.NewCollaboratorRepository()
	userRepo := repositories.NewUserRepository()
	boardRepo := repositories.NewBoardRepository()

	return &CollaboratorHandler{
		ArrayHandler: NewArrayHandler[models.Collaborator, api.CreateCollaboratorRequest, api.UpdateCollaboratorRequest](
			collaboratorRepo.ArrayRepository,
			collaboratorRepo.GetValidator(),
		),
		collaboratorRepo: collaboratorRepo,
		userRepo:         userRepo,
		boardRepo:        boardRepo,
	}
}

// CollaboratorResponseMapper maps collaborator model to response format
func (h *CollaboratorHandler) CollaboratorResponseMapper(collaborator models.Collaborator) interface{} {
	// Convert to API response
	response := utils.ToCollaboratorResponse(collaborator)

	// Try to get user data if user ID exists
	if !collaborator.UserId.IsZero() {
		user, err := h.userRepo.FindOne(bson.M{"_id": collaborator.UserId})
		if err == nil {
			response.Name = user.Name
			response.ImageUrl = &user.ImageUrl
		} else {
			response.Name = "Unknown User"
			response.ImageUrl = nil
		}
	} else {
		response.Name = "Pending Invitation"
		response.ImageUrl = nil
	}

	return response
}

// createCollaboratorFromRequest creates a collaborator from the request data
func (h *CollaboratorHandler) createCollaboratorFromRequest(request api.CreateCollaboratorRequest) models.Collaborator {
	email := request.GetEmail()
	role := request.GetRole()

	// Try to find user by email
	user, err := h.userRepo.GetByEmail(email)

	// Create collaborator with or without user ID
	collaborator := models.Collaborator{
		ID:        primitive.NewObjectID(),
		Email:     email,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err == nil {
		// User exists, set their ID
		collaborator.UserId = user.ID
	} else {
		// User doesn't exist, create placeholder collaborator
		// The user ID will be set when they accept the invitation
		collaborator.UserId = primitive.NilObjectID
	}

	return collaborator
}

// updateCollaboratorFromRequest creates an update document from the request data
func (h *CollaboratorHandler) updateCollaboratorFromRequest(request api.UpdateCollaboratorRequest) bson.M {
	update := bson.M{
		"updatedAt": time.Now(),
	}

	// Role is now a required field (non-pointer)
	if request.Role != "" {
		update["role"] = request.Role
	}

	return update
}

// HandleCreateCollaborator handles CREATE collaborator using array handler
// Requires user to be author, or collaborator with "owner" role
func (h *CollaboratorHandler) HandleCreateCollaborator(c *fiber.Ctx) error {
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

	// Check collaborator management permission
	if !CanManageCollaborators(board, userId) {
		return SendErrorResponse(c, fiber.StatusForbidden, "You do not have permission to manage collaborators on this board", nil)
	}

	return h.HandleCreateInArray(
		c,
		"boardId",
		"collaborators",
		h.createCollaboratorFromRequest,
		func(element interface{}) interface{} {
			if collaborator, ok := element.(models.Collaborator); ok {
				return h.CollaboratorResponseMapper(collaborator)
			}
			return element
		},
		func(c *fiber.Ctx, boardId primitive.ObjectID, request *api.CreateCollaboratorRequest) error {
			// Get the board to check existing collaborators
			board, err := h.boardRepo.FindByID(boardId)
			if err != nil {
				return SendErrorResponse(c, fiber.StatusNotFound, "Board not found", err)
			}

			email := request.GetEmail()

			// Try to find user by email to get userId
			user, userErr := h.userRepo.GetByEmail(email)

			// Check for duplicates by email (for pending invitations)
			for _, existingCollaborator := range board.Collaborators {
				if existingCollaborator.Email == email {
					return SendErrorResponse(c, fiber.StatusBadRequest, "User with this email is already a collaborator on this board", nil)
				}
			}

			// If user exists, check by userId
			if userErr == nil {
				for _, existingCollaborator := range board.Collaborators {
					if existingCollaborator.UserId == user.ID {
						return SendErrorResponse(c, fiber.StatusBadRequest, "User is already a collaborator on this board", nil)
					}
				}
			}
			return nil
		},
	)
}

// HandleUpdateCollaborator handles UPDATE collaborator using array handler
// Requires user to be author, or collaborator with "owner" role
func (h *CollaboratorHandler) HandleUpdateCollaborator(c *fiber.Ctx) error {
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

	// Check collaborator management permission
	if !CanManageCollaborators(board, userId) {
		return SendErrorResponse(c, fiber.StatusForbidden, "You do not have permission to manage collaborators on this board", nil)
	}

	return h.HandleUpdateInArray(
		c,
		"boardId",
		"collaboratorId",
		"collaborators",
		h.updateCollaboratorFromRequest,
		func(element interface{}) interface{} {
			if collaborator, ok := element.(models.Collaborator); ok {
				return h.CollaboratorResponseMapper(collaborator)
			}
			return element
		},
	)
}

// HandleDeleteCollaborator handles DELETE collaborator using array handler
// Requires user to be author, or collaborator with "owner" role
func (h *CollaboratorHandler) HandleDeleteCollaborator(c *fiber.Ctx) error {
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

	// Check collaborator management permission
	if !CanManageCollaborators(board, userId) {
		return SendErrorResponse(c, fiber.StatusForbidden, "You do not have permission to manage collaborators on this board", nil)
	}

	return h.HandleDeleteFromArray(
		c,
		"boardId",        // parentIdParam - this should be the board ID
		"collaboratorId", // elementIdParam
		"collaborators",  // arrayField
	)
}

// HandleGetAllCollaborators handles GET all collaborators using array handler
// Requires user to have access to the board (author or collaborator)
func (h *CollaboratorHandler) HandleGetAllCollaborators(c *fiber.Ctx) error {
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
		"boardId",       // parentIdParam - this should be the board ID
		"collaborators", // arrayField
		func(element interface{}) interface{} {
			if collaborator, ok := element.(models.Collaborator); ok {
				return h.CollaboratorResponseMapper(collaborator)
			}
			return element
		},
	)
}

// HandleGetCollaboratorById handles GET collaborator by ID using array handler
// Requires user to have access to the board (author or collaborator)
func (h *CollaboratorHandler) HandleGetCollaboratorById(c *fiber.Ctx) error {
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
		"boardId",        // parentIdParam - this should be the board ID
		"collaboratorId", // elementIdParam
		"collaborators",  // arrayField
		func(element interface{}) interface{} {
			if collaborator, ok := element.(models.Collaborator); ok {
				return h.CollaboratorResponseMapper(collaborator)
			}
			return element
		},
	)
}
