package repositories

import (
	"testing"

	api "tuneslap/generated"
	"tuneslap/testutils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewKeyRepository(t *testing.T) {
	repo := NewKeyRepository()
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.ArrayRepository)
	assert.NotNil(t, repo.boardRepo)
	assert.NotNil(t, repo.validator)
}

func TestKeyRepository_GetValidator(t *testing.T) {
	repo := NewKeyRepository()
	validator := repo.GetValidator()
	assert.NotNil(t, validator)
}

// Integration tests requiring MongoDB
func TestKeyRepositoryIntegration(t *testing.T) {
	cleanup := testutils.SetupTestMongoDBWithSkip(t)
	defer cleanup()

	redisCleanup := testutils.SetupTestRedisWithSkip(t)
	defer redisCleanup()

	keyRepo := NewKeyRepository()
	boardRepo := NewBoardRepository()
	authorID := primitive.NewObjectID()

	// Create a board for testing
	createBoard := func() primitive.ObjectID {
		req := &api.CreateBoardRequest{
			Name:   "Test Board",
			Layout: "grid",
		}
		board, err := boardRepo.CreateBoard(req, authorID)
		assert.NoError(t, err)
		return board.ID
	}

	t.Run("CreateKey", func(t *testing.T) {
		boardID := createBoard()
		audioMediaID := primitive.NewObjectID()

		req := &api.CreateKeyRequest{
			BoardId:      boardID.Hex(),
			Name:         "Test Key",
			AudioMediaId: audioMediaID.Hex(),
			HotKey:       "A",
		}

		board, err := keyRepo.CreateKey(boardID, req)
		assert.NoError(t, err)
		assert.Len(t, board.Keys, 1)
		assert.Equal(t, "Test Key", board.Keys[0].Name)
		assert.Equal(t, "A", board.Keys[0].HotKey)
	})

	t.Run("CreateKey_WithAllFields", func(t *testing.T) {
		boardID := createBoard()
		audioMediaID := primitive.NewObjectID()
		imageMediaID := primitive.NewObjectID()
		desc := "Test description"

		req := &api.CreateKeyRequest{
			BoardId:      boardID.Hex(),
			Name:         "Full Key",
			Description:  &desc,
			AudioMediaId: audioMediaID.Hex(),
			ImageMediaId: strPtr(imageMediaID.Hex()),
			HotKey:       "B",
		}

		board, err := keyRepo.CreateKey(boardID, req)
		assert.NoError(t, err)
		assert.Len(t, board.Keys, 1)
		assert.Equal(t, "Full Key", board.Keys[0].Name)
		assert.Equal(t, "Test description", board.Keys[0].Description)
		assert.Equal(t, audioMediaID, board.Keys[0].AudioMediaId)
		assert.Equal(t, imageMediaID, board.Keys[0].ImageMediaId)
	})

	t.Run("UpdateKey", func(t *testing.T) {
		boardID := createBoard()
		audioMediaID := primitive.NewObjectID()

		// Create key first
		createReq := &api.CreateKeyRequest{
			BoardId:      boardID.Hex(),
			Name:         "Original Key",
			AudioMediaId: audioMediaID.Hex(),
			HotKey:       "C",
		}

		board, err := keyRepo.CreateKey(boardID, createReq)
		assert.NoError(t, err)
		keyID := board.Keys[0].ID

		// Update key
		newName := "Updated Key"
		newHotKey := "D"
		updateReq := &api.UpdateKeyRequest{
			Name:   &newName,
			HotKey: &newHotKey,
		}

		updatedBoard, err := keyRepo.UpdateKey(boardID, keyID, updateReq)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Key", updatedBoard.Keys[0].Name)
		assert.Equal(t, "D", updatedBoard.Keys[0].HotKey)
	})

	t.Run("DeleteKey", func(t *testing.T) {
		boardID := createBoard()
		audioMediaID := primitive.NewObjectID()

		// Create key first
		createReq := &api.CreateKeyRequest{
			BoardId:      boardID.Hex(),
			Name:         "Delete Key",
			AudioMediaId: audioMediaID.Hex(),
			HotKey:       "E",
		}

		board, err := keyRepo.CreateKey(boardID, createReq)
		assert.NoError(t, err)
		keyID := board.Keys[0].ID

		// Delete key
		updatedBoard, err := keyRepo.DeleteKey(boardID, keyID)
		assert.NoError(t, err)
		assert.Len(t, updatedBoard.Keys, 0)
	})

	t.Run("DeleteKeysByMediaId", func(t *testing.T) {
		boardID := createBoard()
		mediaID := primitive.NewObjectID()

		// Create multiple keys with same media ID
		for i := 0; i < 3; i++ {
			createReq := &api.CreateKeyRequest{
				BoardId:      boardID.Hex(),
				Name:         "Media Key " + string(rune('1'+i)),
				AudioMediaId: mediaID.Hex(),
				HotKey:       string(rune('F' + i)),
			}
			_, err := keyRepo.CreateKey(boardID, createReq)
			assert.NoError(t, err)
		}

		// Delete keys by media ID
		count, err := keyRepo.DeleteKeysByMediaId(mediaID)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(1))

		// Verify keys were deleted
		board, err := boardRepo.FindByID(boardID)
		assert.NoError(t, err)
		for _, key := range board.Keys {
			assert.NotEqual(t, mediaID, key.AudioMediaId)
		}
	})
}

// Helper for tests
func strPtr(s string) *string {
	return &s
}

// Benchmark tests
func BenchmarkKeyRepository_CreateKey(b *testing.B) {
	cleanup, err := testutils.SetupTestMongoDB()
	if err != nil {
		b.Skipf("Skipping benchmark - MongoDB not available: %v", err)
	}
	defer cleanup()

	keyRepo := NewKeyRepository()
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
		audioMediaID := primitive.NewObjectID()
		req := &api.CreateKeyRequest{
			BoardId:      board.ID.Hex(),
			Name:         "Benchmark Key",
			AudioMediaId: audioMediaID.Hex(),
			HotKey:       string(rune('A' + i%26)),
		}
		keyRepo.CreateKey(board.ID, req)
	}
}
