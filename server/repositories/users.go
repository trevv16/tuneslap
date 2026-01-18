package repositories

import (
	"time"
	api "tuneslap/generated"
	"tuneslap/models"
	"tuneslap/validation"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository struct {
	*Repository[models.User]
	validator *validation.UserValidator
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		Repository: NewRepository[models.User]("users"),
		validator:  validation.NewUserValidator(),
	}
}

// GetValidator returns the user validator
func (r *UserRepository) GetValidator() *validation.UserValidator {
	return r.validator
}

// GetMe retrieves a user by their ID
func (r *UserRepository) GetMe(authorId primitive.ObjectID) (models.User, error) {
	return r.FindOne(bson.M{"_id": authorId})
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (models.User, error) {
	return r.FindOne(bson.M{"email": email})
}

// GetByResetToken retrieves a user by reset token
func (r *UserRepository) GetByResetToken(token string) (models.User, error) {
	return r.FindOne(bson.M{"resetToken": token})
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(user models.User) (models.User, error) {
	newUser := models.User{
		ID:           primitive.NewObjectID(),
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return r.Create(newUser)
}

// UpdateMe updates the current user's profile
func (r *UserRepository) UpdateMe(authorId primitive.ObjectID, updateData *api.UpdateMeRequest) (models.User, error) {
	update := bson.M{
		"$set": bson.M{},
	}

	// Only include fields that are not nil
	if updateData.Name != nil {
		update["$set"].(bson.M)["name"] = updateData.GetName()
	}
	if updateData.ImageUrl != nil {
		update["$set"].(bson.M)["imageUrl"] = updateData.GetImageUrl()
	}

	// Only update if there are fields to update
	if len(update["$set"].(bson.M)) == 0 {
		// No fields to update, just return the current user
		return r.FindByID(authorId)
	}

	return r.Update(authorId, update)
}

// UpdateUser updates any user by ID
func (r *UserRepository) UpdateUser(id primitive.ObjectID, updateData *models.UpdateUserRequest) (models.User, error) {
	update := bson.M{
		"$set": bson.M{
			"name":           updateData.Name,
			"imageUrl":       updateData.ImageUrl,
			"resetToken":     updateData.ResetToken,
			"resetExpiresAt": updateData.ResetExpiresAt,
			"passwordHash":   updateData.PasswordHash,
		},
	}
	return r.Update(id, update)
}
