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

// DirPermissions is the default permission for created directories (rwxr-xr-x)
const DirPermissions = 0755

func ProcessImage(media models.Media, params models.ImageProcessingParams) (models.Media, error) {
	// Step 1: Initialize user uploads bucket client
	userUploadsClient, err := storage.GetUserUploadsStorage()
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to get user uploads client: %w", err)
	}

	// get original file upload key
	originalFileUploadKey := storage.GetMediaKey(media.AuthorId.Hex(), media.MediaType, media.FileName)

	// Step 2: Create destination path and ensure directory exists
	downloadedFilePath := filepath.Join(os.TempDir(), originalFileUploadKey)
	downloadDir := filepath.Dir(downloadedFilePath)
	if err := os.MkdirAll(downloadDir, DirPermissions); err != nil {
		return models.Media{}, fmt.Errorf("failed to create download directory: %w", err)
	}

	// Step 3: Download original image from user uploads bucket
	err = userUploadsClient.DownloadFile(context.Background(), originalFileUploadKey, downloadedFilePath)
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to download original image: %w", err)
	}
	// Clean up downloaded file when done
	defer os.Remove(downloadedFilePath)

	// Step 4: Open the file
	file, err := os.Open(downloadedFilePath)
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Step 5: Read the file
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to read file: %w", err)
	}

	// Step 6: Normalize with default options
	// (resize, strip metadata, compress, convert)
	img, err := normalizeDefault(fileBytes)
	if err != nil {
		return models.Media{}, fmt.Errorf("normalization error: %w", err)
	}

	// Step 7: Transform (crop, rotate, blur, grayscale)
	img, err = transform(img, params)
	if err != nil {
		return models.Media{}, fmt.Errorf("transformation error: %w", err)
	}

	// Step 8: Get updated metadata
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
	// Clean up processed file when done
	defer os.Remove(processedFilePath)

	// Step 9: Get processed file upload key
	processedFileUploadKey := storage.GetMediaKey(media.AuthorId.Hex(), media.MediaType, media.FileName)

	// Step 10: Upload processed image to media bucket
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

	// Step 11: Get the uploaded file url
	processedUrl := mediaClient.GetFileURL(processedFileUploadKey)

	// Step 12: Delete original file from user uploads bucket
	err = userUploadsClient.DeleteFile(context.Background(), originalFileUploadKey)
	if err != nil {
		// Log but don't fail - the processed file was uploaded successfully
		fmt.Printf("[Image Process] Warning: failed to delete original file %s from user uploads bucket: %v\n", originalFileUploadKey, err)
	} else {
		fmt.Printf("[Image Process] Deleted original file %s from user uploads bucket\n", originalFileUploadKey)
	}

	// Step 13: Build updated media object
	media.ProcessedUrl = processedUrl
	media.Status = models.ProcessingStatusDone
	media.FileSize = metadata.FileSize
	media.ContentType = metadata.ContentType
	media.Dimensions = metadata.Dimensions

	return media, nil
}
