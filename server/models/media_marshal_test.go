package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestProcessingActivityMarshalJSON(t *testing.T) {
	t.Run("valid dates", func(t *testing.T) {
		now := time.Now()
		activity := ProcessingActivity{
			Status:    ProcessingStatusPending,
			Message:   "Test message",
			CreatedAt: primitive.NewDateTimeFromTime(now),
			UpdatedAt: primitive.NewDateTimeFromTime(now),
		}

		data, err := json.Marshal(activity)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		assert.Equal(t, ProcessingStatusPending, result["status"])
		assert.Equal(t, "Test message", result["message"])
		assert.NotEmpty(t, result["createdAt"])
		assert.NotEmpty(t, result["updatedAt"])
	})

	t.Run("invalid date year too low", func(t *testing.T) {
		// Year 1800 is less than 1900, should be replaced with current time
		invalidTime := time.Date(1800, 1, 1, 0, 0, 0, 0, time.UTC)
		activity := ProcessingActivity{
			Status:    ProcessingStatusDone,
			Message:   "Invalid date",
			CreatedAt: primitive.NewDateTimeFromTime(invalidTime),
			UpdatedAt: primitive.NewDateTimeFromTime(invalidTime),
		}

		data, err := json.Marshal(activity)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		// The dates should be replaced with current time (year >= 1900)
		createdAtStr := result["createdAt"].(string)
		parsedTime, err := time.Parse(time.RFC3339Nano, createdAtStr)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, parsedTime.Year(), 1900)
	})

	t.Run("invalid date year too high", func(t *testing.T) {
		// Year 10000 is greater than 9999, should be replaced with current time
		invalidTime := time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)
		activity := ProcessingActivity{
			Status:    ProcessingStatusError,
			Message:   "Future date",
			CreatedAt: primitive.NewDateTimeFromTime(invalidTime),
			UpdatedAt: primitive.NewDateTimeFromTime(invalidTime),
		}

		data, err := json.Marshal(activity)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		// The dates should be replaced with current time (year <= 9999)
		createdAtStr := result["createdAt"].(string)
		parsedTime, err := time.Parse(time.RFC3339Nano, createdAtStr)
		require.NoError(t, err)
		assert.LessOrEqual(t, parsedTime.Year(), 9999)
		assert.GreaterOrEqual(t, parsedTime.Year(), 1900)
	})

	t.Run("all status values", func(t *testing.T) {
		statuses := []string{
			ProcessingStatusPending,
			ProcessingStatusProcessing,
			ProcessingStatusDone,
			ProcessingStatusError,
		}

		now := time.Now()
		for _, status := range statuses {
			activity := ProcessingActivity{
				Status:    status,
				Message:   "Test",
				CreatedAt: primitive.NewDateTimeFromTime(now),
				UpdatedAt: primitive.NewDateTimeFromTime(now),
			}

			data, err := json.Marshal(activity)
			require.NoError(t, err)

			var result map[string]interface{}
			err = json.Unmarshal(data, &result)
			require.NoError(t, err)

			assert.Equal(t, status, result["status"])
		}
	})
}

func TestMediaMarshalJSON(t *testing.T) {
	validTime := time.Now()
	validObjectID := primitive.NewObjectID()

	t.Run("valid media with valid dates", func(t *testing.T) {
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

		data, err := json.Marshal(media)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		assert.Equal(t, "audio", result["mediaType"])
		assert.Equal(t, "test.mp3", result["fileName"])
		assert.Equal(t, float64(1024), result["fileSize"])

		activities := result["processingActivity"].([]interface{})
		assert.Len(t, activities, 1)
	})

	t.Run("media with invalid dates replaced", func(t *testing.T) {
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

		data, err := json.Marshal(media)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		// Dates should be replaced with current time
		createdAtStr := result["createdAt"].(string)
		parsedTime, err := time.Parse(time.RFC3339Nano, createdAtStr)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, parsedTime.Year(), 1900)
	})

	t.Run("media with empty activities gets default", func(t *testing.T) {
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
			ProcessingActivity: []ProcessingActivity{}, // Empty
		}

		data, err := json.Marshal(media)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		// Should have a default activity
		activities := result["processingActivity"].([]interface{})
		assert.Len(t, activities, 1)
		firstActivity := activities[0].(map[string]interface{})
		assert.Equal(t, ProcessingStatusPending, firstActivity["status"])
		assert.Equal(t, "Queued for processing", firstActivity["message"])
	})

	t.Run("media filters invalid activities", func(t *testing.T) {
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
					CreatedAt: primitive.NewDateTimeFromTime(invalidTime), // Invalid
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

		data, err := json.Marshal(media)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		// Only the valid activity should remain
		activities := result["processingActivity"].([]interface{})
		assert.Len(t, activities, 1)
		firstActivity := activities[0].(map[string]interface{})
		assert.Equal(t, "Valid", firstActivity["message"])
	})

	t.Run("media with all invalid activities gets default", func(t *testing.T) {
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
			Status:       ProcessingStatusPending,
			CreatedAt:    primitive.NewDateTimeFromTime(validTime),
			UpdatedAt:    primitive.NewDateTimeFromTime(validTime),
			ProcessingActivity: []ProcessingActivity{
				{
					Status:    ProcessingStatusPending,
					Message:   "Invalid 1",
					CreatedAt: primitive.NewDateTimeFromTime(invalidTime),
					UpdatedAt: primitive.NewDateTimeFromTime(invalidTime),
				},
				{
					Status:    ProcessingStatusProcessing,
					Message:   "Invalid 2",
					CreatedAt: primitive.NewDateTimeFromTime(invalidTime),
					UpdatedAt: primitive.NewDateTimeFromTime(invalidTime),
				},
			},
		}

		data, err := json.Marshal(media)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		// Should have a default activity since all were invalid
		activities := result["processingActivity"].([]interface{})
		assert.Len(t, activities, 1)
		firstActivity := activities[0].(map[string]interface{})
		assert.Equal(t, ProcessingStatusPending, firstActivity["status"])
		assert.Equal(t, "Queued for processing", firstActivity["message"])
	})

	t.Run("media with dimensions and duration", func(t *testing.T) {
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
			Duration:     0,
			CreatedAt:    primitive.NewDateTimeFromTime(validTime),
			UpdatedAt:    primitive.NewDateTimeFromTime(validTime),
		}

		data, err := json.Marshal(media)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		dims := result["dimensions"].([]interface{})
		assert.Equal(t, float64(1920), dims[0])
		assert.Equal(t, float64(1080), dims[1])
	})
}

// Benchmarks
func BenchmarkMediaMarshalJSON(b *testing.B) {
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
		_, _ = json.Marshal(media)
	}
}
