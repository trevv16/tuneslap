package tasks

import (
	"context"
	"encoding/json"
	"testing"
	"tuneslap/models"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewMediaProcessTask(t *testing.T) {
	tests := []struct {
		name             string
		mediaId          primitive.ObjectID
		userId           primitive.ObjectID
		processingParams models.ProcessingParams
		expectError      bool
	}{
		{
			name:    "successful task creation with audio params",
			mediaId: primitive.NewObjectID(),
			userId:  primitive.NewObjectID(),
			processingParams: models.ProcessingParams{
				Audio: &models.AudioProcessingParams{
					ContentType: "audio/webm",
					TrimStart:   0,
					TrimEnd:     10,
				},
			},
			expectError: false,
		},
		{
			name:    "successful task creation with image params",
			mediaId: primitive.NewObjectID(),
			userId:  primitive.NewObjectID(),
			processingParams: models.ProcessingParams{
				Image: &models.ImageProcessingParams{
					ResizeTo: [2]int{800, 600},
					Format:   1, // 1 for webp format
				},
			},
			expectError: false,
		},
		{
			name:             "successful task creation with empty params",
			mediaId:          primitive.NewObjectID(),
			userId:           primitive.NewObjectID(),
			processingParams: models.ProcessingParams{},
			expectError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := NewMediaProcessTask(tt.mediaId, tt.userId, tt.processingParams)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, task)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, task)
				assert.Equal(t, TypeMediaProcess, task.Type())

				// Verify payload can be unmarshaled
				var payload MediaProcessPayload
				err := json.Unmarshal(task.Payload(), &payload)
				assert.NoError(t, err)
				assert.Equal(t, tt.mediaId, payload.MediaID)
				assert.Equal(t, tt.userId, payload.UserID)
				assert.Equal(t, tt.processingParams, payload.ProcessingParams)
			}
		})
	}
}

func TestMediaProcessPayload_Serialization(t *testing.T) {
	// Test that payload can be properly serialized and deserialized
	originalPayload := MediaProcessPayload{
		MediaID: primitive.NewObjectID(),
		UserID:  primitive.NewObjectID(),
		ProcessingParams: models.ProcessingParams{
			Audio: &models.AudioProcessingParams{
				ContentType: "audio/webm",
				TrimStart:   5,
				TrimEnd:     15,
			},
		},
	}

	// Serialize
	jsonData, err := json.Marshal(originalPayload)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Deserialize
	var deserializedPayload MediaProcessPayload
	err = json.Unmarshal(jsonData, &deserializedPayload)
	assert.NoError(t, err)

	// Verify data integrity
	assert.Equal(t, originalPayload.MediaID, deserializedPayload.MediaID)
	assert.Equal(t, originalPayload.UserID, deserializedPayload.UserID)
	assert.Equal(t, originalPayload.ProcessingParams, deserializedPayload.ProcessingParams)
}

func TestHandleMediaProcessTask_InvalidPayload(t *testing.T) {
	// Test handling of invalid payload
	invalidPayload := []byte(`{"invalid": "json"`)
	task := asynq.NewTask(TypeMediaProcess, invalidPayload)

	err := HandleMediaProcessTask(context.Background(), task)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "json.Unmarshal failed")
}

func TestRegisterMediaProcessTasks(t *testing.T) {
	// Test that task registration doesn't panic
	mux := asynq.NewServeMux()
	assert.NotPanics(t, func() {
		RegisterMediaProcessTasks(mux)
	})
}

func TestTypeMediaProcess_Constant(t *testing.T) {
	// Test that the task type constant is properly defined
	assert.Equal(t, "media:process", TypeMediaProcess)
	assert.NotEmpty(t, TypeMediaProcess)
}

