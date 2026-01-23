//go:build integration
// +build integration

package repositories

import (
	"testing"
	"time"
	api "tuneslap/generated"
	"tuneslap/models"
	"tuneslap/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFullBoardWorkflow(t *testing.T) {
	cleanup := testutils.SetupTestMongoDBWithSkip(t)
	defer cleanup()

	redisCleanup := testutils.SetupTestRedisWithSkip(t)
	defer redisCleanup()

	// Create repositories
	userRepo := NewUserRepository()
	boardRepo := NewBoardRepository()
	keyRepo := NewKeyRepository()
	collabRepo := NewCollaboratorRepository()
	mediaRepo := NewMediaRepository()

	// Create a test user
	user := models.User{
		Name:         "Integration Test User",
		Email:        "integration@test.com",
		PasswordHash: "hashedpassword123",
	}
	createdUser, err := userRepo.CreateUser(user)
	require.NoError(t, err)
	authorID := createdUser.ID

	// Create a second user for collaboration tests
	user2 := models.User{
		Name:         "Collaborator User",
		Email:        "collab@test.com",
		PasswordHash: "hashedpassword456",
	}
	createdUser2, err := userRepo.CreateUser(user2)
	require.NoError(t, err)

	t.Run("create board", func(t *testing.T) {
		name := "Integration Test Board"
		desc := "Testing the full workflow"
		layout := "grid"

		req := &api.CreateBoardRequest{
			Name:        &name,
			Description: &desc,
			Layout:      &layout,
		}

		board, err := boardRepo.CreateBoard(req, authorID)
		require.NoError(t, err)
		assert.Equal(t, "Integration Test Board", board.Name)
		assert.Equal(t, authorID, board.AuthorId)
	})

	t.Run("full workflow with keys and collaborators", func(t *testing.T) {
		// Create board
		name := "Full Workflow Board"
		layout := "list"
		req := &api.CreateBoardRequest{
			Name:   &name,
			Layout: &layout,
		}

		board, err := boardRepo.CreateBoard(req, authorID)
		require.NoError(t, err)
		boardID := board.ID

		// Create media for the key
		mediaType := "audio"
		fileName := "workflow_test.mp3"
		fileUrl := "https://example.com/workflow_test.mp3"
		contentType := "audio/mpeg"
		fileSize := int32(1024)

		mediaReq := &api.CreateMediaRequest{
			MediaType:   &mediaType,
			FileName:    &fileName,
			FileUrl:     &fileUrl,
			ContentType: &contentType,
			FileSize:    &fileSize,
		}

		media, err := mediaRepo.CreateMedia(mediaReq, authorID)
		require.NoError(t, err)

		// Add a key to the board
		keyReq := &api.CreateKeyRequest{
			BoardId:      boardID.Hex(),
			Name:         "Workflow Key",
			AudioMediaId: media.ID.Hex(),
			HotKey:       "W",
		}

		board, err = keyRepo.CreateKey(boardID, keyReq)
		require.NoError(t, err)
		assert.Len(t, board.Keys, 1)
		assert.Equal(t, "Workflow Key", board.Keys[0].Name)

		// Add a collaborator
		email := createdUser2.Email
		role := "editor"
		collabReq := &api.CreateCollaboratorRequest{
			Email: &email,
			Role:  &role,
		}

		board, err = collabRepo.CreateCollaborator(boardID, collabReq, authorID)
		require.NoError(t, err)
		assert.Len(t, board.Collaborators, 1)
		assert.Equal(t, createdUser2.Email, board.Collaborators[0].Email)

		// Verify collaborator can access the board
		accessibleBoard, err := boardRepo.FindByIDWithAccess(boardID, createdUser2.ID)
		require.NoError(t, err)
		assert.Equal(t, boardID, accessibleBoard.ID)

		// Update the key
		newName := "Updated Key Name"
		keyUpdateReq := &api.UpdateKeyRequest{
			Name: &newName,
		}
		board, err = keyRepo.UpdateKey(boardID, board.Keys[0].ID, keyUpdateReq)
		require.NoError(t, err)
		assert.Equal(t, "Updated Key Name", board.Keys[0].Name)

		// Delete the key
		board, err = keyRepo.DeleteKey(boardID, board.Keys[0].ID)
		require.NoError(t, err)
		assert.Len(t, board.Keys, 0)

		// Remove collaborator
		board, err = collabRepo.DeleteCollaborator(boardID, board.Collaborators[0].ID)
		require.NoError(t, err)
		assert.Len(t, board.Collaborators, 0)

		// Delete the board
		err = boardRepo.DeleteBoard(boardID, authorID)
		require.NoError(t, err)

		// Verify board is deleted
		_, err = boardRepo.GetById(boardID, authorID)
		assert.Error(t, err)
	})
}

