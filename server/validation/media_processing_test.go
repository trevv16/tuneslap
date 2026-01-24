package validation

import (
	"testing"

	api "tuneslap/generated"

	"github.com/stretchr/testify/assert"
)

func TestValidateMediaProcessingParamsAudio(t *testing.T) {
	floatPtr := func(v float32) *float32 { return &v }

	tests := []struct {
		name      string
		params    *api.MediaProcessingParamsAudio
		expectErr bool
	}{
		{"nil params", nil, false},
		{"empty params", &api.MediaProcessingParamsAudio{}, false},
		{"valid params", &api.MediaProcessingParamsAudio{
			TrimStart: floatPtr(0.5), TrimEnd: floatPtr(10.0), FadeIn: floatPtr(0.2), FadeOut: floatPtr(0.3), Speed: floatPtr(1.25),
		}, false},
		{"negative trimStart", &api.MediaProcessingParamsAudio{TrimStart: floatPtr(-1.0)}, true},
		{"negative trimEnd", &api.MediaProcessingParamsAudio{TrimEnd: floatPtr(-1.0)}, true},
		{"negative fadeIn", &api.MediaProcessingParamsAudio{FadeIn: floatPtr(-0.5)}, true},
		{"zero speed", &api.MediaProcessingParamsAudio{Speed: floatPtr(0)}, true},
		{"negative speed", &api.MediaProcessingParamsAudio{Speed: floatPtr(-1.0)}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateMediaProcessingParamsAudio(tt.params)
			if tt.expectErr {
				assert.False(t, result.IsValid)
				assert.NotEmpty(t, result.Errors)
			} else {
				assert.True(t, result.IsValid)
			}
		})
	}
}
