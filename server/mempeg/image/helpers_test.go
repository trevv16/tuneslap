package image

import (
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
			input:    "test.png",
			expected: "processed-test.png",
		},
		{
			name:     "filename with spaces",
			input:    "my image.png",
			expected: "processed-my image.png",
		},
		{
			name:     "filename with special chars",
			input:    "image_file-2024.jpg",
			expected: "processed-image_file-2024.jpg",
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

func TestGetContentTypeFromFormat(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{
			name:     "webp format",
			format:   "webp",
			expected: "image/webp",
		},
		{
			name:     "png format",
			format:   "png",
			expected: "image/png",
		},
		{
			name:     "jpeg format",
			format:   "jpeg",
			expected: "image/jpeg",
		},
		{
			name:     "jpg format",
			format:   "jpg",
			expected: "image/jpeg",
		},
		{
			name:     "gif format",
			format:   "gif",
			expected: "image/gif",
		},
		{
			name:     "svg format",
			format:   "svg",
			expected: "image/svg+xml",
		},
		{
			name:     "tiff format",
			format:   "tiff",
			expected: "image/tiff",
		},
		{
			name:     "unknown format defaults to webp",
			format:   "unknown",
			expected: "image/webp",
		},
		{
			name:     "empty format defaults to webp",
			format:   "",
			expected: "image/webp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetContentTypeFromFormat(tt.format)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Benchmark tests
func BenchmarkGetProcessedFileName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetProcessedFileName("test-image-file.png")
	}
}

func BenchmarkGetContentTypeFromFormat(b *testing.B) {
	formats := []string{"webp", "png", "jpeg", "gif", "unknown"}
	for i := 0; i < b.N; i++ {
		GetContentTypeFromFormat(formats[i%len(formats)])
	}
}
