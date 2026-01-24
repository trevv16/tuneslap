package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateIntRange(t *testing.T) {
	tests := []struct {
		name      string
		value     *int32
		min       int64
		max       int64
		expectErr bool
	}{
		{
			name:      "valid within range",
			value:     int32Ptr(50),
			min:       1,
			max:       100,
			expectErr: false,
		},
		{
			name:      "at minimum",
			value:     int32Ptr(1),
			min:       1,
			max:       100,
			expectErr: false,
		},
		{
			name:      "at maximum",
			value:     int32Ptr(100),
			min:       1,
			max:       100,
			expectErr: false,
		},
		{
			name:      "below minimum",
			value:     int32Ptr(0),
			min:       1,
			max:       100,
			expectErr: true,
		},
		{
			name:      "above maximum",
			value:     int32Ptr(101),
			min:       1,
			max:       100,
			expectErr: true,
		},
		{
			name:      "nil value",
			value:     nil,
			min:       1,
			max:       100,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateIntRange("field", tt.value, tt.min, tt.max)
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestValidateNonNegativeFloat(t *testing.T) {
	tests := []struct {
		name      string
		value     *float32
		expectErr bool
	}{
		{
			name:      "nil value",
			value:     nil,
			expectErr: false,
		},
		{
			name:      "positive value",
			value:     float32Ptr(10.5),
			expectErr: false,
		},
		{
			name:      "zero value",
			value:     float32Ptr(0),
			expectErr: false,
		},
		{
			name:      "negative value",
			value:     float32Ptr(-1.5),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateNonNegativeFloat("field", tt.value)
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestValidateDimensions(t *testing.T) {
	tests := []struct {
		name      string
		value     []int32
		expectErr bool
	}{
		{
			name:      "nil value",
			value:     nil,
			expectErr: false,
		},
		{
			name:      "valid dimensions",
			value:     []int32{1920, 1080},
			expectErr: false,
		},
		{
			name:      "single element",
			value:     []int32{1920},
			expectErr: true,
		},
		{
			name:      "three elements",
			value:     []int32{1920, 1080, 24},
			expectErr: true,
		},
		{
			name:      "empty array",
			value:     []int32{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDimensions("field", tt.value)
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

// Helper functions for tests
func int32Ptr(i int32) *int32 {
	return &i
}

func float32Ptr(f float32) *float32 {
	return &f
}
