package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
	"tuneslap/models"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// GCS wraps the Google Cloud Storage client.
type GCS struct {
	Client              *storage.Client
	BucketName          string
	ServiceAccountEmail string
	PrivateKey          []byte
}

var userUploadsClient *GCS
var mediaClient *GCS

type UploadFileRequest struct {
	ObjectName  string
	FilePath    string
	ContentType string
	Metadata    map[string]string
}

// create the gcs client for user uploads
func GetUserUploadsClient() (*GCS, error) {
	if userUploadsClient != nil {
		return userUploadsClient, nil
	}

	// create the gcs client
	var err error
	userUploadsBucket := os.Getenv("USER_UPLOADS_BUCKET")
	if userUploadsBucket == "" {
		return nil, fmt.Errorf("USER_UPLOADS_BUCKET environment variable is required")
	}
	userUploadsClient, err = NewGCSClient(context.Background(), userUploadsBucket, "")
	if err != nil {
		return nil, err
	}

	return userUploadsClient, nil
}

// create the gcs client for media api bucket
func GetMediaClient() (*GCS, error) {
	if mediaClient != nil {
		return mediaClient, nil
	}

	// create the gcs client
	var err error
	mediaBucket := os.Getenv("MEDIA_BUCKET")
	if mediaBucket == "" {
		return nil, fmt.Errorf("MEDIA_BUCKET environment variable is required")
	}
	mediaClient, err = NewGCSClient(context.Background(), mediaBucket, "")
	if err != nil {
		return nil, err
	}

	return mediaClient, nil
}

// NewGCSClient initializes and returns a GCS instance.
// If credentialsFile is empty, will use default credentials.
func NewGCSClient(ctx context.Context, bucketName, credentialsFile string) (*GCS, error) {
	var client *storage.Client
	var err error
	if credentialsFile != "" {
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	} else {
		client, err = storage.NewClient(ctx)
	}
	if err != nil {
		return nil, err
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

	return &GCS{
		Client:              client,
		BucketName:          bucketName,
		ServiceAccountEmail: serviceAccountEmail,
		PrivateKey:          []byte(keyData.PrivateKey),
	}, nil
}

// UploadFile uploads data from a local file to GCS.
func (g *GCS) UploadFile(ctx context.Context, req UploadFileRequest) error {
	f, err := os.Open(req.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	obj := g.Client.Bucket(g.BucketName).Object(req.ObjectName)
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
		return err
	}
	return wc.Close()
}

// DownloadFile downloads data from GCS to a local file.
func (g *GCS) DownloadFile(ctx context.Context, objectName, destPath string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	rc, err := g.Client.Bucket(g.BucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return err
	}
	defer rc.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, rc)
	return err
}

// DeleteFile deletes an object from GCS.
func (g *GCS) DeleteFile(ctx context.Context, objectName string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	return g.Client.Bucket(g.BucketName).Object(objectName).Delete(ctx)
}

func GetMediaKey(media models.Media) string {
	return fmt.Sprintf("%s/%s/%s", media.AuthorId.Hex(), media.MediaType, media.FileName)
}

func GetUploadedFileUrl(bucketName, objectName string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)
}

func DeleteMediaFromS3(media models.Media) error {
	mediaClient, err := GetMediaClient()
	if err != nil {
		return err
	}

	mediaKey := GetMediaKey(media)
	err = mediaClient.DeleteFile(context.Background(), mediaKey)
	if err != nil {
		return err
	}

	return nil
}

// GenerateSignedUploadURL generates a signed URL for direct upload to GCS
func (g *GCS) GenerateSignedUploadURL(ctx context.Context, objectName, contentType string, expires time.Duration) (string, error) {
	opts := &storage.SignedURLOptions{
		GoogleAccessID: g.ServiceAccountEmail,
		PrivateKey:     g.PrivateKey,
		Scheme:         storage.SigningSchemeV4,
		Method:         "PUT",
		Expires:        time.Now().Add(expires),
		Headers: []string{
			fmt.Sprintf("Content-Type: %s", contentType),
		},
	}

	url, err := storage.SignedURL(g.BucketName, objectName, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return url, nil
}

// GenerateSignedDownloadURL generates a signed URL for downloading from GCS
func (g *GCS) GenerateSignedDownloadURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	opts := &storage.SignedURLOptions{
		GoogleAccessID: g.ServiceAccountEmail,
		PrivateKey:     g.PrivateKey,
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		Expires:        time.Now().Add(expires),
	}

	url, err := storage.SignedURL(g.BucketName, objectName, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return url, nil
}
