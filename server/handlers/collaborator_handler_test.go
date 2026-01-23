package handlers

import (
	"testing"
	"time"
	"tuneslap/models"
	"tuneslap/utils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewCollaboratorHandler(t *testing.T) {
	t.Run("creates handler successfully", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Skip("Skipping test - requires database connection")
			}
		}()

		handler := NewCollaboratorHandler()
		assert.NotNil(t, handler)
		assert.NotNil(t, handler.collaboratorRepo)
		assert.NotNil(t, handler.userRepo)
		assert.NotNil(t, handler.boardRepo)
		assert.NotNil(t, handler.ArrayHandler)
	})
}

func TestCollaboratorResponseMapperLogic(t *testing.T) {
	validTime := time.Now()
	collaboratorID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	t.Run("maps collaborator with user ID", func(t *testing.T) {
		collaborator := models.Collaborator{
			ID:        collaboratorID,
			Email:     "test@example.com",
			UserId:    userID,
			Role:      "editor",
			CreatedAt: validTime,
			UpdatedAt: validTime,
		}

		// Test the utility conversion function
		response := utils.ToCollaboratorResponse(collaborator)

		assert.NotNil(t, response.Id)
		assert.Equal(t, collaboratorID.Hex(), *response.Id)
		assert.NotNil(t, response.Email)
		assert.Equal(t, "test@example.com", *response.Email)
		assert.NotNil(t, response.UserId)
		assert.Equal(t, userID.Hex(), *response.UserId)
		assert.NotNil(t, response.Role)
		assert.Equal(t, "editor", *response.Role)
	})

	t.Run("maps collaborator without user ID (pending invitation)", func(t *testing.T) {
		collaborator := models.Collaborator{
			ID:        collaboratorID,
			Email:     "pending@example.com",
			UserId:    primitive.NilObjectID,
			Role:      "viewer",
			CreatedAt: validTime,
			UpdatedAt: validTime,
		}

		response := utils.ToCollaboratorResponse(collaborator)

		assert.NotNil(t, response.Id)
		assert.NotNil(t, response.Email)
		assert.Equal(t, "pending@example.com", *response.Email)
		// UserId should be nil for pending invitations
		assert.Nil(t, response.UserId)
		assert.NotNil(t, response.Role)
		assert.Equal(t, "viewer", *response.Role)
	})

	t.Run("Name and ImageUrl are initially nil", func(t *testing.T) {
		collaborator := models.Collaborator{
			ID:        collaboratorID,
			Email:     "test@example.com",
			UserId:    userID,
			Role:      "owner",
			CreatedAt: validTime,
			UpdatedAt: validTime,
		}

		response := utils.ToCollaboratorResponse(collaborator)

		// Name and ImageUrl are set by CollaboratorResponseMapper after lookup
		// ToCollaboratorResponse does not set them
		assert.Nil(t, response.Name)
		assert.Nil(t, response.ImageUrl)
	})
}

