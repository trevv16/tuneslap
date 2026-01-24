package ffmpeg

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// checkFFmpegAvailable checks if FFmpeg is available
func checkFFmpegAvailable() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}

// checkFFprobeAvailable checks if FFprobe is available
func checkFFprobeAvailable() bool {
	_, err := exec.LookPath("ffprobe")
	return err == nil
}

func TestNormalizeAudio(t *testing.T) {
	if !checkFFmpegAvailable() {
		t.Skip("FFmpeg not available, skipping test")
	}

	// Create a simple test audio file using FFmpeg
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "test_input.wav")
	outputBasePath := filepath.Join(tempDir, "output")

	// Generate a simple test audio (1 second of silence)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-f", "lavfi",
		"-i", "anullsrc=r=44100:cl=mono",
		"-t", "1",
		"-y", inputPath,
	)
	err := cmd.Run()
	if err != nil {
		t.Skipf("Could not create test audio file: %v", err)
	}

	t.Run("normalize valid audio", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		webmPath, err := NormalizeAudio(ctx, inputPath, outputBasePath)
		assert.NoError(t, err)
		assert.Equal(t, outputBasePath+".webm", webmPath)

		// Check that output files exist
		_, err = os.Stat(outputBasePath + ".mp3")
		assert.NoError(t, err)
		_, err = os.Stat(outputBasePath + ".webm")
		assert.NoError(t, err)
	})

	t.Run("invalid input path", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := NormalizeAudio(ctx, "/nonexistent/file.wav", outputBasePath+"_invalid")
		assert.Error(t, err)
	})
}

func TestGetAudioMetadata(t *testing.T) {
	if !checkFFprobeAvailable() {
		t.Skip("FFprobe not available, skipping test")
	}

	// Create a simple test audio file using FFmpeg
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "test_input.wav")

	// Generate a simple test audio (2 seconds of silence)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-f", "lavfi",
		"-i", "anullsrc=r=44100:cl=mono",
		"-t", "2",
		"-y", inputPath,
	)
	err := cmd.Run()
	if err != nil {
		t.Skipf("Could not create test audio file: %v", err)
	}

	t.Run("get metadata from valid file", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		duration, err := GetAudioMetadata(ctx, inputPath)
		assert.NoError(t, err)
		// Duration should be approximately 2 seconds
		assert.InDelta(t, 2.0, duration, 0.5)
	})

	t.Run("invalid file path", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := GetAudioMetadata(ctx, "/nonexistent/file.wav")
		assert.Error(t, err)
	})
}

func TestNormalizeAudio_ContextCancellation(t *testing.T) {
	if !checkFFmpegAvailable() {
		t.Skip("FFmpeg not available, skipping test")
	}

	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "test_input.wav")
	outputBasePath := filepath.Join(tempDir, "output")

	// Generate a simple test audio
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-f", "lavfi",
		"-i", "anullsrc=r=44100:cl=mono",
		"-t", "1",
		"-y", inputPath,
	)
	err := cmd.Run()
	cancel()
	if err != nil {
		t.Skipf("Could not create test audio file: %v", err)
	}

	t.Run("context cancelled immediately", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := NormalizeAudio(ctx, inputPath, outputBasePath)
		assert.Error(t, err)
	})
}

// Benchmark tests (only run if FFmpeg is available)
func BenchmarkGetAudioMetadata(b *testing.B) {
	if !checkFFprobeAvailable() {
		b.Skip("FFprobe not available, skipping benchmark")
	}

	// Create a simple test audio file
	tempDir := b.TempDir()
	inputPath := filepath.Join(tempDir, "test_input.wav")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-f", "lavfi",
		"-i", "anullsrc=r=44100:cl=mono",
		"-t", "1",
		"-y", inputPath,
	)
	if err := cmd.Run(); err != nil {
		b.Skipf("Could not create test audio file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		GetAudioMetadata(ctx, inputPath)
		cancel()
	}
}
