package validation

import (
	"testing"

	api "tuneslap/generated"
	"tuneslap/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewKeyValidator(t *testing.T) {
	v := NewKeyValidator()
	assert.NotNil(t, v)
	assert.NotNil(t, v.Validator)
}

func TestKeyValidator_Validate(t *testing.T) {
	v := NewKeyValidator()

	boardId := primitive.NewObjectID().Hex()
	audioMediaId := primitive.NewObjectID().Hex()
	name := "Test Key"
	hotKey := "A"

	createReq := &api.CreateKeyRequest{
		BoardId:      boardId,
		Name:         name,
		AudioMediaId: audioMediaId,
		HotKey:       hotKey,
	}

	updateName := "Updated Key"
	updateReq := &api.UpdateKeyRequest{
		Name: &updateName,
	}

	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
	}{
		{
			name:      "valid CreateKeyRequest",
			data:      createReq,
			expectErr: false,
		},
		{
			name:      "valid UpdateKeyRequest",
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

func TestValidateCreateKeyRequest(t *testing.T) {
	tests := []struct {
		name      string
		request   *api.CreateKeyRequest
		expectErr bool
		errFields []string
	}{
		{
			name: "valid request",
			request: &api.CreateKeyRequest{
				BoardId:      primitive.NewObjectID().Hex(),
				Name:         "Test Key",
				AudioMediaId: primitive.NewObjectID().Hex(),
				HotKey:       "A",
			},
			expectErr: false,
		},
		{
			name: "valid request with all fields",
			request: func() *api.CreateKeyRequest {
				desc := "A test key description"
				imageMediaId := primitive.NewObjectID().Hex()
				return &api.CreateKeyRequest{
					BoardId:      primitive.NewObjectID().Hex(),
					Name:         "Test Key",
					Description:  &desc,
					AudioMediaId: primitive.NewObjectID().Hex(),
					ImageMediaId: &imageMediaId,
					HotKey:       "B",
				}
			}(),
			expectErr: false,
		},
		{
			name: "missing boardId",
			request: &api.CreateKeyRequest{
				Name:         "Test Key",
				AudioMediaId: primitive.NewObjectID().Hex(),
				HotKey:       "A",
			},
			expectErr: true,
			errFields: []string{"boardId"},
		},
		{
			name: "missing name",
			request: &api.CreateKeyRequest{
				BoardId:      primitive.NewObjectID().Hex(),
				AudioMediaId: primitive.NewObjectID().Hex(),
				HotKey:       "A",
			},
			expectErr: true,
			errFields: []string{"name"},
		},
		{
			name: "name too short",
			request: &api.CreateKeyRequest{
				BoardId:      primitive.NewObjectID().Hex(),
				Name:         "AB",
				AudioMediaId: primitive.NewObjectID().Hex(),
				HotKey:       "A",
			},
			expectErr: true,
			errFields: []string{"name"},
		},
		{
			name: "name too long",
			request: &api.CreateKeyRequest{
				BoardId:      primitive.NewObjectID().Hex(),
				Name:         string(make([]byte, 101)),
				AudioMediaId: primitive.NewObjectID().Hex(),
				HotKey:       "A",
			},
			expectErr: true,
			errFields: []string{"name"},
		},
		{
			name: "name with special characters",
			request: &api.CreateKeyRequest{
				BoardId:      primitive.NewObjectID().Hex(),
				Name:         "Test<Key>",
				AudioMediaId: primitive.NewObjectID().Hex(),
				HotKey:       "A",
			},
			expectErr: true,
			errFields: []string{"name"},
		},
		{
			name: "missing audioMediaId",
			request: &api.CreateKeyRequest{
				BoardId: primitive.NewObjectID().Hex(),
				Name:    "Test Key",
				HotKey:  "A",
			},
			expectErr: true,
			errFields: []string{"audioMediaId"},
		},
		{
			name: "missing hotKey",
			request: &api.CreateKeyRequest{
				BoardId:      primitive.NewObjectID().Hex(),
				Name:         "Test Key",
				AudioMediaId: primitive.NewObjectID().Hex(),
			},
			expectErr: true,
			errFields: []string{"hotKey"},
		},
		{
			name: "hotKey too long",
			request: &api.CreateKeyRequest{
				BoardId:      primitive.NewObjectID().Hex(),
				Name:         "Test Key",
				AudioMediaId: primitive.NewObjectID().Hex(),
				HotKey:       "AB",
			},
			expectErr: true,
			errFields: []string{"hotKey"},
		},
		{
			name: "description too long",
			request: func() *api.CreateKeyRequest {
				desc := string(make([]byte, 501))
				return &api.CreateKeyRequest{
					BoardId:      primitive.NewObjectID().Hex(),
					Name:         "Test Key",
					Description:  &desc,
					AudioMediaId: primitive.NewObjectID().Hex(),
					HotKey:       "A",
				}
			}(),
			expectErr: true,
			errFields: []string{"description"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateCreateKeyRequest(tt.request)
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

func TestValidateUpdateKeyRequest(t *testing.T) {
	tests := []struct {
		name      string
		request   *api.UpdateKeyRequest
		expectErr bool
		errFields []string
	}{
		{
			name:      "valid empty request",
			request:   &api.UpdateKeyRequest{},
			expectErr: false,
		},
		{
			name: "valid request with name",
			request: func() *api.UpdateKeyRequest {
				name := "Updated Key"
				return &api.UpdateKeyRequest{
					Name: &name,
				}
			}(),
			expectErr: false,
		},
		{
			name: "valid request with all fields",
			request: func() *api.UpdateKeyRequest {
				name := "Updated Key"
				desc := "Updated description"
				hotKey := "C"
				audioMediaId := primitive.NewObjectID().Hex()
				imageMediaId := primitive.NewObjectID().Hex()
				return &api.UpdateKeyRequest{
					Name:         &name,
					Description:  &desc,
					HotKey:       &hotKey,
					AudioMediaId: &audioMediaId,
					ImageMediaId: &imageMediaId,
				}
			}(),
			expectErr: false,
		},
		{
			name: "name too short",
			request: func() *api.UpdateKeyRequest {
				name := "AB"
				return &api.UpdateKeyRequest{
					Name: &name,
				}
			}(),
			expectErr: true,
			errFields: []string{"name"},
		},
		{
			name: "name with special characters",
			request: func() *api.UpdateKeyRequest {
				name := "Key{Test}"
				return &api.UpdateKeyRequest{
					Name: &name,
				}
			}(),
			expectErr: true,
			errFields: []string{"name"},
		},
		{
			name: "hotKey too long",
			request: func() *api.UpdateKeyRequest {
				hotKey := "AB"
				return &api.UpdateKeyRequest{
					HotKey: &hotKey,
				}
			}(),
			expectErr: true,
			errFields: []string{"hotKey"},
		},
		{
			name: "description too long",
			request: func() *api.UpdateKeyRequest {
				desc := string(make([]byte, 501))
				return &api.UpdateKeyRequest{
					Description: &desc,
				}
			}(),
			expectErr: true,
			errFields: []string{"description"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateUpdateKeyRequest(tt.request)
			if tt.expectErr {
				assert.False(t, result.IsValid)
				assert.NotEmpty(t, result.Errors)
			} else {
				assert.True(t, result.IsValid)
			}
		})
	}
}

func TestKeyValidator_ValidateHotKeyUniqueness(t *testing.T) {
	v := NewKeyValidator()

	key1ID := primitive.NewObjectID()
	key2ID := primitive.NewObjectID()

	boardKeys := []models.Key{
		{ID: key1ID, HotKey: "A"},
		{ID: key2ID, HotKey: "B"},
	}

	tests := []struct {
		name      string
		hotKey    string
		keys      []models.Key
		excludeId primitive.ObjectID
		expectErr bool
	}{
		{
			name:      "unique hotKey",
			hotKey:    "C",
			keys:      boardKeys,
			excludeId: primitive.NilObjectID,
			expectErr: false,
		},
		{
			name:      "duplicate hotKey",
			hotKey:    "A",
			keys:      boardKeys,
			excludeId: primitive.NilObjectID,
			expectErr: true,
		},
		{
			name:      "same hotKey for same key",
			hotKey:    "A",
			keys:      boardKeys,
			excludeId: key1ID,
			expectErr: false,
		},
		{
			name:      "empty keys",
			hotKey:    "A",
			keys:      []models.Key{},
			excludeId: primitive.NilObjectID,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateHotKeyUniqueness(tt.hotKey, tt.keys, tt.excludeId)
			if tt.expectErr {
				assert.False(t, result.IsValid)
			} else {
				assert.True(t, result.IsValid)
			}
		})
	}
}

func TestKeyValidator_ValidateCreateKey(t *testing.T) {
	v := NewKeyValidator()

	req := &api.CreateKeyRequest{
		BoardId:      primitive.NewObjectID().Hex(),
		Name:         "Test Key",
		AudioMediaId: primitive.NewObjectID().Hex(),
		HotKey:       "A",
	}

	result := v.ValidateCreateKey(req)
	assert.True(t, result.IsValid)
}

func TestKeyValidator_ValidateUpdateKey(t *testing.T) {
	v := NewKeyValidator()

	name := "Updated Key"
	req := &api.UpdateKeyRequest{
		Name: &name,
	}

	result := v.ValidateUpdateKey(req)
	assert.True(t, result.IsValid)
}

// Benchmark tests
func BenchmarkValidateCreateKeyRequest(b *testing.B) {
	req := &api.CreateKeyRequest{
		BoardId:      primitive.NewObjectID().Hex(),
		Name:         "Test Key",
		AudioMediaId: primitive.NewObjectID().Hex(),
		HotKey:       "A",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateCreateKeyRequest(req)
	}
}

func BenchmarkKeyValidator_ValidateHotKeyUniqueness(b *testing.B) {
	v := NewKeyValidator()

	boardKeys := make([]models.Key, 26)
	for i := 0; i < 26; i++ {
		boardKeys[i] = models.Key{
			ID:     primitive.NewObjectID(),
			HotKey: string(rune('a' + i)),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.ValidateHotKeyUniqueness("Z", boardKeys, primitive.NilObjectID)
	}
}
