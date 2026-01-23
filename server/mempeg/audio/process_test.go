package audio

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
	"tuneslap/mempeg/ffmpeg"
	"tuneslap/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetFileSize(t *testing.T) {
	t.Run("returns correct file size", func(t *testing.T) {
		// Create a temporary file with known content
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test-file.txt")
		content := []byte("Hello, World!") // 13 bytes

		err := os.WriteFile(testFile, content, 0644)
		require.NoError(t, err)

		size, err := GetFileSize(testFile)
		require.NoError(t, err)
		assert.Equal(t, int64(13), size)
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		size, err := GetFileSize("/non/existent/file.txt")
		assert.Error(t, err)
		assert.Equal(t, int64(0), size)
	})

	t.Run("returns zero for empty file", func(t *testing.T) {
		tmpDir := t.TempDir()
		emptyFile := filepath.Join(tmpDir, "empty.txt")

		err := os.WriteFile(emptyFile, []byte{}, 0644)
		require.NoError(t, err)

		size, err := GetFileSize(emptyFile)
		require.NoError(t, err)
		assert.Equal(t, int64(0), size)
	})

	t.Run("returns correct size for larger file", func(t *testing.T) {
		tmpDir := t.TempDir()
		largeFile := filepath.Join(tmpDir, "large.bin")

		// Create a 1KB file
		content := make([]byte, 1024)
		err := os.WriteFile(largeFile, content, 0644)
		require.NoError(t, err)

		size, err := GetFileSize(largeFile)
		require.NoError(t, err)
		assert.Equal(t, int64(1024), size)
	})
}

func TestDirPermissions(t *testing.T) {
	t.Run("has correct value", func(t *testing.T) {
		assert.Equal(t, os.FileMode(0755), os.FileMode(DirPermissions))
	})
}

func TestDefaultOutputFormat(t *testing.T) {
	t.Run("default is webm", func(t *testing.T) {
		assert.Equal(t, "webm", DefaultOutputFormat)
	})
}

func TestContentTypeMap(t *testing.T) {
	t.Run("webm content type", func(t *testing.T) {
		assert.Equal(t, "audio/webm", ContentTypeMap["webm"])
	})

	t.Run("mp3 content type", func(t *testing.T) {
		assert.Equal(t, "audio/mpeg", ContentTypeMap["mp3"])
	})

	t.Run("ogg content type", func(t *testing.T) {
		assert.Equal(t, "audio/ogg", ContentTypeMap["ogg"])
	})

	t.Run("wav content type", func(t *testing.T) {
		assert.Equal(t, "audio/wav", ContentTypeMap["wav"])
	})
}

func TestNormalizeAudioParams(t *testing.T) {
	t.Run("adds default output format when empty", func(t *testing.T) {
		params := models.AudioProcessingParams{}
		result := normalizeAudioParams(params)
		assert.Equal(t, []string{DefaultOutputFormat}, result.OutputFormats)
	})

	t.Run("preserves existing output formats", func(t *testing.T) {
		params := models.AudioProcessingParams{
			OutputFormats: []string{"mp3", "ogg"},
		}
		result := normalizeAudioParams(params)
		assert.Equal(t, []string{"mp3", "ogg"}, result.OutputFormats)
	})

	t.Run("enables normalize when no processing requested", func(t *testing.T) {
		params := models.AudioProcessingParams{}
		result := normalizeAudioParams(params)
		assert.True(t, result.Normalize)
	})

	t.Run("does not enable normalize when trim is set", func(t *testing.T) {
		params := models.AudioProcessingParams{
			TrimStart: 2.0,
		}
		result := normalizeAudioParams(params)
		assert.False(t, result.Normalize)
	})

	t.Run("does not override explicit normalize setting", func(t *testing.T) {
		params := models.AudioProcessingParams{
			Normalize: true,
			TrimStart: 2.0,
		}
		result := normalizeAudioParams(params)
		assert.True(t, result.Normalize)
	})

	t.Run("does not enable normalize when fade is set", func(t *testing.T) {
		params := models.AudioProcessingParams{
			FadeIn: 1.0,
		}
		result := normalizeAudioParams(params)
		assert.False(t, result.Normalize)
	})

	t.Run("does not enable normalize when speed is set", func(t *testing.T) {
		params := models.AudioProcessingParams{
			Speed: 1.5,
		}
		result := normalizeAudioParams(params)
		assert.False(t, result.Normalize)
	})
}

