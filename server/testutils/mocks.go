package testutils

import (
	"context"
	"errors"
	"time"
	"tuneslap/services/storage"
)

// Common errors
var (
	ErrNotFound         = errors.New("not found")
	ErrDuplicateKey     = errors.New("duplicate key error")
	ErrInvalidID        = errors.New("invalid object ID")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrRedisUnavailable = errors.New("redis unavailable")
)

// MockStorageClient implements storage.ObjectStorage for testing
type MockStorageClient struct {
	Files           map[string][]byte
	SignedURLs      map[string]string
	BucketName      string
	UploadError     error
	DownloadError   error
	DeleteError     error
	SignedURLError  error
	FileExistsError error
}

// NewMockStorageClient creates a new mock storage client
func NewMockStorageClient(bucketName string) *MockStorageClient {
	return &MockStorageClient{
		Files:      make(map[string][]byte),
		SignedURLs: make(map[string]string),
		BucketName: bucketName,
	}
}

// UploadFile implements storage.ObjectStorage
func (m *MockStorageClient) UploadFile(ctx context.Context, req storage.UploadFileRequest) error {
	if m.UploadError != nil {
		return m.UploadError
	}
	m.Files[req.ObjectName] = []byte("mock file content")
	return nil
}

// DownloadFile implements storage.ObjectStorage
func (m *MockStorageClient) DownloadFile(ctx context.Context, objectName, destPath string) error {
	if m.DownloadError != nil {
		return m.DownloadError
	}
	if _, exists := m.Files[objectName]; !exists {
		return ErrNotFound
	}
	return nil
}

// DeleteFile implements storage.ObjectStorage
func (m *MockStorageClient) DeleteFile(ctx context.Context, objectName string) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	delete(m.Files, objectName)
	return nil
}

// GenerateSignedUploadURL implements storage.ObjectStorage
func (m *MockStorageClient) GenerateSignedUploadURL(ctx context.Context, objectName, contentType string, expires time.Duration) (string, error) {
	if m.SignedURLError != nil {
		return "", m.SignedURLError
	}
	url := "https://mock-storage.example.com/upload/" + objectName
	m.SignedURLs[objectName] = url
	return url, nil
}

// GenerateSignedDownloadURL implements storage.ObjectStorage
func (m *MockStorageClient) GenerateSignedDownloadURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	if m.SignedURLError != nil {
		return "", m.SignedURLError
	}
	return "https://mock-storage.example.com/download/" + objectName, nil
}

// GetBucketName implements storage.ObjectStorage
func (m *MockStorageClient) GetBucketName() string {
	return m.BucketName
}

// GetFileURL implements storage.ObjectStorage
func (m *MockStorageClient) GetFileURL(objectName string) string {
	return "https://mock-storage.example.com/" + m.BucketName + "/" + objectName
}

// FileExists implements storage.ObjectStorage
func (m *MockStorageClient) FileExists(ctx context.Context, objectName string) (bool, error) {
	if m.FileExistsError != nil {
		return false, m.FileExistsError
	}
	_, exists := m.Files[objectName]
	return exists, nil
}

// SetFileContent sets mock file content
func (m *MockStorageClient) SetFileContent(objectName string, content []byte) {
	m.Files[objectName] = content
}

// MockTaskClient mocks the asynq task client
type MockTaskClient struct {
	EnqueuedTasks []MockTask
	EnqueueError  error
}

// MockTask represents an enqueued task
type MockTask struct {
	TypeName string
	Payload  []byte
	Options  []interface{}
}

// NewMockTaskClient creates a new mock task client
func NewMockTaskClient() *MockTaskClient {
	return &MockTaskClient{
		EnqueuedTasks: make([]MockTask, 0),
	}
}

// Enqueue mocks task enqueueing
func (m *MockTaskClient) Enqueue(taskType string, payload []byte, opts ...interface{}) error {
	if m.EnqueueError != nil {
		return m.EnqueueError
	}
	m.EnqueuedTasks = append(m.EnqueuedTasks, MockTask{
		TypeName: taskType,
		Payload:  payload,
		Options:  opts,
	})
	return nil
}

// GetEnqueuedTaskCount returns the number of enqueued tasks
func (m *MockTaskClient) GetEnqueuedTaskCount() int {
	return len(m.EnqueuedTasks)
}

// GetLastEnqueuedTask returns the last enqueued task
func (m *MockTaskClient) GetLastEnqueuedTask() *MockTask {
	if len(m.EnqueuedTasks) == 0 {
		return nil
	}
	return &m.EnqueuedTasks[len(m.EnqueuedTasks)-1]
}

// ClearEnqueuedTasks clears all enqueued tasks
func (m *MockTaskClient) ClearEnqueuedTasks() {
	m.EnqueuedTasks = make([]MockTask, 0)
}

// TestingT is an interface that both *testing.T and *testing.B implement
type TestingT interface {
	Skipf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Logf(format string, args ...interface{})
	Helper()
}
