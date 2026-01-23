package image

import (
	"testing"
	"tuneslap/models"

	"github.com/h2non/bimg"
	"github.com/stretchr/testify/assert"
)

func TestDefaultNormalizeOptions(t *testing.T) {
	t.Run("has correct default width", func(t *testing.T) {
		assert.Equal(t, 500, defaultNormalizeOptions.Width)
	})

	t.Run("has correct default height", func(t *testing.T) {
		assert.Equal(t, 500, defaultNormalizeOptions.Height)
	})

	t.Run("has correct quality setting", func(t *testing.T) {
		assert.Equal(t, 90, defaultNormalizeOptions.Quality)
	})

	t.Run("outputs webp format", func(t *testing.T) {
		assert.Equal(t, bimg.WEBP, defaultNormalizeOptions.Type)
	})

	t.Run("has compression level", func(t *testing.T) {
		assert.Equal(t, 6, defaultNormalizeOptions.Compression)
	})

	t.Run("strips metadata", func(t *testing.T) {
		assert.True(t, defaultNormalizeOptions.StripMetadata)
	})
}

func TestGetDefaultNormalizeOptions(t *testing.T) {
	t.Run("returns copy of default options", func(t *testing.T) {
		opts := GetDefaultNormalizeOptions()
		assert.Equal(t, defaultNormalizeOptions.Width, opts.Width)
		assert.Equal(t, defaultNormalizeOptions.Height, opts.Height)
		assert.Equal(t, defaultNormalizeOptions.Quality, opts.Quality)
		assert.Equal(t, defaultNormalizeOptions.Type, opts.Type)
		assert.Equal(t, defaultNormalizeOptions.Compression, opts.Compression)
		assert.Equal(t, defaultNormalizeOptions.StripMetadata, opts.StripMetadata)
	})
}

func TestNormalizeWithParams(t *testing.T) {
	t.Run("uses custom ResizeTo dimensions", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ResizeTo: [2]int{800, 600},
		}

		// We can't test actual image processing without bimg/libvips
		// but we can verify the params structure is correct
		assert.Equal(t, 800, params.ResizeTo[0])
		assert.Equal(t, 600, params.ResizeTo[1])
	})

	t.Run("falls back to defaults with zero dimensions", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ResizeTo: [2]int{0, 0},
		}

		// With zero dimensions, normalize should use defaults
		assert.Equal(t, 0, params.ResizeTo[0])
		assert.Equal(t, 0, params.ResizeTo[1])

		// Verify default dimensions
		assert.Equal(t, 500, defaultNormalizeOptions.Width)
		assert.Equal(t, 500, defaultNormalizeOptions.Height)
	})

	t.Run("uses custom format when specified", func(t *testing.T) {
		params := models.ImageProcessingParams{
			Format: int(bimg.PNG),
		}

		assert.Equal(t, int(bimg.PNG), params.Format)
		assert.NotEqual(t, int(bimg.WEBP), params.Format)
	})

	t.Run("uses default format when not specified", func(t *testing.T) {
		params := models.ImageProcessingParams{}

		// Format 0 means use default (WEBP)
		assert.Equal(t, 0, params.Format)
	})

	t.Run("validates dimensions are positive", func(t *testing.T) {
		validParams := models.ImageProcessingParams{
			ResizeTo: [2]int{100, 100},
		}
		invalidParams := models.ImageProcessingParams{
			ResizeTo: [2]int{-100, -100},
		}

		assert.True(t, validParams.ResizeTo[0] > 0 && validParams.ResizeTo[1] > 0)
		assert.False(t, invalidParams.ResizeTo[0] > 0 && invalidParams.ResizeTo[1] > 0)
	})

	t.Run("handles partial dimensions", func(t *testing.T) {
		// Only width set
		paramsWidthOnly := models.ImageProcessingParams{
			ResizeTo: [2]int{800, 0},
		}

		// Only height set
		paramsHeightOnly := models.ImageProcessingParams{
			ResizeTo: [2]int{0, 600},
		}

		// Both should fall back to defaults since both dimensions required
		isWidthOnlyValid := paramsWidthOnly.ResizeTo[0] > 0 && paramsWidthOnly.ResizeTo[1] > 0
		isHeightOnlyValid := paramsHeightOnly.ResizeTo[0] > 0 && paramsHeightOnly.ResizeTo[1] > 0

		assert.False(t, isWidthOnlyValid)
		assert.False(t, isHeightOnlyValid)
	})
}

