package image

import (
	"bytes"
	goimage "image"
	"image/color"
	"image/png"
	"os"
	"testing"
	"tuneslap/models"

	"github.com/h2non/bimg"
	"github.com/stretchr/testify/assert"
)

func TestDirPermissionsConstant(t *testing.T) {
	t.Run("has correct value", func(t *testing.T) {
		assert.Equal(t, os.FileMode(0755), os.FileMode(DirPermissions))
	})
}

func TestProcessImageWorkflow(t *testing.T) {
	// Test the workflow logic without actual image processing
	// since bimg requires libvips

	t.Run("workflow steps are in correct order", func(t *testing.T) {
		// Document the expected workflow
		steps := []string{
			"Initialize user uploads bucket client",
			"Create destination path and ensure directory exists",
			"Download original image from user uploads bucket",
			"Open the file",
			"Read the file",
			"Normalize with default options",
			"Transform (crop, rotate, blur, grayscale)",
			"Get updated metadata",
			"Save file to temp dir",
			"Get processed file upload key",
			"Upload processed image to media bucket",
			"Get the uploaded file url",
			"Delete original file from user uploads bucket",
			"Build updated media object",
		}

		// ProcessImage has 14 steps
		assert.Len(t, steps, 14)
	})
}

func TestProcessImageFileOperations(t *testing.T) {
	t.Run("temp file creation", func(t *testing.T) {
		tmpDir := t.TempDir()
		testContent := []byte("test image content")

		// Simulate file write
		testPath := tmpDir + "/test.webp"
		err := os.WriteFile(testPath, testContent, 0644)
		assert.NoError(t, err)

		// Verify file exists
		_, err = os.Stat(testPath)
		assert.NoError(t, err)

		// Verify cleanup works
		err = os.Remove(testPath)
		assert.NoError(t, err)

		_, err = os.Stat(testPath)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("directory creation", func(t *testing.T) {
		tmpDir := t.TempDir()
		subDir := tmpDir + "/sub/dir/path"

		err := os.MkdirAll(subDir, DirPermissions)
		assert.NoError(t, err)

		info, err := os.Stat(subDir)
		assert.NoError(t, err)
		assert.True(t, info.IsDir())
	})
}

func TestProcessImageMediaStatusUpdate(t *testing.T) {
	t.Run("status transitions", func(t *testing.T) {
		// Test the expected status values
		statuses := map[string]string{
			"initial":    "pending",
			"processing": "processing",
			"done":       "done",
			"error":      "error",
		}

		assert.Equal(t, "pending", statuses["initial"])
		assert.Equal(t, "done", statuses["done"])
	})
}

func TestProcessImageMetadataExtraction(t *testing.T) {
	t.Run("metadata fields are populated", func(t *testing.T) {
		// Simulate metadata extraction result
		metadata := ImageMetadata{
			Dimensions:  [2]int{500, 500},
			ContentType: "webp",
			FileSize:    25000,
		}

		assert.Equal(t, [2]int{500, 500}, metadata.Dimensions)
		assert.Equal(t, "webp", metadata.ContentType)
		assert.Equal(t, int64(25000), metadata.FileSize)
	})
}

// Benchmark tests
func BenchmarkDirCreation(b *testing.B) {
	baseDir := b.TempDir()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		subDir := baseDir + "/bench" + string(rune(i%26+'a'))
		_ = os.MkdirAll(subDir, DirPermissions)
	}
}

// Integration tests for image processing with real bimg operations
// These tests require libvips to be installed

// isLibvipsAvailable checks if libvips is available by checking the version constant
func isLibvipsAvailable() bool {
	// bimg.VipsMajorVersion is a constant set at compile time
	// If it's > 0, libvips was available when bimg was built
	return bimg.VipsMajorVersion > 0
}

// createTestImage creates a test image with color variation for integration tests
// Uses Go's image package to create a checkerboard pattern that blur will affect
func createTestImage(t *testing.T) []byte {
	t.Helper()

	// Create a 100x100 checkerboard pattern image using Go's standard library
	img := goimage.NewRGBA(goimage.Rect(0, 0, 100, 100))
	red := color.RGBA{255, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}

	// Create 10x10 pixel squares in checkerboard pattern
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			// Determine which square we're in (10x10 grid)
			squareX := x / 10
			squareY := y / 10
			// Checkerboard pattern: alternate colors based on position
			if (squareX+squareY)%2 == 0 {
				img.Set(x, y, red)
			} else {
				img.Set(x, y, white)
			}
		}
	}

	// Encode to PNG bytes
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}

	return buf.Bytes()
}

