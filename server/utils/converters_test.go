package utils

import (
	"testing"
	"time"

	api "tuneslap/generated"
	"tuneslap/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestObjectIDToString(t *testing.T) {
	tests := []struct {
		name     string
		id       primitive.ObjectID
		expected bool // whether result should be non-nil
	}{
		{
			name:     "valid ObjectID",
			id:       primitive.NewObjectID(),
			expected: true,
		},
		{
			name:     "zero ObjectID",
			id:       primitive.NilObjectID,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := objectIDToString(tt.id)
			if tt.expected {
				assert.NotNil(t, result)
				assert.Equal(t, tt.id.Hex(), *result)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestStringToObjectID(t *testing.T) {
	validID := primitive.NewObjectID()

	tests := []struct {
		name        string
		input       *string
		expectError bool
		expectNil   bool
	}{
		{
			name:        "valid ObjectID string",
			input:       strPtr(validID.Hex()),
			expectError: false,
			expectNil:   false,
		},
		{
			name:        "nil input",
			input:       nil,
			expectError: false,
			expectNil:   true,
		},
		{
			name:        "empty string",
			input:       strPtr(""),
			expectError: false,
			expectNil:   true,
		},
		{
			name:        "invalid ObjectID string",
			input:       strPtr("invalid"),
			expectError: true,
			expectNil:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := stringToObjectID(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.expectNil {
					assert.Equal(t, primitive.NilObjectID, result)
				} else {
					assert.NotEqual(t, primitive.NilObjectID, result)
				}
			}
		})
	}
}

func TestDateTimeToTime(t *testing.T) {
	now := time.Now()
	validDateTime := primitive.NewDateTimeFromTime(now)
	invalidDateTime := primitive.DateTime(0) // Unix epoch, should be valid

	tests := []struct {
		name  string
		input primitive.DateTime
	}{
		{
			name:  "valid DateTime",
			input: validDateTime,
		},
		{
			name:  "zero DateTime",
			input: invalidDateTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dateTimeToTime(tt.input)
			assert.NotNil(t, result)
		})
	}
}

func TestTimeToDateTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name  string
		input *time.Time
	}{
		{
			name:  "valid time",
			input: &now,
		},
		{
			name:  "nil time",
			input: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := timeToDateTime(tt.input)
			assert.NotZero(t, result)
		})
	}
}

