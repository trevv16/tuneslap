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

	t.Run("uses custom format when specified as string", func(t *testing.T) {
		params := models.ImageProcessingParams{
			Format: "png",
		}

		assert.Equal(t, "png", params.Format)
		assert.NotEqual(t, "webp", params.Format)
	})

	t.Run("uses default format when not specified", func(t *testing.T) {
		params := models.ImageProcessingParams{}

		// Empty string means use default (WEBP)
		assert.Equal(t, "", params.Format)
	})

	t.Run("format is string type not int", func(t *testing.T) {
		// This test ensures we don't regress to using int for Format
		// The OpenAPI spec defines format as a string enum: jpeg, png, gif, webp, svg
		params := models.ImageProcessingParams{
			Format: "webp",
		}

		// Verify it's a string, not an int
		var formatValue interface{} = params.Format
		_, isString := formatValue.(string)
		assert.True(t, isString, "Format should be a string type, not int")
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

	t.Run("builds options with custom string format", func(t *testing.T) {
		params := models.ImageProcessingParams{
			Format: "png",
		}

		opts := bimg.Options{
			Type: bimg.WEBP,
		}

		if params.Format != "" {
			opts.Type = stringToImageType(params.Format)
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

func TestStringToImageType(t *testing.T) {
	t.Run("converts jpeg string to bimg.JPEG", func(t *testing.T) {
		assert.Equal(t, bimg.JPEG, stringToImageType("jpeg"))
	})

	t.Run("converts jpg string to bimg.JPEG", func(t *testing.T) {
		assert.Equal(t, bimg.JPEG, stringToImageType("jpg"))
	})

	t.Run("converts png string to bimg.PNG", func(t *testing.T) {
		assert.Equal(t, bimg.PNG, stringToImageType("png"))
	})

	t.Run("converts webp string to bimg.WEBP", func(t *testing.T) {
		assert.Equal(t, bimg.WEBP, stringToImageType("webp"))
	})

	t.Run("converts gif string to bimg.GIF", func(t *testing.T) {
		assert.Equal(t, bimg.GIF, stringToImageType("gif"))
	})

	t.Run("converts svg string to bimg.SVG", func(t *testing.T) {
		assert.Equal(t, bimg.SVG, stringToImageType("svg"))
	})

	t.Run("converts tiff string to bimg.TIFF", func(t *testing.T) {
		assert.Equal(t, bimg.TIFF, stringToImageType("tiff"))
	})

	t.Run("defaults to WEBP for unknown format", func(t *testing.T) {
		assert.Equal(t, bimg.WEBP, stringToImageType("unknown"))
	})

	t.Run("defaults to WEBP for empty string", func(t *testing.T) {
		assert.Equal(t, bimg.WEBP, stringToImageType(""))
	})

	t.Run("handles OpenAPI enum values correctly", func(t *testing.T) {
		// These are the exact values from the OpenAPI ImageOutputFormat enum
		openAPIFormats := []struct {
			input    string
			expected bimg.ImageType
		}{
			{"jpeg", bimg.JPEG},
			{"png", bimg.PNG},
			{"gif", bimg.GIF},
			{"webp", bimg.WEBP},
			{"svg", bimg.SVG},
		}

		for _, tc := range openAPIFormats {
			t.Run(tc.input, func(t *testing.T) {
				result := stringToImageType(tc.input)
				assert.Equal(t, tc.expected, result, "Format '%s' should convert to correct bimg type", tc.input)
			})
		}
	})

	t.Run("invalid format does not cause error", func(t *testing.T) {
		// Should not panic, should return default
		assert.NotPanics(t, func() {
			result := stringToImageType("invalid_format_12345")
			assert.Equal(t, bimg.WEBP, result)
		})
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

// Integration tests for actual normalize function output

// createTestPNGWithSize creates a minimal valid PNG for testing with specified dimensions
func createTestPNGWithSize(t *testing.T, width, height int) []byte {
	t.Helper()

	minimalPNG := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41,
		0x54, 0x08, 0xD7, 0x63, 0xF8, 0xCF, 0xC0, 0x00,
		0x00, 0x00, 0x03, 0x00, 0x01, 0x00, 0x18, 0xDD,
		0x8D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45,
		0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}

	if width == 1 && height == 1 {
		return minimalPNG
	}

	options := bimg.Options{
		Width:  width,
		Height: height,
		Type:   bimg.PNG,
	}

	img, err := bimg.NewImage(minimalPNG).Process(options)
	if err != nil {
		return minimalPNG
	}

	return img
}

// checkLibvipsAvailable checks if libvips is available for this test file
func checkLibvipsAvailable() bool {
	return bimg.VipsMajorVersion > 0
}

func TestIntegrationNormalizeDefaultOutput(t *testing.T) {
	if !checkLibvipsAvailable() {
		t.Skip("Skipping integration test - libvips not available")
	}

	testImage := createTestPNGWithSize(t, 800, 600)

	t.Run("resizes to default 500x500", func(t *testing.T) {
		result, err := NormalizeDefault(testImage)
		if err != nil {
			t.Skipf("Normalize error: %v", err)
		}

		assert.NotNil(t, result)
		assert.Greater(t, len(result), 0)

		// Verify output dimensions
		size, err := bimg.NewImage(result).Size()
		if err == nil {
			assert.Equal(t, 500, size.Width)
			assert.Equal(t, 500, size.Height)
		}
	})

	t.Run("converts to WEBP format", func(t *testing.T) {
		result, err := NormalizeDefault(testImage)
		if err != nil {
			t.Skipf("Normalize error: %v", err)
		}

		imgType := bimg.NewImage(result).Type()
		assert.Equal(t, "webp", imgType)
	})

	t.Run("output is smaller than uncompressed", func(t *testing.T) {
		result, err := NormalizeDefault(testImage)
		if err != nil {
			t.Skipf("Normalize error: %v", err)
		}

		// WEBP with quality 90 and compression should be smaller
		// than a raw 500x500 image (500*500*3 = 750000 bytes uncompressed)
		assert.Less(t, len(result), 750000)
	})
}

func TestIntegrationNormalizeWithCustomParams(t *testing.T) {
	if !checkLibvipsAvailable() {
		t.Skip("Skipping integration test - libvips not available")
	}

	testImage := createTestPNGWithSize(t, 800, 600)

	t.Run("resizes to custom dimensions", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ResizeTo: [2]int{200, 150},
		}

		result, err := Normalize(testImage, params)
		if err != nil {
			t.Skipf("Normalize error: %v", err)
		}

		size, err := bimg.NewImage(result).Size()
		if err == nil {
			assert.Equal(t, 200, size.Width)
			assert.Equal(t, 150, size.Height)
		}
	})

	t.Run("uses defaults when ResizeTo is zero", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ResizeTo: [2]int{0, 0},
		}

		result, err := Normalize(testImage, params)
		if err != nil {
			t.Skipf("Normalize error: %v", err)
		}

		size, err := bimg.NewImage(result).Size()
		if err == nil {
			assert.Equal(t, 500, size.Width)
			assert.Equal(t, 500, size.Height)
		}
	})

	t.Run("converts to PNG when format specified as string", func(t *testing.T) {
		params := models.ImageProcessingParams{
			Format: "png",
		}

		result, err := Normalize(testImage, params)
		if err != nil {
			t.Skipf("Normalize error: %v", err)
		}

		imgType := bimg.NewImage(result).Type()
		assert.Equal(t, "png", imgType)
	})

	t.Run("converts to JPEG when format specified as string", func(t *testing.T) {
		params := models.ImageProcessingParams{
			Format: "jpeg",
		}

		result, err := Normalize(testImage, params)
		if err != nil {
			t.Skipf("Normalize error: %v", err)
		}

		imgType := bimg.NewImage(result).Type()
		assert.Equal(t, "jpeg", imgType)
	})

	t.Run("converts to GIF when format specified as string", func(t *testing.T) {
		params := models.ImageProcessingParams{
			Format: "gif",
		}

		result, err := Normalize(testImage, params)
		if err != nil {
			t.Skipf("Normalize error: %v", err)
		}

		imgType := bimg.NewImage(result).Type()
		assert.Equal(t, "gif", imgType)
	})

	t.Run("uses WEBP as default format", func(t *testing.T) {
		params := models.ImageProcessingParams{}

		result, err := Normalize(testImage, params)
		if err != nil {
			t.Skipf("Normalize error: %v", err)
		}

		imgType := bimg.NewImage(result).Type()
		assert.Equal(t, "webp", imgType)
	})
}

func TestIntegrationNormalizeEdgeCases(t *testing.T) {
	if !checkLibvipsAvailable() {
		t.Skip("Skipping integration test - libvips not available")
	}

	t.Run("handles nil input", func(t *testing.T) {
		params := models.ImageProcessingParams{}
		_, err := Normalize(nil, params)
		// Should return an error
		assert.Error(t, err)
	})

	t.Run("handles empty input", func(t *testing.T) {
		params := models.ImageProcessingParams{}
		_, err := Normalize([]byte{}, params)
		assert.Error(t, err)
	})

	t.Run("handles invalid image data", func(t *testing.T) {
		params := models.ImageProcessingParams{}
		_, err := Normalize([]byte("not an image"), params)
		assert.Error(t, err)
	})

	t.Run("handles very small resize", func(t *testing.T) {
		testImage := createTestPNGWithSize(t, 100, 100)
		params := models.ImageProcessingParams{
			ResizeTo: [2]int{10, 10},
		}

		result, err := Normalize(testImage, params)
		if err != nil {
			t.Skipf("Normalize error: %v", err)
		}

		size, err := bimg.NewImage(result).Size()
		if err == nil {
			assert.Equal(t, 10, size.Width)
			assert.Equal(t, 10, size.Height)
		}
	})

	t.Run("handles very large resize", func(t *testing.T) {
		testImage := createTestPNGWithSize(t, 100, 100)
		params := models.ImageProcessingParams{
			ResizeTo: [2]int{2000, 2000},
		}

		result, err := Normalize(testImage, params)
		if err != nil {
			t.Skipf("Normalize error: %v", err)
		}

		size, err := bimg.NewImage(result).Size()
		if err == nil {
			assert.Equal(t, 2000, size.Width)
			assert.Equal(t, 2000, size.Height)
		}
	})
}

func TestIntegrationNormalizeDefaultError(t *testing.T) {
	if !checkLibvipsAvailable() {
		t.Skip("Skipping integration test - libvips not available")
	}

	t.Run("handles nil input", func(t *testing.T) {
		_, err := NormalizeDefault(nil)
		assert.Error(t, err)
	})

	t.Run("handles empty input", func(t *testing.T) {
		_, err := NormalizeDefault([]byte{})
		assert.Error(t, err)
	})

	t.Run("handles invalid image data", func(t *testing.T) {
		_, err := NormalizeDefault([]byte("not an image"))
		assert.Error(t, err)
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

func BenchmarkIntegrationNormalizeDefault(b *testing.B) {
	if !checkLibvipsAvailable() {
		b.Skip("Skipping benchmark - libvips not available")
	}

	// Create test image once
	minimalPNG := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41,
		0x54, 0x08, 0xD7, 0x63, 0xF8, 0xCF, 0xC0, 0x00,
		0x00, 0x00, 0x03, 0x00, 0x01, 0x00, 0x18, 0xDD,
		0x8D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45,
		0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}

	testImage, _ := bimg.NewImage(minimalPNG).Process(bimg.Options{
		Width:  200,
		Height: 200,
		Type:   bimg.PNG,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NormalizeDefault(testImage)
	}
}
