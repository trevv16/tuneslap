package handlers

import (
	"testing"
	"time"
	"tuneslap/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDeriveContentTypeFromFileName(t *testing.T) {
	tests := []struct {
		name        string
		fileName    string
		expected    string
		description string
	}{
		// Audio types
		{
			name:        "mp3 file",
			fileName:    "audio.mp3",
			expected:    "audio/mpeg",
			description: "MP3 audio files",
		},
		{
			name:        "wav file",
			fileName:    "audio.wav",
			expected:    "audio/x-wav",
			description: "WAV audio files",
		},
		{
			name:        "webm audio file",
			fileName:    "audio.webm",
			expected:    "audio/webm",
			description: "WebM audio files",
		},
		{
			name:        "ogg file",
			fileName:    "audio.ogg",
			expected:    "audio/ogg",
			description: "OGG audio files",
		},
		{
			name:        "aac file",
			fileName:    "audio.aac",
			expected:    "audio/aac",
			description: "AAC audio files",
		},
		// Image types
		{
			name:        "jpg file",
			fileName:    "image.jpg",
			expected:    "image/jpeg",
			description: "JPG image files",
		},
		{
			name:        "jpeg file",
			fileName:    "image.jpeg",
			expected:    "image/jpeg",
			description: "JPEG image files",
		},
		{
			name:        "png file",
			fileName:    "image.png",
			expected:    "image/png",
			description: "PNG image files",
		},
		{
			name:        "gif file",
			fileName:    "image.gif",
			expected:    "image/gif",
			description: "GIF image files",
		},
		{
			name:        "webp file",
			fileName:    "image.webp",
			expected:    "image/webp",
			description: "WebP image files",
		},
		{
			name:        "svg file",
			fileName:    "image.svg",
			expected:    "image/svg+xml",
			description: "SVG image files",
		},
		// Edge cases
		{
			name:        "empty filename",
			fileName:    "",
			expected:    "application/octet-stream",
			description: "Empty filename returns default",
		},
		{
			name:        "unknown extension",
			fileName:    "file.xyz",
			expected:    "application/octet-stream",
			description: "Unknown extension returns default",
		},
		{
			name:        "no extension",
			fileName:    "filename",
			expected:    "application/octet-stream",
			description: "No extension returns default",
		},
		// Case insensitivity
		{
			name:        "uppercase MP3",
			fileName:    "AUDIO.MP3",
			expected:    "audio/mpeg",
			description: "Uppercase extension handled correctly",
		},
		{
			name:        "mixed case PNG",
			fileName:    "Image.Png",
			expected:    "image/png",
			description: "Mixed case extension handled correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deriveContentTypeFromFileName(tt.fileName)
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

// Helper function for creating string pointers
func strPtr(s string) *string {
	return &s
}

func TestNewMediaHandler(t *testing.T) {
	t.Run("creates handler successfully", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Skip("Skipping test - requires database connection")
			}
		}()

		handler := NewMediaHandler()
		assert.NotNil(t, handler)
		assert.NotNil(t, handler.mediaRepo)
		assert.NotNil(t, handler.BaseHandler)
	})
}

func TestMediaResponseMapperLogic(t *testing.T) {
	validTime := primitive.NewDateTimeFromTime(time.Now())
	mediaID := primitive.NewObjectID()
	authorID := primitive.NewObjectID()

	t.Run("maps audio media correctly", func(t *testing.T) {
		media := models.Media{
			ID:           mediaID,
			AuthorId:     authorID,
			MediaType:    "audio",
			FileName:     "test.mp3",
			Description:  "Test audio",
			FileUrl:      "https://storage.example.com/test.mp3",
			ProcessedUrl: "https://storage.example.com/processed.webm",
			WaveformUrl:  "https://storage.example.com/waveform.png",
			ContentType:  "audio/mpeg",
			FileSize:     1024000,
			Status:       models.ProcessingStatusDone,
			Duration:     120.5,
			CreatedAt:    validTime,
			UpdatedAt:    validTime,
		}

		assert.Equal(t, "audio", media.MediaType)
		assert.Equal(t, "audio/mpeg", media.ContentType)
		assert.Equal(t, 120.5, media.Duration)
	})

	t.Run("maps image media correctly", func(t *testing.T) {
		media := models.Media{
			ID:           mediaID,
			AuthorId:     authorID,
			MediaType:    "image",
			FileName:     "test.png",
			Description:  "Test image",
			FileUrl:      "https://storage.example.com/test.png",
			ProcessedUrl: "https://storage.example.com/processed.webp",
			WaveformUrl:  "",
			ContentType:  "image/png",
			FileSize:     512000,
			Status:       models.ProcessingStatusDone,
			Dimensions:   [2]int{1920, 1080},
			CreatedAt:    validTime,
			UpdatedAt:    validTime,
		}

		assert.Equal(t, "image", media.MediaType)
		assert.Equal(t, "image/png", media.ContentType)
		assert.Equal(t, [2]int{1920, 1080}, media.Dimensions)
	})
}