func TestStringPtr(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool // whether result should be non-nil
	}{
		{
			name:     "non-empty string",
			input:    "test",
			expected: true,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stringPtr(tt.input)
			if tt.expected {
				assert.NotNil(t, result)
				assert.Equal(t, tt.input, *result)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestStringVal(t *testing.T) {
	value := "test"

	tests := []struct {
		name     string
		input    *string
		expected string
	}{
		{
			name:     "non-nil pointer",
			input:    &value,
			expected: "test",
		},
		{
			name:     "nil pointer",
			input:    nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stringVal(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFloat64ToFloat32(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float32
	}{
		{
			name:     "positive value",
			input:    10.5,
			expected: 10.5,
		},
		{
			name:     "zero value",
			input:    0,
			expected: 0,
		},
		{
			name:     "negative value",
			input:    -5.25,
			expected: -5.25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := float64ToFloat32(tt.input)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expected, *result)
		})
	}
}

func TestFloat32ToFloat64(t *testing.T) {
	value := float32(10.5)

	tests := []struct {
		name     string
		input    *float32
		expected float64
	}{
		{
			name:     "non-nil value",
			input:    &value,
			expected: 10.5,
		},
		{
			name:     "nil value",
			input:    nil,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := float32ToFloat64(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIntToInt32(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int32
	}{
		{
			name:     "positive value",
			input:    100,
			expected: 100,
		},
		{
			name:     "zero value",
			input:    0,
			expected: 0,
		},
		{
			name:     "negative value",
			input:    -50,
			expected: -50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := intToInt32(tt.input)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expected, *result)
		})
	}
}

func TestInt32ToInt(t *testing.T) {
	value := int32(100)

	tests := []struct {
		name     string
		input    *int32
		expected int
	}{
		{
			name:     "non-nil value",
			input:    &value,
			expected: 100,
		},
		{
			name:     "nil value",
			input:    nil,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := int32ToInt(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIntArrayToInt32Array(t *testing.T) {
	tests := []struct {
		name     string
		input    [2]int
		expected []int32
	}{
		{
			name:     "valid dimensions",
			input:    [2]int{1920, 1080},
			expected: []int32{1920, 1080},
		},
		{
			name:     "zero values",
			input:    [2]int{0, 0},
			expected: []int32{0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := intArrayToInt32Array(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInt32ArrayToIntArray(t *testing.T) {
	tests := []struct {
		name     string
		input    []int32
		expected [2]int
	}{
		{
			name:     "valid dimensions",
			input:    []int32{1920, 1080},
			expected: [2]int{1920, 1080},
		},
		{
			name:     "less than 2 elements",
			input:    []int32{100},
			expected: [2]int{0, 0},
		},
		{
			name:     "more than 2 elements",
			input:    []int32{100, 200, 300},
			expected: [2]int{100, 200},
		},
		{
			name:     "empty array",
			input:    []int32{},
			expected: [2]int{0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := int32ArrayToIntArray(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInt32ArrayToInt4Array(t *testing.T) {
	tests := []struct {
		name     string
		input    []int32
		expected [4]int
	}{
		{
			name:     "valid crop values",
			input:    []int32{10, 20, 100, 200},
			expected: [4]int{10, 20, 100, 200},
		},
		{
			name:     "less than 4 elements",
			input:    []int32{10, 20},
			expected: [4]int{0, 0, 0, 0},
		},
		{
			name:     "empty array",
			input:    []int32{},
			expected: [4]int{0, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := int32ArrayToInt4Array(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

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
	assert.Equal(t, userID.Hex(), *result.Id)
	assert.Equal(t, "Test User", *result.Name)
	assert.Equal(t, "test@example.com", *result.Email)
	assert.Equal(t, "https://example.com/avatar.png", *result.ImageUrl)
	assert.NotNil(t, result.CreatedAt)
	assert.NotNil(t, result.UpdatedAt)
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
	assert.Equal(t, collabID.Hex(), *result.Id)
	assert.Equal(t, userID.Hex(), *result.UserId)
	assert.Equal(t, "collab@example.com", *result.Email)
	assert.Equal(t, "editor", *result.Role)
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
	assert.Equal(t, keyID.Hex(), *result.Id)
	assert.Equal(t, boardID.Hex(), *result.BoardId)
	assert.Equal(t, "Test Key", *result.Name)
	assert.Equal(t, "A test key", *result.Description)
	assert.Equal(t, audioMediaID.Hex(), *result.AudioMediaId)
	assert.Equal(t, "A", *result.HotKey)
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
		Format:       0,
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
	assert.Equal(t, "done", *result.Status)
	assert.Equal(t, "Processing completed", *result.Message)
	assert.NotNil(t, result.CreatedAt)
	assert.NotNil(t, result.UpdatedAt)
}

func TestProcessingActivityFromAPI(t *testing.T) {
	now := time.Now()
	status := "processing"
	message := "Processing in progress"

	activity := &api.MediaProcessingActivity{
		Status:    &status,
		Message:   &message,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	result := ProcessingActivityFromAPI(activity)

	assert.Equal(t, "processing", result.Status)
	assert.Equal(t, "Processing in progress", result.Message)
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
	assert.Equal(t, mediaID.Hex(), *result.Id)
	assert.Equal(t, authorID.Hex(), *result.AuthorId)
	assert.Equal(t, "audio", *result.MediaType)
	assert.Equal(t, "test.mp3", *result.FileName)
	assert.Equal(t, "done", *result.Status)
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
	assert.Equal(t, boardID.Hex(), *result.Id)
	assert.Equal(t, authorID.Hex(), *result.AuthorId)
	assert.Equal(t, "Test Board", *result.Name)
	assert.Equal(t, "grid", *result.Layout)
	assert.Len(t, result.Collaborators, 1)
	assert.Len(t, result.Keys, 1)
}

func TestBoardFromCreateRequest(t *testing.T) {
	authorID := primitive.NewObjectID()
	name := "New Board"
	desc := "A new board"
	layout := "grid"
	imageUrl := "https://example.com/image.png"

	req := &api.CreateBoardRequest{
		Name:        &name,
		Description: &desc,
		Layout:      &layout,
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
	name := "New Board"
	layout := "list"

	req := &api.CreateBoardRequest{
		Name:   &name,
		Layout: &layout,
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
	mediaType := "audio"
	fileName := "test.mp3"
	desc := "Test file"
	fileUrl := "https://example.com/test.mp3"
	contentType := "audio/mpeg"
	fileSize := int32(1024)
	duration := float32(60.0)

	req := &api.CreateMediaRequest{
		MediaType:   &mediaType,
		FileName:    &fileName,
		Description: &desc,
		FileUrl:     &fileUrl,
		ContentType: &contentType,
		FileSize:    &fileSize,
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
	email := "test@example.com"
	role := "editor"

	req := &api.CreateCollaboratorRequest{
		Email: &email,
		Role:  &role,
	}

	result := CollaboratorFromCreateRequest(req)

	assert.NotEqual(t, primitive.NilObjectID, result.ID)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "editor", result.Role)
}

// Helper function for tests
func strPtr(s string) *string {
	return &s
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
