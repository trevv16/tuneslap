package validation

import (
	"testing"

	api "tuneslap/generated"
	"tuneslap/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewMediaValidator(t *testing.T) {
	v := NewMediaValidator()
	assert.NotNil(t, v)
	assert.NotNil(t, v.Validator)
}

func TestMediaValidator_Validate(t *testing.T) {
	v := NewMediaValidator()

	mediaType := "audio"
	fileName := "test.mp3"
	fileUrl := "https://example.com/test.mp3"
	contentType := "audio/mpeg"
	fileSize := int32(1024)

	createReq := &api.CreateMediaRequest{
		MediaType:   &mediaType,
		FileName:    &fileName,
		FileUrl:     &fileUrl,
		ContentType: &contentType,
		FileSize:    &fileSize,
	}

	desc := "Updated description"
	updateReq := &api.UpdateMediaRequest{
		Description: &desc,
	}

	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
	}{
		{
			name:      "valid CreateMediaRequest",
			data:      createReq,
			expectErr: false,
		},
		{
			name:      "valid UpdateMediaRequest",
			data:      updateReq,
			expectErr: false,
		},
		{
			name:      "invalid type",
			data:      "invalid",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.data)
			if tt.expectErr {
				assert.False(t, result.IsValid)
			} else {
				assert.True(t, result.IsValid)
			}
		})
	}
}

