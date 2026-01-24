package validation

import (
	"testing"

	api "tuneslap/generated"
	"tuneslap/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