func TestMediaLifecycle(t *testing.T) {
	cleanup := testutils.SetupTestMongoDBWithSkip(t)
	defer cleanup()

	redisCleanup := testutils.SetupTestRedisWithSkip(t)
	defer redisCleanup()

	mediaRepo := NewMediaRepository()
	authorID := primitive.NewObjectID()

	t.Run("create and process media", func(t *testing.T) {
		// Create media
		mediaType := "audio"
		fileName := "lifecycle_test.mp3"
		fileUrl := "https://example.com/lifecycle_test.mp3"
		contentType := "audio/mpeg"
		fileSize := int32(2048)

		req := &api.CreateMediaRequest{
			MediaType:   &mediaType,
			FileName:    &fileName,
			FileUrl:     &fileUrl,
			ContentType: &contentType,
			FileSize:    &fileSize,
		}

		media, err := mediaRepo.CreateMedia(req, authorID)
		require.NoError(t, err)
		assert.Equal(t, models.ProcessingStatusPending, media.Status)

		// Simulate processing update
		updateData := &models.Media{
			Status:       models.ProcessingStatusProcessing,
			ProcessingActivity: []models.ProcessingActivity{
				{
					Status:    models.ProcessingStatusProcessing,
					Message:   "Processing audio...",
					CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
					UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
				},
			},
		}
		media, err = mediaRepo.UpdateMediaUnscoped(media.ID, updateData)
		require.NoError(t, err)
		assert.Equal(t, models.ProcessingStatusProcessing, media.Status)

		// Complete processing
		updateData = &models.Media{
			Status:       models.ProcessingStatusDone,
			ProcessedUrl: "https://example.com/processed/lifecycle_test.webm",
			Duration:     60.5,
		}
		media, err = mediaRepo.UpdateMediaUnscoped(media.ID, updateData)
		require.NoError(t, err)
		assert.Equal(t, models.ProcessingStatusDone, media.Status)
		assert.Equal(t, float64(60.5), media.Duration)

		// Get stats
		stats, err := mediaRepo.GetMyMediaStats(authorID)
		require.NoError(t, err)
		assert.Equal(t, 1, stats.AudioCount)
		assert.Equal(t, 0, stats.ImageCount)

		// Delete media
		err = mediaRepo.DeleteMedia(media.ID, authorID)
		require.NoError(t, err)
	})
}

func TestConcurrentOperations(t *testing.T) {
	cleanup := testutils.SetupTestMongoDBWithSkip(t)
	defer cleanup()

	redisCleanup := testutils.SetupTestRedisWithSkip(t)
	defer redisCleanup()

	boardRepo := NewBoardRepository()
	authorID := primitive.NewObjectID()

	// Create multiple boards concurrently
	t.Run("concurrent board creation", func(t *testing.T) {
		done := make(chan bool)
		errors := make(chan error, 5)

		for i := 0; i < 5; i++ {
			go func(idx int) {
				name := "Concurrent Board " + string(rune('A'+idx))
				layout := "grid"
				req := &api.CreateBoardRequest{
					Name:   &name,
					Layout: &layout,
				}
				_, err := boardRepo.CreateBoard(req, authorID)
				if err != nil {
					errors <- err
				}
				done <- true
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < 5; i++ {
			<-done
		}
		close(errors)

		// Check for errors
		for err := range errors {
			t.Errorf("Concurrent creation failed: %v", err)
		}

		// Verify all boards were created
		boards, err := boardRepo.GetByAuthor(authorID)
		require.NoError(t, err)
		assert.Len(t, boards, 5)
	})
}
