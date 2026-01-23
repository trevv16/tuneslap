package ffmpeg

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"tuneslap/models"
)

// Audio codec and sample rate constants
const (
	MP3SampleRate  = 44100
	OpusSampleRate = 48000
	DefaultBitrate = "128k"
)

// LoudnormFilter is the default loudness normalization filter
const LoudnormFilter = "loudnorm=I=-16:TP=-1.5:LRA=11"

// SilenceRemoveFilter removes leading silence
const SilenceRemoveFilter = "silenceremove=start_periods=1:start_duration=0.5:start_threshold=-50dB"

// NormalizeAudio applies loudness normalization, resampling, trimming silence, and re-encodes to MP3 and WebM.
func NormalizeAudio(ctx context.Context, inputPath, outputBasePath string) (webmPath string, err error) {
	// Apply loudness normalization, resample, trim silence, and export as MP3
	// MP3 supports 44100 Hz sample rate
	mp3Path := fmt.Sprintf("%s.mp3", outputBasePath)
	cmdMp3 := exec.CommandContext(ctx, "ffmpeg",
		"-i", inputPath,
		"-af", fmt.Sprintf("%s,%s", LoudnormFilter, SilenceRemoveFilter),
		"-ar", fmt.Sprintf("%d", MP3SampleRate),
		"-c:a", "libmp3lame",
		"-b:a", DefaultBitrate,
		"-y", mp3Path,
	)
	if out, err := cmdMp3.CombinedOutput(); err != nil {
		return "", fmt.Errorf("mp3 normalization error: %w\n%s", err, string(out))
	}

	// Apply same processing but encode as WebM (Opus)
	// Opus only supports: 8000, 12000, 16000, 24000, 48000 Hz - using 48000 for best quality
	webmPath = fmt.Sprintf("%s.webm", outputBasePath)
	cmdWebm := exec.CommandContext(ctx, "ffmpeg",
		"-i", inputPath,
		"-af", fmt.Sprintf("%s,%s", LoudnormFilter, SilenceRemoveFilter),
		"-ar", fmt.Sprintf("%d", OpusSampleRate),
		"-c:a", "libopus",
		"-b:a", DefaultBitrate,
		"-y", webmPath,
	)
	if out, err := cmdWebm.CombinedOutput(); err != nil {
		return "", fmt.Errorf("webm normalization error: %w\n%s", err, string(out))
	}

	return webmPath, nil
}

// BuildAudioFilterChain builds an FFmpeg audio filter chain from processing params
// Returns the filter chain string for use with the -af flag
func BuildAudioFilterChain(params models.AudioProcessingParams) string {
	var filters []string

	// Trim filter (atrim)
	if params.TrimStart > 0 || params.TrimEnd > 0 {
		trimFilter := buildTrimFilter(params.TrimStart, params.TrimEnd)
		if trimFilter != "" {
			filters = append(filters, trimFilter)
		}
	}

	// Fade in filter
	if params.FadeIn > 0 {
		filters = append(filters, fmt.Sprintf("afade=t=in:st=0:d=%.2f", params.FadeIn))
	}

	// Fade out filter (requires knowing the duration, so we use -1 to let ffmpeg calculate)
	// This is typically applied after trim, so we use the params to calculate
	if params.FadeOut > 0 {
		// The fade out start time needs to be calculated based on content duration
		// We'll use a placeholder that will be adjusted at processing time
		filters = append(filters, fmt.Sprintf("afade=t=out:st=-1:d=%.2f", params.FadeOut))
	}

	// Speed filter (atempo)
	if params.Speed > 0 && params.Speed != 1.0 {
		speedFilters := buildSpeedFilters(params.Speed)
		filters = append(filters, speedFilters...)
	}

	// Pitch filter (asetrate + aresample + atempo)
	if params.Pitch > 0 && params.Pitch != 1.0 {
		pitchFilters := buildPitchFilters(params.Pitch)
		filters = append(filters, pitchFilters...)
	}

	// Loudness normalization
	if params.Normalize {
		filters = append(filters, LoudnormFilter)
	}

	// If no filters specified, return default normalization
	if len(filters) == 0 {
		return LoudnormFilter
	}

	return strings.Join(filters, ",")
}

// buildTrimFilter builds the atrim filter for start/end trimming
func buildTrimFilter(trimStart, trimEnd float64) string {
	if trimStart > 0 && trimEnd > 0 {
		return fmt.Sprintf("atrim=start=%.2f:end=%.2f", trimStart, trimEnd)
	} else if trimStart > 0 {
		return fmt.Sprintf("atrim=start=%.2f", trimStart)
	} else if trimEnd > 0 {
		return fmt.Sprintf("atrim=end=%.2f", trimEnd)
	}
	return ""
}