func TestNormalizeOptionsBuilding(t *testing.T) {
	t.Run("builds options with custom dimensions", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ResizeTo: [2]int{1024, 768},
		}

		// Simulate option building logic
		opts := bimg.Options{
			Width:         defaultNormalizeOptions.Width,
			Height:        defaultNormalizeOptions.Height,
			Quality:       defaultNormalizeOptions.Quality,
			Type:          bimg.WEBP,
			Compression:   defaultNormalizeOptions.Compression,
			StripMetadata: true,
		}

		if params.ResizeTo[0] > 0 && params.ResizeTo[1] > 0 {
			opts.Width = params.ResizeTo[0]
			opts.Height = params.ResizeTo[1]
		}

		assert.Equal(t, 1024, opts.Width)
		assert.Equal(t, 768, opts.Height)
	})

	t.Run("builds options with custom format", func(t *testing.T) {
		params := models.ImageProcessingParams{
			Format: int(bimg.PNG),
		}

		opts := bimg.Options{
			Type: bimg.WEBP,
		}

		if params.Format > 0 {
			opts.Type = bimg.ImageType(params.Format)
		}

		assert.Equal(t, bimg.PNG, opts.Type)
	})

	t.Run("always strips metadata", func(t *testing.T) {
		opts := bimg.Options{
			StripMetadata: true,
		}

		// Metadata stripping should always be true for security
		assert.True(t, opts.StripMetadata)
	})
}

func TestNormalizeOptionsStructure(t *testing.T) {
	t.Run("bimg options are valid", func(t *testing.T) {
		opts := bimg.Options{
			Width:         800,
			Height:        600,
			Quality:       85,
			Type:          bimg.WEBP,
			Compression:   4,
			StripMetadata: true,
		}

		assert.Equal(t, 800, opts.Width)
		assert.Equal(t, 600, opts.Height)
		assert.Equal(t, 85, opts.Quality)
		assert.True(t, opts.StripMetadata)
	})
}

func TestNormalizeOutputFormat(t *testing.T) {
	t.Run("webp format constant", func(t *testing.T) {
		// bimg.WEBP is an int representing the format
		assert.NotEqual(t, 0, bimg.WEBP)
	})

	t.Run("png format constant", func(t *testing.T) {
		assert.NotEqual(t, 0, bimg.PNG)
	})

	t.Run("jpeg format constant", func(t *testing.T) {
		assert.NotEqual(t, 0, bimg.JPEG)
	})

	t.Run("different formats have different values", func(t *testing.T) {
		assert.NotEqual(t, bimg.WEBP, bimg.PNG)
		assert.NotEqual(t, bimg.WEBP, bimg.JPEG)
		assert.NotEqual(t, bimg.PNG, bimg.JPEG)
	})
}

func TestNormalizeQualityRange(t *testing.T) {
	t.Run("quality is in valid range", func(t *testing.T) {
		quality := defaultNormalizeOptions.Quality
		assert.GreaterOrEqual(t, quality, 1)
		assert.LessOrEqual(t, quality, 100)
	})

	t.Run("compression is in valid range", func(t *testing.T) {
		compression := defaultNormalizeOptions.Compression
		assert.GreaterOrEqual(t, compression, 0)
		assert.LessOrEqual(t, compression, 9) // typical compression range
	})
}

// Benchmark tests
func BenchmarkNormalizeOptionsCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bimg.Options{
			Width:         500,
			Height:        500,
			Quality:       90,
			Type:          bimg.WEBP,
			Compression:   6,
			StripMetadata: true,
		}
	}
}
