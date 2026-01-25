package utils

import (
	"testing"
	"time"

	"tuneslap/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestToUserResponse(t *testing.T) {
	now := time.Now()
	userID := primitive.NewObjectID()

	user := models.User{
		ID:        userID,
		Name:      "Test User",
		Email:     "test@example.com",
		ImageUrl:  "https://example.com/avatar.png",
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := ToUserResponse(user)

	assert.NotNil(t, result)
	assert.Equal(t, userID.Hex(), result.Id)
	assert.Equal(t, "Test User", result.Name)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "https://example.com/avatar.png", *result.ImageUrl)
	assert.NotZero(t, result.CreatedAt)
	assert.NotZero(t, result.UpdatedAt)
}

func TestToCollaboratorResponse(t *testing.T) {
	now := time.Now()
	collabID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	collaborator := models.Collaborator{
		ID:        collabID,
		UserId:    userID,
		Email:     "collab@example.com",
		Role:      "editor",
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := ToCollaboratorResponse(collaborator)

	assert.NotNil(t, result)
	assert.Equal(t, collabID.Hex(), result.Id)
	assert.Equal(t, userID.Hex(), result.UserId)
	assert.Equal(t, "collab@example.com", result.Email)
	assert.Equal(t, "editor", result.Role)
}

func TestToKeyResponse(t *testing.T) {
	now := time.Now()
	keyID := primitive.NewObjectID()
	boardID := primitive.NewObjectID()
	audioMediaID := primitive.NewObjectID()
	imageMediaID := primitive.NewObjectID()

	key := models.Key{
		ID:           keyID,
		BoardId:      boardID,
		Name:         "Test Key",
		Description:  "A test key",
		AudioMediaId: audioMediaID,
		AudioUrl:     "https://example.com/audio.mp3",
		ImageMediaId: imageMediaID,
		ImageUrl:     "https://example.com/image.png",
		HotKey:       "A",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	result := ToKeyResponse(key)

	assert.NotNil(t, result)
	assert.Equal(t, keyID.Hex(), result.Id)
	assert.Equal(t, boardID.Hex(), result.BoardId)
	assert.Equal(t, "Test Key", result.Name)
	assert.Equal(t, "A test key", *result.Description)
	assert.Equal(t, audioMediaID.Hex(), result.AudioMediaId)
	assert.Equal(t, "A", result.HotKey)
}

func TestToMediaProcessingParamsAudio(t *testing.T) {
	params := models.AudioProcessingParams{
		ContentType:   "audio/webm",
		TrimStart:     1.5,
		TrimEnd:       10.0,
		Normalize:     true,
		FadeIn:        0.5,
		FadeOut:       0.5,
		Speed:         1.25,
		Pitch:         1.0,
		OutputFormats: []string{"webm", "mp3"},
	}

	result := ToMediaProcessingParamsAudio(params)

	assert.NotNil(t, result)
	assert.Equal(t, "audio/webm", *result.ContentType)
	assert.Equal(t, float32(1.5), *result.TrimStart)
	assert.Equal(t, float32(10.0), *result.TrimEnd)
	assert.True(t, *result.Normalize)
	assert.Equal(t, float32(0.5), *result.FadeIn)
	assert.Equal(t, float32(0.5), *result.FadeOut)
	assert.Equal(t, float32(1.25), *result.Speed)
}

func TestToMediaProcessingParamsImage(t *testing.T) {
	params := models.ImageProcessingParams{
		ResizeTo:     [2]int{1920, 1080},
		Format:       "webp",
		Crop:         [4]int{10, 20, 100, 200},
		AspectRatio:  "16:9",
		ApplyFilters: "grayscale",
	}

	result := ToMediaProcessingParamsImage(params)

	assert.NotNil(t, result)
	assert.Equal(t, []int32{1920, 1080}, result.ResizeTo)
	assert.Equal(t, []int32{10, 20, 100, 200}, result.Crop)
	assert.Equal(t, "16:9", *result.AspectRatio)
	assert.Contains(t, result.ApplyFilters, "grayscale")
}

func TestToMediaProcessingParams(t *testing.T) {
	audioParams := &models.AudioProcessingParams{
		ContentType: "audio/webm",
		TrimStart:   1.0,
	}
	imageParams := &models.ImageProcessingParams{
		ResizeTo: [2]int{800, 600},
	}

	tests := []struct {
		name   string
		params models.ProcessingParams
	}{
		{
			name:   "empty params",
			params: models.ProcessingParams{},
		},
		{
			name: "audio only",
			params: models.ProcessingParams{
				Audio: audioParams,
			},
		},
		{
			name: "image only",
			params: models.ProcessingParams{
				Image: imageParams,
			},
		},
		{
			name: "both audio and image",
			params: models.ProcessingParams{
				Audio: audioParams,
				Image: imageParams,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToMediaProcessingParams(tt.params)
			assert.NotNil(t, result)

			if tt.params.Audio != nil {
				assert.NotNil(t, result.Audio)
			}
			if tt.params.Image != nil {
				assert.NotNil(t, result.Image)
			}
		})
	}
}

func TestToMediaProcessingActivity(t *testing.T) {
	now := time.Now()
	activity := models.ProcessingActivity{
		Status:    "done",
		Message:   "Processing completed",
		CreatedAt: primitive.NewDateTimeFromTime(now),
		UpdatedAt: primitive.NewDateTimeFromTime(now),
	}

	result := ToMediaProcessingActivity(activity)

	assert.NotNil(t, result)
	assert.Equal(t, "done", result.Status)
	assert.Equal(t, "Processing completed", result.Message)
	assert.NotZero(t, result.CreatedAt)
	assert.NotZero(t, result.UpdatedAt)
}

func TestToMediaResponse(t *testing.T) {
	now := primitive.NewDateTimeFromTime(time.Now())
	mediaID := primitive.NewObjectID()
	authorID := primitive.NewObjectID()

	media := models.Media{
		ID:           mediaID,
		AuthorId:     authorID,
		MediaType:    "audio",
		FileName:     "test.mp3",
		Description:  "Test audio file",
		FileUrl:      "https://example.com/test.mp3",
		ProcessedUrl: "https://example.com/processed/test.webm",
		ContentType:  "audio/mpeg",
		FileSize:     1024000,
		Status:       "done",
		Dimensions:   [2]int{0, 0},
		Duration:     120.5,
		CreatedAt:    now,
		UpdatedAt:    now,
		ProcessingActivity: []models.ProcessingActivity{
			{
				Status:    "done",
				Message:   "Processing completed",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}

	result := ToMediaResponse(media)

	assert.NotNil(t, result)
	assert.Equal(t, mediaID.Hex(), result.Id)
	assert.Equal(t, authorID.Hex(), result.AuthorId)
	assert.Equal(t, "audio", result.MediaType)
	assert.Equal(t, "test.mp3", result.FileName)
	assert.Equal(t, "done", result.Status)
	assert.Equal(t, float32(120.5), *result.Duration)
}

func TestToBoardResponse(t *testing.T) {
	now := time.Now()
	boardID := primitive.NewObjectID()
	authorID := primitive.NewObjectID()

	board := models.Board{
		ID:          boardID,
		AuthorId:    authorID,
		Name:        "Test Board",
		Description: "A test board",
		Layout:      models.GridLayout,
		ImageUrl:    "https://example.com/board.png",
		Collaborators: []models.Collaborator{
			{
				ID:        primitive.NewObjectID(),
				UserId:    primitive.NewObjectID(),
				Email:     "collab@example.com",
				Role:      "editor",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		Keys: []models.Key{
			{
				ID:           primitive.NewObjectID(),
				BoardId:      boardID,
				Name:         "Key 1",
				HotKey:       "A",
				AudioMediaId: primitive.NewObjectID(),
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := ToBoardResponse(board)

	assert.NotNil(t, result)
	assert.Equal(t, boardID.Hex(), result.Id)
	assert.Equal(t, authorID.Hex(), result.AuthorId)
	assert.Equal(t, "Test Board", result.Name)
	assert.Equal(t, "grid", result.Layout)
	assert.Len(t, result.Collaborators, 1)
	assert.Len(t, result.Keys, 1)
}

// Benchmark tests
func BenchmarkToMediaResponse(b *testing.B) {
	now := primitive.NewDateTimeFromTime(time.Now())
	media := models.Media{
		ID:          primitive.NewObjectID(),
		AuthorId:    primitive.NewObjectID(),
		MediaType:   "audio",
		FileName:    "test.mp3",
		Description: "Test audio file",
		FileUrl:     "https://example.com/test.mp3",
		ContentType: "audio/mpeg",
		FileSize:    1024000,
		Status:      "done",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToMediaResponse(media)
	}
}

func BenchmarkToBoardResponse(b *testing.B) {
	now := time.Now()
	board := models.Board{
		ID:          primitive.NewObjectID(),
		AuthorId:    primitive.NewObjectID(),
		Name:        "Test Board",
		Description: "A test board",
		Layout:      models.GridLayout,
		Collaborators: []models.Collaborator{
			{ID: primitive.NewObjectID(), Email: "test@example.com", Role: "editor", CreatedAt: now, UpdatedAt: now},
		},
		Keys: []models.Key{
			{ID: primitive.NewObjectID(), Name: "Key 1", HotKey: "A", CreatedAt: now, UpdatedAt: now},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToBoardResponse(board)
	}
}
