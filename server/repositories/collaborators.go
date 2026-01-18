package repositories

import (
	"time"
	api "tuneslap/generated"
	"tuneslap/models"
	"tuneslap/utils"
	"tuneslap/validation"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CollaboratorRepository struct {
	*ArrayRepository[models.Collaborator]
	boardRepo *BoardRepository
	userRepo  *UserRepository
	validator *validation.CollaboratorValidator
}

func NewCollaboratorRepository() *CollaboratorRepository {
	return &CollaboratorRepository{
		ArrayRepository: NewArrayRepository[models.Collaborator]("boards"),
		boardRepo:       NewBoardRepository(),
		userRepo:        NewUserRepository(),
		validator:       validation.NewCollaboratorValidator(),
	}
}

// GetValidator returns the collaborator validator
func (r *CollaboratorRepository) GetValidator() *validation.CollaboratorValidator {
	return r.validator
}

// GetBoardRepo returns the board repository
func (r *CollaboratorRepository) GetBoardRepo() *BoardRepository {
	return r.boardRepo
}

// CreateCollaborator adds a new collaborator to a board's collaborators array
func (r *CollaboratorRepository) CreateCollaborator(
	boardId primitive.ObjectID,
	createData *api.CreateCollaboratorRequest,
	authorId primitive.ObjectID,
) (models.Board, error) {
	// Validate the request
	validationResult := r.validator.ValidateCreateCollaborator(createData)
	if !validationResult.IsValid {
		return models.Board{}, validation.NewValidationError(validationResult.Errors)
	}

	email := createData.GetEmail()

	// Try to find user by email
	user, userErr := r.userRepo.GetByEmail(email)

	// Create the new collaborator using converter
	newCollaborator := utils.CollaboratorFromCreateRequest(createData)

	if userErr == nil {
		// User exists, set their ID
		newCollaborator.UserId = user.ID
	} else {
		// User doesn't exist, create placeholder collaborator
		// The user ID will be set when they accept the invitation
		newCollaborator.UserId = primitive.NilObjectID
	}

	// Add to the board's collaborators array
	_, createErr := r.CreateInArray(boardId, "collaborators", newCollaborator)
	if createErr != nil {
		return models.Board{}, createErr
	}

	// Return the updated board
	updatedBoard, err := r.boardRepo.FindByID(boardId)
	if err != nil {
		return models.Board{}, err
	}

	return updatedBoard, nil
}

// UpdateCollaborator updates a specific collaborator in a board's collaborators array
func (r *CollaboratorRepository) UpdateCollaborator(
	boardId primitive.ObjectID,
	collaboratorId primitive.ObjectID,
	updateData *api.UpdateCollaboratorRequest,
) (models.Board, error) {
	// Validate the request
	validationResult := r.validator.ValidateUpdateCollaborator(updateData)
	if !validationResult.IsValid {
		return models.Board{}, validation.NewValidationError(validationResult.Errors)
	}

	// Update the collaborator in the array
	update := bson.M{
		"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
	}

	if updateData.Role != nil {
		update["role"] = updateData.GetRole()
	}

	err := r.UpdateInArray(boardId, "collaborators", collaboratorId, update)
	if err != nil {
		return models.Board{}, err
	}

	// Return the updated board
	board, err := r.boardRepo.FindByID(boardId)
	if err != nil {
		return models.Board{}, err
	}

	return board, nil
}

// DeleteCollaborator removes a collaborator from a board's collaborators array
func (r *CollaboratorRepository) DeleteCollaborator(
	boardId primitive.ObjectID,
	collaboratorId primitive.ObjectID,
) (models.Board, error) {
	// Remove the collaborator from the array
	err := r.DeleteFromArray(boardId, "collaborators", collaboratorId)
	if err != nil {
		return models.Board{}, err
	}

	// Return the updated board
	board, err := r.boardRepo.FindByID(boardId)
	if err != nil {
		return models.Board{}, err
	}

	return board, nil
}
