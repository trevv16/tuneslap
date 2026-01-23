package audio

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProcessedFileName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple filename",
			input:    "test.mp3",
			expected: "processed-test.mp3",
		},
		{
			name:     "filename with spaces",
			input:    "my audio.mp3",
			expected: "processed-my audio.mp3",
		},
		{
			name:     "filename with special chars",
			input:    "audio_file-2024.mp3",
			expected: "processed-audio_file-2024.mp3",
		},
		{
			name:     "empty filename",
			input:    "",
			expected: "processed-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetProcessedFileName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetAudioFileDataFromPath(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "test.mp3")
	testContent := []byte("test audio data content")

	err := os.WriteFile(testFilePath, testContent, 0644)
	assert.NoError(t, err)

	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{
			name:        "valid file",
			path:        testFilePath,
			expectError: false,
		},
		{
			name:        "non-existent file",
			path:        "/nonexistent/path/file.mp3",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := GetAudioFileDataFromPath(tt.path)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testContent, data)
			}
		})
	}
}

// Benchmark tests
func BenchmarkGetProcessedFileName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetProcessedFileName("test-audio-file.mp3")
	}
}

func BenchmarkGetAudioFileDataFromPath(b *testing.B) {
	// Create a temporary test file
	tempDir := b.TempDir()
	testFilePath := filepath.Join(tempDir, "test.mp3")
	testContent := make([]byte, 1024) // 1KB test file
	os.WriteFile(testFilePath, testContent, 0644)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetAudioFileDataFromPath(testFilePath)
	}
}
