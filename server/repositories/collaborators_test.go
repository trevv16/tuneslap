package repositories

import (
	"testing"
	api "tuneslap/generated"
	"tuneslap/testutils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewCollaboratorRepository(t *testing.T) {
	repo := NewCollaboratorRepository()
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.ArrayRepository)
	assert.NotNil(t, repo.boardRepo)
	assert.NotNil(t, repo.userRepo)
	assert.NotNil(t, repo.validator)
}

func TestCollaboratorRepository_GetValidator(t *testing.T) {
	repo := NewCollaboratorRepository()
	validator := repo.GetValidator()
	assert.NotNil(t, validator)
}

func TestCollaboratorRepository_GetBoardRepo(t *testing.T) {
	repo := NewCollaboratorRepository()
	boardRepo := repo.GetBoardRepo()
	assert.NotNil(t, boardRepo)
}

// Integration tests requiring MongoDB
func TestCollaboratorRepositoryIntegration(t *testing.T) {
	cleanup := testutils.SetupTestMongoDBWithSkip(t)
	defer cleanup()

	redisCleanup := testutils.SetupTestRedisWithSkip(t)
	defer redisCleanup()

	collabRepo := NewCollaboratorRepository()
	boardRepo := NewBoardRepository()
	authorID := primitive.NewObjectID()

	// Create a board for testing
	createBoard := func() primitive.ObjectID {
		name := "Collab Test Board"
		layout := "grid"
		req := &api.CreateBoardRequest{
			Name:   &name,
			Layout: &layout,
		}
		board, err := boardRepo.CreateBoard(req, authorID)
		assert.NoError(t, err)
		return board.ID
	}

	t.Run("CreateCollaborator", func(t *testing.T) {
		boardID := createBoard()
		email := "newcollab@test.com"
		role := "editor"

		req := &api.CreateCollaboratorRequest{
			Email: &email,
			Role:  &role,
		}

		board, err := collabRepo.CreateCollaborator(boardID, req, authorID)
		assert.NoError(t, err)
		assert.Len(t, board.Collaborators, 1)
		assert.Equal(t, "newcollab@test.com", board.Collaborators[0].Email)
		assert.Equal(t, "editor", board.Collaborators[0].Role)
	})

	t.Run("CreateCollaborator_ViewerRole", func(t *testing.T) {
		boardID := createBoard()
		email := "viewer@test.com"
		role := "viewer"

		req := &api.CreateCollaboratorRequest{
			Email: &email,
			Role:  &role,
		}

		board, err := collabRepo.CreateCollaborator(boardID, req, authorID)
		assert.NoError(t, err)
		assert.Len(t, board.Collaborators, 1)
		assert.Equal(t, "viewer", board.Collaborators[0].Role)
	})

	t.Run("UpdateCollaborator", func(t *testing.T) {
		boardID := createBoard()
		email := "updatecollab@test.com"
		role := "viewer"

		// Create collaborator first
		createReq := &api.CreateCollaboratorRequest{
			Email: &email,
			Role:  &role,
		}

		board, err := collabRepo.CreateCollaborator(boardID, createReq, authorID)
		assert.NoError(t, err)
		collabID := board.Collaborators[0].ID

		// Update role
		newRole := "editor"
		updateReq := &api.UpdateCollaboratorRequest{
			Role: &newRole,
		}

		updatedBoard, err := collabRepo.UpdateCollaborator(boardID, collabID, updateReq)
		assert.NoError(t, err)
		assert.Equal(t, "editor", updatedBoard.Collaborators[0].Role)
	})

	t.Run("DeleteCollaborator", func(t *testing.T) {
		boardID := createBoard()
		email := "deletecollab@test.com"
		role := "editor"

		// Create collaborator first
		createReq := &api.CreateCollaboratorRequest{
			Email: &email,
			Role:  &role,
		}

		board, err := collabRepo.CreateCollaborator(boardID, createReq, authorID)
		assert.NoError(t, err)
		collabID := board.Collaborators[0].ID

		// Delete collaborator
		updatedBoard, err := collabRepo.DeleteCollaborator(boardID, collabID)
		assert.NoError(t, err)
		assert.Len(t, updatedBoard.Collaborators, 0)
	})

	t.Run("CreateCollaborator_InvalidValidation", func(t *testing.T) {
		boardID := createBoard()
		email := "invalid-email"
		role := "editor"

		req := &api.CreateCollaboratorRequest{
			Email: &email,
			Role:  &role,
		}

		_, err := collabRepo.CreateCollaborator(boardID, req, authorID)
		assert.Error(t, err)
	})

	t.Run("CreateCollaborator_InvalidRole", func(t *testing.T) {
		boardID := createBoard()
		email := "valid@test.com"
		role := "invalid_role"

		req := &api.CreateCollaboratorRequest{
			Email: &email,
			Role:  &role,
		}

		_, err := collabRepo.CreateCollaborator(boardID, req, authorID)
		assert.Error(t, err)
	})
}

// Benchmark tests
func BenchmarkCollaboratorRepository_CreateCollaborator(b *testing.B) {
	cleanup, err := testutils.SetupTestMongoDB()
	if err != nil {
		b.Skipf("Skipping benchmark - MongoDB not available: %v", err)
	}
	defer cleanup()

	collabRepo := NewCollaboratorRepository()
	boardRepo := NewBoardRepository()
	authorID := primitive.NewObjectID()

	// Create a board
	name := "Benchmark Board"
	layout := "grid"
	boardReq := &api.CreateBoardRequest{
		Name:   &name,
		Layout: &layout,
	}
	board, _ := boardRepo.CreateBoard(boardReq, authorID)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		email := "bench" + string(rune('0'+i%10)) + "@test.com"
		role := "editor"
		req := &api.CreateCollaboratorRequest{
			Email: &email,
			Role:  &role,
		}
		collabRepo.CreateCollaborator(board.ID, req, authorID)
	}
}