func TestValidateCreateMediaRequest(t *testing.T) {
	tests := []struct {
		name      string
		request   *api.CreateMediaRequest
		expectErr bool
		errFields []string
	}{
		{
			name: "valid audio request",
			request: func() *api.CreateMediaRequest {
				mediaType := "audio"
				fileName := "test.mp3"
				fileUrl := "https://example.com/test.mp3"
				contentType := "audio/mpeg"
				fileSize := int32(1024)
				return &api.CreateMediaRequest{
					MediaType:   &mediaType,
					FileName:    &fileName,
					FileUrl:     &fileUrl,
					ContentType: &contentType,
					FileSize:    &fileSize,
				}
			}(),
			expectErr: false,
		},
		{
			name: "valid image request",
			request: func() *api.CreateMediaRequest {
				mediaType := "image"
				fileName := "test.png"
				fileUrl := "https://example.com/test.png"
				contentType := "image/png"
				fileSize := int32(512000)
				return &api.CreateMediaRequest{
					MediaType:   &mediaType,
					FileName:    &fileName,
					FileUrl:     &fileUrl,
					ContentType: &contentType,
					FileSize:    &fileSize,
				}
			}(),
			expectErr: false,
		},
		{
			name: "valid request with optional fields",
			request: func() *api.CreateMediaRequest {
				mediaType := "audio"
				fileName := "test.mp3"
				fileUrl := "https://example.com/test.mp3"
				contentType := "audio/mpeg"
				fileSize := int32(1024)
				desc := "Test description"
				status := "pending"
				duration := float32(120.5)
				return &api.CreateMediaRequest{
					MediaType:   &mediaType,
					FileName:    &fileName,
					FileUrl:     &fileUrl,
					ContentType: &contentType,
					FileSize:    &fileSize,
					Description: &desc,
					Status:      &status,
					Duration:    &duration,
				}
			}(),
			expectErr: false,
		},
		{
			name: "missing mediaType",
			request: func() *api.CreateMediaRequest {
				fileName := "test.mp3"
				fileUrl := "https://example.com/test.mp3"
				fileSize := int32(1024)
				return &api.CreateMediaRequest{
					FileName: &fileName,
					FileUrl:  &fileUrl,
					FileSize: &fileSize,
				}
			}(),
			expectErr: true,
			errFields: []string{"mediaType"},
		},
		{
			name: "invalid mediaType",
			request: func() *api.CreateMediaRequest {
				mediaType := "video"
				fileName := "test.mp4"
				fileUrl := "https://example.com/test.mp4"
				fileSize := int32(1024)
				return &api.CreateMediaRequest{
					MediaType: &mediaType,
					FileName:  &fileName,
					FileUrl:   &fileUrl,
					FileSize:  &fileSize,
				}
			}(),
			expectErr: true,
			errFields: []string{"mediaType"},
		},
		{
			name: "missing fileName",
			request: func() *api.CreateMediaRequest {
				mediaType := "audio"
				fileUrl := "https://example.com/test.mp3"
				fileSize := int32(1024)
				return &api.CreateMediaRequest{
					MediaType: &mediaType,
					FileUrl:   &fileUrl,
					FileSize:  &fileSize,
				}
			}(),
			expectErr: true,
			errFields: []string{"fileName"},
		},
		{
			name: "fileName too short",
			request: func() *api.CreateMediaRequest {
				mediaType := "audio"
				fileName := "AB"
				fileUrl := "https://example.com/test.mp3"
				fileSize := int32(1024)
				return &api.CreateMediaRequest{
					MediaType: &mediaType,
					FileName:  &fileName,
					FileUrl:   &fileUrl,
					FileSize:  &fileSize,
				}
			}(),
			expectErr: true,
			errFields: []string{"fileName"},
		},
		{
			name: "fileName with invalid characters",
			request: func() *api.CreateMediaRequest {
				mediaType := "audio"
				fileName := "test<file>.mp3"
				fileUrl := "https://example.com/test.mp3"
				fileSize := int32(1024)
				return &api.CreateMediaRequest{
					MediaType: &mediaType,
					FileName:  &fileName,
					FileUrl:   &fileUrl,
					FileSize:  &fileSize,
				}
			}(),
			expectErr: true,
			errFields: []string{"fileName"},
		},
		{
			name: "missing fileUrl",
			request: func() *api.CreateMediaRequest {
				mediaType := "audio"
				fileName := "test.mp3"
				fileSize := int32(1024)
				return &api.CreateMediaRequest{
					MediaType: &mediaType,
					FileName:  &fileName,
					FileSize:  &fileSize,
				}
			}(),
			expectErr: true,
			errFields: []string{"fileUrl"},
		},
		{
			name: "invalid fileUrl",
			request: func() *api.CreateMediaRequest {
				mediaType := "audio"
				fileName := "test.mp3"
				fileUrl := "not-a-url"
				fileSize := int32(1024)
				return &api.CreateMediaRequest{
					MediaType: &mediaType,
					FileName:  &fileName,
					FileUrl:   &fileUrl,
					FileSize:  &fileSize,
				}
			}(),
			expectErr: true,
			errFields: []string{"fileUrl"},
		},
		{
			name: "missing fileSize",
			request: func() *api.CreateMediaRequest {
				mediaType := "audio"
				fileName := "test.mp3"
				fileUrl := "https://example.com/test.mp3"
				return &api.CreateMediaRequest{
					MediaType: &mediaType,
					FileName:  &fileName,
					FileUrl:   &fileUrl,
				}
			}(),
			expectErr: true,
			errFields: []string{"fileSize"},
		},
		{
			name: "fileSize too small",
			request: func() *api.CreateMediaRequest {
				mediaType := "audio"
				fileName := "test.mp3"
				fileUrl := "https://example.com/test.mp3"
				fileSize := int32(0)
				return &api.CreateMediaRequest{
					MediaType: &mediaType,
					FileName:  &fileName,
					FileUrl:   &fileUrl,
					FileSize:  &fileSize,
				}
			}(),
			expectErr: true,
			errFields: []string{"fileSize"},
		},
		{
			name: "invalid status",
			request: func() *api.CreateMediaRequest {
				mediaType := "audio"
				fileName := "test.mp3"
				fileUrl := "https://example.com/test.mp3"
				fileSize := int32(1024)
				status := "invalid"
				return &api.CreateMediaRequest{
					MediaType: &mediaType,
					FileName:  &fileName,
					FileUrl:   &fileUrl,
					FileSize:  &fileSize,
					Status:    &status,
				}
			}(),
			expectErr: true,
			errFields: []string{"status"},
		},
		{
			name: "invalid dimensions count",
			request: func() *api.CreateMediaRequest {
				mediaType := "image"
				fileName := "test.png"
				fileUrl := "https://example.com/test.png"
				fileSize := int32(1024)
				dimensions := []int32{100}
				return &api.CreateMediaRequest{
					MediaType:  &mediaType,
					FileName:   &fileName,
					FileUrl:    &fileUrl,
					FileSize:   &fileSize,
					Dimensions: dimensions,
				}
			}(),
			expectErr: true,
			errFields: []string{"dimensions"},
		},
		{
			name: "negative duration",
			request: func() *api.CreateMediaRequest {
				mediaType := "audio"
				fileName := "test.mp3"
				fileUrl := "https://example.com/test.mp3"
				fileSize := int32(1024)
				duration := float32(-5)
				return &api.CreateMediaRequest{
					MediaType: &mediaType,
					FileName:  &fileName,
					FileUrl:   &fileUrl,
					FileSize:  &fileSize,
					Duration:  &duration,
				}
			}(),
			expectErr: true,
			errFields: []string{"duration"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateCreateMediaRequest(tt.request)
			if tt.expectErr {
				assert.False(t, result.IsValid)
				assert.NotEmpty(t, result.Errors)
				for _, expectedField := range tt.errFields {
					found := false
					for _, err := range result.Errors {
						if err.Field == expectedField {
							found = true
							break
						}
					}
					assert.True(t, found, "Expected error for field %s", expectedField)
				}
			} else {
				assert.True(t, result.IsValid)
				assert.Empty(t, result.Errors)
			}
		})
	}
}

