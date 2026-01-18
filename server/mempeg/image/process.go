package image

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"tuneslap/models"
	"tuneslap/services/storage"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ProcessImage(media models.Media, params models.ImageProcessingParams) (models.Media, error) {
	// Step 1: Initialize user uploads bucket client
	userUploadsClient, err := storage.GetUserUploadsStorage()
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to get user uploads client: %w", err)
	}

	// get original file upload key
	originalFileUploadKey := storage.GetMediaKey(media.AuthorId.Hex(), media.MediaType, media.FileName)

	// Step 2: Download original image from user uploads bucket
	downloadedFilePath := filepath.Join(os.TempDir(), originalFileUploadKey)
	err = userUploadsClient.DownloadFile(context.Background(), originalFileUploadKey, downloadedFilePath)
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to download original image: %w", err)
	}

	// Step 3: Open the file
	file, err := os.Open(downloadedFilePath)
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Step 4: Read the file
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to read file: %w", err)
	}

	// Step 5: Normalize with default options
	// (resize, strip metadata, compress, convert)
	img, err := normalizeDefault(fileBytes)
	if err != nil {
		return models.Media{}, fmt.Errorf("normalization error: %w", err)
	}

	// Step 6: Transform (crop, rotate, blur, grayscale)
	img, err = transform(img, params)
	if err != nil {
		return models.Media{}, fmt.Errorf("transformation error: %w", err)
	}

	// Step 7: Get updated metadata
	metadata, err := getMetadata(img)
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to get metadata: %w", err)
	}

	// save file to temp dir
	processedFilePath := filepath.Join(os.TempDir(), primitive.NewObjectID().Hex())
	err = os.WriteFile(processedFilePath, img, 0644)
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to save file: %w", err)
	}

	// Step 8: Get processed file upload key
	processedFileUploadKey := storage.GetMediaKey(media.AuthorId.Hex(), media.MediaType, media.FileName)

	// Step 9: Upload processed image to media bucket
	mediaClient, err := storage.GetMediaStorage()
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to get media client: %w", err)
	}
	err = mediaClient.UploadFile(context.Background(), storage.UploadFileRequest{
		ObjectName:  processedFileUploadKey,
		FilePath:    processedFilePath,
		ContentType: GetContentTypeFromFormat(params.Format),
	})
	if err != nil {
		return models.Media{}, fmt.Errorf("upload error: %w", err)
	}

	// Step 7: Get the uploaded file url
	processedUrl := mediaClient.GetFileURL(processedFileUploadKey)

	// Step 8: Build updated media object
	media.ProcessedUrl = processedUrl
	media.Status = models.ProcessingStatusDone
	media.FileSize = metadata.FileSize
	media.ContentType = metadata.ContentType
	media.Dimensions = metadata.Dimensions

	return media, nil
}
