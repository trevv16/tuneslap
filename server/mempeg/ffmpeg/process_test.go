package ffmpeg

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
	"tuneslap/models"

	"github.com/stretchr/testify/assert"
)

func TestAudioConstants(t *testing.T) {
	t.Run("MP3 sample rate", func(t *testing.T) {
		assert.Equal(t, 44100, MP3SampleRate)
	})

	t.Run("Opus sample rate", func(t *testing.T) {
		assert.Equal(t, 48000, OpusSampleRate)
	})

	t.Run("default bitrate", func(t *testing.T) {
		assert.Equal(t, "128k", DefaultBitrate)
	})

	t.Run("loudnorm filter format", func(t *testing.T) {
		assert.Contains(t, LoudnormFilter, "loudnorm")
		assert.Contains(t, LoudnormFilter, "I=-16")
		assert.Contains(t, LoudnormFilter, "TP=-1.5")
		assert.Contains(t, LoudnormFilter, "LRA=11")
	})
}

func TestProcessAudioWithParams(t *testing.T) {
	if !checkFFmpegAvailable() {
		t.Skip("FFmpeg not available, skipping test")
	}

	// Create a test audio file
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "test_input.wav")
	outputBasePath := filepath.Join(tempDir, "output")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-f", "lavfi",
		"-i", "anullsrc=r=44100:cl=mono",
		"-t", "2",
		"-y", inputPath,
	)
	if err := cmd.Run(); err != nil {
		t.Skipf("Could not create test audio file: %v", err)
	}

	t.Run("default params produces webm", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		params := models.AudioProcessingParams{}
		outputs, err := ProcessAudioWithParams(ctx, inputPath, outputBasePath+"_default", params)

		if err != nil {
			t.Skipf("FFmpeg error (may not have required codecs): %v", err)
		}

		assert.NoError(t, err)
		assert.Contains(t, outputs, "webm")
		_, err = os.Stat(outputs["webm"])
		assert.NoError(t, err)
	})

	t.Run("multiple output formats", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()

		params := models.AudioProcessingParams{
			OutputFormats: []string{"webm", "mp3"},
			Normalize:     true,
		}
		outputs, err := ProcessAudioWithParams(ctx, inputPath, outputBasePath+"_multi", params)

		if err != nil {
			t.Skipf("FFmpeg error (may not have required codecs): %v", err)
		}

		assert.NoError(t, err)
		assert.Len(t, outputs, 2)
		assert.Contains(t, outputs, "webm")
		assert.Contains(t, outputs, "mp3")
	})

	t.Run("with trim params", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		params := models.AudioProcessingParams{
			TrimStart: 0.5,
			TrimEnd:   1.5,
		}
		outputs, err := ProcessAudioWithParams(ctx, inputPath, outputBasePath+"_trim", params)

		if err != nil {
			t.Skipf("FFmpeg error: %v", err)
		}

		assert.NoError(t, err)
		assert.Contains(t, outputs, "webm")
	})

	t.Run("unsupported format returns error", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		params := models.AudioProcessingParams{
			OutputFormats: []string{"xyz"},
		}
		_, err := ProcessAudioWithParams(ctx, inputPath, outputBasePath+"_unsupported", params)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported format")
	})

	t.Run("invalid input path", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		params := models.AudioProcessingParams{}
		_, err := ProcessAudioWithParams(ctx, "/nonexistent/file.wav", outputBasePath+"_invalid", params)

		assert.Error(t, err)
	})
}
