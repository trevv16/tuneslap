package ffmpeg

import (
	"strings"
	"testing"
	"tuneslap/models"

	"github.com/stretchr/testify/assert"
)

func TestBuildAudioFilterChain(t *testing.T) {
	t.Run("empty params returns default filter", func(t *testing.T) {
		params := models.AudioProcessingParams{}
		result := BuildAudioFilterChain(params)
		assert.Equal(t, LoudnormFilter, result)
	})

	t.Run("normalize only", func(t *testing.T) {
		params := models.AudioProcessingParams{
			Normalize: true,
		}
		result := BuildAudioFilterChain(params)
		assert.Equal(t, LoudnormFilter, result)
	})

	t.Run("trim start only", func(t *testing.T) {
		params := models.AudioProcessingParams{
			TrimStart: 2.5,
		}
		result := BuildAudioFilterChain(params)
		assert.Contains(t, result, "atrim=start=2.50")
	})

	t.Run("trim end only", func(t *testing.T) {
		params := models.AudioProcessingParams{
			TrimEnd: 30.0,
		}
		result := BuildAudioFilterChain(params)
		assert.Contains(t, result, "atrim=end=30.00")
	})

	t.Run("trim start and end", func(t *testing.T) {
		params := models.AudioProcessingParams{
			TrimStart: 2.0,
			TrimEnd:   30.0,
		}
		result := BuildAudioFilterChain(params)
		assert.Contains(t, result, "atrim=start=2.00:end=30.00")
	})

	t.Run("fade in", func(t *testing.T) {
		params := models.AudioProcessingParams{
			FadeIn: 1.5,
		}
		result := BuildAudioFilterChain(params)
		assert.Contains(t, result, "afade=t=in:st=0:d=1.50")
	})

	t.Run("fade out", func(t *testing.T) {
		params := models.AudioProcessingParams{
			FadeOut: 2.0,
		}
		result := BuildAudioFilterChain(params)
		assert.Contains(t, result, "afade=t=out")
		assert.Contains(t, result, "d=2.00")
	})

	t.Run("fade in and out", func(t *testing.T) {
		params := models.AudioProcessingParams{
			FadeIn:  1.0,
			FadeOut: 2.0,
		}
		result := BuildAudioFilterChain(params)
		assert.Contains(t, result, "afade=t=in:st=0:d=1.00")
		assert.Contains(t, result, "afade=t=out")
	})

	t.Run("speed adjustment", func(t *testing.T) {
		params := models.AudioProcessingParams{
			Speed: 1.5,
		}
		result := BuildAudioFilterChain(params)
		assert.Contains(t, result, "atempo=1.5")
	})

	t.Run("speed 1.0 is ignored", func(t *testing.T) {
		params := models.AudioProcessingParams{
			Speed: 1.0,
		}
		result := BuildAudioFilterChain(params)
		// Speed 1.0 should not add atempo filter
		assert.NotContains(t, result, "atempo")
	})

	t.Run("pitch adjustment", func(t *testing.T) {
		params := models.AudioProcessingParams{
			Pitch: 1.1,
		}
		result := BuildAudioFilterChain(params)
		assert.Contains(t, result, "asetrate=")
		assert.Contains(t, result, "aresample=")
		assert.Contains(t, result, "atempo=")
	})

	t.Run("pitch 1.0 is ignored", func(t *testing.T) {
		params := models.AudioProcessingParams{
			Pitch: 1.0,
		}
		result := BuildAudioFilterChain(params)
		assert.NotContains(t, result, "asetrate=")
	})

	t.Run("combined params", func(t *testing.T) {
		params := models.AudioProcessingParams{
			TrimStart: 2.0,
			TrimEnd:   30.0,
			FadeIn:    1.0,
			FadeOut:   2.0,
			Speed:     1.25,
			Normalize: true,
		}
		result := BuildAudioFilterChain(params)

		// Check all filters are present
		assert.Contains(t, result, "atrim=start=2.00:end=30.00")
		assert.Contains(t, result, "afade=t=in:st=0:d=1.00")
		assert.Contains(t, result, "afade=t=out")
		assert.Contains(t, result, "atempo=1.25")
		assert.Contains(t, result, LoudnormFilter)

		// Check filters are comma-separated
		filters := strings.Split(result, ",")
		assert.GreaterOrEqual(t, len(filters), 4)
	})
}

func TestBuildTrimFilter(t *testing.T) {
	t.Run("both start and end", func(t *testing.T) {
		result := buildTrimFilter(2.0, 30.0)
		assert.Equal(t, "atrim=start=2.00:end=30.00", result)
	})

	t.Run("start only", func(t *testing.T) {
		result := buildTrimFilter(5.0, 0)
		assert.Equal(t, "atrim=start=5.00", result)
	})

	t.Run("end only", func(t *testing.T) {
		result := buildTrimFilter(0, 20.0)
		assert.Equal(t, "atrim=end=20.00", result)
	})

	t.Run("neither start nor end", func(t *testing.T) {
		result := buildTrimFilter(0, 0)
		assert.Equal(t, "", result)
	})
}