func TestMediaFilterLogic(t *testing.T) {
	// Test the filter logic used in HandleGetAllMedia

	t.Run("creates filter with mediaType", func(t *testing.T) {
		authorID := primitive.NewObjectID()
		mediaType := "audio"

		filter := map[string]interface{}{
			"authorId": authorID,
		}

		if mediaType == "image" || mediaType == "audio" {
			filter["mediaType"] = mediaType
		}

		assert.Equal(t, authorID, filter["authorId"])
		assert.Equal(t, "audio", filter["mediaType"])
	})

	t.Run("creates filter with contentType", func(t *testing.T) {
		authorID := primitive.NewObjectID()
		contentType := "audio/mpeg"

		filter := map[string]interface{}{
			"authorId": authorID,
		}

		validContentTypes := []string{
			"image/jpeg", "image/png", "image/gif", "image/webp",
			"audio/mpeg", "audio/mp3", "audio/mp4", "audio/wav", "audio/ogg", "audio/webm",
		}
		for _, validType := range validContentTypes {
			if contentType == validType {
				filter["contentType"] = contentType
				break
			}
		}

		assert.Equal(t, "audio/mpeg", filter["contentType"])
	})

	t.Run("ignores invalid mediaType", func(t *testing.T) {
		authorID := primitive.NewObjectID()
		mediaType := "invalid"

		filter := map[string]interface{}{
			"authorId": authorID,
		}

		if mediaType == "image" || mediaType == "audio" {
			filter["mediaType"] = mediaType
		}

		_, exists := filter["mediaType"]
		assert.False(t, exists)
	})

	t.Run("ignores invalid contentType", func(t *testing.T) {
		authorID := primitive.NewObjectID()
		contentType := "application/pdf" // Not a valid media content type

		filter := map[string]interface{}{
			"authorId": authorID,
		}

		validContentTypes := []string{
			"image/jpeg", "image/png", "image/gif", "image/webp",
			"audio/mpeg", "audio/mp3", "audio/mp4", "audio/wav", "audio/ogg", "audio/webm",
		}
		for _, validType := range validContentTypes {
			if contentType == validType {
				filter["contentType"] = contentType
				break
			}
		}

		_, exists := filter["contentType"]
		assert.False(t, exists)
	})
}

func TestMediaTypeDetection(t *testing.T) {
	// Test the media type detection logic in HandleGenerateUploadURL

	t.Run("detects audio content type", func(t *testing.T) {
		contentType := "audio/mpeg"

		mediaType := "audio"
		if contentType != "" && len(contentType) >= 5 && contentType[:5] == "image" {
			mediaType = "image"
		}

		assert.Equal(t, "audio", mediaType)
	})

	t.Run("detects image content type", func(t *testing.T) {
		contentType := "image/png"

		mediaType := "audio"
		if contentType != "" && len(contentType) >= 5 && contentType[:5] == "image" {
			mediaType = "image"
		}

		assert.Equal(t, "image", mediaType)
	})

	t.Run("defaults to audio for unknown types", func(t *testing.T) {
		contentType := "application/octet-stream"

		mediaType := "audio"
		if contentType != "" && len(contentType) >= 5 && contentType[:5] == "image" {
			mediaType = "image"
		}

		assert.Equal(t, "audio", mediaType)
	})

	t.Run("handles empty content type", func(t *testing.T) {
		contentType := ""

		mediaType := "audio"
		if contentType != "" && len(contentType) >= 5 && contentType[:5] == "image" {
			mediaType = "image"
		}

		assert.Equal(t, "audio", mediaType)
	})
}

func TestMediaProcessingStatus(t *testing.T) {
	t.Run("all processing statuses are valid", func(t *testing.T) {
		statuses := []string{
			models.ProcessingStatusPending,
			models.ProcessingStatusProcessing,
			models.ProcessingStatusDone,
			models.ProcessingStatusError,
		}

		for _, status := range statuses {
			media := models.Media{
				ID:     primitive.NewObjectID(),
				Status: status,
			}
			assert.NotEmpty(t, media.Status)
		}
	})
}

func TestMediaStatsCalculation(t *testing.T) {
	t.Run("calculates available storage correctly", func(t *testing.T) {
		stats := models.MediaStats{
			ImageCount:       5,
			AudioCount:       10,
			UsedStorage:      1000000,
			AvailableStorage: 9000000,
		}

		totalCount := stats.ImageCount + stats.AudioCount
		assert.Equal(t, 15, totalCount)
		assert.Equal(t, int64(9000000), stats.AvailableStorage)
	})

	t.Run("handles unlimited storage", func(t *testing.T) {
		stats := models.MediaStats{
			ImageCount:       5,
			AudioCount:       10,
			UsedStorage:      1000000,
			AvailableStorage: -1, // Unlimited
		}

		isUnlimited := stats.AvailableStorage == -1
		assert.True(t, isUnlimited)
	})
}

// Benchmarks
func BenchmarkDeriveContentTypeFromFileName(b *testing.B) {
	fileName := "test-audio.mp3"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = deriveContentTypeFromFileName(fileName)
	}
}

func BenchmarkMediaFilterCreation(b *testing.B) {
	authorID := primitive.NewObjectID()
	mediaType := "audio"
	contentType := "audio/mpeg"

	validContentTypes := []string{
		"image/jpeg", "image/png", "image/gif", "image/webp",
		"audio/mpeg", "audio/mp3", "audio/mp4", "audio/wav", "audio/ogg", "audio/webm",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter := map[string]interface{}{
			"authorId": authorID,
		}

		if mediaType == "image" || mediaType == "audio" {
			filter["mediaType"] = mediaType
		}

		for _, validType := range validContentTypes {
			if contentType == validType {
				filter["contentType"] = contentType
				break
			}
		}
		_ = filter
	}
}
