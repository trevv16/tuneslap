package storage

import (
	"context"
	"time"
)

// UploadFileRequest contains the parameters for uploading a file
type UploadFileRequest struct {
	ObjectName  string
	FilePath    string
	ContentType string
	Metadata    map[string]string
}

// ObjectStorage defines the interface for object storage operations
type ObjectStorage interface {
	// UploadFile uploads a file from the local filesystem to object storage
	UploadFile(ctx context.Context, req UploadFileRequest) error

	// DownloadFile downloads a file from object storage to the local filesystem
	DownloadFile(ctx context.Context, objectName, destPath string) error

	// DeleteFile deletes a file from object storage
	DeleteFile(ctx context.Context, objectName string) error

	// GenerateSignedUploadURL generates a signed URL for direct upload to object storage
	GenerateSignedUploadURL(ctx context.Context, objectName, contentType string, expires time.Duration) (string, error)

	// GenerateSignedDownloadURL generates a signed URL for downloading from object storage
	GenerateSignedDownloadURL(ctx context.Context, objectName string, expires time.Duration) (string, error)

	// GetBucketName returns the bucket name for this storage instance
	GetBucketName() string

	// GetFileURL returns the public URL for a file (if applicable)
	GetFileURL(objectName string) string
}
