package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateIntRange(t *testing.T) {
	intPtr := func(v int32) *int32 { return &v }

	tests := []struct {
		name      string
		value     *int32
		min, max  int64
		expectErr bool
	}{
		{"valid within range", intPtr(50), 1, 100, false},
		{"at minimum", intPtr(1), 1, 100, false},
		{"at maximum", intPtr(100), 1, 100, false},
		{"below minimum", intPtr(0), 1, 100, true},
		{"above maximum", intPtr(101), 1, 100, true},
		{"nil value", nil, 1, 100, true},
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
	floatPtr := func(v float32) *float32 { return &v }

	tests := []struct {
		name      string
		value     *float32
		expectErr bool
	}{
		{"nil value", nil, false},
		{"positive value", floatPtr(10.5), false},
		{"zero value", floatPtr(0), false},
		{"negative value", floatPtr(-1.5), true},
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
		{"nil value", nil, false},
		{"valid dimensions", []int32{1920, 1080}, false},
		{"single element", []int32{1920}, true},
		{"three elements", []int32{1920, 1080, 24}, true},
		{"empty array", []int32{}, true},
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
