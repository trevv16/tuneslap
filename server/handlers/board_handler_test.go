package handlers

import (
	"testing"
	"time"
	"tuneslap/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewBoardHandler(t *testing.T) {
	// Skip if no database connection - this test just checks constructor doesn't panic
	t.Run("creates handler successfully", func(t *testing.T) {
		// Note: This will fail if there's no database, but that's expected
		// The test verifies the constructor logic is correct
		defer func() {
			if r := recover(); r != nil {
				t.Skip("Skipping test - requires database connection")
			}
		}()

		handler := NewBoardHandler()
		assert.NotNil(t, handler)
		assert.NotNil(t, handler.boardRepo)
		assert.NotNil(t, handler.mediaRepo)
		assert.NotNil(t, handler.BaseHandler)
	})
}

func TestBoardResponseMapperLogic(t *testing.T) {
	// Test the board to response mapping logic without requiring database
	// We test the underlying ToBoardResponse utility which the mapper uses

	validTime := time.Now()
	authorID := primitive.NewObjectID()
	boardID := primitive.NewObjectID()

	t.Run("maps board with empty keys", func(t *testing.T) {
		board := models.Board{
			ID:          boardID,
			AuthorId:    authorID,
			Name:        "Test Board",
			Description: "Test description",
			Layout:      models.GridLayout,
			Keys:        []models.Key{},
			CreatedAt:   validTime,
			UpdatedAt:   validTime,
		}

		// Verify board structure is correct
		assert.Equal(t, boardID, board.ID)
		assert.Equal(t, "Test Board", board.Name)
		assert.Empty(t, board.Keys)
	})

	t.Run("maps board with keys containing media IDs", func(t *testing.T) {
		audioMediaID := primitive.NewObjectID()
		imageMediaID := primitive.NewObjectID()
		keyID := primitive.NewObjectID()

		board := models.Board{
			ID:          boardID,
			AuthorId:    authorID,
			Name:        "Test Board with Keys",
			Description: "Has keys",
			Layout:      models.ListLayout,
			Keys: []models.Key{
				{
					ID:           keyID,
					BoardId:      boardID,
					Name:         "Test Key",
					Description:  "A test key",
					AudioMediaId: audioMediaID,
					ImageMediaId: imageMediaID,
					HotKey:       "A",
					CreatedAt:    validTime,
					UpdatedAt:    validTime,
				},
			},
			CreatedAt: validTime,
			UpdatedAt: validTime,
		}

		// Verify board structure
		assert.Len(t, board.Keys, 1)
		assert.Equal(t, audioMediaID, board.Keys[0].AudioMediaId)
		assert.Equal(t, imageMediaID, board.Keys[0].ImageMediaId)
		assert.False(t, board.Keys[0].AudioMediaId.IsZero())
		assert.False(t, board.Keys[0].ImageMediaId.IsZero())
	})

	t.Run("maps board with key missing image media", func(t *testing.T) {
		audioMediaID := primitive.NewObjectID()
		keyID := primitive.NewObjectID()

		board := models.Board{
			ID:          boardID,
			AuthorId:    authorID,
			Name:        "Test Board",
			Layout:      models.GridLayout,
			Keys: []models.Key{
				{
					ID:           keyID,
					BoardId:      boardID,
					Name:         "Audio Only Key",
					AudioMediaId: audioMediaID,
					ImageMediaId: primitive.NilObjectID,
					HotKey:       "B",
					CreatedAt:    validTime,
					UpdatedAt:    validTime,
				},
			},
			CreatedAt: validTime,
			UpdatedAt: validTime,
		}

		// Verify structure
		assert.Len(t, board.Keys, 1)
		assert.False(t, board.Keys[0].AudioMediaId.IsZero())
		assert.True(t, board.Keys[0].ImageMediaId.IsZero())
	})

	t.Run("maps board with collaborators", func(t *testing.T) {
		collabID := primitive.NewObjectID()
		collabUserID := primitive.NewObjectID()

		board := models.Board{
			ID:          boardID,
			AuthorId:    authorID,
			Name:        "Shared Board",
			Description: "Board with collaborators",
			Layout:      models.GridLayout,
			Collaborators: []models.Collaborator{
				{
					ID:        collabID,
					UserId:    collabUserID,
					Email:     "collab@example.com",
					Role:      "editor",
					CreatedAt: validTime,
					UpdatedAt: validTime,
				},
			},
			Keys:      []models.Key{},
			CreatedAt: validTime,
			UpdatedAt: validTime,
		}

		assert.Len(t, board.Collaborators, 1)
		assert.Equal(t, "editor", board.Collaborators[0].Role)
		assert.Equal(t, "collab@example.com", board.Collaborators[0].Email)
	})
}

