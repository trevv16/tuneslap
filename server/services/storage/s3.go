package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Storage implements ObjectStorage using AWS S3-compatible storage (S3, MinIO, Cloudflare R2, etc.)
type S3Storage struct {
	client           *s3.Client
	bucketName       string
	endpoint         string
	externalEndpoint string // External endpoint for presigned URLs (for browser access)
	region           string
	accessKey        string
	secretKey        string
}

// NewS3Storage creates a new S3Storage instance
func NewS3Storage(ctx context.Context, bucketName, endpoint, region, accessKey, secretKey string) (*S3Storage, error) {
	// Build AWS config
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// If endpoint is provided, use it for S3-compatible services (MinIO, etc.)
	if endpoint != "" {
		cfg.BaseEndpoint = aws.String(endpoint)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.UsePathStyle = true // Required for MinIO
		}
	})

	// Get external endpoint for presigned URLs (for browser access)
	// Defaults to endpoint if not set, but can be overridden for Docker internal/external URL differences
	externalEndpoint := os.Getenv("S3_EXTERNAL_ENDPOINT")
	if externalEndpoint == "" {
		externalEndpoint = endpoint
	}

	return &S3Storage{
		client:           client,
		bucketName:       bucketName,
		endpoint:         endpoint,
		externalEndpoint: externalEndpoint,
		region:           region,
		accessKey:        accessKey,
		secretKey:        secretKey,
	}, nil
}

// UploadFile uploads a file from the local filesystem to S3
func (s *S3Storage) UploadFile(ctx context.Context, req UploadFileRequest) error {
	file, err := os.Open(req.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	uploader := manager.NewUploader(s.client)
	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(req.ObjectName),
		Body:        file,
		ContentType: aws.String(req.ContentType),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

// DownloadFile downloads a file from S3 to the local filesystem
func (s *S3Storage) DownloadFile(ctx context.Context, objectName, destPath string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		return fmt.Errorf("failed to get object: %w", err)
	}
	defer result.Body.Close()

	outFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, result.Body)
	if err != nil {
		return fmt.Errorf("failed to copy object data: %w", err)
	}

	return nil
}

// DeleteFile deletes a file from S3
func (s *S3Storage) DeleteFile(ctx context.Context, objectName string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

// FileExists checks if a file exists in S3 storage
func (s *S3Storage) FileExists(ctx context.Context, objectName string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		// Check if it's a "not found" error
		var nsk *types.NoSuchKey
		var nskb *types.NotFound
		if errors.As(err, &nsk) || errors.As(err, &nskb) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if object exists: %w", err)
	}

	return true, nil
}

// GenerateSignedUploadURL generates a presigned URL for uploading to S3
func (s *S3Storage) GenerateSignedUploadURL(ctx context.Context, objectName, contentType string, expires time.Duration) (string, error) {
	// Use external endpoint for presigning if different from internal endpoint
	// This ensures the signature is valid for the external URL
	var presignClient *s3.PresignClient
	if s.externalEndpoint != "" && s.endpoint != "" && s.externalEndpoint != s.endpoint {
		// Create a separate config and client for presigning with external endpoint
		cfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(s.region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s.accessKey, s.secretKey, "")),
		)
		if err != nil {
			return "", fmt.Errorf("failed to load AWS config for presigning: %w", err)
		}

		cfg.BaseEndpoint = aws.String(s.externalEndpoint)
		externalClient := s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.UsePathStyle = true // Required for MinIO
		})
		presignClient = s3.NewPresignClient(externalClient)
	} else {
		presignClient = s3.NewPresignClient(s.client)
	}

	// Note: We don't include ContentType or ACL in the signed request because:
	// 1. AWS SDK Go v2 doesn't automatically include ContentType in signed headers
	// 2. Browsers will send Content-Type as an unsigned header, which MinIO accepts
	// 3. ACL headers (x-amz-acl) shouldn't be signed for browser uploads as browsers don't send them
	// The contentType parameter is kept for API compatibility but not used in signing
	// ACL can be set server-side after upload if needed
	request, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectName),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return request.URL, nil
}

// GenerateSignedDownloadURL generates a presigned URL for downloading from S3
func (s *S3Storage) GenerateSignedDownloadURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	// Use external endpoint for presigning if different from internal endpoint
	// This ensures the signature is valid for the external URL
	var presignClient *s3.PresignClient
	if s.externalEndpoint != "" && s.endpoint != "" && s.externalEndpoint != s.endpoint {
		// Create a separate config and client for presigning with external endpoint
		cfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(s.region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s.accessKey, s.secretKey, "")),
		)
		if err != nil {
			return "", fmt.Errorf("failed to load AWS config for presigning: %w", err)
		}

		cfg.BaseEndpoint = aws.String(s.externalEndpoint)
		externalClient := s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.UsePathStyle = true // Required for MinIO
		})
		presignClient = s3.NewPresignClient(externalClient)
	} else {
		presignClient = s3.NewPresignClient(s.client)
	}

	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectName),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return request.URL, nil
}

// GetBucketName returns the bucket name
func (s *S3Storage) GetBucketName() string {
	return s.bucketName
}

// GetFileURL returns the public URL for a file
func (s *S3Storage) GetFileURL(objectName string) string {
	endpoint := s.externalEndpoint
	if endpoint == "" {
		endpoint = s.endpoint
	}

	if endpoint != "" {
		// For MinIO or custom endpoints, construct URL manually
		return fmt.Sprintf("%s/%s/%s", endpoint, s.bucketName, objectName)
	}
	// For AWS S3, use standard S3 URL format
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucketName, objectName)
}
