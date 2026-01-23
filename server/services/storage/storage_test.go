package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMediaKey(t *testing.T) {
	tests := []struct {
		name      string
		authorId  string
		mediaType string
		fileName  string
		expected  string
	}{
		{
			name:      "audio file",
			authorId:  "abc123",
			mediaType: "audio",
			fileName:  "test.mp3",
			expected:  "abc123/audio/test.mp3",
		},
		{
			name:      "image file",
			authorId:  "xyz789",
			mediaType: "image",
			fileName:  "photo.png",
			expected:  "xyz789/image/photo.png",
		},
		{
			name:      "file with spaces",
			authorId:  "user1",
			mediaType: "audio",
			fileName:  "my file.mp3",
			expected:  "user1/audio/my file.mp3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMediaKey(tt.authorId, tt.mediaType, tt.fileName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProviderType(t *testing.T) {
	assert.Equal(t, ProviderType("s3"), ProviderS3)
	assert.Equal(t, ProviderType("gcs"), ProviderGCS)
}

// Note: Actual storage tests require either:
// 1. A running S3-compatible service (like MinIO)
// 2. GCS credentials and a test bucket
// These should be run as integration tests with appropriate build tags
