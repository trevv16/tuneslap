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
			request: func() *api.SignupRequest {
				name := "Test User"
				email := "test@example.com"
				password := "password123"
				return &api.SignupRequest{
					Name:     &name,
					Email:    &email,
					Password: &password,
				}
			}(),
			expectErr: false,
		},
		{
			name: "missing name",
			request: func() *api.SignupRequest {
				email := "test@example.com"
				password := "password123"
				return &api.SignupRequest{
					Email:    &email,
					Password: &password,
				}
			}(),
			expectErr: true,
			errFields: []string{"name"},
		},
		{
			name: "missing email",
			request: func() *api.SignupRequest {
				name := "Test User"
				password := "password123"
				return &api.SignupRequest{
					Name:     &name,
					Password: &password,
				}
			}(),
			expectErr: true,
			errFields: []string{"email"},
		},
		{
			name: "missing password",
			request: func() *api.SignupRequest {
				name := "Test User"
				email := "test@example.com"
				return &api.SignupRequest{
					Name:  &name,
					Email: &email,
				}
			}(),
			expectErr: true,
			errFields: []string{"password"},
		},
		{
			name: "name too short",
			request: func() *api.SignupRequest {
				name := "AB"
				email := "test@example.com"
				password := "password123"
				return &api.SignupRequest{
					Name:     &name,
					Email:    &email,
					Password: &password,
				}
			}(),
			expectErr: true,
			errFields: []string{"name"},
		},
		{
			name: "name too long",
			request: func() *api.SignupRequest {
				name := string(make([]byte, 101))
				email := "test@example.com"
				password := "password123"
				return &api.SignupRequest{
					Name:     &name,
					Email:    &email,
					Password: &password,
				}
			}(),
			expectErr: true,
			errFields: []string{"name"},
		},
		{
			name: "invalid email",
			request: func() *api.SignupRequest {
				name := "Test User"
				email := "invalid-email"
				password := "password123"
				return &api.SignupRequest{
					Name:     &name,
					Email:    &email,
					Password: &password,
				}
			}(),
			expectErr: true,
			errFields: []string{"email"},
		},
		{
			name: "password too short",
			request: func() *api.SignupRequest {
				name := "Test User"
				email := "test@example.com"
				password := "short"
				return &api.SignupRequest{
					Name:     &name,
					Email:    &email,
					Password: &password,
				}
			}(),
			expectErr: true,
			errFields: []string{"password"},
		},
		{
			name: "password too long",
			request: func() *api.SignupRequest {
				name := "Test User"
				email := "test@example.com"
				password := string(make([]byte, 129))
				return &api.SignupRequest{
					Name:     &name,
					Email:    &email,
					Password: &password,
				}
			}(),
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
			request: func() *api.SigninRequest {
				email := "test@example.com"
				password := "password123"
				return &api.SigninRequest{
					Email:    &email,
					Password: &password,
				}
			}(),
			expectErr: false,
		},
		{
			name: "missing email",
			request: func() *api.SigninRequest {
				password := "password123"
				return &api.SigninRequest{
					Password: &password,
				}
			}(),
			expectErr: true,
			errFields: []string{"email"},
		},
		{
			name: "missing password",
			request: func() *api.SigninRequest {
				email := "test@example.com"
				return &api.SigninRequest{
					Email: &email,
				}
			}(),
			expectErr: true,
			errFields: []string{"password"},
		},
		{
			name: "invalid email",
			request: func() *api.SigninRequest {
				email := "invalid"
				password := "password123"
				return &api.SigninRequest{
					Email:    &email,
					Password: &password,
				}
			}(),
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
			request: func() *api.ForgotRequest {
				email := "test@example.com"
				return &api.ForgotRequest{
					Email: &email,
				}
			}(),
			expectErr: false,
		},
		{
			name: "missing email",
			request: func() *api.ForgotRequest {
				return &api.ForgotRequest{}
			}(),
			expectErr: true,
			errFields: []string{"email"},
		},
		{
			name: "invalid email",
			request: func() *api.ForgotRequest {
				email := "invalid"
				return &api.ForgotRequest{
					Email: &email,
				}
			}(),
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
			request: func() *api.ResetRequest {
				token := "12345678901234567890123456789012" // 32 chars
				password := "newpassword123"
				return &api.ResetRequest{
					Token:    &token,
					Password: &password,
				}
			}(),
			expectErr: false,
		},
		{
			name: "missing token",
			request: func() *api.ResetRequest {
				password := "newpassword123"
				return &api.ResetRequest{
					Password: &password,
				}
			}(),
			expectErr: true,
			errFields: []string{"token"},
		},
		{
			name: "token too short",
			request: func() *api.ResetRequest {
				token := "short"
				password := "newpassword123"
				return &api.ResetRequest{
					Token:    &token,
					Password: &password,
				}
			}(),
			expectErr: true,
			errFields: []string{"token"},
		},
		{
			name: "missing password",
			request: func() *api.ResetRequest {
				token := "12345678901234567890123456789012"
				return &api.ResetRequest{
					Token: &token,
				}
			}(),
			expectErr: true,
			errFields: []string{"password"},
		},
		{
			name: "password too short",
			request: func() *api.ResetRequest {
				token := "12345678901234567890123456789012"
				password := "short"
				return &api.ResetRequest{
					Token:    &token,
					Password: &password,
				}
			}(),
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
	name := "Test User"
	email := "test@example.com"
	password := "password123"
	req := &api.SignupRequest{
		Name:     &name,
		Email:    &email,
		Password: &password,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateSignupRequest(req)
	}
}

func BenchmarkValidateSigninRequest(b *testing.B) {
	email := "test@example.com"
	password := "password123"
	req := &api.SigninRequest{
		Email:    &email,
		Password: &password,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateSigninRequest(req)
	}
}
