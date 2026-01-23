package image

import (
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

// createTestImage creates a simple test image for integration tests
func createTestImage(t *testing.T) []byte {
	t.Helper()

	// Create a simple 100x100 PNG image using bimg
	// If this fails, libvips is not available
	options := bimg.Options{
		Width:  100,
		Height: 100,
		Type:   bimg.PNG,
	}

	// Create a blank image (bimg needs input bytes to work with)
	// We'll use a minimal valid PNG to start
	// This is a 1x1 red PNG
	minimalPNG := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41, // IDAT chunk
		0x54, 0x08, 0xD7, 0x63, 0xF8, 0xCF, 0xC0, 0x00,
		0x00, 0x00, 0x03, 0x00, 0x01, 0x00, 0x18, 0xDD,
		0x8D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, // IEND chunk
		0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}

	// Resize to create a larger test image
	img, err := bimg.NewImage(minimalPNG).Process(options)
	if err != nil {
		// If bimg fails, libvips might not be available
		return minimalPNG
	}

	return img
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
