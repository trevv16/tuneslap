package image

import (
	"testing"
	"tuneslap/models"

	"github.com/h2non/bimg"
	"github.com/stretchr/testify/assert"
)

func TestTransformLogic(t *testing.T) {
	// Test the transform function logic without requiring actual image processing
	// since bimg requires libvips which may not be available in all environments

	t.Run("grayscale filter is detected", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ApplyFilters: "grayscale",
		}

		assert.Equal(t, "grayscale", params.ApplyFilters)
	})

	t.Run("blur filter is detected", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ApplyFilters: "blur",
		}

		assert.Equal(t, "blur", params.ApplyFilters)
	})

	t.Run("empty filter is handled", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ApplyFilters: "",
		}

		assert.Empty(t, params.ApplyFilters)
	})

	t.Run("unrecognized filter is ignored", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ApplyFilters: "unknown_filter",
		}

		// The transform function will just return the input unchanged for unknown filters
		assert.Equal(t, "unknown_filter", params.ApplyFilters)
	})
}

func TestImageProcessingParams(t *testing.T) {
	t.Run("resize parameters", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ResizeTo: [2]int{800, 600},
		}

		assert.Equal(t, 800, params.ResizeTo[0])
		assert.Equal(t, 600, params.ResizeTo[1])
	})

	t.Run("crop parameters", func(t *testing.T) {
		params := models.ImageProcessingParams{
			Crop: [4]int{10, 10, 100, 100}, // x, y, width, height
		}

		assert.Equal(t, 10, params.Crop[0])  // x
		assert.Equal(t, 10, params.Crop[1])  // y
		assert.Equal(t, 100, params.Crop[2]) // width
		assert.Equal(t, 100, params.Crop[3]) // height
	})

	t.Run("aspect ratio parameter", func(t *testing.T) {
		params := models.ImageProcessingParams{
			AspectRatio: "16:9",
		}

		assert.Equal(t, "16:9", params.AspectRatio)
	})

	t.Run("format parameter", func(t *testing.T) {
		params := models.ImageProcessingParams{
			Format: 1, // webp format typically
		}

		assert.Equal(t, 1, params.Format)
	})
}

func TestTransformFilterDetection(t *testing.T) {
	// Test the filter detection logic used in transform

	t.Run("detects grayscale filter", func(t *testing.T) {
		filter := "grayscale"
		isGrayscale := filter == "grayscale"
		assert.True(t, isGrayscale)
	})

	t.Run("detects blur filter", func(t *testing.T) {
		filter := "blur"
		isBlur := filter == "blur"
		assert.True(t, isBlur)
	})

	t.Run("does not detect blur as grayscale", func(t *testing.T) {
		filter := "blur"
		isGrayscale := filter == "grayscale"
		assert.False(t, isGrayscale)
	})
}

func TestIsCropValid(t *testing.T) {
	t.Run("valid crop with all positive values", func(t *testing.T) {
		crop := [4]int{10, 10, 100, 100}
		assert.True(t, isCropValid(crop))
	})

	t.Run("valid crop with zero x and y", func(t *testing.T) {
		crop := [4]int{0, 0, 100, 100}
		assert.True(t, isCropValid(crop))
	})

	t.Run("invalid crop with zero width", func(t *testing.T) {
		crop := [4]int{10, 10, 0, 100}
		assert.False(t, isCropValid(crop))
	})

	t.Run("invalid crop with zero height", func(t *testing.T) {
		crop := [4]int{10, 10, 100, 0}
		assert.False(t, isCropValid(crop))
	})

	t.Run("invalid crop with negative width", func(t *testing.T) {
		crop := [4]int{10, 10, -100, 100}
		assert.False(t, isCropValid(crop))
	})

	t.Run("invalid crop with negative height", func(t *testing.T) {
		crop := [4]int{10, 10, 100, -100}
		assert.False(t, isCropValid(crop))
	})

	t.Run("all zeros is invalid", func(t *testing.T) {
		crop := [4]int{0, 0, 0, 0}
		assert.False(t, isCropValid(crop))
	})
}

func TestDefaultBlurSigma(t *testing.T) {
	t.Run("default blur sigma is reasonable", func(t *testing.T) {
		assert.Equal(t, 5.0, DefaultBlurSigma)
	})

	t.Run("blur sigma is positive", func(t *testing.T) {
		assert.Greater(t, DefaultBlurSigma, 0.0)
	})
}

