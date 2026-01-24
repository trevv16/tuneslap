package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRequiredString(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		value     *string
		expectErr bool
	}{
		{
			name:      "valid string",
			field:     "name",
			value:     strPtr("test"),
			expectErr: false,
		},
		{
			name:      "nil value",
			field:     "name",
			value:     nil,
			expectErr: true,
		},
		{
			name:      "empty string",
			field:     "name",
			value:     strPtr(""),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRequiredString(tt.field, tt.value)
			if tt.expectErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.field, err.Field)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestValidateOptionalString(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		value     *string
		min       int
		max       int
		rules     []string
		expectErr bool
	}{
		{
			name:      "nil value",
			field:     "desc",
			value:     nil,
			expectErr: false,
		},
		{
			name:      "empty value",
			field:     "desc",
			value:     strPtr(""),
			expectErr: false,
		},
		{
			name:      "valid within limits",
			field:     "desc",
			value:     strPtr("test"),
			min:       1,
			max:       10,
			expectErr: false,
		},
		{
			name:      "too short",
			field:     "desc",
			value:     strPtr("ab"),
			min:       5,
			max:       100,
			expectErr: true,
		},
		{
			name:      "too long",
			field:     "desc",
			value:     strPtr("this is a very long string"),
			min:       1,
			max:       10,
			expectErr: true,
		},
		{
			name:      "valid with url rule",
			field:     "url",
			value:     strPtr("https://example.com"),
			rules:     []string{"url"},
			expectErr: false,
		},
		{
			name:      "invalid url",
			field:     "url",
			value:     strPtr("not-a-url"),
			rules:     []string{"url"},
			expectErr: true,
		},
		{
			name:      "valid email",
			field:     "email",
			value:     strPtr("test@example.com"),
			rules:     []string{"email"},
			expectErr: false,
		},
		{
			name:      "invalid email",
			field:     "email",
			value:     strPtr("invalid-email"),
			rules:     []string{"email"},
			expectErr: true,
		},
		{
			name:      "valid alphanumspace",
			field:     "name",
			value:     strPtr("Test Name 123"),
			rules:     []string{"alphanumspace"},
			expectErr: false,
		},
		{
			name:      "invalid alphanumspace",
			field:     "name",
			value:     strPtr("Test@Name"),
			rules:     []string{"alphanumspace"},
			expectErr: true,
		},
		{
			name:      "valid excludesall",
			field:     "name",
			value:     strPtr("Test Name"),
			rules:     []string{"excludesall"},
			expectErr: false,
		},
		{
			name:      "invalid excludesall",
			field:     "name",
			value:     strPtr("Test<Name>"),
			rules:     []string{"excludesall"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOptionalString(tt.field, tt.value, tt.min, tt.max, tt.rules...)
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestValidateStringLength(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		value     string
		min       int
		max       int
		expectErr bool
	}{
		{
			name:      "valid length",
			field:     "name",
			value:     "test",
			min:       1,
			max:       10,
			expectErr: false,
		},
		{
			name:      "at minimum",
			field:     "name",
			value:     "abc",
			min:       3,
			max:       10,
			expectErr: false,
		},
		{
			name:      "at maximum",
			field:     "name",
			value:     "1234567890",
			min:       1,
			max:       10,
			expectErr: false,
		},
		{
			name:      "below minimum",
			field:     "name",
			value:     "ab",
			min:       3,
			max:       10,
			expectErr: true,
		},
		{
			name:      "above maximum",
			field:     "name",
			value:     "12345678901",
			min:       1,
			max:       10,
			expectErr: true,
		},
		{
			name:      "no minimum check",
			field:     "name",
			value:     "a",
			min:       0,
			max:       10,
			expectErr: false,
		},
		{
			name:      "no maximum check",
			field:     "name",
			value:     "a very long string that should pass",
			min:       1,
			max:       0,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateStringLength(tt.field, tt.value, tt.min, tt.max)
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestValidateAlphaNumSpaceHelper(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{
			name:      "only letters",
			value:     "HelloWorld",
			expectErr: false,
		},
		{
			name:      "letters and numbers",
			value:     "Test123",
			expectErr: false,
		},
		{
			name:      "letters numbers and spaces",
			value:     "Test 123 Name",
			expectErr: false,
		},
		{
			name:      "empty string",
			value:     "",
			expectErr: false,
		},
		{
			name:      "with hyphen",
			value:     "test-name",
			expectErr: true,
		},
		{
			name:      "with underscore",
			value:     "test_name",
			expectErr: true,
		},
		{
			name:      "with special chars",
			value:     "test@name",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAlphaNumSpaceHelper("field", tt.value)
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestValidateExcludesAllHelper(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		excluded  string
		expectErr bool
	}{
		{
			name:      "no excluded chars",
			value:     "hello world",
			excluded:  "<>{}[]",
			expectErr: false,
		},
		{
			name:      "contains <",
			value:     "hello<world",
			excluded:  "<>",
			expectErr: true,
		},
		{
			name:      "contains >",
			value:     "hello>world",
			excluded:  "<>",
			expectErr: true,
		},
		{
			name:      "contains {",
			value:     "hello{world",
			excluded:  "{}",
			expectErr: true,
		},
		{
			name:      "contains }",
			value:     "hello}world",
			excluded:  "{}",
			expectErr: true,
		},
		{
			name:      "empty string",
			value:     "",
			excluded:  "<>{}",
			expectErr: false,
		},
		{
			name:      "empty excluded",
			value:     "<hello>",
			excluded:  "",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateExcludesAllHelper("field", tt.value, tt.excluded)
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

// Helper functions for tests
func strPtr(s string) *string {
	return &s
}

// Benchmark tests
func BenchmarkValidateRequiredString(b *testing.B) {
	value := strPtr("test value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateRequiredString("field", value)
	}
}
