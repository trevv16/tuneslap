package image

import (
	"testing"

	"github.com/h2non/bimg"
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
		format   int
		expected string
	}{
		{
			name:     "WEBP format",
			format:   int(bimg.WEBP),
			expected: "image/webp",
		},
		{
			name:     "PNG format",
			format:   int(bimg.PNG),
			expected: "image/png",
		},
		{
			name:     "unknown format defaults to PNG",
			format:   999,
			expected: "image/png",
		},
		{
			name:     "zero format defaults to PNG",
			format:   0,
			expected: "image/png",
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
	formats := []int{int(bimg.WEBP), int(bimg.PNG), 999}
	for i := 0; i < b.N; i++ {
		GetContentTypeFromFormat(formats[i%len(formats)])
	}
}