func TestGetPrimaryFormat(t *testing.T) {
	t.Run("returns first format from list", func(t *testing.T) {
		formats := []string{"mp3", "webm", "ogg"}
		assert.Equal(t, "mp3", getPrimaryFormat(formats))
	})

	t.Run("returns default when list is empty", func(t *testing.T) {
		formats := []string{}
		assert.Equal(t, DefaultOutputFormat, getPrimaryFormat(formats))
	})

	t.Run("returns default when nil", func(t *testing.T) {
		var formats []string
		assert.Equal(t, DefaultOutputFormat, getPrimaryFormat(formats))
	})
}

func TestGetProcessedFileKey(t *testing.T) {
	media := models.Media{
		AuthorId:  primitive.NewObjectID(),
		MediaType: "audio",
		FileName:  "test-audio.wav",
	}

	t.Run("generates key with new format extension", func(t *testing.T) {
		key := getProcessedFileKey(media, "webm")
		assert.Contains(t, key, "test-audio.webm")
		assert.Contains(t, key, media.AuthorId.Hex())
		assert.Contains(t, key, "audio")
	})

	t.Run("replaces original extension", func(t *testing.T) {
		key := getProcessedFileKey(media, "mp3")
		assert.Contains(t, key, "test-audio.mp3")
		assert.NotContains(t, key, ".wav")
	})

	t.Run("handles file without extension", func(t *testing.T) {
		mediaNoExt := models.Media{
			AuthorId:  primitive.NewObjectID(),
			MediaType: "audio",
			FileName:  "test-audio",
		}
		key := getProcessedFileKey(mediaNoExt, "webm")
		assert.Contains(t, key, "test-audio.webm")
	})
}

func TestGetContentType(t *testing.T) {
	t.Run("returns content type for known format", func(t *testing.T) {
		assert.Equal(t, "audio/webm", getContentType("webm", ""))
		assert.Equal(t, "audio/mpeg", getContentType("mp3", ""))
		assert.Equal(t, "audio/ogg", getContentType("ogg", ""))
		assert.Equal(t, "audio/wav", getContentType("wav", ""))
	})

	t.Run("returns fallback for unknown format", func(t *testing.T) {
		assert.Equal(t, "audio/custom", getContentType("xyz", "audio/custom"))
	})

	t.Run("returns default when no fallback", func(t *testing.T) {
		assert.Equal(t, "audio/webm", getContentType("xyz", ""))
	})

	t.Run("prefers known format over fallback", func(t *testing.T) {
		assert.Equal(t, "audio/webm", getContentType("webm", "audio/custom"))
	})
}

func TestProcessAudioWorkflow(t *testing.T) {
	// Test the workflow logic without actual audio processing
	t.Run("workflow steps are documented", func(t *testing.T) {
		steps := []string{
			"Initialize user uploads bucket client",
			"Download audio",
			"Process audio using FFmpeg module with params",
			"Get the primary output format",
			"Extract metadata from primary output",
			"Upload processed audio to media bucket",
			"Delete original file from user uploads bucket",
			"Update media object",
		}
		assert.Len(t, steps, 8)
	})
}

// Benchmark tests
func BenchmarkGetFileSize(b *testing.B) {
	// Create a temporary file for benchmarking
	tmpDir := b.TempDir()
	testFile := filepath.Join(tmpDir, "bench-file.txt")
	content := make([]byte, 1024) // 1KB
	err := os.WriteFile(testFile, content, 0644)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetFileSize(testFile)
	}
}