func TestValidateUpdateMediaRequest(t *testing.T) {
	tests := []struct {
		name      string
		request   *api.UpdateMediaRequest
		expectErr bool
		errFields []string
	}{
		{
			name:      "valid empty request",
			request:   &api.UpdateMediaRequest{},
			expectErr: false,
		},
		{
			name: "valid request with description",
			request: func() *api.UpdateMediaRequest {
				desc := "Updated description"
				return &api.UpdateMediaRequest{
					Description: &desc,
				}
			}(),
			expectErr: false,
		},
		{
			name: "description too long",
			request: func() *api.UpdateMediaRequest {
				desc := string(make([]byte, 1001))
				return &api.UpdateMediaRequest{
					Description: &desc,
				}
			}(),
			expectErr: true,
			errFields: []string{"description"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateUpdateMediaRequest(tt.request)
			if tt.expectErr {
				assert.False(t, result.IsValid)
				assert.NotEmpty(t, result.Errors)
			} else {
				assert.True(t, result.IsValid)
			}
		})
	}
}

func TestMediaValidator_ValidateMediaOwnership(t *testing.T) {
	v := NewMediaValidator()

	authorID := primitive.NewObjectID()
	otherUserID := primitive.NewObjectID()

	media := models.Media{
		ID:       primitive.NewObjectID(),
		AuthorId: authorID,
	}

	tests := []struct {
		name      string
		media     models.Media
		userID    primitive.ObjectID
		expectErr bool
	}{
		{
			name:      "owner can access",
			media:     media,
			userID:    authorID,
			expectErr: false,
		},
		{
			name:      "non-owner cannot access",
			media:     media,
			userID:    otherUserID,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateMediaOwnership(tt.media, tt.userID)
			if tt.expectErr {
				assert.False(t, result.IsValid)
			} else {
				assert.True(t, result.IsValid)
			}
		})
	}
}

