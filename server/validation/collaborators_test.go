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

	createReq := &api.CreateCollaboratorRequest{
		Email: "test@example.com",
		Role:  "editor",
	}

	updateReq := &api.UpdateCollaboratorRequest{
		Role: "editor",
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
			request: &api.CreateCollaboratorRequest{
				Email: "test@example.com",
				Role:  "editor",
			},
			expectErr: false,
		},
		{
			name: "valid request with viewer role",
			request: &api.CreateCollaboratorRequest{
				Email: "test@example.com",
				Role:  "viewer",
			},
			expectErr: false,
		},
		{
			name: "valid request with owner role",
			request: &api.CreateCollaboratorRequest{
				Email: "test@example.com",
				Role:  "owner",
			},
			expectErr: false,
		},
		{
			name: "missing email",
			request: &api.CreateCollaboratorRequest{
				Email: "",
				Role:  "editor",
			},
			expectErr: true,
			errFields: []string{"email"},
		},
		{
			name: "empty email",
			request: &api.CreateCollaboratorRequest{
				Email: "",
				Role:  "editor",
			},
			expectErr: true,
			errFields: []string{"email"},
		},
		{
			name: "invalid email format",
			request: &api.CreateCollaboratorRequest{
				Email: "not-an-email",
				Role:  "editor",
			},
			expectErr: true,
			errFields: []string{"email"},
		},
		{
			name: "missing role",
			request: &api.CreateCollaboratorRequest{
				Email: "test@example.com",
				Role:  "",
			},
			expectErr: true,
			errFields: []string{"role"},
		},
		{
			name: "empty role",
			request: &api.CreateCollaboratorRequest{
				Email: "test@example.com",
				Role:  "",
			},
			expectErr: true,
			errFields: []string{"role"},
		},
		{
			name: "invalid role",
			request: &api.CreateCollaboratorRequest{
				Email: "test@example.com",
				Role:  "admin",
			},
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
			request: &api.UpdateCollaboratorRequest{
				Role: "editor",
			},
			expectErr: false,
		},
		{
			name: "valid request with viewer role",
			request: &api.UpdateCollaboratorRequest{
				Role: "viewer",
			},
			expectErr: false,
		},
		{
			name: "valid request with owner role",
			request: &api.UpdateCollaboratorRequest{
				Role: "owner",
			},
			expectErr: false,
		},
		{
			name: "missing role",
			request: &api.UpdateCollaboratorRequest{
				Role: "",
			},
			expectErr: true,
			errFields: []string{"role"},
		},
		{
			name: "empty role",
			request: &api.UpdateCollaboratorRequest{
				Role: "",
			},
			expectErr: true,
			errFields: []string{"role"},
		},
		{
			name: "invalid role",
			request: &api.UpdateCollaboratorRequest{
				Role: "admin",
			},
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

	req := &api.CreateCollaboratorRequest{
		Email: "test@example.com",
		Role:  "editor",
	}

	result := v.ValidateCreateCollaborator(req)
	assert.True(t, result.IsValid)
}

func TestCollaboratorValidator_ValidateUpdateCollaborator(t *testing.T) {
	v := NewCollaboratorValidator()

	req := &api.UpdateCollaboratorRequest{
		Role: "viewer",
	}

	result := v.ValidateUpdateCollaborator(req)
	assert.True(t, result.IsValid)
}

// Benchmark tests
func BenchmarkValidateCreateCollaboratorRequest(b *testing.B) {
	req := &api.CreateCollaboratorRequest{
		Email: "test@example.com",
		Role:  "editor",
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
