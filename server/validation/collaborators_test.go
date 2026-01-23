package validation

import (
	"testing"

	api "tuneslap/generated"
	"tuneslap/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewCollaboratorValidator(t *testing.T) {
	v := NewCollaboratorValidator()
	assert.NotNil(t, v)
	assert.NotNil(t, v.Validator)
}

func TestCollaboratorValidator_Validate(t *testing.T) {
	v := NewCollaboratorValidator()

	email := "test@example.com"
	role := "editor"

	createReq := &api.CreateCollaboratorRequest{
		Email: &email,
		Role:  &role,
	}

	updateReq := &api.UpdateCollaboratorRequest{
		Role: &role,
	}

	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
	}{
		{
			name:      "valid CreateCollaboratorRequest",
			data:      createReq,
			expectErr: false,
		},
		{
			name:      "valid UpdateCollaboratorRequest",
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

func TestValidateCreateCollaboratorRequest(t *testing.T) {
	tests := []struct {
		name      string
		request   *api.CreateCollaboratorRequest
		expectErr bool
		errFields []string
	}{
		{
			name: "valid request with editor role",
			request: func() *api.CreateCollaboratorRequest {
				email := "test@example.com"
				role := "editor"
				return &api.CreateCollaboratorRequest{
					Email: &email,
					Role:  &role,
				}
			}(),
			expectErr: false,
		},
		{
			name: "valid request with viewer role",
			request: func() *api.CreateCollaboratorRequest {
				email := "test@example.com"
				role := "viewer"
				return &api.CreateCollaboratorRequest{
					Email: &email,
					Role:  &role,
				}
			}(),
			expectErr: false,
		},
		{
			name: "valid request with owner role",
			request: func() *api.CreateCollaboratorRequest {
				email := "test@example.com"
				role := "owner"
				return &api.CreateCollaboratorRequest{
					Email: &email,
					Role:  &role,
				}
			}(),
			expectErr: false,
		},
		{
			name: "missing email",
			request: func() *api.CreateCollaboratorRequest {
				role := "editor"
				return &api.CreateCollaboratorRequest{
					Role: &role,
				}
			}(),
			expectErr: true,
			errFields: []string{"email"},
		},
		{
			name: "empty email",
			request: func() *api.CreateCollaboratorRequest {
				email := ""
				role := "editor"
				return &api.CreateCollaboratorRequest{
					Email: &email,
					Role:  &role,
				}
			}(),
			expectErr: true,
			errFields: []string{"email"},
		},
		{
			name: "invalid email format",
			request: func() *api.CreateCollaboratorRequest {
				email := "not-an-email"
				role := "editor"
				return &api.CreateCollaboratorRequest{
					Email: &email,
					Role:  &role,
				}
			}(),
			expectErr: true,
			errFields: []string{"email"},
		},
		{
			name: "missing role",
			request: func() *api.CreateCollaboratorRequest {
				email := "test@example.com"
				return &api.CreateCollaboratorRequest{
					Email: &email,
				}
			}(),
			expectErr: true,
			errFields: []string{"role"},
		},
		{
			name: "empty role",
			request: func() *api.CreateCollaboratorRequest {
				email := "test@example.com"
				role := ""
				return &api.CreateCollaboratorRequest{
					Email: &email,
					Role:  &role,
				}
			}(),
			expectErr: true,
			errFields: []string{"role"},
		},
		{
			name: "invalid role",
			request: func() *api.CreateCollaboratorRequest {
				email := "test@example.com"
				role := "admin"
				return &api.CreateCollaboratorRequest{
					Email: &email,
					Role:  &role,
				}
			}(),
			expectErr: true,
			errFields: []string{"role"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateCreateCollaboratorRequest(tt.request)
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

func TestValidateUpdateCollaboratorRequest(t *testing.T) {
	tests := []struct {
		name      string
		request   *api.UpdateCollaboratorRequest
		expectErr bool
		errFields []string
	}{
		{
			name: "valid request with editor role",
			request: func() *api.UpdateCollaboratorRequest {
				role := "editor"
				return &api.UpdateCollaboratorRequest{
					Role: &role,
				}
			}(),
			expectErr: false,
		},
		{
			name: "valid request with viewer role",
			request: func() *api.UpdateCollaboratorRequest {
				role := "viewer"
				return &api.UpdateCollaboratorRequest{
					Role: &role,
				}
			}(),
			expectErr: false,
		},
		{
			name: "valid request with owner role",
			request: func() *api.UpdateCollaboratorRequest {
				role := "owner"
				return &api.UpdateCollaboratorRequest{
					Role: &role,
				}
			}(),
			expectErr: false,
		},
		{
			name:      "missing role",
			request:   &api.UpdateCollaboratorRequest{},
			expectErr: true,
			errFields: []string{"role"},
		},
		{
			name: "empty role",
			request: func() *api.UpdateCollaboratorRequest {
				role := ""
				return &api.UpdateCollaboratorRequest{
					Role: &role,
				}
			}(),
			expectErr: true,
			errFields: []string{"role"},
		},
		{
			name: "invalid role",
			request: func() *api.UpdateCollaboratorRequest {
				role := "admin"
				return &api.UpdateCollaboratorRequest{
					Role: &role,
				}
			}(),
			expectErr: true,
			errFields: []string{"role"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateUpdateCollaboratorRequest(tt.request)
			if tt.expectErr {
				assert.False(t, result.IsValid)
				assert.NotEmpty(t, result.Errors)
			} else {
				assert.True(t, result.IsValid)
			}
		})
	}
}

func TestCollaboratorValidator_ValidateCollaboratorUniqueness(t *testing.T) {
	v := NewCollaboratorValidator()

	user1ID := primitive.NewObjectID()
	user2ID := primitive.NewObjectID()
	newUserID := primitive.NewObjectID()

	collaborators := []models.Collaborator{
		{ID: primitive.NewObjectID(), UserId: user1ID, Email: "user1@example.com"},
		{ID: primitive.NewObjectID(), UserId: user2ID, Email: "user2@example.com"},
	}

	tests := []struct {
		name          string
		userId        primitive.ObjectID
		collaborators []models.Collaborator
		expectErr     bool
	}{
		{
			name:          "new user is unique",
			userId:        newUserID,
			collaborators: collaborators,
			expectErr:     false,
		},
		{
			name:          "user already collaborator",
			userId:        user1ID,
			collaborators: collaborators,
			expectErr:     true,
		},
		{
			name:          "empty collaborators list",
			userId:        user1ID,
			collaborators: []models.Collaborator{},
			expectErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateCollaboratorUniqueness(tt.userId, tt.collaborators)
			if tt.expectErr {
				assert.False(t, result.IsValid)
			} else {
				assert.True(t, result.IsValid)
			}
		})
	}
}

func TestCollaboratorValidator_ValidateCreateCollaborator(t *testing.T) {
	v := NewCollaboratorValidator()

	email := "test@example.com"
	role := "editor"
	req := &api.CreateCollaboratorRequest{
		Email: &email,
		Role:  &role,
	}

	result := v.ValidateCreateCollaborator(req)
	assert.True(t, result.IsValid)
}

func TestCollaboratorValidator_ValidateUpdateCollaborator(t *testing.T) {
	v := NewCollaboratorValidator()

	role := "viewer"
	req := &api.UpdateCollaboratorRequest{
		Role: &role,
	}

	result := v.ValidateUpdateCollaborator(req)
	assert.True(t, result.IsValid)
}

// Benchmark tests
func BenchmarkValidateCreateCollaboratorRequest(b *testing.B) {
	email := "test@example.com"
	role := "editor"
	req := &api.CreateCollaboratorRequest{
		Email: &email,
		Role:  &role,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateCreateCollaboratorRequest(req)
	}
}

func BenchmarkCollaboratorValidator_ValidateCollaboratorUniqueness(b *testing.B) {
	v := NewCollaboratorValidator()

	collaborators := make([]models.Collaborator, 50)
	for i := 0; i < 50; i++ {
		collaborators[i] = models.Collaborator{
			ID:     primitive.NewObjectID(),
			UserId: primitive.NewObjectID(),
			Email:  "user" + string(rune('0'+i)) + "@example.com",
		}
	}

	newUserID := primitive.NewObjectID()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.ValidateCollaboratorUniqueness(newUserID, collaborators)
	}
}
