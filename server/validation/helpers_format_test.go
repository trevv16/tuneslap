package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

// Benchmark tests
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