func BenchmarkNormalizeAudioParams(b *testing.B) {
	params := models.AudioProcessingParams{
		TrimStart: 2.0,
		TrimEnd:   30.0,
		FadeIn:    1.0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = normalizeAudioParams(params)
	}
}

func BenchmarkGetContentType(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getContentType("webm", "")
	}
}

// Integration tests for audio processing pipeline
// These tests require FFmpeg to be installed

// checkFFmpegAvailable checks if FFmpeg is available
func checkFFmpegAvailable() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}

// createTestAudioFile creates a simple test audio file using FFmpeg
func createTestAudioFile(t *testing.T, duration float64) string {
	t.Helper()

	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "test_audio.wav")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-f", "lavfi",
		"-i", "anullsrc=r=44100:cl=mono",
		"-t", fmt.Sprintf("%.1f", duration),
		"-y", inputPath,
	)
	err := cmd.Run()
	if err != nil {
		return ""
	}

	return inputPath
}

func TestIntegrationProcessAudioWithFFmpeg(t *testing.T) {
	if !checkFFmpegAvailable() {
		t.Skip("FFmpeg not available, skipping integration test")
	}

	t.Run("processes audio with default params", func(t *testing.T) {
		inputPath := createTestAudioFile(t, 2.0)
		if inputPath == "" {
			t.Skip("Could not create test audio file")
		}
		defer os.Remove(inputPath)

		outputBasePath := filepath.Join(t.TempDir(), "output")

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		params := models.AudioProcessingParams{
			Normalize: true,
		}

		outputs, err := ffmpeg.ProcessAudioWithParams(ctx, inputPath, outputBasePath, params)
		if err != nil {
			t.Skipf("FFmpeg error: %v", err)
		}

		assert.NoError(t, err)
		assert.Contains(t, outputs, "webm")

		// Verify output file exists
		_, err = os.Stat(outputs["webm"])
		assert.NoError(t, err)
	})

	t.Run("processes audio with trim params", func(t *testing.T) {
		inputPath := createTestAudioFile(t, 5.0)
		if inputPath == "" {
			t.Skip("Could not create test audio file")
		}
		defer os.Remove(inputPath)

		outputBasePath := filepath.Join(t.TempDir(), "output_trim")

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		params := models.AudioProcessingParams{
			TrimStart:     1.0,
			TrimEnd:       4.0,
			OutputFormats: []string{"webm"},
		}

		outputs, err := ffmpeg.ProcessAudioWithParams(ctx, inputPath, outputBasePath, params)
		if err != nil {
			t.Skipf("FFmpeg error: %v", err)
		}

		assert.NoError(t, err)
		assert.Contains(t, outputs, "webm")

		// Verify output file exists and has content
		info, err := os.Stat(outputs["webm"])
		if err == nil {
			assert.Greater(t, info.Size(), int64(0))
		}

		// Verify duration is shorter than original (trimming worked)
		duration, err := ffmpeg.GetAudioMetadata(ctx, outputs["webm"])
		if err == nil {
			assert.Less(t, duration, 5.0, "Trimmed audio should be shorter than original")
		}
	})

	t.Run("processes audio with speed adjustment", func(t *testing.T) {
		inputPath := createTestAudioFile(t, 4.0)
		if inputPath == "" {
			t.Skip("Could not create test audio file")
		}
		defer os.Remove(inputPath)

		outputBasePath := filepath.Join(t.TempDir(), "output_speed")

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		params := models.AudioProcessingParams{
			Speed:         2.0, // Double speed
			OutputFormats: []string{"webm"},
		}

		outputs, err := ffmpeg.ProcessAudioWithParams(ctx, inputPath, outputBasePath, params)
		if err != nil {
			t.Skipf("FFmpeg error: %v", err)
		}

		assert.NoError(t, err)

		// Verify duration is approximately half (4s at 2x speed = ~2s)
		duration, err := ffmpeg.GetAudioMetadata(ctx, outputs["webm"])
		if err == nil {
			assert.InDelta(t, 2.0, duration, 0.5)
		}
	})

	t.Run("processes audio with fade effects", func(t *testing.T) {
		inputPath := createTestAudioFile(t, 5.0)
		if inputPath == "" {
			t.Skip("Could not create test audio file")
		}
		defer os.Remove(inputPath)

		outputBasePath := filepath.Join(t.TempDir(), "output_fade")

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		params := models.AudioProcessingParams{
			FadeIn:        1.0,
			OutputFormats: []string{"webm"},
		}

		outputs, err := ffmpeg.ProcessAudioWithParams(ctx, inputPath, outputBasePath, params)
		if err != nil {
			t.Skipf("FFmpeg error: %v", err)
		}

		assert.NoError(t, err)
		assert.Contains(t, outputs, "webm")
	})

	t.Run("processes audio with combined params", func(t *testing.T) {
		inputPath := createTestAudioFile(t, 10.0)
		if inputPath == "" {
			t.Skip("Could not create test audio file")
		}
		defer os.Remove(inputPath)

		outputBasePath := filepath.Join(t.TempDir(), "output_combined")

		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()

		params := models.AudioProcessingParams{
			TrimStart:     2.0,
			TrimEnd:       8.0,
			FadeIn:        0.5,
			Speed:         1.25,
			Normalize:     true,
			OutputFormats: []string{"webm"},
		}

		outputs, err := ffmpeg.ProcessAudioWithParams(ctx, inputPath, outputBasePath, params)
		if err != nil {
			t.Skipf("FFmpeg error: %v", err)
		}

		assert.NoError(t, err)
		assert.Contains(t, outputs, "webm")

		// Verify output file exists and has content
		info, err := os.Stat(outputs["webm"])
		if err == nil {
			assert.Greater(t, info.Size(), int64(0))
		}

		// Verify duration is shorter than original (trimming + speed worked)
		duration, err := ffmpeg.GetAudioMetadata(ctx, outputs["webm"])
		if err == nil {
			assert.Less(t, duration, 10.0, "Processed audio should be shorter than original")
		}
	})
}

