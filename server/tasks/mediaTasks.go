package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"tuneslap/mempeg/audio"
	"tuneslap/mempeg/image"
	"tuneslap/models"
	"tuneslap/repositories"

	"github.com/hibiken/asynq"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TypeMediaProcess = "media:process"
)

type MediaProcessPayload struct {
	MediaID          primitive.ObjectID      `json:"mediaId"`
	UserID           primitive.ObjectID      `json:"userId"`
	ProcessingParams models.ProcessingParams `json:"processingParams"`
}

func NewMediaProcessTask(mediaId primitive.ObjectID, userId primitive.ObjectID, processingParams models.ProcessingParams) (*asynq.Task, error) {
	payload, err := json.Marshal(MediaProcessPayload{
		MediaID:          mediaId,
		UserID:           userId,
		ProcessingParams: processingParams,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeMediaProcess, payload), nil
}

// function used by asynq
func HandleMediaProcessTask(ctx context.Context, task *asynq.Task) error {
	log.Printf("[Task] Received media:process task")

	var payload MediaProcessPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		log.Printf("[Task] Failed to unmarshal payload: %v", err)
		return fmt.Errorf("json.Unmarshal failed: %v", err)
	}

	log.Printf("[Task] Processing media ID: %s, User ID: %s", payload.MediaID.Hex(), payload.UserID.Hex())

	// get the media
	mediaRepo := repositories.NewMediaRepository()
	media, err := mediaRepo.GetByIdUnscoped(payload.MediaID)
	if err != nil {
		log.Printf("[Task] Failed to get media: %v", err)
		return fmt.Errorf("failed to get media: %v", err)
	}

	log.Printf("[Task] Found media: %s (type: %s)", media.FileName, media.MediaType)

	// update the media status to processing
	_, err = mediaRepo.UpdateMediaUnscoped(payload.MediaID, &models.Media{
		Status: models.ProcessingStatusProcessing,
		ProcessingActivity: append(media.ProcessingActivity, models.ProcessingActivity{
			Status:    models.ProcessingStatusProcessing,
			Message:   "Processing started",
			CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
			UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
		}),
	})
	if err != nil {
		return fmt.Errorf("failed to update media before processing: %v", err)
	}

	updateData := models.Media{
		Status: models.ProcessingStatusDone,
		ProcessingActivity: append(media.ProcessingActivity, models.ProcessingActivity{
			Status:    models.ProcessingStatusDone,
			Message:   "Processing completed",
			CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
			UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
		}),
	}

	// process the media
	switch media.MediaType {
	case "audio":
		// process the audio
		processedAudio, err := audio.ProcessAudio(media, *payload.ProcessingParams.Audio)
		if err != nil {
			return fmt.Errorf("failed to process audio: %v", err)
		}
		updateData.Status = models.ProcessingStatusDone
		updateData.FileSize = processedAudio.FileSize
		updateData.ContentType = processedAudio.ContentType
		updateData.ProcessedUrl = processedAudio.ProcessedUrl
		updateData.WaveformUrl = processedAudio.WaveformUrl
		updateData.Duration = processedAudio.Duration
	case "image":
		// process the image
		processedImage, err := image.ProcessImage(media, *payload.ProcessingParams.Image)
		if err != nil {
			return fmt.Errorf("failed to process image: %v", err)
		}
		updateData.Status = models.ProcessingStatusDone
		updateData.FileSize = processedImage.FileSize
		updateData.ContentType = processedImage.ContentType
		updateData.ProcessedUrl = processedImage.ProcessedUrl
		updateData.Dimensions = processedImage.Dimensions
	default:
		return fmt.Errorf("unsupported media type: %s", media.MediaType)
	}

	// update the media metadata
	_, err = mediaRepo.UpdateMediaUnscoped(payload.MediaID, &updateData)
	if err != nil {
		log.Printf("[Task] Failed to update media after processing: %v", err)
		return fmt.Errorf("failed to update media after processing: %v", err)
	}

	log.Printf("[Task] Successfully processed media ID: %s", payload.MediaID.Hex())
	return nil
}

// RegisterMediaProcessTasks registers media processing task handlers to the provided mux
func RegisterMediaProcessTasks(mux *asynq.ServeMux) {
	log.Println("[Worker] Registering media:process task handler")
	mux.HandleFunc(TypeMediaProcess, HandleMediaProcessTask)
}
