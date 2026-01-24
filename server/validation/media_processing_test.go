package validation

import (
	"testing"

	api "tuneslap/generated"

	"github.com/stretchr/testify/assert"
)

func TestValidateMediaProcessingParamsAudio(t *testing.T) {
	tests := []struct {
		name      string
		params    *api.MediaProcessingParamsAudio
		expectErr bool
		errFields []string
	}{
		{
			name:      "nil params",
			params:    nil,
			expectErr: false,
		},
		{
			name:      "empty params",
			params:    &api.MediaProcessingParamsAudio{},
			expectErr: false,
		},
		{
			name: "valid params",
			params: func() *api.MediaProcessingParamsAudio {
				trimStart := float32(0.5)
				trimEnd := float32(10.0)
				fadeIn := float32(0.2)
				fadeOut := float32(0.3)
				speed := float32(1.25)
				return &api.MediaProcessingParamsAudio{
					TrimStart: &trimStart,
					TrimEnd:   &trimEnd,
					FadeIn:    &fadeIn,
					FadeOut:   &fadeOut,
					Speed:     &speed,
				}
			}(),
			expectErr: false,
		},
		{
			name: "negative trimStart",
			params: func() *api.MediaProcessingParamsAudio {
				trimStart := float32(-1.0)
				return &api.MediaProcessingParamsAudio{
					TrimStart: &trimStart,
				}
			}(),
			expectErr: true,
			errFields: []string{"trimStart"},
		},
		{
			name: "negative trimEnd",
			params: func() *api.MediaProcessingParamsAudio {
				trimEnd := float32(-1.0)
				return &api.MediaProcessingParamsAudio{
					TrimEnd: &trimEnd,
				}
			}(),
			expectErr: true,
			errFields: []string{"trimEnd"},
		},
		{
			name: "negative fadeIn",
			params: func() *api.MediaProcessingParamsAudio {
				fadeIn := float32(-0.5)
				return &api.MediaProcessingParamsAudio{
					FadeIn: &fadeIn,
				}
			}(),
			expectErr: true,
			errFields: []string{"fadeIn"},
		},
		{
			name: "zero speed",
			params: func() *api.MediaProcessingParamsAudio {
				speed := float32(0)
				return &api.MediaProcessingParamsAudio{
					Speed: &speed,
				}
			}(),
			expectErr: true,
			errFields: []string{"speed"},
		},
		{
			name: "negative speed",
			params: func() *api.MediaProcessingParamsAudio {
				speed := float32(-1.0)
				return &api.MediaProcessingParamsAudio{
					Speed: &speed,
				}
			}(),
			expectErr: true,
			errFields: []string{"speed"},
		},
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