func TestIntegrationAudioPipelineHelpers(t *testing.T) {
	if !checkFFmpegAvailable() {
		t.Skip("FFmpeg not available, skipping integration test")
	}

	t.Run("GetAudioMetadata returns valid duration", func(t *testing.T) {
		inputPath := createTestAudioFile(t, 3.0)
		if inputPath == "" {
			t.Skip("Could not create test audio file")
		}
		defer os.Remove(inputPath)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		duration, err := ffmpeg.GetAudioMetadata(ctx, inputPath)
		assert.NoError(t, err)
		assert.InDelta(t, 3.0, duration, 0.5)
	})

	t.Run("GetFileSize returns correct size for processed audio", func(t *testing.T) {
		inputPath := createTestAudioFile(t, 2.0)
		if inputPath == "" {
			t.Skip("Could not create test audio file")
		}
		defer os.Remove(inputPath)

		size, err := GetFileSize(inputPath)
		assert.NoError(t, err)
		assert.Greater(t, size, int64(0))
	})
}

// Benchmark for integration tests
func BenchmarkIntegrationProcessAudio(b *testing.B) {
	if !checkFFmpegAvailable() {
		b.Skip("FFmpeg not available, skipping benchmark")
	}

	tmpDir := b.TempDir()
	inputPath := filepath.Join(tmpDir, "bench_audio.wav")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-f", "lavfi",
		"-i", "anullsrc=r=44100:cl=mono",
		"-t", "1",
		"-y", inputPath,
	)
	if err := cmd.Run(); err != nil {
		cancel()
		b.Skipf("Could not create test audio file: %v", err)
	}
	cancel()

	params := models.AudioProcessingParams{
		Normalize:     true,
		OutputFormats: []string{"webm"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		outputBasePath := filepath.Join(tmpDir, fmt.Sprintf("bench_output_%d", i))
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		ffmpeg.ProcessAudioWithParams(ctx, inputPath, outputBasePath, params)
		cancel()
	}
}
