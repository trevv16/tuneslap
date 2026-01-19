package audio

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"tuneslap/mempeg/ffmpeg"
	"tuneslap/models"
	"tuneslap/services/storage"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ProcessAudio(media models.Media, params models.AudioProcessingParams) (models.Media, error) {
	// Step 1: Initialize user uploads bucket client
	userUploadsClient, err := storage.GetUserUploadsStorage()
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to get user uploads client: %w", err)
	}

	// Step 2: Download audio
	inputFilePath, err := DownloadAudio(userUploadsClient, media)
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to download audio: %w", err)
	}
	// Clean up downloaded file when done
	defer os.Remove(inputFilePath)

	// Step 3: Process audio using FFmpeg module
	processedOutputFileBasePath := filepath.Join(os.TempDir(), primitive.NewObjectID().Hex())
	webmPath, err := ffmpeg.NormalizeAudio(context.Background(), inputFilePath, processedOutputFileBasePath)
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to process audio: %w", err)
	}
	// Clean up processed file when done
	defer os.Remove(webmPath)

	// Step 4: Extract metadata
	audioDuration, err := ffmpeg.GetAudioMetadata(context.Background(), webmPath)
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to extract metadata: %w", err)
	}
	webmFileSize, err := GetFileSize(webmPath)
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to get file size: %w", err)
	}

	// Step 5: Upload processed audio to storage
	processedFileKey := storage.GetMediaKey(media.AuthorId.Hex(), media.MediaType, media.FileName)

	// Step 6: Upload processed audio to media bucket
	mediaClient, err := storage.GetMediaStorage()
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to get media client: %w", err)
	}

	err = mediaClient.UploadFile(context.Background(), storage.UploadFileRequest{
		ObjectName:  processedFileKey,
		FilePath:    webmPath,
		ContentType: params.ContentType,
	})
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to upload processed audio: %w", err)
	}

	// Step 7: Get the uploaded file url
	processedUrl := mediaClient.GetFileURL(processedFileKey)

	// Step 8: Delete original file from user uploads bucket
	originalFileKey := storage.GetMediaKey(media.AuthorId.Hex(), media.MediaType, media.FileName)
	err = userUploadsClient.DeleteFile(context.Background(), originalFileKey)
	if err != nil {
		// Log but don't fail - the processed file was uploaded successfully
		fmt.Printf("[Audio Process] Warning: failed to delete original file %s from user uploads bucket: %v\n", originalFileKey, err)
	} else {
		fmt.Printf("[Audio Process] Deleted original file %s from user uploads bucket\n", originalFileKey)
	}

	// Step 9: Update media object
	media.ProcessedUrl = processedUrl
	media.Status = models.ProcessingStatusDone
	media.FileSize = webmFileSize
	media.ContentType = params.ContentType
	media.Duration = audioDuration
	media.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	return media, nil
}

func GetFileSize(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil // size in bytes
}
