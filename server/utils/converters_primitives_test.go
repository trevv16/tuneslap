package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestObjectIDToString(t *testing.T) {
	tests := []struct {
		name     string
		id       primitive.ObjectID
		expected bool // whether result should be non-nil
	}{
		{
			name:     "valid ObjectID",
			id:       primitive.NewObjectID(),
			expected: true,
		},
		{
			name:     "zero ObjectID",
			id:       primitive.NilObjectID,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := objectIDToString(tt.id)
			if tt.expected {
				assert.NotNil(t, result)
				assert.Equal(t, tt.id.Hex(), *result)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestStringToObjectID(t *testing.T) {
	validID := primitive.NewObjectID()

	tests := []struct {
		name        string
		input       *string
		expectError bool
		expectNil   bool
	}{
		{
			name:        "valid ObjectID string",
			input:       strPtr(validID.Hex()),
			expectError: false,
			expectNil:   false,
		},
		{
			name:        "nil input",
			input:       nil,
			expectError: false,
			expectNil:   true,
		},
		{
			name:        "empty string",
			input:       strPtr(""),
			expectError: false,
			expectNil:   true,
		},
		{
			name:        "invalid ObjectID string",
			input:       strPtr("invalid"),
			expectError: true,
			expectNil:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := stringToObjectID(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.expectNil {
					assert.Equal(t, primitive.NilObjectID, result)
				} else {
					assert.NotEqual(t, primitive.NilObjectID, result)
				}
			}
		})
	}
}

func TestDateTimeToTime(t *testing.T) {
	now := time.Now()
	validDateTime := primitive.NewDateTimeFromTime(now)
	invalidDateTime := primitive.DateTime(0) // Unix epoch, should be valid

	tests := []struct {
		name  string
		input primitive.DateTime
	}{
		{
			name:  "valid DateTime",
			input: validDateTime,
		},
		{
			name:  "zero DateTime",
			input: invalidDateTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dateTimeToTime(tt.input)
			assert.NotNil(t, result)
		})
	}
}

func TestTimeToDateTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name  string
		input *time.Time
	}{
		{
			name:  "valid time",
			input: &now,
		},
		{
			name:  "nil time",
			input: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := timeToDateTime(tt.input)
			assert.NotZero(t, result)
		})
	}
}

func TestStringPtr(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool // whether result should be non-nil
	}{
		{
			name:     "non-empty string",
			input:    "test",
			expected: true,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stringPtr(tt.input)
			if tt.expected {
				assert.NotNil(t, result)
				assert.Equal(t, tt.input, *result)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestStringVal(t *testing.T) {
	value := "test"

	tests := []struct {
		name     string
		input    *string
		expected string
	}{
		{
			name:     "non-nil pointer",
			input:    &value,
			expected: "test",
		},
		{
			name:     "nil pointer",
			input:    nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stringVal(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFloat64ToFloat32(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float32
	}{
		{
			name:     "positive value",
			input:    10.5,
			expected: 10.5,
		},
		{
			name:     "zero value",
			input:    0,
			expected: 0,
		},
		{
			name:     "negative value",
			input:    -5.25,
			expected: -5.25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := float64ToFloat32(tt.input)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expected, *result)
		})
	}
}

func TestFloat32ToFloat64(t *testing.T) {
	value := float32(10.5)

	tests := []struct {
		name     string
		input    *float32
		expected float64
	}{
		{
			name:     "non-nil value",
			input:    &value,
			expected: 10.5,
		},
		{
			name:     "nil value",
			input:    nil,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := float32ToFloat64(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIntToInt32(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int32
	}{
		{
			name:     "positive value",
			input:    100,
			expected: 100,
		},
		{
			name:     "zero value",
			input:    0,
			expected: 0,
		},
		{
			name:     "negative value",
			input:    -50,
			expected: -50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := intToInt32(tt.input)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expected, *result)
		})
	}
}

func TestInt32ToInt(t *testing.T) {
	value := int32(100)

	tests := []struct {
		name     string
		input    *int32
		expected int
	}{
		{
			name:     "non-nil value",
			input:    &value,
			expected: 100,
		},
		{
			name:     "nil value",
			input:    nil,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := int32ToInt(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIntArrayToInt32Array(t *testing.T) {
	tests := []struct {
		name     string
		input    [2]int
		expected []int32
	}{
		{
			name:     "valid dimensions",
			input:    [2]int{1920, 1080},
			expected: []int32{1920, 1080},
		},
		{
			name:     "zero values",
			input:    [2]int{0, 0},
			expected: []int32{0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := intArrayToInt32Array(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInt32ArrayToIntArray(t *testing.T) {
	tests := []struct {
		name     string
		input    []int32
		expected [2]int
	}{
		{
			name:     "valid dimensions",
			input:    []int32{1920, 1080},
			expected: [2]int{1920, 1080},
		},
		{
			name:     "less than 2 elements",
			input:    []int32{100},
			expected: [2]int{0, 0},
		},
		{
			name:     "more than 2 elements",
			input:    []int32{100, 200, 300},
			expected: [2]int{100, 200},
		},
		{
			name:     "empty array",
			input:    []int32{},
			expected: [2]int{0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := int32ArrayToIntArray(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInt32ArrayToInt4Array(t *testing.T) {
	tests := []struct {
		name     string
		input    []int32
		expected [4]int
	}{
		{
			name:     "valid crop values",
			input:    []int32{10, 20, 100, 200},
			expected: [4]int{10, 20, 100, 200},
		},
		{
			name:     "less than 4 elements",
			input:    []int32{10, 20},
			expected: [4]int{0, 0, 0, 0},
		},
		{
			name:     "empty array",
			input:    []int32{},
			expected: [4]int{0, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := int32ArrayToInt4Array(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function for tests
func strPtr(s string) *string {
	return &s
}
