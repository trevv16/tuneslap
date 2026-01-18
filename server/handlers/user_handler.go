package handlers

import (
	api "tuneslap/generated"
	"tuneslap/models"
	"tuneslap/repositories"
	"tuneslap/utils"
	"tuneslap/validation"

	"github.com/gofiber/fiber/v2"
)

// UserHandler provides user-specific operations using the base handler
type UserHandler struct {
	*BaseHandler[models.User, api.UpdateMeRequest, api.UpdateMeRequest]
	userRepo *repositories.UserRepository
}

// NewUserHandler creates a new user handler
func NewUserHandler() *UserHandler {
	userRepo := repositories.NewUserRepository()
	userValidator := validation.NewUserValidator()

	return &UserHandler{
		BaseHandler: NewBaseHandler[models.User, api.UpdateMeRequest, api.UpdateMeRequest](
			userRepo.Repository,
			userValidator,
		),
		userRepo: userRepo,
	}
}

// UserResponseMapper maps user model to response format
func (h *UserHandler) UserResponseMapper(user models.User) interface{} {
	return utils.ToUserResponse(user)
}

// HandleGetMe handles GET current user
func (h *UserHandler) HandleGetMe(c *fiber.Ctx) error {
	return h.HandleGetCurrentUser(c, h.UserResponseMapper)
}

// HandleUpdateMe handles UPDATE current user
func (h *UserHandler) HandleUpdateMe(c *fiber.Ctx) error {
	authorId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}

	// Parse request body using generated type
	var body api.UpdateMeRequest
	if err := c.BodyParser(&body); err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	validationResult := h.validator.Validate(&body)
	if !validationResult.IsValid {
		return SendValidationErrorResponse(c, validationResult)
	}

	// Update the user using the repository
	updatedUser, err := h.userRepo.UpdateMe(authorId, &body)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "User failed to update", err)
	}

	return SendSuccessResponse(c, fiber.StatusOK, "User updated", h.UserResponseMapper(updatedUser))
}