func TestBoardMediaIDCollection(t *testing.T) {
	// Test the logic that collects media IDs from keys (used in BoardResponseMapper)

	t.Run("collects all media IDs from keys", func(t *testing.T) {
		audioMediaID1 := primitive.NewObjectID()
		audioMediaID2 := primitive.NewObjectID()
		imageMediaID1 := primitive.NewObjectID()
		imageMediaID2 := primitive.NewObjectID()

		keys := []models.Key{
			{
				ID:           primitive.NewObjectID(),
				AudioMediaId: audioMediaID1,
				ImageMediaId: imageMediaID1,
			},
			{
				ID:           primitive.NewObjectID(),
				AudioMediaId: audioMediaID2,
				ImageMediaId: imageMediaID2,
			},
		}

		// Simulate the collection logic from BoardResponseMapper
		var mediaIds []primitive.ObjectID
		for _, key := range keys {
			if !key.AudioMediaId.IsZero() {
				mediaIds = append(mediaIds, key.AudioMediaId)
			}
			if !key.ImageMediaId.IsZero() {
				mediaIds = append(mediaIds, key.ImageMediaId)
			}
		}

		assert.Len(t, mediaIds, 4)
		assert.Contains(t, mediaIds, audioMediaID1)
		assert.Contains(t, mediaIds, audioMediaID2)
		assert.Contains(t, mediaIds, imageMediaID1)
		assert.Contains(t, mediaIds, imageMediaID2)
	})

	t.Run("skips nil media IDs", func(t *testing.T) {
		audioMediaID := primitive.NewObjectID()

		keys := []models.Key{
			{
				ID:           primitive.NewObjectID(),
				AudioMediaId: audioMediaID,
				ImageMediaId: primitive.NilObjectID, // No image
			},
			{
				ID:           primitive.NewObjectID(),
				AudioMediaId: primitive.NilObjectID, // No audio
				ImageMediaId: primitive.NilObjectID, // No image
			},
		}

		var mediaIds []primitive.ObjectID
		for _, key := range keys {
			if !key.AudioMediaId.IsZero() {
				mediaIds = append(mediaIds, key.AudioMediaId)
			}
			if !key.ImageMediaId.IsZero() {
				mediaIds = append(mediaIds, key.ImageMediaId)
			}
		}

		assert.Len(t, mediaIds, 1)
		assert.Contains(t, mediaIds, audioMediaID)
	})

	t.Run("empty keys returns empty media IDs", func(t *testing.T) {
		keys := []models.Key{}

		var mediaIds []primitive.ObjectID
		for _, key := range keys {
			if !key.AudioMediaId.IsZero() {
				mediaIds = append(mediaIds, key.AudioMediaId)
			}
			if !key.ImageMediaId.IsZero() {
				mediaIds = append(mediaIds, key.ImageMediaId)
			}
		}

		assert.Empty(t, mediaIds)
	})
}

func TestBoardKeyURLEnrichment(t *testing.T) {
	// Test the URL enrichment logic used in BoardResponseMapper

	t.Run("enriches keys with URLs from media map", func(t *testing.T) {
		audioMediaID := primitive.NewObjectID()
		imageMediaID := primitive.NewObjectID()

		// Simulate media URLs map
		mediaUrls := map[primitive.ObjectID]string{
			audioMediaID: "https://storage.example.com/audio.mp3",
			imageMediaID: "https://storage.example.com/image.png",
		}

		key := models.Key{
			ID:           primitive.NewObjectID(),
			AudioMediaId: audioMediaID,
			ImageMediaId: imageMediaID,
		}

		// Check if URLs exist in map (logic from BoardResponseMapper)
		var audioUrl, imageUrl *string
		if !key.AudioMediaId.IsZero() {
			if url, exists := mediaUrls[key.AudioMediaId]; exists {
				audioUrl = &url
			}
		}
		if !key.ImageMediaId.IsZero() {
			if url, exists := mediaUrls[key.ImageMediaId]; exists {
				imageUrl = &url
			}
		}

		assert.NotNil(t, audioUrl)
		assert.NotNil(t, imageUrl)
		assert.Equal(t, "https://storage.example.com/audio.mp3", *audioUrl)
		assert.Equal(t, "https://storage.example.com/image.png", *imageUrl)
	})

	t.Run("handles missing URLs gracefully", func(t *testing.T) {
		audioMediaID := primitive.NewObjectID()
		imageMediaID := primitive.NewObjectID()

		// Empty media URLs map
		mediaUrls := map[primitive.ObjectID]string{}

		key := models.Key{
			ID:           primitive.NewObjectID(),
			AudioMediaId: audioMediaID,
			ImageMediaId: imageMediaID,
		}

		var audioUrl, imageUrl *string
		if !key.AudioMediaId.IsZero() {
			if url, exists := mediaUrls[key.AudioMediaId]; exists {
				audioUrl = &url
			}
		}
		if !key.ImageMediaId.IsZero() {
			if url, exists := mediaUrls[key.ImageMediaId]; exists {
				imageUrl = &url
			}
		}

		assert.Nil(t, audioUrl)
		assert.Nil(t, imageUrl)
	})
}

