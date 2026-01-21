package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// GCSStorage implements ObjectStorage using Google Cloud Storage
type GCSStorage struct {
	client              *storage.Client
	bucketName          string
	serviceAccountEmail string
	privateKey          []byte
}

// NewGCSStorage creates a new GCSStorage instance
func NewGCSStorage(ctx context.Context, bucketName, credentialsFile string) (*GCSStorage, error) {
	var client *storage.Client
	var err error
	if credentialsFile != "" {
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	} else {
		client, err = storage.NewClient(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}

	// Get service account email - try environment variable first, then key file
	serviceAccountEmail := os.Getenv("GOOGLE_SERVICE_ACCOUNT_EMAIL")

	// Get private key file path for signed URLs
	privateKeyPath := os.Getenv("GOOGLE_PRIVATE_KEY_PATH")
	if privateKeyPath == "" {
		return nil, fmt.Errorf("GOOGLE_PRIVATE_KEY_PATH environment variable is required for signed URLs")
	}

	// Read and parse the service account key file
	keyFileData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	// Parse the JSON to extract private key and client email
	var keyData struct {
		PrivateKey  string `json:"private_key"`
		ClientEmail string `json:"client_email"`
	}

	if err := json.Unmarshal(keyFileData, &keyData); err != nil {
		return nil, fmt.Errorf("failed to parse service account key file: %w", err)
	}

	if keyData.PrivateKey == "" {
		return nil, fmt.Errorf("private_key not found in service account key file")
	}

	// Use environment variable if set, otherwise use the one from key file
	if serviceAccountEmail == "" {
		if keyData.ClientEmail == "" {
			return nil, fmt.Errorf("client_email not found in service account key file and GOOGLE_SERVICE_ACCOUNT_EMAIL not set")
		}
		serviceAccountEmail = keyData.ClientEmail
	}

	return &GCSStorage{
		client:              client,
		bucketName:          bucketName,
		serviceAccountEmail: serviceAccountEmail,
		privateKey:          []byte(keyData.PrivateKey),
	}, nil
}

// UploadFile uploads a file from the local filesystem to GCS
func (g *GCSStorage) UploadFile(ctx context.Context, req UploadFileRequest) error {
	f, err := os.Open(req.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	obj := g.client.Bucket(g.bucketName).Object(req.ObjectName)
	wc := obj.NewWriter(ctx)

	if req.ContentType != "" {
		wc.ContentType = req.ContentType
	}
	if req.Metadata != nil {
		wc.Metadata = req.Metadata
	}
	wc.PredefinedACL = "publicRead"

	if _, err := io.Copy(wc, f); err != nil {
		_ = wc.Close()
		return fmt.Errorf("failed to copy file data: %w", err)
	}
	return wc.Close()
}

// DownloadFile downloads a file from GCS to the local filesystem
func (g *GCSStorage) DownloadFile(ctx context.Context, objectName, destPath string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	rc, err := g.client.Bucket(g.bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("failed to get object reader: %w", err)
	}
	defer rc.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, rc)
	if err != nil {
		return fmt.Errorf("failed to copy object data: %w", err)
	}
	return nil
}

// DeleteFile deletes a file from GCS
// FileExists checks if a file exists in GCS storage
func (g *GCSStorage) FileExists(ctx context.Context, objectName string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	obj := g.client.Bucket(g.bucketName).Object(objectName)
	_, err := obj.Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if object exists: %w", err)
	}

	return true, nil
}

func (g *GCSStorage) DeleteFile(ctx context.Context, objectName string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	if err := g.client.Bucket(g.bucketName).Object(objectName).Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

// GenerateSignedUploadURL generates a signed URL for direct upload to GCS
func (g *GCSStorage) GenerateSignedUploadURL(ctx context.Context, objectName, contentType string, expires time.Duration) (string, error) {
	opts := &storage.SignedURLOptions{
		GoogleAccessID: g.serviceAccountEmail,
		PrivateKey:     g.privateKey,
		Scheme:         storage.SigningSchemeV4,
		Method:         "PUT",
		Expires:        time.Now().Add(expires),
		Headers: []string{
			fmt.Sprintf("Content-Type: %s", contentType),
		},
	}

	url, err := storage.SignedURL(g.bucketName, objectName, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return url, nil
}

// GenerateSignedDownloadURL generates a signed URL for downloading from GCS
func (g *GCSStorage) GenerateSignedDownloadURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	opts := &storage.SignedURLOptions{
		GoogleAccessID: g.serviceAccountEmail,
		PrivateKey:     g.privateKey,
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		Expires:        time.Now().Add(expires),
	}

	url, err := storage.SignedURL(g.bucketName, objectName, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return url, nil
}

// GetBucketName returns the bucket name
func (g *GCSStorage) GetBucketName() string {
	return g.bucketName
}

// GetFileURL returns the public URL for a file
func (g *GCSStorage) GetFileURL(objectName string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.bucketName, objectName)
}