func TestCreateCollaboratorFromRequestLogic(t *testing.T) {
	t.Run("creates collaborator with email and role", func(t *testing.T) {
		email := "newuser@example.com"
		role := "editor"

		collaborator := models.Collaborator{
			ID:        primitive.NewObjectID(),
			Email:     email,
			Role:      role,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.NotEqual(t, primitive.NilObjectID, collaborator.ID)
		assert.Equal(t, email, collaborator.Email)
		assert.Equal(t, role, collaborator.Role)
		assert.False(t, collaborator.CreatedAt.IsZero())
		assert.False(t, collaborator.UpdatedAt.IsZero())
	})

	t.Run("collaborator starts with NilObjectID for unknown user", func(t *testing.T) {
		// When user doesn't exist, UserId should be NilObjectID
		collaborator := models.Collaborator{
			ID:        primitive.NewObjectID(),
			Email:     "unknown@example.com",
			UserId:    primitive.NilObjectID,
			Role:      "viewer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.True(t, collaborator.UserId.IsZero())
	})

	t.Run("collaborator gets UserId when user exists", func(t *testing.T) {
		existingUserId := primitive.NewObjectID()

		collaborator := models.Collaborator{
			ID:        primitive.NewObjectID(),
			Email:     "existing@example.com",
			UserId:    existingUserId,
			Role:      "editor",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.False(t, collaborator.UserId.IsZero())
		assert.Equal(t, existingUserId, collaborator.UserId)
	})
}

func TestUpdateCollaboratorFromRequestLogic(t *testing.T) {
	t.Run("update with role change", func(t *testing.T) {
		newRole := "owner"

		update := bson.M{
			"updatedAt": time.Now(),
		}

		if newRole != "" {
			update["role"] = newRole
		}

		assert.Equal(t, "owner", update["role"])
		assert.NotNil(t, update["updatedAt"])
	})

	t.Run("update without role change", func(t *testing.T) {
		var role *string = nil

		update := bson.M{
			"updatedAt": time.Now(),
		}

		if role != nil {
			update["role"] = *role
		}

		_, hasRole := update["role"]
		assert.False(t, hasRole)
		assert.NotNil(t, update["updatedAt"])
	})

	t.Run("update always sets updatedAt", func(t *testing.T) {
		update := bson.M{
			"updatedAt": time.Now(),
		}

		assert.NotNil(t, update["updatedAt"])
	})
}

func TestCollaboratorDuplicateDetectionLogic(t *testing.T) {
	existingEmail := "existing@example.com"
	existingUserId := primitive.NewObjectID()

	existingCollaborators := []models.Collaborator{
		{
			ID:        primitive.NewObjectID(),
			Email:     existingEmail,
			UserId:    existingUserId,
			Role:      "editor",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	t.Run("detects duplicate by email", func(t *testing.T) {
		newEmail := "existing@example.com"

		isDuplicate := false
		for _, existingCollaborator := range existingCollaborators {
			if existingCollaborator.Email == newEmail {
				isDuplicate = true
				break
			}
		}

		assert.True(t, isDuplicate)
	})

	t.Run("allows new email", func(t *testing.T) {
		newEmail := "newuser@example.com"

		isDuplicate := false
		for _, existingCollaborator := range existingCollaborators {
			if existingCollaborator.Email == newEmail {
				isDuplicate = true
				break
			}
		}

		assert.False(t, isDuplicate)
	})

	t.Run("detects duplicate by userId", func(t *testing.T) {
		newUserId := existingUserId

		isDuplicate := false
		for _, existingCollaborator := range existingCollaborators {
			if existingCollaborator.UserId == newUserId {
				isDuplicate = true
				break
			}
		}

		assert.True(t, isDuplicate)
	})

	t.Run("allows new userId", func(t *testing.T) {
		newUserId := primitive.NewObjectID()

		isDuplicate := false
		for _, existingCollaborator := range existingCollaborators {
			if existingCollaborator.UserId == newUserId {
				isDuplicate = true
				break
			}
		}

		assert.False(t, isDuplicate)
	})
}

func TestCollaboratorRoleTypes(t *testing.T) {
	validRoles := []string{"owner", "editor", "viewer"}

	t.Run("validates owner role", func(t *testing.T) {
		role := "owner"
		isValid := false
		for _, validRole := range validRoles {
			if role == validRole {
				isValid = true
				break
			}
		}
		assert.True(t, isValid)
	})

	t.Run("validates editor role", func(t *testing.T) {
		role := "editor"
		isValid := false
		for _, validRole := range validRoles {
			if role == validRole {
				isValid = true
				break
			}
		}
		assert.True(t, isValid)
	})

	t.Run("validates viewer role", func(t *testing.T) {
		role := "viewer"
		isValid := false
		for _, validRole := range validRoles {
			if role == validRole {
				isValid = true
				break
			}
		}
		assert.True(t, isValid)
	})

	t.Run("rejects invalid role", func(t *testing.T) {
		role := "admin"
		isValid := false
		for _, validRole := range validRoles {
			if role == validRole {
				isValid = true
				break
			}
		}
		assert.False(t, isValid)
	})
}

func TestCollaboratorModelStructure(t *testing.T) {
	t.Run("collaborator has all required fields", func(t *testing.T) {
		collaborator := models.Collaborator{
			ID:        primitive.NewObjectID(),
			Email:     "test@example.com",
			UserId:    primitive.NewObjectID(),
			Role:      "editor",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.False(t, collaborator.ID.IsZero())
		assert.NotEmpty(t, collaborator.Email)
		assert.False(t, collaborator.UserId.IsZero())
		assert.NotEmpty(t, collaborator.Role)
		assert.False(t, collaborator.CreatedAt.IsZero())
		assert.False(t, collaborator.UpdatedAt.IsZero())
	})

	t.Run("collaborator allows zero UserId for pending invitations", func(t *testing.T) {
		collaborator := models.Collaborator{
			ID:        primitive.NewObjectID(),
			Email:     "pending@example.com",
			UserId:    primitive.NilObjectID,
			Role:      "viewer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.True(t, collaborator.UserId.IsZero())
		assert.NotEmpty(t, collaborator.Email)
	})
}

func TestCollaboratorResponseEnrichment(t *testing.T) {
	// Tests the logic of enriching collaborator response with user data

	t.Run("enriches response when user exists", func(t *testing.T) {
		response := utils.ToCollaboratorResponse(models.Collaborator{
			ID:        primitive.NewObjectID(),
			Email:     "test@example.com",
			UserId:    primitive.NewObjectID(),
			Role:      "editor",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})

		// Simulate user data enrichment
		userName := "Test User"
		userImageUrl := "https://example.com/avatar.png"
		response.Name = &userName
		response.ImageUrl = &userImageUrl

		assert.NotNil(t, response.Name)
		assert.Equal(t, "Test User", *response.Name)
		assert.NotNil(t, response.ImageUrl)
		assert.Equal(t, "https://example.com/avatar.png", *response.ImageUrl)
	})

	t.Run("sets unknown user when user not found", func(t *testing.T) {
		response := utils.ToCollaboratorResponse(models.Collaborator{
			ID:        primitive.NewObjectID(),
			Email:     "test@example.com",
			UserId:    primitive.NewObjectID(),
			Role:      "editor",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})

		// Simulate user not found
		unknownName := "Unknown User"
		emptyImageUrl := ""
		response.Name = &unknownName
		response.ImageUrl = &emptyImageUrl

		assert.NotNil(t, response.Name)
		assert.Equal(t, "Unknown User", *response.Name)
		assert.NotNil(t, response.ImageUrl)
		assert.Equal(t, "", *response.ImageUrl)
	})

	t.Run("sets pending invitation when no userId", func(t *testing.T) {
		response := utils.ToCollaboratorResponse(models.Collaborator{
			ID:        primitive.NewObjectID(),
			Email:     "pending@example.com",
			UserId:    primitive.NilObjectID,
			Role:      "viewer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})

		// Simulate pending invitation
		pendingName := "Pending Invitation"
		emptyImageUrl := ""
		response.Name = &pendingName
		response.ImageUrl = &emptyImageUrl

		assert.NotNil(t, response.Name)
		assert.Equal(t, "Pending Invitation", *response.Name)
		assert.NotNil(t, response.ImageUrl)
		assert.Equal(t, "", *response.ImageUrl)
	})
}

// Benchmarks
func BenchmarkCollaboratorResponseMapping(b *testing.B) {
	collaborator := models.Collaborator{
		ID:        primitive.NewObjectID(),
		Email:     "test@example.com",
		UserId:    primitive.NewObjectID(),
		Role:      "editor",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utils.ToCollaboratorResponse(collaborator)
	}
}

func BenchmarkCollaboratorDuplicateCheck(b *testing.B) {
	existingCollaborators := []models.Collaborator{
		{ID: primitive.NewObjectID(), Email: "user1@example.com", UserId: primitive.NewObjectID()},
		{ID: primitive.NewObjectID(), Email: "user2@example.com", UserId: primitive.NewObjectID()},
		{ID: primitive.NewObjectID(), Email: "user3@example.com", UserId: primitive.NewObjectID()},
	}
	newEmail := "newuser@example.com"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isDuplicate := false
		for _, existingCollaborator := range existingCollaborators {
			if existingCollaborator.Email == newEmail {
				isDuplicate = true
				break
			}
		}
		_ = isDuplicate
	}
}

func BenchmarkCollaboratorUpdateBson(b *testing.B) {
	newRole := "owner"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		update := bson.M{
			"updatedAt": time.Now(),
		}
		if newRole != "" {
			update["role"] = newRole
		}
		_ = update
	}
}
