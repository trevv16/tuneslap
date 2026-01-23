package audio

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
	"tuneslap/models"
	"tuneslap/services/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockStorageClient implements storage.ObjectStorage for testing
type MockStorageClient struct {
	DownloadFunc       func(ctx context.Context, objectName, destPath string) error
	UploadFunc         func(ctx context.Context, req storage.UploadFileRequest) error
	DeleteFunc         func(ctx context.Context, objectName string) error
	FileExistsFunc     func(ctx context.Context, objectName string) (bool, error)
	GenerateUploadURL  func(ctx context.Context, objectName, contentType string, expires time.Duration) (string, error)
	GenerateDownloadURL func(ctx context.Context, objectName string, expires time.Duration) (string, error)
	GetBucketNameFunc  func() string
	GetFileURLFunc     func(objectName string) string
}

func (m *MockStorageClient) UploadFile(ctx context.Context, req storage.UploadFileRequest) error {
	if m.UploadFunc != nil {
		return m.UploadFunc(ctx, req)
	}
	return nil
}

func (m *MockStorageClient) DownloadFile(ctx context.Context, objectName, destPath string) error {
	if m.DownloadFunc != nil {
		return m.DownloadFunc(ctx, objectName, destPath)
	}
	return nil
}

func (m *MockStorageClient) DeleteFile(ctx context.Context, objectName string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, objectName)
	}
	return nil
}

func (m *MockStorageClient) FileExists(ctx context.Context, objectName string) (bool, error) {
	if m.FileExistsFunc != nil {
		return m.FileExistsFunc(ctx, objectName)
	}
	return true, nil
}

func (m *MockStorageClient) GenerateSignedUploadURL(ctx context.Context, objectName, contentType string, expires time.Duration) (string, error) {
	if m.GenerateUploadURL != nil {
		return m.GenerateUploadURL(ctx, objectName, contentType, expires)
	}
	return "https://example.com/upload", nil
}

func (m *MockStorageClient) GenerateSignedDownloadURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	if m.GenerateDownloadURL != nil {
		return m.GenerateDownloadURL(ctx, objectName, expires)
	}
	return "https://example.com/download", nil
}

func (m *MockStorageClient) GetBucketName() string {
	if m.GetBucketNameFunc != nil {
		return m.GetBucketNameFunc()
	}
	return "test-bucket"
}

func (m *MockStorageClient) GetFileURL(objectName string) string {
	if m.GetFileURLFunc != nil {
		return m.GetFileURLFunc(objectName)
	}
	return "https://example.com/" + objectName
}

func TestDownloadAudio(t *testing.T) {
	authorID := primitive.NewObjectID()

	t.Run("successfully downloads audio file", func(t *testing.T) {
		tmpDir := t.TempDir()

		mockClient := &MockStorageClient{
			DownloadFunc: func(ctx context.Context, objectName, destPath string) error {
				// Create the file to simulate download
				err := os.MkdirAll(filepath.Dir(destPath), DirPermissions)
				if err != nil {
					return err
				}
				return os.WriteFile(destPath, []byte("audio content"), 0644)
			},
		}

		media := models.Media{
			ID:        primitive.NewObjectID(),
			AuthorId:  authorID,
			MediaType: "audio",
			FileName:  "test.mp3",
		}

		// Temporarily override temp dir for predictable paths
		originalTempDir := os.TempDir()
		err := os.Setenv("TMPDIR", tmpDir)
		require.NoError(t, err)
		defer os.Setenv("TMPDIR", originalTempDir)

		filePath, err := DownloadAudio(mockClient, media)
		require.NoError(t, err)
		assert.NotEmpty(t, filePath)

		// Verify file exists
		_, err = os.Stat(filePath)
		assert.NoError(t, err)
	})

	t.Run("returns error when download fails", func(t *testing.T) {
		mockClient := &MockStorageClient{
			DownloadFunc: func(ctx context.Context, objectName, destPath string) error {
				return errors.New("download failed")
			},
		}

		media := models.Media{
			ID:        primitive.NewObjectID(),
			AuthorId:  authorID,
			MediaType: "audio",
			FileName:  "test.mp3",
		}

		filePath, err := DownloadAudio(mockClient, media)
		assert.Error(t, err)
		assert.Empty(t, filePath)
		assert.Contains(t, err.Error(), "failed to download original audio")
	})

	t.Run("creates directory if not exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		dirCreated := false

		mockClient := &MockStorageClient{
			DownloadFunc: func(ctx context.Context, objectName, destPath string) error {
				// Check if directory was created
				dir := filepath.Dir(destPath)
				if _, err := os.Stat(dir); err == nil {
					dirCreated = true
				}
				return os.WriteFile(destPath, []byte("content"), 0644)
			},
		}

		media := models.Media{
			ID:        primitive.NewObjectID(),
			AuthorId:  authorID,
			MediaType: "audio",
			FileName:  "test.mp3",
		}

		// Temporarily override temp dir
		originalTempDir := os.TempDir()
		err := os.Setenv("TMPDIR", tmpDir)
		require.NoError(t, err)
		defer os.Setenv("TMPDIR", originalTempDir)

		_, err = DownloadAudio(mockClient, media)
		require.NoError(t, err)
		assert.True(t, dirCreated)
	})
}

func TestDownloadAudioFileKeyGeneration(t *testing.T) {
	authorID := primitive.NewObjectID()

	t.Run("generates correct file key", func(t *testing.T) {
		var capturedObjectName string

		mockClient := &MockStorageClient{
			DownloadFunc: func(ctx context.Context, objectName, destPath string) error {
				capturedObjectName = objectName
				// Create minimal file
				os.MkdirAll(filepath.Dir(destPath), DirPermissions)
				return os.WriteFile(destPath, []byte("x"), 0644)
			},
		}

		media := models.Media{
			ID:        primitive.NewObjectID(),
			AuthorId:  authorID,
			MediaType: "audio",
			FileName:  "my-song.mp3",
		}

		_, err := DownloadAudio(mockClient, media)
		require.NoError(t, err)

		// Verify the object key matches expected format
		expectedKey := storage.GetMediaKey(authorID.Hex(), "audio", "my-song.mp3")
		assert.Equal(t, expectedKey, capturedObjectName)
	})
}

// Benchmark tests
func BenchmarkDownloadAudio(b *testing.B) {
	authorID := primitive.NewObjectID()
	mockClient := &MockStorageClient{
		DownloadFunc: func(ctx context.Context, objectName, destPath string) error {
			os.MkdirAll(filepath.Dir(destPath), DirPermissions)
			return os.WriteFile(destPath, []byte("content"), 0644)
		},
	}

	media := models.Media{
		ID:        primitive.NewObjectID(),
		AuthorId:  authorID,
		MediaType: "audio",
		FileName:  "bench.mp3",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filePath, _ := DownloadAudio(mockClient, media)
		if filePath != "" {
			os.Remove(filePath)
		}
	}
}