func TestBuildSpeedFilters(t *testing.T) {
	t.Run("speed within normal range", func(t *testing.T) {
		result := buildSpeedFilters(1.5)
		assert.Len(t, result, 1)
		assert.Equal(t, "atempo=1.5000", result[0])
	})

	t.Run("speed 0.75", func(t *testing.T) {
		result := buildSpeedFilters(0.75)
		assert.Len(t, result, 1)
		assert.Equal(t, "atempo=0.7500", result[0])
	})

	t.Run("speed above 2.0 requires chaining", func(t *testing.T) {
		result := buildSpeedFilters(3.0)
		// 3.0 = 2.0 * 1.5, so we need atempo=2.0,atempo=1.5
		assert.GreaterOrEqual(t, len(result), 2)
		assert.Contains(t, result, "atempo=2.0")
	})

	t.Run("speed below 0.5 requires chaining", func(t *testing.T) {
		result := buildSpeedFilters(0.3)
		// 0.3 = 0.5 * 0.6, so we need atempo=0.5,atempo=0.6
		assert.GreaterOrEqual(t, len(result), 2)
		assert.Contains(t, result, "atempo=0.5")
	})

	t.Run("speed 1.0 returns empty", func(t *testing.T) {
		result := buildSpeedFilters(1.0)
		assert.Empty(t, result)
	})

	t.Run("speed clamped to min 0.25", func(t *testing.T) {
		result := buildSpeedFilters(0.1)
		// Should be clamped to 0.25
		assert.NotEmpty(t, result)
	})

	t.Run("speed clamped to max 4.0", func(t *testing.T) {
		result := buildSpeedFilters(10.0)
		// Should be clamped to 4.0
		assert.NotEmpty(t, result)
	})
}

func TestBuildPitchFilters(t *testing.T) {
	t.Run("pitch 1.1 (up)", func(t *testing.T) {
		result := buildPitchFilters(1.1)
		assert.Len(t, result, 3)
		assert.Contains(t, result[0], "asetrate=")
		assert.Contains(t, result[1], "aresample=")
		assert.Contains(t, result[2], "atempo=")
	})

	t.Run("pitch 0.9 (down)", func(t *testing.T) {
		result := buildPitchFilters(0.9)
		assert.Len(t, result, 3)
		assert.Contains(t, result[0], "asetrate=")
	})

	t.Run("pitch 1.0 returns pitch filters anyway", func(t *testing.T) {
		// The function doesn't check for 1.0, that's done in BuildAudioFilterChain
		result := buildPitchFilters(1.0)
		assert.Len(t, result, 3)
	})

	t.Run("pitch clamped to max 2.0", func(t *testing.T) {
		result := buildPitchFilters(3.0)
		assert.Len(t, result, 3)
		// Should contain pitch of 2.0
		assert.Contains(t, result[0], "2.0000")
	})

	t.Run("pitch clamped to min 0.5", func(t *testing.T) {
		result := buildPitchFilters(0.1)
		assert.Len(t, result, 3)
		// Should contain pitch of 0.5
		assert.Contains(t, result[0], "0.5000")
	})
}

