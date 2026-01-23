package handlers

import (
	"testing"

	api "tuneslap/generated"

	"github.com/stretchr/testify/assert"
)

func TestValidateSignupRequest(t *testing.T) {
	tests := []struct {
		name      string
		request   *api.SignupRequest
		expectErr bool
		errFields []string
	}{
		{
			name: "valid request",
			request: &api.SignupRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			expectErr: false,
		},
		{
			name: "missing name",
			request: &api.SignupRequest{
				Name:     "",
				Email:    "test@example.com",
				Password: "password123",
			},
			expectErr: true,
			errFields: []string{"name"},
		},
		{
			name: "missing email",
			request: &api.SignupRequest{
				Name:     "Test User",
				Email:    "",
				Password: "password123",
			},
			expectErr: true,
			errFields: []string{"email"},
		},
		{
			name: "missing password",
			request: &api.SignupRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "",
			},
			expectErr: true,
			errFields: []string{"password"},
		},
		{
			name: "name too short",
			request: &api.SignupRequest{
				Name:     "AB",
				Email:    "test@example.com",
				Password: "password123",
			},
			expectErr: true,
			errFields: []string{"name"},
		},
		{
			name: "name too long",
			request: &api.SignupRequest{
				Name:     string(make([]byte, 101)),
				Email:    "test@example.com",
				Password: "password123",
			},
			expectErr: true,
			errFields: []string{"name"},
		},
		{
			name: "invalid email",
			request: &api.SignupRequest{
				Name:     "Test User",
				Email:    "invalid-email",
				Password: "password123",
			},
			expectErr: true,
			errFields: []string{"email"},
		},
		{
			name: "password too short",
			request: &api.SignupRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "short",
			},
			expectErr: true,
			errFields: []string{"password"},
		},
		{
			name: "password too long",
			request: &api.SignupRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: string(make([]byte, 129)),
			},
			expectErr: true,
			errFields: []string{"password"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateSignupRequest(tt.request)
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

func TestValidateSigninRequest(t *testing.T) {
	tests := []struct {
		name      string
		request   *api.SigninRequest
		expectErr bool
		errFields []string
	}{
		{
			name: "valid request",
			request: &api.SigninRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			expectErr: false,
		},
		{
			name: "missing email",
			request: &api.SigninRequest{
				Email:    "",
				Password: "password123",
			},
			expectErr: true,
			errFields: []string{"email"},
		},
		{
			name: "missing password",
			request: &api.SigninRequest{
				Email:    "test@example.com",
				Password: "",
			},
			expectErr: true,
			errFields: []string{"password"},
		},
		{
			name: "invalid email",
			request: &api.SigninRequest{
				Email:    "invalid",
				Password: "password123",
			},
			expectErr: true,
			errFields: []string{"email"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateSigninRequest(tt.request)
			if tt.expectErr {
				assert.False(t, result.IsValid)
			} else {
				assert.True(t, result.IsValid)
			}
		})
	}
}

func TestValidateForgotRequest(t *testing.T) {
	tests := []struct {
		name      string
		request   *api.ForgotRequest
		expectErr bool
		errFields []string
	}{
		{
			name: "valid request",
			request: &api.ForgotRequest{
				Email: "test@example.com",
			},
			expectErr: false,
		},
		{
			name: "missing email",
			request: &api.ForgotRequest{
				Email: "",
			},
			expectErr: true,
			errFields: []string{"email"},
		},
		{
			name: "invalid email",
			request: &api.ForgotRequest{
				Email: "invalid",
			},
			expectErr: true,
			errFields: []string{"email"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateForgotRequest(tt.request)
			if tt.expectErr {
				assert.False(t, result.IsValid)
			} else {
				assert.True(t, result.IsValid)
			}
		})
	}
}

func TestValidateResetRequest(t *testing.T) {
	tests := []struct {
		name      string
		request   *api.ResetRequest
		expectErr bool
		errFields []string
	}{
		{
			name: "valid request",
			request: &api.ResetRequest{
				Token:    "12345678901234567890123456789012", // 32 chars
				Password: "newpassword123",
			},
			expectErr: false,
		},
		{
			name: "missing token",
			request: &api.ResetRequest{
				Token:    "",
				Password: "newpassword123",
			},
			expectErr: true,
			errFields: []string{"token"},
		},
		{
			name: "token too short",
			request: &api.ResetRequest{
				Token:    "short",
				Password: "newpassword123",
			},
			expectErr: true,
			errFields: []string{"token"},
		},
		{
			name: "missing password",
			request: &api.ResetRequest{
				Token:    "12345678901234567890123456789012",
				Password: "",
			},
			expectErr: true,
			errFields: []string{"password"},
		},
		{
			name: "password too short",
			request: &api.ResetRequest{
				Token:    "12345678901234567890123456789012",
				Password: "short",
			},
			expectErr: true,
			errFields: []string{"password"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateResetRequest(tt.request)
			if tt.expectErr {
				assert.False(t, result.IsValid)
			} else {
				assert.True(t, result.IsValid)
			}
		})
	}
}

func TestNewAuthHandler(t *testing.T) {
	handler := NewAuthHandler()
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.userRepo)
	assert.NotNil(t, handler.validator)
}

// Benchmark tests
func BenchmarkValidateSignupRequest(b *testing.B) {
	req := &api.SignupRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateSignupRequest(req)
	}
}

func BenchmarkValidateSigninRequest(b *testing.B) {
	req := &api.SigninRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateSigninRequest(req)
	}
}
