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

// DefaultOutputFormat is the default audio output format if none specified
const DefaultOutputFormat = "webm"

// ContentTypeMap maps format extensions to MIME types
var ContentTypeMap = map[string]string{
	"webm": "audio/webm",
	"mp3":  "audio/mpeg",
	"ogg":  "audio/ogg",
	"wav":  "audio/wav",
}

// ProcessAudio processes audio media using FFmpeg with the specified params
// Supports custom processing params including trim, fade, speed, pitch, and normalization
func ProcessAudio(media models.Media, params models.AudioProcessingParams) (models.Media, error) {
	ctx := context.Background()

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

	// Step 3: Process audio using FFmpeg module with params
	processedOutputFileBasePath := filepath.Join(os.TempDir(), primitive.NewObjectID().Hex())

	// Set default params if needed
	processParams := normalizeAudioParams(params)

	outputs, err := ffmpeg.ProcessAudioWithParams(ctx, inputFilePath, processedOutputFileBasePath, processParams)
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to process audio: %w", err)
	}

	// Clean up all processed files when done
	defer func() {
		for _, outputPath := range outputs {
			os.Remove(outputPath)
		}
	}()

	// Step 4: Get the primary output format (first in the list or webm)
	primaryFormat := getPrimaryFormat(processParams.OutputFormats)
	primaryOutputPath := outputs[primaryFormat]

	// Step 5: Extract metadata from primary output
	audioDuration, err := ffmpeg.GetAudioMetadata(ctx, primaryOutputPath)
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to extract metadata: %w", err)
	}
	primaryFileSize, err := GetFileSize(primaryOutputPath)
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to get file size: %w", err)
	}

	// Step 6: Upload processed audio to media bucket
	mediaClient, err := storage.GetMediaStorage()
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to get media client: %w", err)
	}

	// Upload all output formats
	var processedUrl string
	for format, outputPath := range outputs {
		fileKey := getProcessedFileKey(media, format)
		contentType := getContentType(format, params.ContentType)

		err = mediaClient.UploadFile(ctx, storage.UploadFileRequest{
			ObjectName:  fileKey,
			FilePath:    outputPath,
			ContentType: contentType,
		})
		if err != nil {
			return models.Media{}, fmt.Errorf("failed to upload processed audio (%s): %w", format, err)
		}

		// Set the primary format URL as the processed URL
		if format == primaryFormat {
			processedUrl = mediaClient.GetFileURL(fileKey)
		}
	}

	// Step 7: Delete original file from user uploads bucket
	originalFileKey := storage.GetMediaKey(media.AuthorId.Hex(), media.MediaType, media.FileName)
	err = userUploadsClient.DeleteFile(ctx, originalFileKey)
	if err != nil {
		// Log but don't fail - the processed file was uploaded successfully
		fmt.Printf("[Audio Process] Warning: failed to delete original file %s from user uploads bucket: %v\n", originalFileKey, err)
	} else {
		fmt.Printf("[Audio Process] Deleted original file %s from user uploads bucket\n", originalFileKey)
	}

	// Step 8: Update media object
	media.ProcessedUrl = processedUrl
	media.Status = models.ProcessingStatusDone
	media.FileSize = primaryFileSize
	media.ContentType = getContentType(primaryFormat, params.ContentType)
	media.Duration = audioDuration
	media.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	return media, nil
}

// normalizeAudioParams ensures params have sensible defaults
func normalizeAudioParams(params models.AudioProcessingParams) models.AudioProcessingParams {
	// If no output formats specified, default to webm
	if len(params.OutputFormats) == 0 {
		params.OutputFormats = []string{DefaultOutputFormat}
	}

	// If normalize not explicitly set but no other processing requested, enable it
	if !params.Normalize && params.TrimStart == 0 && params.TrimEnd == 0 &&
		params.FadeIn == 0 && params.FadeOut == 0 && params.Speed == 0 && params.Pitch == 0 {
		params.Normalize = true
	}

	return params
}

// getPrimaryFormat returns the first output format or default
func getPrimaryFormat(formats []string) string {
	if len(formats) > 0 {
		return formats[0]
	}
	return DefaultOutputFormat
}

// getProcessedFileKey generates the storage key for a processed file
func getProcessedFileKey(media models.Media, format string) string {
	// Use the base file name with the new format extension
	baseName := media.FileName
	// Remove existing extension if present
	ext := filepath.Ext(baseName)
	if ext != "" {
		baseName = baseName[:len(baseName)-len(ext)]
	}
	newFileName := fmt.Sprintf("%s.%s", baseName, format)
	return storage.GetMediaKey(media.AuthorId.Hex(), media.MediaType, newFileName)
}

// getContentType returns the content type for a format
func getContentType(format string, fallback string) string {
	if ct, ok := ContentTypeMap[format]; ok {
		return ct
	}
	if fallback != "" {
		return fallback
	}
	return "audio/webm"
}

func GetFileSize(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil // size in bytes
}