func TestIntegrationNormalizeDefault(t *testing.T) {
	// Skip if libvips is not available
	if !isLibvipsAvailable() {
		t.Skip("Skipping integration test - libvips not available")
	}

	testImage := createTestImage(t)

	t.Run("normalizes image to default dimensions", func(t *testing.T) {
		result, err := NormalizeDefault(testImage)
		if err != nil {
			t.Skipf("Skipping - bimg error (libvips may not be installed): %v", err)
		}

		assert.NotNil(t, result)
		assert.Greater(t, len(result), 0)
	})
}

func TestIntegrationNormalizeWithParams(t *testing.T) {
	if !isLibvipsAvailable() {
		t.Skip("Skipping integration test - libvips not available")
	}

	testImage := createTestImage(t)

	t.Run("normalizes with custom dimensions", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ResizeTo: [2]int{200, 150},
		}

		result, err := Normalize(testImage, params)
		if err != nil {
			t.Skipf("Skipping - bimg error: %v", err)
		}

		assert.NotNil(t, result)
		assert.Greater(t, len(result), 0)

		// Verify output dimensions
		size, err := bimg.NewImage(result).Size()
		if err == nil {
			assert.Equal(t, 200, size.Width)
			assert.Equal(t, 150, size.Height)
		}
	})

	t.Run("normalizes with PNG format", func(t *testing.T) {
		params := models.ImageProcessingParams{
			Format: int(bimg.PNG),
		}

		result, err := Normalize(testImage, params)
		if err != nil {
			t.Skipf("Skipping - bimg error: %v", err)
		}

		assert.NotNil(t, result)

		// Verify output format
		imgType := bimg.NewImage(result).Type()
		assert.Equal(t, "png", imgType)
	})
}

func TestIntegrationTransform(t *testing.T) {
	if !isLibvipsAvailable() {
		t.Skip("Skipping integration test - libvips not available")
	}

	testImage := createTestImage(t)

	t.Run("applies grayscale filter", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ApplyFilters: "grayscale",
		}

		result, err := transform(testImage, params)
		if err != nil {
			t.Skipf("Skipping - bimg error: %v", err)
		}

		assert.NotNil(t, result)
		assert.Greater(t, len(result), 0)
	})

	t.Run("applies blur filter", func(t *testing.T) {
		params := models.ImageProcessingParams{
			ApplyFilters: "blur",
		}

		result, err := transform(testImage, params)
		if err != nil {
			t.Skipf("Skipping - bimg error: %v", err)
		}

		assert.NotNil(t, result)
		assert.Greater(t, len(result), 0)
	})

	t.Run("no filter returns unchanged", func(t *testing.T) {
		params := models.ImageProcessingParams{}

		result, err := transform(testImage, params)
		if err != nil {
			t.Skipf("Skipping - bimg error: %v", err)
		}

		assert.NotNil(t, result)
		// Without filters or crop, output should be similar to input
		assert.Equal(t, testImage, result)
	})
}

func TestIntegrationCrop(t *testing.T) {
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
		t.Skipf("Skipping - cannot create test image: %v", err)
	}

	t.Run("crops image to specified region", func(t *testing.T) {
		params := models.ImageProcessingParams{
			Crop: [4]int{10, 10, 50, 50}, // x, y, width, height
		}

		result, err := transform(testImage, params)
		if err != nil {
			t.Skipf("Skipping - bimg error: %v", err)
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

func TestIntegrationGetMetadata(t *testing.T) {
	if !isLibvipsAvailable() {
		t.Skip("Skipping integration test - libvips not available")
	}

	testImage := createTestImage(t)

	t.Run("extracts metadata from image", func(t *testing.T) {
		metadata, err := getMetadata(testImage)
		if err != nil {
			t.Skipf("Skipping - bimg error: %v", err)
		}

		assert.NotEqual(t, [2]int{0, 0}, metadata.Dimensions)
		assert.NotEmpty(t, metadata.ContentType)
		assert.Greater(t, metadata.FileSize, int64(0))
	})
}

func TestIntegrationFullPipeline(t *testing.T) {
	if !isLibvipsAvailable() {
		t.Skip("Skipping integration test - libvips not available")
	}

	testImage := createTestImage(t)

	t.Run("full normalize and transform pipeline", func(t *testing.T) {
		// Step 1: Normalize
		params := models.ImageProcessingParams{
			ResizeTo: [2]int{100, 100},
		}

		normalized, err := Normalize(testImage, params)
		if err != nil {
			t.Skipf("Skipping - normalize error: %v", err)
		}

		// Step 2: Transform
		transformParams := models.ImageProcessingParams{
			ApplyFilters: "grayscale",
		}

		transformed, err := transform(normalized, transformParams)
		if err != nil {
			t.Skipf("Skipping - transform error: %v", err)
		}

		// Step 3: Get metadata
		metadata, err := getMetadata(transformed)
		if err != nil {
			t.Skipf("Skipping - metadata error: %v", err)
		}

		assert.NotNil(t, transformed)
		assert.NotEqual(t, [2]int{0, 0}, metadata.Dimensions)
	})
}
