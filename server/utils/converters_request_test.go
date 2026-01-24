package utils

import (
	"testing"
	"time"

	api "tuneslap/generated"
	"tuneslap/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestProcessingParamsFromAPI(t *testing.T) {
	trimStart := float32(1.0)
	audioParams := &api.MediaProcessingParamsAudio{
		TrimStart: &trimStart,
	}
	imageParams := &api.MediaProcessingParamsImage{
		ResizeTo: []int32{800, 600},
	}

	tests := []struct {
		name   string
		params *api.MediaProcessingParams
	}{
		{
			name:   "nil params",
			params: nil,
		},
		{
			name:   "empty params",
			params: &api.MediaProcessingParams{},
		},
		{
			name: "audio only",
			params: &api.MediaProcessingParams{
				Audio: audioParams,
			},
		},
		{
			name: "image only",
			params: &api.MediaProcessingParams{
				Image: imageParams,
			},
		},
		{
			name: "both audio and image",
			params: &api.MediaProcessingParams{
				Audio: audioParams,
				Image: imageParams,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessingParamsFromAPI(tt.params)
			// Result should always be a valid struct
			assert.NotNil(t, &result)

			if tt.params != nil && tt.params.Audio != nil {
				assert.NotNil(t, result.Audio)
			}
			if tt.params != nil && tt.params.Image != nil {
				assert.NotNil(t, result.Image)
			}
		})
	}
}

func TestProcessingActivityFromAPI(t *testing.T) {
	now := time.Now()

	activity := &api.MediaProcessingActivity{
		Status:    "processing",
		Message:   "Processing in progress",
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := ProcessingActivityFromAPI(activity)

	assert.Equal(t, "processing", result.Status)
	assert.Equal(t, "Processing in progress", result.Message)
	assert.NotZero(t, result.CreatedAt)
	assert.NotZero(t, result.UpdatedAt)
}

func TestBoardFromCreateRequest(t *testing.T) {
	authorID := primitive.NewObjectID()
	desc := "A new board"
	imageUrl := "https://example.com/image.png"

	req := &api.CreateBoardRequest{
		Name:        "New Board",
		Description: &desc,
		Layout:      "grid",
		ImageUrl:    &imageUrl,
	}

	result := BoardFromCreateRequest(req, authorID)

	assert.NotEqual(t, primitive.NilObjectID, result.ID)
	assert.Equal(t, authorID, result.AuthorId)
	assert.Equal(t, "New Board", result.Name)
	assert.Equal(t, "A new board", result.Description)
	assert.Equal(t, models.GridLayout, result.Layout)
	assert.Equal(t, "https://example.com/image.png", result.ImageUrl)
	assert.Empty(t, result.Collaborators)
	assert.Empty(t, result.Keys)
}

func TestBoardFromCreateRequest_ListLayout(t *testing.T) {
	authorID := primitive.NewObjectID()

	req := &api.CreateBoardRequest{
		Name:   "New Board",
		Layout: "list",
	}

	result := BoardFromCreateRequest(req, authorID)

	assert.Equal(t, models.ListLayout, result.Layout)
}

func TestBoardFromUpdateRequest(t *testing.T) {
	now := time.Now()
	boardID := primitive.NewObjectID()

	existingBoard := models.Board{
		ID:          boardID,
		AuthorId:    primitive.NewObjectID(),
		Name:        "Old Name",
		Description: "Old description",
		Layout:      models.GridLayout,
		ImageUrl:    "https://example.com/old.png",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	newName := "New Name"
	newDesc := "New description"
	newLayout := "list"
	newImageUrl := "https://example.com/new.png"

	req := &api.UpdateBoardRequest{
		Name:        &newName,
		Description: &newDesc,
		Layout:      &newLayout,
		ImageUrl:    &newImageUrl,
	}

	result := BoardFromUpdateRequest(existingBoard, req)

	assert.Equal(t, boardID, result.ID)
	assert.Equal(t, "New Name", result.Name)
	assert.Equal(t, "New description", result.Description)
	assert.Equal(t, models.ListLayout, result.Layout)
	assert.Equal(t, "https://example.com/new.png", result.ImageUrl)
}

func TestMediaFromCreateRequest(t *testing.T) {
	authorID := primitive.NewObjectID()
	desc := "Test file"
	contentType := "audio/mpeg"
	duration := float32(60.0)

	req := &api.CreateMediaRequest{
		MediaType:   "audio",
		FileName:    "test.mp3",
		Description: &desc,
		FileUrl:     "https://example.com/test.mp3",
		ContentType: &contentType,
		FileSize:    1024,
		Duration:    &duration,
	}

	result := MediaFromCreateRequest(req, authorID)

	assert.NotEqual(t, primitive.NilObjectID, result.ID)
	assert.Equal(t, authorID, result.AuthorId)
	assert.Equal(t, "audio", result.MediaType)
	assert.Equal(t, "test.mp3", result.FileName)
	assert.Equal(t, "Test file", result.Description)
	assert.Equal(t, "https://example.com/test.mp3", result.FileUrl)
	assert.Equal(t, "audio/mpeg", result.ContentType)
	assert.Equal(t, int64(1024), result.FileSize)
	assert.Equal(t, float64(60.0), result.Duration)
	assert.Equal(t, models.ProcessingStatusPending, result.Status)
	assert.NotEmpty(t, result.ProcessingActivity)
}

func TestKeyFromCreateRequest(t *testing.T) {
	boardID := primitive.NewObjectID()
	audioMediaID := primitive.NewObjectID()
	imageMediaID := primitive.NewObjectID()
	desc := "Test key"

	req := &api.CreateKeyRequest{
		Name:         "Test Key",
		Description:  &desc,
		AudioMediaId: audioMediaID.Hex(),
		ImageMediaId: strPtr(imageMediaID.Hex()),
		HotKey:       "A",
	}

	result := KeyFromCreateRequest(req, boardID)

	assert.NotEqual(t, primitive.NilObjectID, result.ID)
	assert.Equal(t, boardID, result.BoardId)
	assert.Equal(t, "Test Key", result.Name)
	assert.Equal(t, "Test key", result.Description)
	assert.Equal(t, audioMediaID, result.AudioMediaId)
	assert.Equal(t, imageMediaID, result.ImageMediaId)
	assert.Equal(t, "A", result.HotKey)
}

func TestCollaboratorFromCreateRequest(t *testing.T) {
	req := &api.CreateCollaboratorRequest{
		Email: "test@example.com",
		Role:  "editor",
	}

	result := CollaboratorFromCreateRequest(req)

	assert.NotEqual(t, primitive.NilObjectID, result.ID)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "editor", result.Role)
}