// buildSpeedFilters builds atempo filters for speed adjustment
// atempo only supports values between 0.5 and 2.0, so we chain multiple filters for wider range
func buildSpeedFilters(speed float64) []string {
	var filters []string

	// Clamp speed to reasonable range
	if speed < 0.25 {
		speed = 0.25
	} else if speed > 4.0 {
		speed = 4.0
	}

	// atempo filter range: 0.5 to 2.0
	// For values outside this range, chain multiple filters
	remaining := speed
	for remaining < 0.5 || remaining > 2.0 {
		if remaining < 0.5 {
			filters = append(filters, "atempo=0.5")
			remaining /= 0.5
		} else if remaining > 2.0 {
			filters = append(filters, "atempo=2.0")
			remaining /= 2.0
		}
	}

	// Add the final atempo filter for the remaining value
	if remaining != 1.0 {
		filters = append(filters, fmt.Sprintf("atempo=%.4f", remaining))
	}

	return filters
}

// buildPitchFilters builds filters for pitch shifting without changing speed
// Uses asetrate to change pitch, then aresample to restore sample rate,
// and atempo to restore duration
func buildPitchFilters(pitch float64) []string {
	// Clamp pitch to reasonable range (0.5 to 2.0)
	if pitch < 0.5 {
		pitch = 0.5
	} else if pitch > 2.0 {
		pitch = 2.0
	}

	// Calculate the inverse tempo adjustment to maintain duration
	inverseTempo := 1.0 / pitch

	return []string{
		fmt.Sprintf("asetrate=%d*%.4f", OpusSampleRate, pitch),
		fmt.Sprintf("aresample=%d", OpusSampleRate),
		fmt.Sprintf("atempo=%.4f", inverseTempo),
	}
}

// ProcessAudioWithParams processes audio using the specified params
// Returns a map of format -> output path
func ProcessAudioWithParams(ctx context.Context, inputPath, outputBasePath string, params models.AudioProcessingParams) (outputs map[string]string, err error) {
	outputs = make(map[string]string)

	// Build the filter chain
	filterChain := BuildAudioFilterChain(params)

	// Determine output formats
	formats := params.OutputFormats
	if len(formats) == 0 {
		// Default to webm if no formats specified
		formats = []string{"webm"}
	}

	// Process each output format
	for _, format := range formats {
		outputPath, err := processFormat(ctx, inputPath, outputBasePath, format, filterChain)
		if err != nil {
			return nil, fmt.Errorf("failed to process %s: %w", format, err)
		}
		outputs[format] = outputPath
	}

	return outputs, nil
}

// processFormat processes audio to a specific format
func processFormat(ctx context.Context, inputPath, outputBasePath, format, filterChain string) (string, error) {
	outputPath := fmt.Sprintf("%s.%s", outputBasePath, format)

	var cmd *exec.Cmd
	switch format {
	case "mp3":
		cmd = exec.CommandContext(ctx, "ffmpeg",
			"-i", inputPath,
			"-af", filterChain,
			"-ar", fmt.Sprintf("%d", MP3SampleRate),
			"-c:a", "libmp3lame",
			"-b:a", DefaultBitrate,
			"-y", outputPath,
		)
	case "webm":
		cmd = exec.CommandContext(ctx, "ffmpeg",
			"-i", inputPath,
			"-af", filterChain,
			"-ar", fmt.Sprintf("%d", OpusSampleRate),
			"-c:a", "libopus",
			"-b:a", DefaultBitrate,
			"-y", outputPath,
		)
	case "ogg":
		cmd = exec.CommandContext(ctx, "ffmpeg",
			"-i", inputPath,
			"-af", filterChain,
			"-ar", fmt.Sprintf("%d", OpusSampleRate),
			"-c:a", "libvorbis",
			"-b:a", DefaultBitrate,
			"-y", outputPath,
		)
	case "wav":
		cmd = exec.CommandContext(ctx, "ffmpeg",
			"-i", inputPath,
			"-af", filterChain,
			"-ar", fmt.Sprintf("%d", OpusSampleRate),
			"-c:a", "pcm_s16le",
			"-y", outputPath,
		)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}

	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("ffmpeg error: %w\n%s", err, string(out))
	}

	return outputPath, nil
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
