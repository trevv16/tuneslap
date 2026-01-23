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
		req := &api.CreateBoardRequest{
			Name:   "Collab Test Board",
			Layout: "grid",
		}
		board, err := boardRepo.CreateBoard(req, authorID)
		assert.NoError(t, err)
		return board.ID
	}

	t.Run("CreateCollaborator", func(t *testing.T) {
		boardID := createBoard()

		req := &api.CreateCollaboratorRequest{
			Email: "newcollab@test.com",
			Role:  "editor",
		}

		board, err := collabRepo.CreateCollaborator(boardID, req, authorID)
		assert.NoError(t, err)
		assert.Len(t, board.Collaborators, 1)
		assert.Equal(t, "newcollab@test.com", board.Collaborators[0].Email)
		assert.Equal(t, "editor", board.Collaborators[0].Role)
	})

	t.Run("CreateCollaborator_ViewerRole", func(t *testing.T) {
		boardID := createBoard()

		req := &api.CreateCollaboratorRequest{
			Email: "viewer@test.com",
			Role:  "viewer",
		}

		board, err := collabRepo.CreateCollaborator(boardID, req, authorID)
		assert.NoError(t, err)
		assert.Len(t, board.Collaborators, 1)
		assert.Equal(t, "viewer", board.Collaborators[0].Role)
	})

	t.Run("UpdateCollaborator", func(t *testing.T) {
		boardID := createBoard()

		// Create collaborator first
		createReq := &api.CreateCollaboratorRequest{
			Email: "updatecollab@test.com",
			Role:  "viewer",
		}

		board, err := collabRepo.CreateCollaborator(boardID, createReq, authorID)
		assert.NoError(t, err)
		collabID := board.Collaborators[0].ID

		// Update role
		updateReq := &api.UpdateCollaboratorRequest{
			Role: "editor",
		}

		updatedBoard, err := collabRepo.UpdateCollaborator(boardID, collabID, updateReq)
		assert.NoError(t, err)
		assert.Equal(t, "editor", updatedBoard.Collaborators[0].Role)
	})

	t.Run("DeleteCollaborator", func(t *testing.T) {
		boardID := createBoard()
		// Create collaborator first
		createReq := &api.CreateCollaboratorRequest{
			Email: "deletecollab@test.com",
			Role:  "editor",
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

		req := &api.CreateCollaboratorRequest{
			Email: "invalid-email",
			Role:  "editor",
		}

		_, err := collabRepo.CreateCollaborator(boardID, req, authorID)
		assert.Error(t, err)
	})

	t.Run("CreateCollaborator_InvalidRole", func(t *testing.T) {
		boardID := createBoard()

		req := &api.CreateCollaboratorRequest{
			Email: "valid@test.com",
			Role:  "invalid_role",
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
	boardReq := &api.CreateBoardRequest{
		Name:   "Benchmark Board",
		Layout: "grid",
	}
	board, _ := boardRepo.CreateBoard(boardReq, authorID)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := &api.CreateCollaboratorRequest{
			Email: "bench" + string(rune('0'+i%10)) + "@test.com",
			Role:  "editor",
		}
		collabRepo.CreateCollaborator(board.ID, req, authorID)
	}
}
