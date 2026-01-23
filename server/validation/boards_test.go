package validation

import (
	"testing"

	api "tuneslap/generated"
	"tuneslap/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewBoardValidator(t *testing.T) {
	v := NewBoardValidator()
	assert.NotNil(t, v)
	assert.NotNil(t, v.Validator)
}

func TestBoardValidator_Validate(t *testing.T) {
	v := NewBoardValidator()

	name := "Test Board"
	createReq := &api.CreateBoardRequest{
		Name:   "Test Board",
		Layout: "grid",
	}

	updateReq := &api.UpdateBoardRequest{
		Name: &name,
	}

	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
	}{
		{
			name:      "valid CreateBoardRequest",
			data:      createReq,
			expectErr: false,
		},
		{
			name:      "valid UpdateBoardRequest",
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

func TestValidateCreateBoardRequest(t *testing.T) {
	tests := []struct {
		name      string
		request   *api.CreateBoardRequest
		expectErr bool
		errFields []string
	}{
		{
			name: "valid request",
			request: &api.CreateBoardRequest{
				Name:   "Test Board",
				Layout: "grid",
			},
			expectErr: false,
		},
		{
			name: "valid request with all fields",
			request: func() *api.CreateBoardRequest {
				desc := "A test board description"
				imageUrl := "https://example.com/image.png"
				return &api.CreateBoardRequest{
					Name:        "Test Board",
					Description: &desc,
					Layout:      "list",
					ImageUrl:    &imageUrl,
				}
			}(),
			expectErr: false,
		},
		{
			name: "missing name",
			request: &api.CreateBoardRequest{
				Name:   "",
				Layout: "grid",
			},
			expectErr: true,
			errFields: []string{"name"},
		},
		{
			name: "missing layout",
			request: &api.CreateBoardRequest{
				Name:   "Test Board",
				Layout: "",
			},
			expectErr: true,
			errFields: []string{"layout"},
		},
		{
			name: "name too short",
			request: &api.CreateBoardRequest{
				Name:   "AB",
				Layout: "grid",
			},
			expectErr: true,
			errFields: []string{"name"},
		},
		{
			name: "name too long",
			request: &api.CreateBoardRequest{
				Name:   string(make([]byte, 101)),
				Layout: "grid",
			},
			expectErr: true,
			errFields: []string{"name"},
		},
		{
			name: "invalid layout",
			request: &api.CreateBoardRequest{
				Name:   "Test Board",
				Layout: "invalid",
			},
			expectErr: true,
			errFields: []string{"layout"},
		},
		{
			name: "name with special characters",
			request: &api.CreateBoardRequest{
				Name:   "Test<Board>",
				Layout: "grid",
			},
			expectErr: true,
			errFields: []string{"name"},
		},
		{
			name: "description too long",
			request: func() *api.CreateBoardRequest {
				desc := string(make([]byte, 501))
				return &api.CreateBoardRequest{
					Name:        "Test Board",
					Description: &desc,
					Layout:      "grid",
				}
			}(),
			expectErr: true,
			errFields: []string{"description"},
		},
		{
			name: "invalid image URL",
			request: func() *api.CreateBoardRequest {
				imageUrl := "not-a-url"
				return &api.CreateBoardRequest{
					Name:     "Test Board",
					Layout:   "grid",
					ImageUrl: &imageUrl,
				}
			}(),
			expectErr: true,
			errFields: []string{"imageUrl"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateCreateBoardRequest(tt.request)
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

func TestValidateUpdateBoardRequest(t *testing.T) {
	tests := []struct {
		name      string
		request   *api.UpdateBoardRequest
		expectErr bool
		errFields []string
	}{
		{
			name:      "valid empty request",
			request:   &api.UpdateBoardRequest{},
			expectErr: false,
		},
		{
			name: "valid request with name",
			request: func() *api.UpdateBoardRequest {
				name := "Updated Board"
				return &api.UpdateBoardRequest{
					Name: &name,
				}
			}(),
			expectErr: false,
		},
		{
			name: "valid request with all fields",
			request: func() *api.UpdateBoardRequest {
				name := "Updated Board"
				desc := "Updated description"
				layout := "list"
				imageUrl := "https://example.com/new-image.png"
				return &api.UpdateBoardRequest{
					Name:        &name,
					Description: &desc,
					Layout:      &layout,
					ImageUrl:    &imageUrl,
				}
			}(),
			expectErr: false,
		},
		{
			name: "name too short",
			request: func() *api.UpdateBoardRequest {
				name := "AB"
				return &api.UpdateBoardRequest{
					Name: &name,
				}
			}(),
			expectErr: true,
			errFields: []string{"name"},
		},
		{
			name: "invalid layout",
			request: func() *api.UpdateBoardRequest {
				layout := "invalid"
				return &api.UpdateBoardRequest{
					Layout: &layout,
				}
			}(),
			expectErr: true,
			errFields: []string{"layout"},
		},
		{
			name: "name with special characters",
			request: func() *api.UpdateBoardRequest {
				name := "Board{Test}"
				return &api.UpdateBoardRequest{
					Name: &name,
				}
			}(),
			expectErr: true,
			errFields: []string{"name"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateUpdateBoardRequest(tt.request)
			if tt.expectErr {
				assert.False(t, result.IsValid)
				assert.NotEmpty(t, result.Errors)
			} else {
				assert.True(t, result.IsValid)
			}
		})
	}
}

func TestBoardValidator_ValidateBoardOwnership(t *testing.T) {
	v := NewBoardValidator()

	authorID := primitive.NewObjectID()
	otherUserID := primitive.NewObjectID()

	board := models.Board{
		ID:       primitive.NewObjectID(),
		AuthorId: authorID,
		Name:     "Test Board",
	}

	tests := []struct {
		name      string
		board     models.Board
		userID    primitive.ObjectID
		expectErr bool
	}{
		{
			name:      "owner can access",
			board:     board,
			userID:    authorID,
			expectErr: false,
		},
		{
			name:      "non-owner cannot access",
			board:     board,
			userID:    otherUserID,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateBoardOwnership(tt.board, tt.userID)
			if tt.expectErr {
				assert.False(t, result.IsValid)
			} else {
				assert.True(t, result.IsValid)
			}
		})
	}
}

func TestBoardValidator_ValidateCreateBoard(t *testing.T) {
	v := NewBoardValidator()

	req := &api.CreateBoardRequest{
		Name:   "Test Board",
		Layout: "grid",
	}

	result := v.ValidateCreateBoard(req)
	assert.True(t, result.IsValid)
}

func TestBoardValidator_ValidateUpdateBoard(t *testing.T) {
	v := NewBoardValidator()

	name := "Updated Board"
	req := &api.UpdateBoardRequest{
		Name: &name,
	}

	result := v.ValidateUpdateBoard(req)
	assert.True(t, result.IsValid)
}

// Benchmark tests
func BenchmarkValidateCreateBoardRequest(b *testing.B) {
	req := &api.CreateBoardRequest{
		Name:   "Test Board",
		Layout: "grid",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateCreateBoardRequest(req)
	}
}

func BenchmarkValidateUpdateBoardRequest(b *testing.B) {
	name := "Updated Board"
	req := &api.UpdateBoardRequest{
		Name: &name,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateUpdateBoardRequest(req)
	}
}
