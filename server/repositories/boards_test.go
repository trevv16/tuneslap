package repositories

import (
	"testing"
	"time"
	api "tuneslap/generated"
	"tuneslap/models"
	"tuneslap/testutils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewBoardRepository(t *testing.T) {
	repo := NewBoardRepository()
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.Repository)
	assert.NotNil(t, repo.validator)
}

func TestBoardRepository_GetValidator(t *testing.T) {
	repo := NewBoardRepository()
	validator := repo.GetValidator()
	assert.NotNil(t, validator)
}

// Integration tests requiring MongoDB
func TestBoardRepositoryIntegration(t *testing.T) {
	cleanup := testutils.SetupTestMongoDBWithSkip(t)
	defer cleanup()

	redisCleanup := testutils.SetupTestRedisWithSkip(t)
	defer redisCleanup()

	repo := NewBoardRepository()

	// Create test user
	authorID := primitive.NewObjectID()

	t.Run("CreateBoard", func(t *testing.T) {
		desc := "Test Description"

		req := &api.CreateBoardRequest{
			Name:        "Test Board",
			Description: &desc,
			Layout:      "grid",
		}

		created, err := repo.CreateBoard(req, authorID)
		assert.NoError(t, err)
		assert.NotEqual(t, primitive.NilObjectID, created.ID)
		assert.Equal(t, "Test Board", created.Name)
		assert.Equal(t, "Test Description", created.Description)
		assert.Equal(t, models.GridLayout, created.Layout)
		assert.Equal(t, authorID, created.AuthorId)
	})

	t.Run("GetByAuthor", func(t *testing.T) {
		// Create some boards
		for i := 0; i < 3; i++ {
			req := &api.CreateBoardRequest{
				Name:   "Author Board " + string(rune('1'+i)),
				Layout: "grid",
			}
			_, err := repo.CreateBoard(req, authorID)
			assert.NoError(t, err)
		}

		// Get boards by author
		boards, err := repo.GetByAuthor(authorID)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(boards), 3)
	})

	t.Run("GetById", func(t *testing.T) {
		// Create a board
		req := &api.CreateBoardRequest{
			Name:   "GetById Board",
			Layout: "list",
		}

		created, err := repo.CreateBoard(req, authorID)
		assert.NoError(t, err)

		// Get by ID
		found, err := repo.GetById(created.ID, authorID)
		assert.NoError(t, err)
		assert.Equal(t, "GetById Board", found.Name)
		assert.Equal(t, models.ListLayout, found.Layout)
	})

	t.Run("GetById_WrongAuthor", func(t *testing.T) {
		// Create a board
		req := &api.CreateBoardRequest{
			Name:   "Wrong Author Board",
			Layout: "grid",
		}

		created, err := repo.CreateBoard(req, authorID)
		assert.NoError(t, err)

		// Try to get with different author
		wrongAuthorID := primitive.NewObjectID()
		_, err = repo.GetById(created.ID, wrongAuthorID)
		assert.Error(t, err)
	})

	t.Run("UpdateBoard", func(t *testing.T) {
		// Create a board
		req := &api.CreateBoardRequest{
			Name:   "Update Board",
			Layout: "grid",
		}

		created, err := repo.CreateBoard(req, authorID)
		assert.NoError(t, err)

		// Update the board
		newName := "Updated Board Name"
		newDesc := "New Description"
		updateReq := &api.UpdateBoardRequest{
			Name:        &newName,
			Description: &newDesc,
		}

		updated, err := repo.UpdateBoard(created.ID, authorID, updateReq)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Board Name", updated.Name)
		assert.Equal(t, "New Description", updated.Description)
	})

	t.Run("DeleteBoard", func(t *testing.T) {
		// Create a board
		req := &api.CreateBoardRequest{
			Name:   "Delete Board",
			Layout: "grid",
		}

		created, err := repo.CreateBoard(req, authorID)
		assert.NoError(t, err)

		// Delete the board
		err = repo.DeleteBoard(created.ID, authorID)
		assert.NoError(t, err)

		// Verify deletion
		_, err = repo.GetById(created.ID, authorID)
		assert.Error(t, err)
	})

	t.Run("FindByIDWithAccess_AsAuthor", func(t *testing.T) {
		// Create a board
		req := &api.CreateBoardRequest{
			Name:   "Access Board",
			Layout: "grid",
		}

		created, err := repo.CreateBoard(req, authorID)
		assert.NoError(t, err)

		// Find with access as author
		found, err := repo.FindByIDWithAccess(created.ID, authorID)
		assert.NoError(t, err)
		assert.Equal(t, "Access Board", found.Name)
	})

	t.Run("FindByIDWithAccess_AsCollaborator", func(t *testing.T) {
		collaboratorID := primitive.NewObjectID()

		// Create a board with collaborator
		req := &api.CreateBoardRequest{
			Name:   "Collab Access Board",
			Layout: "grid",
		}

		created, err := repo.CreateBoard(req, authorID)
		assert.NoError(t, err)

		// Add collaborator manually
		err = repo.AddToNestedArray(created.ID, "collaborators", models.Collaborator{
			ID:        primitive.NewObjectID(),
			UserId:    collaboratorID,
			Email:     "collab@test.com",
			Role:      "editor",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		assert.NoError(t, err)

		// Find with access as collaborator
		found, err := repo.FindByIDWithAccess(created.ID, collaboratorID)
		assert.NoError(t, err)
		assert.Equal(t, "Collab Access Board", found.Name)
	})

	t.Run("FindAllWithAccess", func(t *testing.T) {
		testAuthorID := primitive.NewObjectID()

		// Create boards as author
		for i := 0; i < 2; i++ {
			req := &api.CreateBoardRequest{
				Name:   "Access Test Board " + string(rune('1'+i)),
				Layout: "grid",
			}
			_, err := repo.CreateBoard(req, testAuthorID)
			assert.NoError(t, err)
		}

		// Find all with access
		boards, err := repo.FindAllWithAccess(testAuthorID)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(boards), 2)
	})
}

// Benchmark tests
func BenchmarkBoardRepository_CreateBoard(b *testing.B) {
	cleanup, err := testutils.SetupTestMongoDB()
	if err != nil {
		b.Skipf("Skipping benchmark - MongoDB not available: %v", err)
	}
	defer cleanup()

	repo := NewBoardRepository()
	authorID := primitive.NewObjectID()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := &api.CreateBoardRequest{
			Name:   "Benchmark Board",
			Layout: "grid",
		}
		repo.CreateBoard(req, authorID)
	}
}

func BenchmarkBoardRepository_GetByAuthor(b *testing.B) {
	cleanup, err := testutils.SetupTestMongoDB()
	if err != nil {
		b.Skipf("Skipping benchmark - MongoDB not available: %v", err)
	}
	defer cleanup()

	repo := NewBoardRepository()
	authorID := primitive.NewObjectID()

	// Create some boards
	for i := 0; i < 10; i++ {
		req := &api.CreateBoardRequest{
			Name:   "Benchmark Board " + string(rune('0'+i)),
			Layout: "grid",
		}
		repo.CreateBoard(req, authorID)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.GetByAuthor(authorID)
	}
}
