package repositories

import (
	api "tuneslap/generated"
	"tuneslap/models"
	"tuneslap/utils"
	"tuneslap/validation"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BoardRepository struct {
	*Repository[models.Board]
	validator *validation.BoardValidator
}

func NewBoardRepository() *BoardRepository {
	return &BoardRepository{
		Repository: NewRepository[models.Board]("boards"),
		validator:  validation.NewBoardValidator(),
	}
}

// GetValidator returns the board validator
func (r *BoardRepository) GetValidator() *validation.BoardValidator {
	return r.validator
}

// GetByAuthor retrieves all boards for a specific author
func (r *BoardRepository) GetByAuthor(authorId primitive.ObjectID) ([]models.Board, error) {
	return r.FindAll(bson.M{"authorId": authorId})
}

// GetById retrieves a board by ID for a specific author
func (r *BoardRepository) GetById(id primitive.ObjectID, authorId primitive.ObjectID) (models.Board, error) {
	return r.FindOne(bson.M{
		"_id":      id,
		"authorId": authorId,
	})
}

// CreateBoard creates a new board
func (r *BoardRepository) CreateBoard(board *api.CreateBoardRequest, authorId primitive.ObjectID) (models.Board, error) {
	newBoard := utils.BoardFromCreateRequest(board, authorId)
	return r.Create(newBoard)
}

// UpdateBoard updates a board by ID for a specific author
func (r *BoardRepository) UpdateBoard(id primitive.ObjectID, authorId primitive.ObjectID, updateData *api.UpdateBoardRequest) (models.Board, error) {
	update := bson.M{
		"$set": bson.M{},
	}

	// Only include fields that are not nil
	if updateData.Name != nil {
		update["$set"].(bson.M)["name"] = updateData.GetName()
	}
	if updateData.Description != nil {
		update["$set"].(bson.M)["description"] = updateData.GetDescription()
	}
	if updateData.Layout != nil {
		update["$set"].(bson.M)["layout"] = models.LayoutType(updateData.GetLayout())
	}
	if updateData.ImageUrl != nil {
		update["$set"].(bson.M)["imageUrl"] = updateData.GetImageUrl()
	}

	// Only update if there are fields to update
	if len(update["$set"].(bson.M)) == 0 {
		// No fields to update, just return the current board
		return r.FindOne(bson.M{
			"_id":      id,
			"authorId": authorId,
		})
	}

	// First update the document
	_, err := r.Update(id, update)
	if err != nil {
		return models.Board{}, err
	}

	// Then fetch the updated document with author filter
	return r.FindOne(bson.M{
		"_id":      id,
		"authorId": authorId,
	})
}

// DeleteBoard deletes a board by ID for a specific author
func (r *BoardRepository) DeleteBoard(id primitive.ObjectID, authorId primitive.ObjectID) error {
	// We need to check if the board exists and belongs to the author first
	_, err := r.FindOne(bson.M{
		"_id":      id,
		"authorId": authorId,
	})
	if err != nil {
		return err
	}

	return r.Delete(id)
}

// AggregateBoards performs aggregation on boards
func (r *BoardRepository) AggregateBoards(pipeline interface{}, authorId primitive.ObjectID) ([]models.Board, int64, error) {
	return r.AggregateWithCount(pipeline, bson.M{"authorId": authorId})
}