// Benchmark tests
func BenchmarkNewMediaProcessTask(b *testing.B) {
	mediaId := primitive.NewObjectID()
	userId := primitive.NewObjectID()
	params := models.ProcessingParams{
		Audio: &models.AudioProcessingParams{
			ContentType: "audio/webm",
			TrimStart:   0,
			TrimEnd:     10,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewMediaProcessTask(mediaId, userId, params)
	}
}

func BenchmarkMediaProcessPayload_Serialization(b *testing.B) {
	payload := MediaProcessPayload{
		MediaID: primitive.NewObjectID(),
		UserID:  primitive.NewObjectID(),
		ProcessingParams: models.ProcessingParams{
			Audio: &models.AudioProcessingParams{
				ContentType: "audio/webm",
				TrimStart:   5,
				TrimEnd:     15,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(payload)
	}
}

func TestNewMediaProcessTask_ErrorScenarios(t *testing.T) {
	tests := []struct {
		name        string
		mediaId     primitive.ObjectID
		userId      primitive.ObjectID
		params      models.ProcessingParams
		expectError bool
		errorType   string
	}{
		{
			name:        "nil object IDs",
			mediaId:     primitive.NilObjectID,
			userId:      primitive.NilObjectID,
			params:      models.ProcessingParams{},
			expectError: false, // Should not error, just create task with nil IDs
		},
		{
			name:    "complex audio params",
			mediaId: primitive.NewObjectID(),
			userId:  primitive.NewObjectID(),
			params: models.ProcessingParams{
				Audio: &models.AudioProcessingParams{
					ContentType: "audio/mp3",
					TrimStart:   1.5,
					TrimEnd:     99.9,
					Speed:       1.2,
					Pitch:       1.1,
				},
			},
			expectError: false,
		},
		{
			name:    "complex image params",
			mediaId: primitive.NewObjectID(),
			userId:  primitive.NewObjectID(),
			params: models.ProcessingParams{
				Image: &models.ImageProcessingParams{
					ResizeTo: [2]int{1920, 1080},
					Format:   2, // 2 for png format
					Crop:     [4]int{0, 0, 1920, 1080},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := NewMediaProcessTask(tt.mediaId, tt.userId, tt.params)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, task)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, task)
				assert.Equal(t, TypeMediaProcess, task.Type())

				// Verify payload structure
				var payload MediaProcessPayload
				err = json.Unmarshal(task.Payload(), &payload)
				assert.NoError(t, err)
				assert.Equal(t, tt.mediaId, payload.MediaID)
				assert.Equal(t, tt.userId, payload.UserID)
				assert.Equal(t, tt.params, payload.ProcessingParams)
			}
		})
	}
}

func TestMediaProcessPayload_Validation(t *testing.T) {
	tests := []struct {
		name        string
		payload     MediaProcessPayload
		expectValid bool
	}{
		{
			name: "valid payload with audio params",
			payload: MediaProcessPayload{
				MediaID: primitive.NewObjectID(),
				UserID:  primitive.NewObjectID(),
				ProcessingParams: models.ProcessingParams{
					Audio: &models.AudioProcessingParams{
						ContentType: "audio/webm",
						TrimStart:   0,
						TrimEnd:     10,
					},
				},
			},
			expectValid: true,
		},
		{
			name: "valid payload with image params",
			payload: MediaProcessPayload{
				MediaID: primitive.NewObjectID(),
				UserID:  primitive.NewObjectID(),
				ProcessingParams: models.ProcessingParams{
					Image: &models.ImageProcessingParams{
						ResizeTo: [2]int{800, 600},
						Format:   1,
					},
				},
			},
			expectValid: true,
		},
		{
			name: "payload with nil object IDs",
			payload: MediaProcessPayload{
				MediaID:          primitive.NilObjectID,
				UserID:           primitive.NilObjectID,
				ProcessingParams: models.ProcessingParams{},
			},
			expectValid: true, // Should be valid, just with nil IDs
		},
		{
			name: "payload with empty processing params",
			payload: MediaProcessPayload{
				MediaID:          primitive.NewObjectID(),
				UserID:           primitive.NewObjectID(),
				ProcessingParams: models.ProcessingParams{},
			},
			expectValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test serialization
			jsonData, err := json.Marshal(tt.payload)
			if tt.expectValid {
				assert.NoError(t, err)
				assert.NotEmpty(t, jsonData)
			} else {
				assert.Error(t, err)
				return
			}

			// Test deserialization
			var deserializedPayload MediaProcessPayload
			err = json.Unmarshal(jsonData, &deserializedPayload)
			assert.NoError(t, err)

			// Verify data integrity
			assert.Equal(t, tt.payload.MediaID, deserializedPayload.MediaID)
			assert.Equal(t, tt.payload.UserID, deserializedPayload.UserID)
			assert.Equal(t, tt.payload.ProcessingParams, deserializedPayload.ProcessingParams)
		})
	}
}

func TestTaskConstants(t *testing.T) {
	// Test that constants are properly defined
	assert.Equal(t, "media:process", TypeMediaProcess)
	assert.NotEmpty(t, TypeMediaProcess)

	// Test that constant is used consistently
	task, err := NewMediaProcessTask(primitive.NewObjectID(), primitive.NewObjectID(), models.ProcessingParams{})
	assert.NoError(t, err)
	assert.Equal(t, TypeMediaProcess, task.Type())
}

func TestRegisterMediaProcessTasks_NoPanic(t *testing.T) {
	// Test that registration doesn't panic on a fresh mux
	t.Run("registers on fresh mux", func(t *testing.T) {
		mux := asynq.NewServeMux()
		assert.NotPanics(t, func() {
			RegisterMediaProcessTasks(mux)
		})
	})

	// Test multiple registrations on different mux instances
	t.Run("registers on multiple mux instances", func(t *testing.T) {
		assert.NotPanics(t, func() {
			mux1 := asynq.NewServeMux()
			mux2 := asynq.NewServeMux()
			RegisterMediaProcessTasks(mux1)
			RegisterMediaProcessTasks(mux2)
		})
	})
}
