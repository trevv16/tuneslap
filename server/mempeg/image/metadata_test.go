package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageMetadataStruct(t *testing.T) {
	t.Run("creates metadata with all fields", func(t *testing.T) {
		metadata := ImageMetadata{
			Dimensions:  [2]int{1920, 1080},
			ContentType: "image/webp",
			FileSize:    1024000,
		}

		assert.Equal(t, 1920, metadata.Dimensions[0])
		assert.Equal(t, 1080, metadata.Dimensions[1])
		assert.Equal(t, "image/webp", metadata.ContentType)
		assert.Equal(t, int64(1024000), metadata.FileSize)
	})

	t.Run("handles square dimensions", func(t *testing.T) {
		metadata := ImageMetadata{
			Dimensions:  [2]int{500, 500},
			ContentType: "image/png",
			FileSize:    512000,
		}

		assert.Equal(t, metadata.Dimensions[0], metadata.Dimensions[1])
	})

	t.Run("handles portrait dimensions", func(t *testing.T) {
		metadata := ImageMetadata{
			Dimensions:  [2]int{1080, 1920},
			ContentType: "image/jpeg",
			FileSize:    2048000,
		}

		// Portrait: height > width
		assert.Greater(t, metadata.Dimensions[1], metadata.Dimensions[0])
	})

	t.Run("handles landscape dimensions", func(t *testing.T) {
		metadata := ImageMetadata{
			Dimensions:  [2]int{1920, 1080},
			ContentType: "image/png",
			FileSize:    3072000,
		}

		// Landscape: width > height
		assert.Greater(t, metadata.Dimensions[0], metadata.Dimensions[1])
	})
}

func TestImageMetadataContentTypes(t *testing.T) {
	validContentTypes := []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
		"jpeg", // bimg returns short type names
		"png",
		"gif",
		"webp",
	}

	for _, contentType := range validContentTypes {
		t.Run("content type: "+contentType, func(t *testing.T) {
			metadata := ImageMetadata{
				ContentType: contentType,
			}
			assert.NotEmpty(t, metadata.ContentType)
		})
	}
}

func TestImageMetadataFileSize(t *testing.T) {
	t.Run("zero file size", func(t *testing.T) {
		metadata := ImageMetadata{
			FileSize: 0,
		}
		assert.Equal(t, int64(0), metadata.FileSize)
	})

	t.Run("large file size", func(t *testing.T) {
		// 100MB
		metadata := ImageMetadata{
			FileSize: 100 * 1024 * 1024,
		}
		assert.Equal(t, int64(104857600), metadata.FileSize)
	})
}

// Benchmark tests
func BenchmarkImageMetadataCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ImageMetadata{
			Dimensions:  [2]int{1920, 1080},
			ContentType: "image/webp",
			FileSize:    1024000,
		}
	}
}
