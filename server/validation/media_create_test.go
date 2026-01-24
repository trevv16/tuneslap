package validation

import (
	"testing"

	api "tuneslap/generated"

	"github.com/stretchr/testify/assert"
)

func TestNewMediaValidator(t *testing.T) {
	v := NewMediaValidator()
	assert.NotNil(t, v)
	assert.NotNil(t, v.Validator)
}

func TestMediaValidator_Validate(t *testing.T) {
	v := NewMediaValidator()

	contentType := "audio/mpeg"

	createReq := &api.CreateMediaRequest{
		MediaType:   "audio",
		FileName:    "test.mp3",
		FileUrl:     "https://example.com/test.mp3",
		ContentType: &contentType,
		FileSize:    1024,
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
				contentType := "audio/mpeg"
				return &api.CreateMediaRequest{
					MediaType:   "audio",
					FileName:    "test.mp3",
					FileUrl:     "https://example.com/test.mp3",
					ContentType: &contentType,
					FileSize:    1024,
				}
			}(),
			expectErr: false,
		},
		{
			name: "valid image request",
			request: func() *api.CreateMediaRequest {
				contentType := "image/png"
				return &api.CreateMediaRequest{
					MediaType:   "image",
					FileName:    "test.png",
					FileUrl:     "https://example.com/test.png",
					ContentType: &contentType,
					FileSize:    512000,
				}
			}(),
			expectErr: false,
		},
		{
			name: "valid request with optional fields",
			request: func() *api.CreateMediaRequest {
				contentType := "audio/mpeg"
				desc := "Test description"
				status := "pending"
				duration := float32(120.5)
				return &api.CreateMediaRequest{
					MediaType:   "audio",
					FileName:    "test.mp3",
					FileUrl:     "https://example.com/test.mp3",
					ContentType: &contentType,
					FileSize:    1024,
					Description: &desc,
					Status:      &status,
					Duration:    &duration,
				}
			}(),
			expectErr: false,
		},
		{
			name: "missing mediaType",
			request: &api.CreateMediaRequest{
				MediaType: "",
				FileName:  "test.mp3",
				FileUrl:   "https://example.com/test.mp3",
				FileSize:  1024,
			},
			expectErr: true,
			errFields: []string{"mediaType"},
		},
		{
			name: "invalid mediaType",
			request: &api.CreateMediaRequest{
				MediaType: "video",
				FileName:  "test.mp4",
				FileUrl:   "https://example.com/test.mp4",
				FileSize:  1024,
			},
			expectErr: true,
			errFields: []string{"mediaType"},
		},
		{
			name: "missing fileName",
			request: &api.CreateMediaRequest{
				MediaType: "audio",
				FileName:  "",
				FileUrl:   "https://example.com/test.mp3",
				FileSize:  1024,
			},
			expectErr: true,
			errFields: []string{"fileName"},
		},
		{
			name: "fileName too short",
			request: &api.CreateMediaRequest{
				MediaType: "audio",
				FileName:  "AB",
				FileUrl:   "https://example.com/test.mp3",
				FileSize:  1024,
			},
			expectErr: true,
			errFields: []string{"fileName"},
		},
		{
			name: "fileName with invalid characters",
			request: &api.CreateMediaRequest{
				MediaType: "audio",
				FileName:  "test<file>.mp3",
				FileUrl:   "https://example.com/test.mp3",
				FileSize:  1024,
			},
			expectErr: true,
			errFields: []string{"fileName"},
		},
		{
			name: "missing fileUrl",
			request: &api.CreateMediaRequest{
				MediaType: "audio",
				FileName:  "test.mp3",
				FileUrl:   "",
				FileSize:  1024,
			},
			expectErr: true,
			errFields: []string{"fileUrl"},
		},
		{
			name: "invalid fileUrl",
			request: &api.CreateMediaRequest{
				MediaType: "audio",
				FileName:  "test.mp3",
				FileUrl:   "not-a-url",
				FileSize:  1024,
			},
			expectErr: true,
			errFields: []string{"fileUrl"},
		},
		{
			name: "missing fileSize",
			request: &api.CreateMediaRequest{
				MediaType: "audio",
				FileName:  "test.mp3",
				FileUrl:   "https://example.com/test.mp3",
				FileSize:  0,
			},
			expectErr: true,
			errFields: []string{"fileSize"},
		},
		{
			name: "fileSize too small",
			request: &api.CreateMediaRequest{
				MediaType: "audio",
				FileName:  "test.mp3",
				FileUrl:   "https://example.com/test.mp3",
				FileSize:  0,
			},
			expectErr: true,
			errFields: []string{"fileSize"},
		},
		{
			name: "invalid status",
			request: func() *api.CreateMediaRequest {
				status := "invalid"
				return &api.CreateMediaRequest{
					MediaType: "audio",
					FileName:  "test.mp3",
					FileUrl:   "https://example.com/test.mp3",
					FileSize:  1024,
					Status:    &status,
				}
			}(),
			expectErr: true,
			errFields: []string{"status"},
		},
		{
			name: "invalid dimensions count",
			request: &api.CreateMediaRequest{
				MediaType:  "image",
				FileName:   "test.png",
				FileUrl:    "https://example.com/test.png",
				FileSize:   1024,
				Dimensions: []int32{100},
			},
			expectErr: true,
			errFields: []string{"dimensions"},
		},
		{
			name: "negative duration",
			request: func() *api.CreateMediaRequest {
				duration := float32(-5)
				return &api.CreateMediaRequest{
					MediaType: "audio",
					FileName:  "test.mp3",
					FileUrl:   "https://example.com/test.mp3",
					FileSize:  1024,
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

// Benchmark tests
func BenchmarkValidateCreateMediaRequest(b *testing.B) {
	contentType := "audio/mpeg"
	req := &api.CreateMediaRequest{
		MediaType:   "audio",
		FileName:    "test.mp3",
		FileUrl:     "https://example.com/test.mp3",
		ContentType: &contentType,
		FileSize:    1024,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateCreateMediaRequest(req)
	}
}
