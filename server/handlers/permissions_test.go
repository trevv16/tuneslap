package handlers

import (
	"testing"
	"time"
	"tuneslap/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetUserBoardRole(t *testing.T) {
	authorID := primitive.NewObjectID()
	collabID := primitive.NewObjectID()
	otherUserID := primitive.NewObjectID()

	tests := []struct {
		name         string
		board        models.Board
		userID       primitive.ObjectID
		expectedRole string
	}{
		{
			name: "user is author",
			board: models.Board{
				ID:            primitive.NewObjectID(),
				AuthorId:      authorID,
				Collaborators: []models.Collaborator{},
			},
			userID:       authorID,
			expectedRole: "author",
		},
		{
			name: "user is editor collaborator",
			board: models.Board{
				ID:       primitive.NewObjectID(),
				AuthorId: authorID,
				Collaborators: []models.Collaborator{
					{
						ID:     primitive.NewObjectID(),
						UserId: collabID,
						Role:   "editor",
					},
				},
			},
			userID:       collabID,
			expectedRole: "editor",
		},
		{
			name: "user is viewer collaborator",
			board: models.Board{
				ID:       primitive.NewObjectID(),
				AuthorId: authorID,
				Collaborators: []models.Collaborator{
					{
						ID:     primitive.NewObjectID(),
						UserId: collabID,
						Role:   "viewer",
					},
				},
			},
			userID:       collabID,
			expectedRole: "viewer",
		},
		{
			name: "user has no access",
			board: models.Board{
				ID:            primitive.NewObjectID(),
				AuthorId:      authorID,
				Collaborators: []models.Collaborator{},
			},
			userID:       otherUserID,
			expectedRole: "",
		},
		{
			name: "user in multiple collaborators",
			board: models.Board{
				ID:       primitive.NewObjectID(),
				AuthorId: authorID,
				Collaborators: []models.Collaborator{
					{ID: primitive.NewObjectID(), UserId: primitive.NewObjectID(), Role: "viewer"},
					{ID: primitive.NewObjectID(), UserId: collabID, Role: "editor"},
					{ID: primitive.NewObjectID(), UserId: primitive.NewObjectID(), Role: "viewer"},
				},
			},
			userID:       collabID,
			expectedRole: "editor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role := GetUserBoardRole(tt.board, tt.userID)
			assert.Equal(t, tt.expectedRole, role)
		})
	}
}