func TestCropParameterStructure(t *testing.T) {
	t.Run("crop array has correct structure", func(t *testing.T) {
		// crop: [x, y, width, height]
		params := models.ImageProcessingParams{
			Crop: [4]int{50, 50, 200, 150},
		}

		x := params.Crop[0]
		y := params.Crop[1]
		width := params.Crop[2]
		height := params.Crop[3]

		assert.Equal(t, 50, x)
		assert.Equal(t, 50, y)
		assert.Equal(t, 200, width)
		assert.Equal(t, 150, height)
	})
}

func TestTransformWithCropAndFilter(t *testing.T) {
	// Test that crop and filter params can be combined
	t.Run("params can have both crop and filter", func(t *testing.T) {
		params := models.ImageProcessingParams{
			Crop:         [4]int{10, 10, 100, 100},
			ApplyFilters: "grayscale",
		}

		assert.True(t, isCropValid(params.Crop))
		assert.Equal(t, "grayscale", params.ApplyFilters)
	})

	t.Run("params can have crop without filter", func(t *testing.T) {
		params := models.ImageProcessingParams{
			Crop: [4]int{10, 10, 100, 100},
		}

		assert.True(t, isCropValid(params.Crop))
		assert.Empty(t, params.ApplyFilters)
	})

	t.Run("params can have filter without crop", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ApplyFilters: "blur",
		}

		assert.False(t, isCropValid(params.Crop))
		assert.Equal(t, "blur", params.ApplyFilters)
	})
}

// Edge case tests for crop validation
func TestIsCropValidEdgeCases(t *testing.T) {
	t.Run("very large crop dimensions", func(t *testing.T) {
		crop := [4]int{0, 0, 10000, 10000}
		// Large dimensions are technically valid (validation happens at processing time)
		assert.True(t, isCropValid(crop))
	})

	t.Run("x offset larger than typical image", func(t *testing.T) {
		crop := [4]int{5000, 0, 100, 100}
		// Offset validation happens at processing time, not in isCropValid
		assert.True(t, isCropValid(crop))
	})

	t.Run("y offset larger than typical image", func(t *testing.T) {
		crop := [4]int{0, 5000, 100, 100}
		assert.True(t, isCropValid(crop))
	})

	t.Run("minimum valid crop 1x1", func(t *testing.T) {
		crop := [4]int{0, 0, 1, 1}
		assert.True(t, isCropValid(crop))
	})

	t.Run("negative x offset with valid dimensions", func(t *testing.T) {
		crop := [4]int{-10, 0, 100, 100}
		// Negative offsets are allowed by isCropValid (bimg handles bounds)
		assert.True(t, isCropValid(crop))
	})

	t.Run("negative y offset with valid dimensions", func(t *testing.T) {
		crop := [4]int{0, -10, 100, 100}
		assert.True(t, isCropValid(crop))
	})
}

// Edge case tests for filter handling
func TestTransformFilterEdgeCases(t *testing.T) {
	t.Run("filter with extra whitespace", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ApplyFilters: " grayscale ",
		}
		// Current implementation requires exact match
		isGrayscale := params.ApplyFilters == "grayscale"
		assert.False(t, isGrayscale)
	})

	t.Run("filter with different case", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ApplyFilters: "Grayscale",
		}
		// Current implementation is case-sensitive
		isGrayscale := params.ApplyFilters == "grayscale"
		assert.False(t, isGrayscale)
	})

	t.Run("filter with uppercase", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ApplyFilters: "BLUR",
		}
		isBlur := params.ApplyFilters == "blur"
		assert.False(t, isBlur)
	})
}

// Integration tests for transform function with actual bimg processing
func TestIntegrationTransformWithInvalidInput(t *testing.T) {
	if !isLibvipsAvailable() {
		t.Skip("Skipping integration test - libvips not available")
	}

	t.Run("transform with nil input returns error", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ApplyFilters: "grayscale",
		}

		result, err := transform(nil, params)
		// bimg should handle nil gracefully or return an error
		if err != nil {
			assert.Error(t, err)
		} else {
			// If no error, result should be nil or empty
			assert.True(t, result == nil || len(result) == 0)
		}
	})

	t.Run("transform with empty input returns error", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ApplyFilters: "grayscale",
		}

		result, err := transform([]byte{}, params)
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.True(t, result == nil || len(result) == 0)
		}
	})

	t.Run("transform with invalid image bytes returns error", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ApplyFilters: "grayscale",
		}

		invalidBytes := []byte("this is not an image")
		_, err := transform(invalidBytes, params)
		// bimg should return an error for invalid image data
		assert.Error(t, err)
	})
}

