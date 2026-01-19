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

type KeyRepository struct {
	*ArrayRepository[models.Key]
	boardRepo *BoardRepository
	validator *validation.KeyValidator
}

func NewKeyRepository() *KeyRepository {
	return &KeyRepository{
		ArrayRepository: NewArrayRepository[models.Key]("boards"),
		boardRepo:       NewBoardRepository(),
		validator:       validation.NewKeyValidator(),
	}
}

// GetValidator returns the key validator
func (r *KeyRepository) GetValidator() *validation.KeyValidator {
	return r.validator
}

// CreateKey adds a new key to a board's keys array
func (r *KeyRepository) CreateKey(
	boardId primitive.ObjectID,
	createData *api.CreateKeyRequest,
) (models.Board, error) {
	// Validate the request
	validationResult := r.validator.ValidateCreateKey(createData)
	if !validationResult.IsValid {
		return models.Board{}, validation.NewValidationError(validationResult.Errors)
	}

	// Create the new key using converter
	newKey := utils.KeyFromCreateRequest(createData, boardId)

	// Add to the board's keys array
	_, err := r.CreateInArray(boardId, "keys", newKey)
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

// UpdateKey updates a specific key in a board's keys array
func (r *KeyRepository) UpdateKey(
	boardId primitive.ObjectID,
	keyId primitive.ObjectID,
	updateData *api.UpdateKeyRequest,
) (models.Board, error) {
	// Validate the request
	validationResult := r.validator.ValidateUpdateKey(updateData)
	if !validationResult.IsValid {
		return models.Board{}, validation.NewValidationError(validationResult.Errors)
	}

	// Build update document from pointer fields
	update := bson.M{
		"updatedAt": time.Now(),
	}

	if updateData.Name != nil {
		update["name"] = *updateData.Name
	}
	if updateData.Description != nil {
		update["description"] = *updateData.Description
	}
	if updateData.HotKey != nil {
		update["hotKey"] = *updateData.HotKey
	}
	if updateData.AudioMediaId != nil && *updateData.AudioMediaId != "" {
		audioMediaId, _ := primitive.ObjectIDFromHex(*updateData.AudioMediaId)
		update["audioMediaId"] = audioMediaId
	}
	if updateData.ImageMediaId != nil && *updateData.ImageMediaId != "" {
		imageMediaId, _ := primitive.ObjectIDFromHex(*updateData.ImageMediaId)
		update["imageMediaId"] = imageMediaId
	}

	err := r.UpdateInArray(boardId, "keys", keyId, update)
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

// DeleteKey removes a key from a board's keys array
func (r *KeyRepository) DeleteKey(
	boardId primitive.ObjectID,
	keyId primitive.ObjectID,
) (models.Board, error) {
	// Remove the key from the array
	err := r.DeleteFromArray(boardId, "keys", keyId)
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

// DeleteKeysByMediaId removes all keys that reference a specific media ID
// (either as audioMediaId or imageMediaId) from all boards
func (r *KeyRepository) DeleteKeysByMediaId(mediaId primitive.ObjectID) (int64, error) {
	collection := r.getCollection()
	ctx, cancel := r.getContext()
	defer cancel()

	// Pull keys from all boards where the key references this media ID
	update := bson.M{
		"$pull": bson.M{
			"keys": bson.M{
				"$or": []bson.M{
					{"audioMediaId": mediaId},
					{"imageMediaId": mediaId},
				},
			},
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}

	// Update all boards that have keys referencing this media
	result, err := collection.UpdateMany(ctx, bson.M{}, update)
	if err != nil {
		return 0, err
	}

	// Clear cache after deletion
	r.clearCache()

	return result.ModifiedCount, nil
}
