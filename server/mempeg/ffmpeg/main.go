package ffmpeg

import (
	"context"
	"fmt"
	"os/exec"
)

// NormalizeAudio applies loudness normalization, resampling, trimming silence, and re-encodes to MP3 and WebM.
func NormalizeAudio(ctx context.Context, inputPath, outputBasePath string) (webmPath string, err error) {
	// Apply loudness normalization, resample, trim silence, and export as MP3
	mp3Path := fmt.Sprintf("%s.mp3", outputBasePath)
	cmdMp3 := exec.CommandContext(ctx, "ffmpeg",
		"-i", inputPath,
		"-af", "loudnorm=I=-16:TP=-1.5:LRA=11,silenceremove=start_periods=1:start_duration=0.5:start_threshold=-50dB",
		"-ar", "44100",
		"-c:a", "libmp3lame",
		"-b:a", "128k",
		"-y", mp3Path,
	)
	if out, err := cmdMp3.CombinedOutput(); err != nil {
		return "", fmt.Errorf("mp3 normalization error: %w\n%s", err, string(out))
	}

	// Apply same processing but encode as WebM (Opus)
	webmPath = fmt.Sprintf("%s.webm", outputBasePath)
	cmdWebm := exec.CommandContext(ctx, "ffmpeg",
		"-i", inputPath,
		"-af", "loudnorm=I=-16:TP=-1.5:LRA=11,silenceremove=start_periods=1:start_duration=0.5:start_threshold=-50dB",
		"-ar", "44100",
		"-c:a", "libopus",
		"-b:a", "128k",
		"-y", webmPath,
	)
	if out, err := cmdWebm.CombinedOutput(); err != nil {
		return "", fmt.Errorf("webm normalization error: %w\n%s", err, string(out))
	}

	return webmPath, nil
}

// GetAudioMetadata uses ffprobe to retrieve duration and codec details
func GetAudioMetadata(ctx context.Context, inputPath string) (duration float64, err error) {
	cmd := exec.CommandContext(ctx, "ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		inputPath,
	)
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("metadata extraction error: %w", err)
	}
	fmt.Sscanf(string(output), "%f", &duration)

	return duration, nil
}