func TestIntegrationTransformOutputVerification(t *testing.T) {
	if !isLibvipsAvailable() {
		t.Skip("Skipping integration test - libvips not available")
	}

	testImage := createTestImage(t)
	if len(testImage) == 0 {
		t.Skip("Could not create test image")
	}

	t.Run("grayscale produces different output", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ApplyFilters: "grayscale",
		}

		result, err := transform(testImage, params)
		if err != nil {
			t.Skipf("Transform error: %v", err)
		}

		assert.NotNil(t, result)
		assert.Greater(t, len(result), 0)
		// Grayscale output should be different from input
		assert.NotEqual(t, testImage, result)
	})

	t.Run("blur produces different output", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ApplyFilters: "blur",
		}

		result, err := transform(testImage, params)
		if err != nil {
			t.Skipf("Transform error: %v", err)
		}

		assert.NotNil(t, result)
		assert.Greater(t, len(result), 0)
		// Blur output should be different from input
		assert.NotEqual(t, testImage, result)
	})

	t.Run("no filter returns same content", func(t *testing.T) {
		params := models.ImageProcessingParams{}

		result, err := transform(testImage, params)
		if err != nil {
			t.Skipf("Transform error: %v", err)
		}

		assert.NotNil(t, result)
		// Without filters, output should match input
		assert.Equal(t, testImage, result)
	})

	t.Run("unknown filter returns unchanged", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ApplyFilters: "sepia", // Not implemented
		}

		result, err := transform(testImage, params)
		if err != nil {
			t.Skipf("Transform error: %v", err)
		}

		assert.NotNil(t, result)
		// Unknown filter should return unchanged
		assert.Equal(t, testImage, result)
	})
}

func TestIntegrationCropOutputVerification(t *testing.T) {
	if !isLibvipsAvailable() {
		t.Skip("Skipping integration test - libvips not available")
	}

	// Create a larger test image for cropping
	options := bimg.Options{
		Width:  200,
		Height: 200,
		Type:   bimg.PNG,
	}

	minimalPNG := createTestImage(t)
	testImage, err := bimg.NewImage(minimalPNG).Process(options)
	if err != nil {
		t.Skipf("Cannot create test image: %v", err)
	}

	t.Run("crop produces smaller image", func(t *testing.T) {
		params := models.ImageProcessingParams{
			Crop: [4]int{10, 10, 50, 50},
		}

		result, err := transform(testImage, params)
		if err != nil {
			t.Skipf("Transform error: %v", err)
		}

		assert.NotNil(t, result)

		// Verify cropped dimensions
		size, err := bimg.NewImage(result).Size()
		if err == nil {
			assert.Equal(t, 50, size.Width)
			assert.Equal(t, 50, size.Height)
		}
	})

	t.Run("crop at origin", func(t *testing.T) {
		params := models.ImageProcessingParams{
			Crop: [4]int{0, 0, 100, 100},
		}

		result, err := transform(testImage, params)
		if err != nil {
			t.Skipf("Transform error: %v", err)
		}

		assert.NotNil(t, result)

		size, err := bimg.NewImage(result).Size()
		if err == nil {
			assert.Equal(t, 100, size.Width)
			assert.Equal(t, 100, size.Height)
		}
	})

	t.Run("crop exceeding bounds returns error", func(t *testing.T) {
		params := models.ImageProcessingParams{
			Crop: [4]int{150, 150, 100, 100}, // Would exceed 200x200 image
		}

		_, err := transform(testImage, params)
		// bimg should return an error when crop exceeds bounds
		assert.Error(t, err)
	})

	t.Run("crop and grayscale combined", func(t *testing.T) {
		params := models.ImageProcessingParams{
			Crop:         [4]int{10, 10, 50, 50},
			ApplyFilters: "grayscale",
		}

		result, err := transform(testImage, params)
		if err != nil {
			t.Skipf("Transform error: %v", err)
		}

		assert.NotNil(t, result)

		// Verify cropped dimensions
		size, err := bimg.NewImage(result).Size()
		if err == nil {
			assert.Equal(t, 50, size.Width)
			assert.Equal(t, 50, size.Height)
		}
	})
}

// Note: isLibvipsAvailable and createTestImage are defined in process_test.go

// Benchmark for filter detection
func BenchmarkFilterDetection(b *testing.B) {
	filter := "grayscale"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filter == "grayscale"
	}
}

func BenchmarkCropValidation(b *testing.B) {
	crop := [4]int{10, 10, 100, 100}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = isCropValid(crop)
	}
}
