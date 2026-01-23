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

func TestValidateURLHelper(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{
			name:      "valid https url",
			value:     "https://example.com",
			expectErr: false,
		},
		{
			name:      "valid http url",
			value:     "http://example.com",
			expectErr: false,
		},
		{
			name:      "valid url with path",
			value:     "https://example.com/path/to/resource",
			expectErr: false,
		},
		{
			name:      "valid url with query",
			value:     "https://example.com?key=value",
			expectErr: false,
		},
		{
			name:      "invalid url",
			value:     "not-a-url",
			expectErr: true,
		},
		{
			name:      "missing protocol",
			value:     "example.com",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURLHelper("field", tt.value)
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestValidateEmailHelper(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{
			name:      "valid email",
			value:     "test@example.com",
			expectErr: false,
		},
		{
			name:      "valid email with subdomain",
			value:     "test@mail.example.com",
			expectErr: false,
		},
		{
			name:      "valid email with plus",
			value:     "test+tag@example.com",
			expectErr: false,
		},
		{
			name:      "invalid email no @",
			value:     "testexample.com",
			expectErr: true,
		},
		{
			name:      "invalid email no domain",
			value:     "test@",
			expectErr: true,
		},
		{
			name:      "invalid email no local part",
			value:     "@example.com",
			expectErr: true,
		},
		{
			name:      "empty string",
			value:     "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmailHelper("field", tt.value)
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestValidateOneOf(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		options   []string
		expectErr bool
	}{
		{
			name:      "valid first option",
			value:     "grid",
			options:   []string{"grid", "list"},
			expectErr: false,
		},
		{
			name:      "valid second option",
			value:     "list",
			options:   []string{"grid", "list"},
			expectErr: false,
		},
		{
			name:      "invalid option",
			value:     "table",
			options:   []string{"grid", "list"},
			expectErr: true,
		},
		{
			name:      "empty value",
			value:     "",
			options:   []string{"grid", "list"},
			expectErr: true,
		},
		{
			name:      "case sensitive mismatch",
			value:     "GRID",
			options:   []string{"grid", "list"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOneOf("field", tt.value, tt.options)
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

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
func strPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}

func float32Ptr(f float32) *float32 {
	return &f
}

// Benchmark tests
func BenchmarkValidateRequiredString(b *testing.B) {
	value := strPtr("test value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateRequiredString("field", value)
	}
}

func BenchmarkValidateEmailHelper(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateEmailHelper("email", "test@example.com")
	}
}

func BenchmarkValidateURLHelper(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateURLHelper("url", "https://example.com/path/to/resource")
	}
}

func BenchmarkValidateOneOf(b *testing.B) {
	options := []string{"grid", "list", "table", "card", "timeline"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateOneOf("layout", "card", options)
	}
}
