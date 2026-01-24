package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMediaToJSONResponse(t *testing.T) {
	validTime := time.Now()
	validObjectID := primitive.NewObjectID()

	t.Run("converts valid media to response", func(t *testing.T) {
		media := Media{
			ID:           validObjectID,
			AuthorId:     validObjectID,
			MediaType:    "audio",
			FileName:     "test.mp3",
			Description:  "Test description",
			FileUrl:      "https://example.com/file.mp3",
			ProcessedUrl: "https://example.com/processed.mp3",
			WaveformUrl:  "https://example.com/waveform.png",
			ContentType:  "audio/mpeg",
			FileSize:     1024,
			Status:       ProcessingStatusDone,
			Duration:     120.5,
			CreatedAt:    primitive.NewDateTimeFromTime(validTime),
			UpdatedAt:    primitive.NewDateTimeFromTime(validTime),
			ProcessingActivity: []ProcessingActivity{
				{
					Status:    ProcessingStatusDone,
					Message:   "Completed",
					CreatedAt: primitive.NewDateTimeFromTime(validTime),
					UpdatedAt: primitive.NewDateTimeFromTime(validTime),
				},
			},
		}

		response := media.ToJSONResponse()

		assert.Equal(t, validObjectID, response.ID)
		assert.Equal(t, "audio", response.MediaType)
		assert.Equal(t, "test.mp3", response.FileName)
		assert.Equal(t, "Test description", response.Description)
		assert.Equal(t, int64(1024), response.FileSize)
		assert.Equal(t, 120.5, response.Duration)
		assert.Len(t, response.ProcessingActivity, 1)
	})

	t.Run("replaces invalid dates with current time", func(t *testing.T) {
		invalidTime := time.Date(1800, 1, 1, 0, 0, 0, 0, time.UTC)
		media := Media{
			ID:           validObjectID,
			AuthorId:     validObjectID,
			MediaType:    "image",
			FileName:     "test.jpg",
			FileUrl:      "https://example.com/file.jpg",
			ProcessedUrl: "https://example.com/processed.jpg",
			WaveformUrl:  "https://example.com/waveform.png",
			ContentType:  "image/jpeg",
			FileSize:     2048,
			Status:       ProcessingStatusPending,
			CreatedAt:    primitive.NewDateTimeFromTime(invalidTime),
			UpdatedAt:    primitive.NewDateTimeFromTime(invalidTime),
		}

		response := media.ToJSONResponse()

		// Dates should be replaced with current time
		assert.GreaterOrEqual(t, response.CreatedAt.Year(), 1900)
		assert.GreaterOrEqual(t, response.UpdatedAt.Year(), 1900)
	})

	t.Run("creates default activity when empty", func(t *testing.T) {
		media := Media{
			ID:                 validObjectID,
			AuthorId:           validObjectID,
			MediaType:          "audio",
			FileName:           "test.mp3",
			FileUrl:            "https://example.com/file.mp3",
			ProcessedUrl:       "https://example.com/processed.mp3",
			WaveformUrl:        "https://example.com/waveform.png",
			ContentType:        "audio/mpeg",
			FileSize:           1024,
			Status:             ProcessingStatusPending,
			CreatedAt:          primitive.NewDateTimeFromTime(validTime),
			UpdatedAt:          primitive.NewDateTimeFromTime(validTime),
			ProcessingActivity: []ProcessingActivity{},
		}

		response := media.ToJSONResponse()

		assert.Len(t, response.ProcessingActivity, 1)
		assert.Equal(t, ProcessingStatusPending, response.ProcessingActivity[0].Status)
		assert.Equal(t, "Queued for processing", response.ProcessingActivity[0].Message)
	})

	t.Run("filters invalid activities", func(t *testing.T) {
		invalidTime := time.Date(1800, 1, 1, 0, 0, 0, 0, time.UTC)
		media := Media{
			ID:           validObjectID,
			AuthorId:     validObjectID,
			MediaType:    "audio",
			FileName:     "test.mp3",
			FileUrl:      "https://example.com/file.mp3",
			ProcessedUrl: "https://example.com/processed.mp3",
			WaveformUrl:  "https://example.com/waveform.png",
			ContentType:  "audio/mpeg",
			FileSize:     1024,
			Status:       ProcessingStatusDone,
			CreatedAt:    primitive.NewDateTimeFromTime(validTime),
			UpdatedAt:    primitive.NewDateTimeFromTime(validTime),
			ProcessingActivity: []ProcessingActivity{
				{
					Status:    ProcessingStatusPending,
					Message:   "Invalid",
					CreatedAt: primitive.NewDateTimeFromTime(invalidTime),
					UpdatedAt: primitive.NewDateTimeFromTime(validTime),
				},
				{
					Status:    ProcessingStatusDone,
					Message:   "Valid",
					CreatedAt: primitive.NewDateTimeFromTime(validTime),
					UpdatedAt: primitive.NewDateTimeFromTime(validTime),
				},
			},
		}

		response := media.ToJSONResponse()

		assert.Len(t, response.ProcessingActivity, 1)
		assert.Equal(t, "Valid", response.ProcessingActivity[0].Message)
	})

	t.Run("preserves dimensions", func(t *testing.T) {
		media := Media{
			ID:           validObjectID,
			AuthorId:     validObjectID,
			MediaType:    "image",
			FileName:     "test.jpg",
			FileUrl:      "https://example.com/file.jpg",
			ProcessedUrl: "https://example.com/processed.jpg",
			WaveformUrl:  "https://example.com/waveform.png",
			ContentType:  "image/jpeg",
			FileSize:     2048,
			Status:       ProcessingStatusDone,
			Dimensions:   [2]int{1920, 1080},
			CreatedAt:    primitive.NewDateTimeFromTime(validTime),
			UpdatedAt:    primitive.NewDateTimeFromTime(validTime),
		}

		response := media.ToJSONResponse()

		assert.Equal(t, [2]int{1920, 1080}, response.Dimensions)
	})
}

