package storage

import (
	"context"
	"fmt"
	"os"
	"tuneslap/models"
)

// ProviderType represents the storage provider type
type ProviderType string

const (
	ProviderS3  ProviderType = "s3"
	ProviderGCS ProviderType = "gcs"
)

var (
	userUploadsStorage ObjectStorage
	mediaStorage       ObjectStorage
)

// GetUserUploadsStorage returns the storage client for user uploads
func GetUserUploadsStorage() (ObjectStorage, error) {
	if userUploadsStorage != nil {
		return userUploadsStorage, nil
	}

	bucketName := os.Getenv("USER_UPLOADS_BUCKET")
	if bucketName == "" {
		return nil, fmt.Errorf("USER_UPLOADS_BUCKET environment variable is required")
	}

	storage, err := createStorage(context.Background(), bucketName)
	if err != nil {
		return nil, err
	}

	userUploadsStorage = storage
	return userUploadsStorage, nil
}

// GetMediaStorage returns the storage client for media
func GetMediaStorage() (ObjectStorage, error) {
	if mediaStorage != nil {
		return mediaStorage, nil
	}

	bucketName := os.Getenv("MEDIA_BUCKET")
	if bucketName == "" {
		return nil, fmt.Errorf("MEDIA_BUCKET environment variable is required")
	}

	storage, err := createStorage(context.Background(), bucketName)
	if err != nil {
		return nil, err
	}

	mediaStorage = storage
	return mediaStorage, nil
}

// createStorage creates a storage instance based on the STORAGE_PROVIDER environment variable
func createStorage(ctx context.Context, bucketName string) (ObjectStorage, error) {
	provider := ProviderType(os.Getenv("STORAGE_PROVIDER"))
	if provider == "" {
		provider = ProviderS3 // Default to S3 for self-hosted
	}

	switch provider {
	case ProviderS3:
		endpoint := os.Getenv("S3_ENDPOINT")
		region := os.Getenv("S3_REGION")
		if region == "" {
			region = "us-east-1"
		}
		accessKey := os.Getenv("S3_ACCESS_KEY")
		secretKey := os.Getenv("S3_SECRET_KEY")

		if accessKey == "" || secretKey == "" {
			return nil, fmt.Errorf("S3_ACCESS_KEY and S3_SECRET_KEY environment variables are required for S3 storage")
		}

		return NewS3Storage(ctx, bucketName, endpoint, region, accessKey, secretKey)

	case ProviderGCS:
		credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		return NewGCSStorage(ctx, bucketName, credentialsFile)

	default:
		return nil, fmt.Errorf("unsupported storage provider: %s (supported: s3, gcs)", provider)
	}
}

// GetMediaKey generates a storage key for a media file
func GetMediaKey(authorId, mediaType, fileName string) string {
	return fmt.Sprintf("%s/%s/%s", authorId, mediaType, fileName)
}

// DeleteMedia deletes a media file from storage
func DeleteMedia(media models.Media) error {
	mediaStorage, err := GetMediaStorage()
	if err != nil {
		return fmt.Errorf("failed to get media storage: %w", err)
	}

	mediaKey := GetMediaKey(media.AuthorId.Hex(), media.MediaType, media.FileName)
	err = mediaStorage.DeleteFile(context.Background(), mediaKey)
	if err != nil {
		return fmt.Errorf("failed to delete media file: %w", err)
	}

	return nil
}