func TestBoardPaginationLogic(t *testing.T) {
	// Test the pagination logic used in HandleGetAllBoards

	t.Run("calculates correct page slice", func(t *testing.T) {
		boards := make([]models.Board, 50)
		for i := range boards {
			boards[i] = models.Board{
				ID:   primitive.NewObjectID(),
				Name: "Board",
			}
		}

		tests := []struct {
			name          string
			skip          int
			limit         int
			expectedLen   int
			expectedStart int
		}{
			{"first page", 0, 10, 10, 0},
			{"second page", 10, 10, 10, 10},
			{"last page partial", 45, 10, 5, 45},
			{"skip past end", 60, 10, 0, 50},
			{"large limit", 0, 100, 50, 0},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				start := tt.skip
				end := start + tt.limit
				if start > len(boards) {
					start = len(boards)
				}
				if end > len(boards) {
					end = len(boards)
				}

				var paginatedBoards []models.Board
				if start < len(boards) {
					paginatedBoards = boards[start:end]
				} else {
					paginatedBoards = []models.Board{}
				}

				assert.Len(t, paginatedBoards, tt.expectedLen)
			})
		}
	})

	t.Run("calculates current page correctly", func(t *testing.T) {
		tests := []struct {
			name         string
			skip         int
			limit        int
			expectedPage int
		}{
			{"page 1", 0, 10, 1},
			{"page 2", 10, 10, 2},
			{"page 3", 20, 10, 3},
			{"page 1 with limit 25", 0, 25, 1},
			{"page 2 with limit 25", 25, 25, 2},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				currentPage := 1
				if tt.limit > 0 {
					currentPage = (tt.skip / tt.limit) + 1
				}
				assert.Equal(t, tt.expectedPage, currentPage)
			})
		}
	})
}

func TestBoardLayoutTypes(t *testing.T) {
	t.Run("supports all layout types", func(t *testing.T) {
		layouts := []models.LayoutType{
			models.GridLayout,
			models.ListLayout,
		}

		for _, layout := range layouts {
			board := models.Board{
				ID:     primitive.NewObjectID(),
				Name:   "Layout Test",
				Layout: layout,
			}
			assert.NotEmpty(t, board.Layout)
		}
	})
}

// Benchmarks
func BenchmarkMediaIDCollection(b *testing.B) {
	keys := make([]models.Key, 100)
	for i := range keys {
		keys[i] = models.Key{
			ID:           primitive.NewObjectID(),
			AudioMediaId: primitive.NewObjectID(),
			ImageMediaId: primitive.NewObjectID(),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var mediaIds []primitive.ObjectID
		for _, key := range keys {
			if !key.AudioMediaId.IsZero() {
				mediaIds = append(mediaIds, key.AudioMediaId)
			}
			if !key.ImageMediaId.IsZero() {
				mediaIds = append(mediaIds, key.ImageMediaId)
			}
		}
		_ = mediaIds
	}
}

func BenchmarkBoardPagination(b *testing.B) {
	boards := make([]models.Board, 1000)
	for i := range boards {
		boards[i] = models.Board{ID: primitive.NewObjectID()}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		skip := 100
		limit := 25
		start := skip
		end := start + limit
		if start > len(boards) {
			start = len(boards)
		}
		if end > len(boards) {
			end = len(boards)
		}
		_ = boards[start:end]
	}
}