func TestMediaValidator_ValidateFileSize(t *testing.T) {
	v := NewMediaValidator()

	tests := []struct {
		name      string
		fileSize  int64
		maxSize   int64
		expectErr bool
	}{
		{
			name:      "file size within limit",
			fileSize:  1024,
			maxSize:   MaxFileSize,
			expectErr: false,
		},
		{
			name:      "file size at limit",
			fileSize:  MaxFileSize,
			maxSize:   MaxFileSize,
			expectErr: false,
		},
		{
			name:      "file size exceeds limit",
			fileSize:  MaxFileSize + 1,
			maxSize:   MaxFileSize,
			expectErr: true,
		},
		{
			name:      "zero file size",
			fileSize:  0,
			maxSize:   MaxFileSize,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateFileSize(tt.fileSize, tt.maxSize)
			if tt.expectErr {
				assert.False(t, result.IsValid)
			} else {
				assert.True(t, result.IsValid)
			}
		})
	}
}

func TestValidateMediaProcessingParamsAudio(t *testing.T) {
	tests := []struct {
		name      string
		params    *api.MediaProcessingParamsAudio
		expectErr bool
		errFields []string
	}{
		{
			name:      "nil params",
			params:    nil,
			expectErr: false,
		},
		{
			name:      "empty params",
			params:    &api.MediaProcessingParamsAudio{},
			expectErr: false,
		},
		{
			name: "valid params",
			params: func() *api.MediaProcessingParamsAudio {
				trimStart := float32(0.5)
				trimEnd := float32(10.0)
				fadeIn := float32(0.2)
				fadeOut := float32(0.3)
				speed := float32(1.25)
				return &api.MediaProcessingParamsAudio{
					TrimStart: &trimStart,
					TrimEnd:   &trimEnd,
					FadeIn:    &fadeIn,
					FadeOut:   &fadeOut,
					Speed:     &speed,
				}
			}(),
			expectErr: false,
		},
		{
			name: "negative trimStart",
			params: func() *api.MediaProcessingParamsAudio {
				trimStart := float32(-1.0)
				return &api.MediaProcessingParamsAudio{
					TrimStart: &trimStart,
				}
			}(),
			expectErr: true,
			errFields: []string{"trimStart"},
		},
		{
			name: "negative trimEnd",
			params: func() *api.MediaProcessingParamsAudio {
				trimEnd := float32(-1.0)
				return &api.MediaProcessingParamsAudio{
					TrimEnd: &trimEnd,
				}
			}(),
			expectErr: true,
			errFields: []string{"trimEnd"},
		},
		{
			name: "negative fadeIn",
			params: func() *api.MediaProcessingParamsAudio {
				fadeIn := float32(-0.5)
				return &api.MediaProcessingParamsAudio{
					FadeIn: &fadeIn,
				}
			}(),
			expectErr: true,
			errFields: []string{"fadeIn"},
		},
		{
			name: "zero speed",
			params: func() *api.MediaProcessingParamsAudio {
				speed := float32(0)
				return &api.MediaProcessingParamsAudio{
					Speed: &speed,
				}
			}(),
			expectErr: true,
			errFields: []string{"speed"},
		},
		{
			name: "negative speed",
			params: func() *api.MediaProcessingParamsAudio {
				speed := float32(-1.0)
				return &api.MediaProcessingParamsAudio{
					Speed: &speed,
				}
			}(),
			expectErr: true,
			errFields: []string{"speed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateMediaProcessingParamsAudio(tt.params)
			if tt.expectErr {
				assert.False(t, result.IsValid)
				assert.NotEmpty(t, result.Errors)
			} else {
				assert.True(t, result.IsValid)
			}
		})
	}
}

// Benchmark tests
func BenchmarkValidateCreateMediaRequest(b *testing.B) {
	mediaType := "audio"
	fileName := "test.mp3"
	fileUrl := "https://example.com/test.mp3"
	contentType := "audio/mpeg"
	fileSize := int32(1024)
	req := &api.CreateMediaRequest{
		MediaType:   &mediaType,
		FileName:    &fileName,
		FileUrl:     &fileUrl,
		ContentType: &contentType,
		FileSize:    &fileSize,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateCreateMediaRequest(req)
	}
}