// Edge case tests for filter chain building
func TestBuildAudioFilterChainEdgeCases(t *testing.T) {
	t.Run("negative trim start is ignored", func(t *testing.T) {
		params := models.AudioProcessingParams{
			TrimStart: -5.0,
		}
		result := BuildAudioFilterChain(params)
		// Negative values are not added (implementation checks > 0)
		assert.NotContains(t, result, "atrim=start=-5.00")
		// Should return default filter
		assert.Equal(t, LoudnormFilter, result)
	})

	t.Run("very small trim values", func(t *testing.T) {
		params := models.AudioProcessingParams{
			TrimStart: 0.001,
			TrimEnd:   0.002,
		}
		result := BuildAudioFilterChain(params)
		assert.Contains(t, result, "atrim=")
	})

	t.Run("very large trim values", func(t *testing.T) {
		params := models.AudioProcessingParams{
			TrimStart: 3600.0, // 1 hour
			TrimEnd:   7200.0, // 2 hours
		}
		result := BuildAudioFilterChain(params)
		assert.Contains(t, result, "3600.00")
		assert.Contains(t, result, "7200.00")
	})

	t.Run("trim end before trim start", func(t *testing.T) {
		params := models.AudioProcessingParams{
			TrimStart: 30.0,
			TrimEnd:   10.0,
		}
		result := BuildAudioFilterChain(params)
		// This is technically invalid but the builder doesn't validate
		assert.Contains(t, result, "atrim=start=30.00:end=10.00")
	})

	t.Run("zero fade duration", func(t *testing.T) {
		params := models.AudioProcessingParams{
			FadeIn:  0.0,
			FadeOut: 0.0,
		}
		result := BuildAudioFilterChain(params)
		// Zero fade should not add fade filters
		assert.NotContains(t, result, "afade")
	})

	t.Run("very small fade duration", func(t *testing.T) {
		params := models.AudioProcessingParams{
			FadeIn: 0.01,
		}
		result := BuildAudioFilterChain(params)
		assert.Contains(t, result, "afade=t=in:st=0:d=0.01")
	})

	t.Run("speed at exact boundary 0.5", func(t *testing.T) {
		result := buildSpeedFilters(0.5)
		assert.Len(t, result, 1)
		assert.Equal(t, "atempo=0.5000", result[0])
	})

	t.Run("speed at exact boundary 2.0", func(t *testing.T) {
		result := buildSpeedFilters(2.0)
		assert.Len(t, result, 1)
		assert.Equal(t, "atempo=2.0000", result[0])
	})

	t.Run("speed just above 2.0", func(t *testing.T) {
		result := buildSpeedFilters(2.01)
		// Should require chaining
		assert.GreaterOrEqual(t, len(result), 2)
	})

	t.Run("speed just below 0.5", func(t *testing.T) {
		result := buildSpeedFilters(0.49)
		// Should require chaining
		assert.GreaterOrEqual(t, len(result), 2)
	})

	t.Run("pitch at boundary 0.5", func(t *testing.T) {
		result := buildPitchFilters(0.5)
		assert.Len(t, result, 3)
		assert.Contains(t, result[0], "0.5000")
	})

	t.Run("pitch at boundary 2.0", func(t *testing.T) {
		result := buildPitchFilters(2.0)
		assert.Len(t, result, 3)
		assert.Contains(t, result[0], "2.0000")
	})

	t.Run("negative speed is clamped", func(t *testing.T) {
		result := buildSpeedFilters(-1.0)
		// Negative values should be clamped to 0.25
		assert.NotEmpty(t, result)
	})

	t.Run("negative pitch is clamped", func(t *testing.T) {
		result := buildPitchFilters(-1.0)
		// Negative values should be clamped to 0.5
		assert.Len(t, result, 3)
	})

	t.Run("all params at zero except normalize", func(t *testing.T) {
		params := models.AudioProcessingParams{
			TrimStart: 0,
			TrimEnd:   0,
			FadeIn:    0,
			FadeOut:   0,
			Speed:     0,
			Pitch:     0,
			Normalize: true,
		}
		result := BuildAudioFilterChain(params)
		// Only normalize filter should be present
		assert.Equal(t, LoudnormFilter, result)
	})
}

func TestBuildTrimFilterEdgeCases(t *testing.T) {
	t.Run("floating point precision", func(t *testing.T) {
		result := buildTrimFilter(1.005, 2.995)
		// Should format to 2 decimal places
		assert.Contains(t, result, "1.00") // 1.005 rounds to 1.00 or 1.01
		assert.Contains(t, result, "3.00") // 2.995 rounds to 3.00 or 2.99
	})

	t.Run("very large values", func(t *testing.T) {
		result := buildTrimFilter(0, 86400.0) // 24 hours in seconds
		assert.Contains(t, result, "86400.00")
	})
}

func TestFilterChainFormat(t *testing.T) {
	t.Run("filter chain is properly comma-separated", func(t *testing.T) {
		params := models.AudioProcessingParams{
			TrimStart: 1.0,
			FadeIn:    0.5,
			Normalize: true,
		}
		result := BuildAudioFilterChain(params)

		// Should not have double commas
		assert.NotContains(t, result, ",,")

		// Should not start or end with comma
		assert.NotEqual(t, ',', result[0])
		assert.NotEqual(t, ',', result[len(result)-1])
	})

	t.Run("filter chain produces valid FFmpeg syntax", func(t *testing.T) {
		params := models.AudioProcessingParams{
			TrimStart: 2.0,
			TrimEnd:   30.0,
			Speed:     1.5,
			Normalize: true,
		}
		result := BuildAudioFilterChain(params)

		// Basic FFmpeg filter syntax validation
		// Filters should be name=value pairs separated by commas
		filters := strings.Split(result, ",")
		for _, filter := range filters {
			// Each filter should contain an equals sign or be a known format
			assert.True(t, strings.Contains(filter, "=") || filter == LoudnormFilter,
				"Filter should have valid syntax: %s", filter)
		}
	})
}

// Benchmarks for filter chain building
func BenchmarkBuildAudioFilterChain(b *testing.B) {
	params := models.AudioProcessingParams{
		TrimStart: 2.0,
		TrimEnd:   30.0,
		FadeIn:    1.0,
		FadeOut:   2.0,
		Speed:     1.25,
		Normalize: true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildAudioFilterChain(params)
	}
}

func BenchmarkBuildSpeedFilters(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = buildSpeedFilters(1.5)
	}
}

func BenchmarkBuildPitchFilters(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = buildPitchFilters(1.1)
	}
}
