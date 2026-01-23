package validation

import (
	"testing"

	api "tuneslap/generated"

	"github.com/stretchr/testify/assert"
)

func TestNewUserValidator(t *testing.T) {
	v := NewUserValidator()
	assert.NotNil(t, v)
	assert.NotNil(t, v.Validator)
}

func TestUserValidator_Validate(t *testing.T) {
	v := NewUserValidator()

	name := "Test User"
	updateReq := &api.UpdateMeRequest{
		Name: &name,
	}

	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
	}{
		{
			name:      "valid UpdateMeRequest",
			data:      updateReq,
			expectErr: false,
		},
		{
			name:      "invalid type",
			data:      "invalid",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.data)
			if tt.expectErr {
				assert.False(t, result.IsValid)
			} else {
				assert.True(t, result.IsValid)
			}
		})
	}
}

func TestValidateUpdateMeRequest(t *testing.T) {
	tests := []struct {
		name      string
		request   *api.UpdateMeRequest
		expectErr bool
		errFields []string
	}{
		{
			name:      "valid empty request",
			request:   &api.UpdateMeRequest{},
			expectErr: false,
		},
		{
			name: "valid request with name",
			request: func() *api.UpdateMeRequest {
				name := "Updated Name"
				return &api.UpdateMeRequest{
					Name: &name,
				}
			}(),
			expectErr: false,
		},
		{
			name: "valid request with all fields",
			request: func() *api.UpdateMeRequest {
				name := "Updated Name"
				imageUrl := "https://example.com/avatar.png"
				return &api.UpdateMeRequest{
					Name:     &name,
					ImageUrl: &imageUrl,
				}
			}(),
			expectErr: false,
		},
		{
			name: "name too short",
			request: func() *api.UpdateMeRequest {
				name := "AB"
				return &api.UpdateMeRequest{
					Name: &name,
				}
			}(),
			expectErr: true,
			errFields: []string{"name"},
		},
		{
			name: "name too long",
			request: func() *api.UpdateMeRequest {
				name := string(make([]byte, 101))
				return &api.UpdateMeRequest{
					Name: &name,
				}
			}(),
			expectErr: true,
			errFields: []string{"name"},
		},
		{
			name: "name with special characters",
			request: func() *api.UpdateMeRequest {
				name := "User<Test>"
				return &api.UpdateMeRequest{
					Name: &name,
				}
			}(),
			expectErr: true,
			errFields: []string{"name"},
		},
		{
			name: "invalid image URL",
			request: func() *api.UpdateMeRequest {
				imageUrl := "not-a-url"
				return &api.UpdateMeRequest{
					ImageUrl: &imageUrl,
				}
			}(),
			expectErr: true,
			errFields: []string{"imageUrl"},
		},
		{
			name: "image URL with special characters",
			request: func() *api.UpdateMeRequest {
				imageUrl := "https://example.com/image<>.png"
				return &api.UpdateMeRequest{
					ImageUrl: &imageUrl,
				}
			}(),
			expectErr: true,
			errFields: []string{"imageUrl"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateUpdateMeRequest(tt.request)
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
			}
		})
	}
}

func TestUserValidator_ValidateEmail(t *testing.T) {
	v := NewUserValidator()

	tests := []struct {
		name           string
		email          string
		existingEmails []string
		expectErr      bool
	}{
		{
			name:           "valid email no duplicates",
			email:          "test@example.com",
			existingEmails: []string{"other@example.com"},
			expectErr:      false,
		},
		{
			name:           "valid email empty list",
			email:          "test@example.com",
			existingEmails: []string{},
			expectErr:      false,
		},
		{
			name:           "invalid email format",
			email:          "invalid-email",
			existingEmails: []string{},
			expectErr:      true,
		},
		{
			name:           "empty email",
			email:          "",
			existingEmails: []string{},
			expectErr:      true,
		},
		{
			name:           "duplicate email",
			email:          "test@example.com",
			existingEmails: []string{"test@example.com"},
			expectErr:      true,
		},
		{
			name:           "duplicate email case insensitive",
			email:          "Test@Example.com",
			existingEmails: []string{"test@example.com"},
			expectErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateEmail(tt.email, tt.existingEmails)
			if tt.expectErr {
				assert.False(t, result.IsValid)
			} else {
				assert.True(t, result.IsValid)
			}
		})
	}
}

func TestUserValidator_ValidateUpdateMe(t *testing.T) {
	v := NewUserValidator()

	name := "Updated Name"
	req := &api.UpdateMeRequest{
		Name: &name,
	}

	result := v.ValidateUpdateMe(req)
	assert.True(t, result.IsValid)
}

// Benchmark tests
func BenchmarkValidateUpdateMeRequest(b *testing.B) {
	name := "Updated Name"
	imageUrl := "https://example.com/avatar.png"
	req := &api.UpdateMeRequest{
		Name:     &name,
		ImageUrl: &imageUrl,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateUpdateMeRequest(req)
	}
}

func BenchmarkUserValidator_ValidateEmail(b *testing.B) {
	v := NewUserValidator()
	existingEmails := []string{"user1@example.com", "user2@example.com", "user3@example.com"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.ValidateEmail("newuser@example.com", existingEmails)
	}
}
