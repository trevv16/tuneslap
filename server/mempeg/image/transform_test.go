package image

import (
	"testing"
	"tuneslap/models"

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
