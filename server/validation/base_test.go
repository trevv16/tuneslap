package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValidator(t *testing.T) {
	v := NewValidator()
	assert.NotNil(t, v)
	assert.NotNil(t, v.validate)
}

func TestValidator_Validate(t *testing.T) {
	v := NewValidator()

	type TestStruct struct {
		Name  string `validate:"required,min=3,max=50"`
		Email string `validate:"required,email"`
		Age   int    `validate:"min=0,max=150"`
	}

	tests := []struct {
		name      string
		data      TestStruct
		expectErr bool
		errFields []string
	}{
		{
			name: "valid struct",
			data: TestStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   30,
			},
			expectErr: false,
		},
		{
			name: "missing required name",
			data: TestStruct{
				Name:  "",
				Email: "john@example.com",
				Age:   30,
			},
			expectErr: true,
			errFields: []string{"Name"},
		},
		{
			name: "invalid email",
			data: TestStruct{
				Name:  "John Doe",
				Email: "invalid-email",
				Age:   30,
			},
			expectErr: true,
			errFields: []string{"Email"},
		},
		{
			name: "name too short",
			data: TestStruct{
				Name:  "Jo",
				Email: "john@example.com",
				Age:   30,
			},
			expectErr: true,
			errFields: []string{"Name"},
		},
		{
			name: "multiple errors",
			data: TestStruct{
				Name:  "",
				Email: "invalid",
				Age:   -1,
			},
			expectErr: true,
			errFields: []string{"Name", "Email", "Age"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.data)

			if tt.expectErr {
				assert.False(t, result.IsValid)
				assert.NotEmpty(t, result.Errors)
				for _, expectedField := range tt.errFields {
					found := false
					for _, err := range result.Errors {
						if err.Field == expectedField {
							found = true
							break
						}
					}
					assert.True(t, found, "Expected error for field %s", expectedField)
				}
			} else {
				assert.True(t, result.IsValid)
				assert.Empty(t, result.Errors)
			}
		})
	}
}

func TestValidateAlphaNumSpaceLogic(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{
			name:      "valid alphanumeric",
			input:     "abc123",
			expectErr: false,
		},
		{
			name:      "valid with spaces",
			input:     "Hello World 123",
			expectErr: false,
		},
		{
			name:      "empty string",
			input:     "",
			expectErr: false,
		},
		{
			name:      "invalid special characters",
			input:     "hello@world",
			expectErr: true,
		},
		{
			name:      "invalid underscore",
			input:     "hello_world",
			expectErr: true,
		},
		{
			name:      "invalid hyphen",
			input:     "hello-world",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAlphaNumSpaceHelper("field", tt.input)
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestValidateExcludesAllLogic(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		excluded  string
		expectErr bool
	}{
		{
			name:      "valid no excluded chars",
			input:     "hello world",
			excluded:  "<>{}",
			expectErr: false,
		},
		{
			name:      "invalid contains <",
			input:     "hello<world",
			excluded:  "<>{}",
			expectErr: true,
		},
		{
			name:      "invalid contains >",
			input:     "hello>world",
			excluded:  "<>{}",
			expectErr: true,
		},
		{
			name:      "empty string",
			input:     "",
			excluded:  "<>{}",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateExcludesAllHelper("field", tt.input, tt.excluded)
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestValidateFileNameHelper(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{
			name:      "valid filename",
			input:     "document.pdf",
			expectErr: false,
		},
		{
			name:      "valid filename with spaces",
			input:     "my document.pdf",
			expectErr: false,
		},
		{
			name:      "valid filename with hyphens",
			input:     "my-document.pdf",
			expectErr: false,
		},
		{
			name:      "valid filename with underscores",
			input:     "my_document.pdf",
			expectErr: false,
		},
		{
			name:      "invalid filename with special chars",
			input:     "document<>.pdf",
			expectErr: true,
		},
		{
			name:      "invalid filename with brackets",
			input:     "document[1].pdf",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFileNameHelper("fileName", tt.input)
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestValidateContentTypeHelper(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{
			name:      "valid audio mpeg",
			input:     "audio/mpeg",
			expectErr: false,
		},
		{
			name:      "valid audio wav",
			input:     "audio/x-wav",
			expectErr: false,
		},
		{
			name:      "valid image jpeg",
			input:     "image/jpeg",
			expectErr: false,
		},
		{
			name:      "valid image png",
			input:     "image/png",
			expectErr: false,
		},
		{
			name:      "valid application json",
			input:     "application/json",
			expectErr: false,
		},
		{
			name:      "invalid with special chars",
			input:     "audio/<mpeg>",
			expectErr: true,
		},
		{
			name:      "invalid with brackets",
			input:     "audio/[mpeg]",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateContentTypeHelper("contentType", tt.input)
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestNewValidationError(t *testing.T) {
	tests := []struct {
		name      string
		errors    []ValidationError
		expectErr bool
	}{
		{
			name:      "empty errors",
			errors:    []ValidationError{},
			expectErr: false,
		},
		{
			name: "single error",
			errors: []ValidationError{
				{Field: "name", Message: "name is required"},
			},
			expectErr: true,
		},
		{
			name: "multiple errors",
			errors: []ValidationError{
				{Field: "name", Message: "name is required"},
				{Field: "email", Message: "email is invalid"},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewValidationError(tt.errors)
			if tt.expectErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), "validation failed")
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetErrorMessage(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name     string
		field    string
		tag      string
		param    string
		expected string
	}{
		{
			name:     "required error",
			field:    "name",
			tag:      "required",
			param:    "",
			expected: "name is required",
		},
		{
			name:     "min error",
			field:    "name",
			tag:      "min",
			param:    "3",
			expected: "name must be at least 3 characters long",
		},
		{
			name:     "max error",
			field:    "name",
			tag:      "max",
			param:    "100",
			expected: "name must be at most 100 characters long",
		},
		{
			name:     "email error",
			field:    "email",
			tag:      "email",
			param:    "",
			expected: "email must be a valid email address",
		},
		{
			name:     "url error",
			field:    "website",
			tag:      "url",
			param:    "",
			expected: "website must be a valid URL",
		},
		{
			name:     "oneof error",
			field:    "role",
			tag:      "oneof",
			param:    "admin user",
			expected: "role must be one of: admin user",
		},
		{
			name:     "unknown tag",
			field:    "field",
			tag:      "unknown",
			param:    "",
			expected: "field is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := v.getErrorMessage(tt.field, tt.tag, tt.param)
			assert.Equal(t, tt.expected, msg)
		})
	}
}

// Benchmark tests
func BenchmarkValidator_Validate(b *testing.B) {
	v := NewValidator()
	type TestStruct struct {
		Name  string `validate:"required,min=3,max=50"`
		Email string `validate:"required,email"`
	}

	data := TestStruct{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Validate(data)
	}
}

func BenchmarkValidateAlphaNumSpaceHelper(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateAlphaNumSpaceHelper("field", "Hello World 123")
	}
}