func TestMediaStatsStruct(t *testing.T) {
	stats := MediaStats{
		ImageCount:       5,
		AudioCount:       10,
		UsedStorage:      1024000,
		AvailableStorage: 9999000,
	}

	data, err := json.Marshal(stats)
	require.NoError(t, err)

	var result MediaStats
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, 5, result.ImageCount)
	assert.Equal(t, 10, result.AudioCount)
	assert.Equal(t, int64(1024000), result.UsedStorage)
	assert.Equal(t, int64(9999000), result.AvailableStorage)
}

func TestProcessingStatusConstants(t *testing.T) {
	assert.Equal(t, "pending", ProcessingStatusPending)
	assert.Equal(t, "processing", ProcessingStatusProcessing)
	assert.Equal(t, "done", ProcessingStatusDone)
	assert.Equal(t, "error", ProcessingStatusError)
}

func BenchmarkMediaToJSONResponse(b *testing.B) {
	validTime := time.Now()
	validObjectID := primitive.NewObjectID()
	media := Media{
		ID:           validObjectID,
		AuthorId:     validObjectID,
		MediaType:    "audio",
		FileName:     "test.mp3",
		Description:  "Test description",
		FileUrl:      "https://example.com/file.mp3",
		ProcessedUrl: "https://example.com/processed.mp3",
		WaveformUrl:  "https://example.com/waveform.png",
		ContentType:  "audio/mpeg",
		FileSize:     1024,
		Status:       ProcessingStatusDone,
		Duration:     120.5,
		CreatedAt:    primitive.NewDateTimeFromTime(validTime),
		UpdatedAt:    primitive.NewDateTimeFromTime(validTime),
		ProcessingActivity: []ProcessingActivity{
			{
				Status:    ProcessingStatusDone,
				Message:   "Completed",
				CreatedAt: primitive.NewDateTimeFromTime(validTime),
				UpdatedAt: primitive.NewDateTimeFromTime(validTime),
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = media.ToJSONResponse()
	}
}