func TestCanEditBoard(t *testing.T) {
	authorID := primitive.NewObjectID()
	editorID := primitive.NewObjectID()
	viewerID := primitive.NewObjectID()
	otherUserID := primitive.NewObjectID()

	board := models.Board{
		ID:       primitive.NewObjectID(),
		AuthorId: authorID,
		Collaborators: []models.Collaborator{
			{ID: primitive.NewObjectID(), UserId: editorID, Role: "editor"},
			{ID: primitive.NewObjectID(), UserId: viewerID, Role: "viewer"},
		},
	}

	tests := []struct {
		name     string
		userID   primitive.ObjectID
		expected bool
	}{
		{
			name:     "author can edit",
			userID:   authorID,
			expected: true,
		},
		{
			name:     "editor can edit",
			userID:   editorID,
			expected: true,
		},
		{
			name:     "viewer cannot edit",
			userID:   viewerID,
			expected: false,
		},
		{
			name:     "other user cannot edit",
			userID:   otherUserID,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CanEditBoard(board, tt.userID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCanDeleteBoard(t *testing.T) {
	authorID := primitive.NewObjectID()
	editorID := primitive.NewObjectID()
	viewerID := primitive.NewObjectID()

	board := models.Board{
		ID:       primitive.NewObjectID(),
		AuthorId: authorID,
		Collaborators: []models.Collaborator{
			{ID: primitive.NewObjectID(), UserId: editorID, Role: "editor"},
			{ID: primitive.NewObjectID(), UserId: viewerID, Role: "viewer"},
		},
	}

	tests := []struct {
		name     string
		userID   primitive.ObjectID
		expected bool
	}{
		{
			name:     "author can delete",
			userID:   authorID,
			expected: true,
		},
		{
			name:     "editor cannot delete",
			userID:   editorID,
			expected: false,
		},
		{
			name:     "viewer cannot delete",
			userID:   viewerID,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CanDeleteBoard(board, tt.userID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCanManageCollaborators(t *testing.T) {
	authorID := primitive.NewObjectID()
	editorID := primitive.NewObjectID()
	viewerID := primitive.NewObjectID()

	board := models.Board{
		ID:       primitive.NewObjectID(),
		AuthorId: authorID,
		Collaborators: []models.Collaborator{
			{ID: primitive.NewObjectID(), UserId: editorID, Role: "editor"},
			{ID: primitive.NewObjectID(), UserId: viewerID, Role: "viewer"},
		},
	}

	tests := []struct {
		name     string
		userID   primitive.ObjectID
		expected bool
	}{
		{
			name:     "author can manage collaborators",
			userID:   authorID,
			expected: true,
		},
		{
			name:     "editor cannot manage collaborators",
			userID:   editorID,
			expected: false,
		},
		{
			name:     "viewer cannot manage collaborators",
			userID:   viewerID,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CanManageCollaborators(board, tt.userID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCanManageKeys(t *testing.T) {
	authorID := primitive.NewObjectID()
	editorID := primitive.NewObjectID()
	viewerID := primitive.NewObjectID()

	board := models.Board{
		ID:       primitive.NewObjectID(),
		AuthorId: authorID,
		Collaborators: []models.Collaborator{
			{ID: primitive.NewObjectID(), UserId: editorID, Role: "editor"},
			{ID: primitive.NewObjectID(), UserId: viewerID, Role: "viewer"},
		},
	}

	tests := []struct {
		name     string
		userID   primitive.ObjectID
		expected bool
	}{
		{
			name:     "author can manage keys",
			userID:   authorID,
			expected: true,
		},
		{
			name:     "editor can manage keys",
			userID:   editorID,
			expected: true,
		},
		{
			name:     "viewer cannot manage keys",
			userID:   viewerID,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CanManageKeys(board, tt.userID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBoardWithCollaboratorOwner(t *testing.T) {
	authorID := primitive.NewObjectID()
	ownerCollabID := primitive.NewObjectID()

	board := models.Board{
		ID:       primitive.NewObjectID(),
		AuthorId: authorID,
		Collaborators: []models.Collaborator{
			{
				ID:        primitive.NewObjectID(),
				UserId:    ownerCollabID,
				Role:      "owner",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	// Collaborator with owner role should be able to do owner-level things
	assert.True(t, CanEditBoard(board, ownerCollabID))
	assert.True(t, CanDeleteBoard(board, ownerCollabID))
	assert.True(t, CanManageCollaborators(board, ownerCollabID))
	assert.True(t, CanManageKeys(board, ownerCollabID))
}

// Benchmark tests
func BenchmarkGetUserBoardRole(b *testing.B) {
	authorID := primitive.NewObjectID()
	board := models.Board{
		ID:       primitive.NewObjectID(),
		AuthorId: authorID,
		Collaborators: []models.Collaborator{
			{ID: primitive.NewObjectID(), UserId: primitive.NewObjectID(), Role: "editor"},
			{ID: primitive.NewObjectID(), UserId: primitive.NewObjectID(), Role: "viewer"},
			{ID: primitive.NewObjectID(), UserId: primitive.NewObjectID(), Role: "editor"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetUserBoardRole(board, authorID)
	}
}

func BenchmarkCanEditBoard(b *testing.B) {
	authorID := primitive.NewObjectID()
	board := models.Board{
		ID:       primitive.NewObjectID(),
		AuthorId: authorID,
		Collaborators: []models.Collaborator{
			{ID: primitive.NewObjectID(), UserId: primitive.NewObjectID(), Role: "editor"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CanEditBoard(board, authorID)
	}
}
