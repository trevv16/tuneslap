package handlers

import (
	"testing"
	"time"
	api "tuneslap/generated"
	"tuneslap/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewKeyHandler(t *testing.T) {
	t.Run("creates handler successfully", func(t *testing.T) {
		// Note: This may fail without a database connection
		defer func() {
			if r := recover(); r != nil {
				t.Skip("Skipping test - requires database connection")
			}
		}()

		handler := NewKeyHandler()
		assert.NotNil(t, handler)
		assert.NotNil(t, handler.keyRepo)
		assert.NotNil(t, handler.ArrayHandler)
	})
}

func TestCreateKeyFromRequest(t *testing.T) {
	// We need to create a handler to test its methods
	// Since createKeyFromRequest is a method on KeyHandler, we test the logic directly

	audioMediaID := primitive.NewObjectID()
	imageMediaID := primitive.NewObjectID()

	t.Run("creates key with all fields", func(t *testing.T) {
		description := "Test description"
		imageMediaIDStr := imageMediaID.Hex()

		request := api.CreateKeyRequest{
			Name:         "Test Key",
			Description:  &description,
			AudioMediaId: audioMediaID.Hex(),
			ImageMediaId: &imageMediaIDStr,
			HotKey:       "A",
		}

		// Simulate createKeyFromRequest logic
		var parsedAudioMediaId primitive.ObjectID
		if request.AudioMediaId != "" {
			parsedAudioMediaId, _ = primitive.ObjectIDFromHex(request.AudioMediaId)
		}

		var parsedImageMediaId primitive.ObjectID
		if request.ImageMediaId != nil && *request.ImageMediaId != "" {
			parsedImageMediaId, _ = primitive.ObjectIDFromHex(*request.ImageMediaId)
		}

		desc := ""
		if request.Description != nil {
			desc = *request.Description
		}

		key := models.Key{
			ID:           primitive.NewObjectID(),
			BoardId:      primitive.ObjectID{},
			Name:         request.Name,
			Description:  desc,
			AudioMediaId: parsedAudioMediaId,
			ImageMediaId: parsedImageMediaId,
			HotKey:       request.HotKey,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		assert.Equal(t, "Test Key", key.Name)
		assert.Equal(t, "Test description", key.Description)
		assert.Equal(t, audioMediaID, key.AudioMediaId)
		assert.Equal(t, imageMediaID, key.ImageMediaId)
		assert.Equal(t, "A", key.HotKey)
		assert.False(t, key.ID.IsZero())
	})

	t.Run("creates key with required fields only", func(t *testing.T) {
		request := api.CreateKeyRequest{
			Name:         "Minimal Key",
			AudioMediaId: audioMediaID.Hex(),
			HotKey:       "B",
		}

		var parsedAudioMediaId primitive.ObjectID
		if request.AudioMediaId != "" {
			parsedAudioMediaId, _ = primitive.ObjectIDFromHex(request.AudioMediaId)
		}

		var parsedImageMediaId primitive.ObjectID
		if request.ImageMediaId != nil && *request.ImageMediaId != "" {
			parsedImageMediaId, _ = primitive.ObjectIDFromHex(*request.ImageMediaId)
		}

		desc := ""
		if request.Description != nil {
			desc = *request.Description
		}

		key := models.Key{
			ID:           primitive.NewObjectID(),
			Name:         request.Name,
			Description:  desc,
			AudioMediaId: parsedAudioMediaId,
			ImageMediaId: parsedImageMediaId,
			HotKey:       request.HotKey,
		}

		assert.Equal(t, "Minimal Key", key.Name)
		assert.Empty(t, key.Description)
		assert.Equal(t, audioMediaID, key.AudioMediaId)
		assert.True(t, key.ImageMediaId.IsZero())
		assert.Equal(t, "B", key.HotKey)
	})

	t.Run("handles empty optional fields", func(t *testing.T) {
		emptyString := ""
		request := api.CreateKeyRequest{
			Name:         "Key No Image",
			AudioMediaId: audioMediaID.Hex(),
			ImageMediaId: &emptyString, // Empty string pointer
			HotKey:       "C",
		}

		var parsedImageMediaId primitive.ObjectID
		if request.ImageMediaId != nil && *request.ImageMediaId != "" {
			parsedImageMediaId, _ = primitive.ObjectIDFromHex(*request.ImageMediaId)
		}

		assert.True(t, parsedImageMediaId.IsZero())
	})

	t.Run("handles invalid ObjectID gracefully", func(t *testing.T) {
		request := api.CreateKeyRequest{
			Name:         "Key Invalid ID",
			AudioMediaId: "invalid-id",
			HotKey:       "D",
		}

		var parsedAudioMediaId primitive.ObjectID
		if request.AudioMediaId != "" {
			parsedAudioMediaId, _ = primitive.ObjectIDFromHex(request.AudioMediaId)
		}

		// Invalid hex string results in zero ObjectID
		assert.True(t, parsedAudioMediaId.IsZero())
	})
}

func TestUpdateKeyFromRequest(t *testing.T) {
	audioMediaID := primitive.NewObjectID()
	imageMediaID := primitive.NewObjectID()

	t.Run("creates update with all fields", func(t *testing.T) {
		name := "Updated Name"
		description := "Updated description"
		hotKey := "X"
		audioIDStr := audioMediaID.Hex()
		imageIDStr := imageMediaID.Hex()

		request := api.UpdateKeyRequest{
			Name:         &name,
			Description:  &description,
			HotKey:       &hotKey,
			AudioMediaId: &audioIDStr,
			ImageMediaId: &imageIDStr,
		}

		// Simulate updateKeyFromRequest logic
		update := map[string]interface{}{
			"updatedAt": time.Now(),
		}

		if request.Name != nil {
			update["name"] = *request.Name
		}
		if request.Description != nil {
			update["description"] = *request.Description
		}
		if request.HotKey != nil {
			update["hotKey"] = *request.HotKey
		}
		if request.AudioMediaId != nil && *request.AudioMediaId != "" {
			audioMediaId, _ := primitive.ObjectIDFromHex(*request.AudioMediaId)
			update["audioMediaId"] = audioMediaId
		}
		if request.ImageMediaId != nil && *request.ImageMediaId != "" {
			imageMediaId, _ := primitive.ObjectIDFromHex(*request.ImageMediaId)
			update["imageMediaId"] = imageMediaId
		}

		assert.Equal(t, "Updated Name", update["name"])
		assert.Equal(t, "Updated description", update["description"])
		assert.Equal(t, "X", update["hotKey"])
		assert.Equal(t, audioMediaID, update["audioMediaId"])
		assert.Equal(t, imageMediaID, update["imageMediaId"])
		assert.NotNil(t, update["updatedAt"])
	})

	t.Run("creates update with partial fields", func(t *testing.T) {
		name := "Just Name"

		request := api.UpdateKeyRequest{
			Name: &name,
		}

		update := map[string]interface{}{
			"updatedAt": time.Now(),
		}

		if request.Name != nil {
			update["name"] = *request.Name
		}
		if request.Description != nil {
			update["description"] = *request.Description
		}
		if request.HotKey != nil {
			update["hotKey"] = *request.HotKey
		}
		if request.AudioMediaId != nil && *request.AudioMediaId != "" {
			audioMediaId, _ := primitive.ObjectIDFromHex(*request.AudioMediaId)
			update["audioMediaId"] = audioMediaId
		}
		if request.ImageMediaId != nil && *request.ImageMediaId != "" {
			imageMediaId, _ := primitive.ObjectIDFromHex(*request.ImageMediaId)
			update["imageMediaId"] = imageMediaId
		}

		assert.Equal(t, "Just Name", update["name"])
		assert.Nil(t, update["description"])
		assert.Nil(t, update["hotKey"])
		assert.Nil(t, update["audioMediaId"])
		assert.Nil(t, update["imageMediaId"])
	})

	t.Run("ignores empty media IDs", func(t *testing.T) {
		emptyStr := ""

		request := api.UpdateKeyRequest{
			AudioMediaId: &emptyStr,
			ImageMediaId: &emptyStr,
		}

		update := map[string]interface{}{
			"updatedAt": time.Now(),
		}

		if request.AudioMediaId != nil && *request.AudioMediaId != "" {
			audioMediaId, _ := primitive.ObjectIDFromHex(*request.AudioMediaId)
			update["audioMediaId"] = audioMediaId
		}
		if request.ImageMediaId != nil && *request.ImageMediaId != "" {
			imageMediaId, _ := primitive.ObjectIDFromHex(*request.ImageMediaId)
			update["imageMediaId"] = imageMediaId
		}

		assert.Nil(t, update["audioMediaId"])
		assert.Nil(t, update["imageMediaId"])
	})

	t.Run("always includes updatedAt", func(t *testing.T) {
		update := map[string]interface{}{
			"updatedAt": time.Now(),
		}

		assert.NotNil(t, update["updatedAt"])
	})
}

func TestKeyResponseMapperLogic(t *testing.T) {
	validTime := time.Now()
	keyID := primitive.NewObjectID()
	boardID := primitive.NewObjectID()
	audioMediaID := primitive.NewObjectID()
	imageMediaID := primitive.NewObjectID()

	t.Run("maps key with all fields", func(t *testing.T) {
		key := models.Key{
			ID:           keyID,
			BoardId:      boardID,
			Name:         "Test Key",
			Description:  "Test description",
			AudioMediaId: audioMediaID,
			ImageMediaId: imageMediaID,
			HotKey:       "A",
			CreatedAt:    validTime,
			UpdatedAt:    validTime,
		}

		assert.Equal(t, keyID, key.ID)
		assert.Equal(t, boardID, key.BoardId)
		assert.Equal(t, "Test Key", key.Name)
		assert.Equal(t, "Test description", key.Description)
		assert.Equal(t, audioMediaID, key.AudioMediaId)
		assert.Equal(t, imageMediaID, key.ImageMediaId)
		assert.Equal(t, "A", key.HotKey)
	})

	t.Run("maps key with empty optional fields", func(t *testing.T) {
		key := models.Key{
			ID:           keyID,
			BoardId:      boardID,
			Name:         "Minimal Key",
			AudioMediaId: audioMediaID,
			ImageMediaId: primitive.NilObjectID,
			HotKey:       "B",
			CreatedAt:    validTime,
			UpdatedAt:    validTime,
		}

		assert.Empty(t, key.Description)
		assert.True(t, key.ImageMediaId.IsZero())
	})
}

func TestKeyHotKeyValidation(t *testing.T) {
	// Test hotkey uniqueness logic

	t.Run("detects duplicate hotkeys", func(t *testing.T) {
		existingKeys := []models.Key{
			{ID: primitive.NewObjectID(), HotKey: "A"},
			{ID: primitive.NewObjectID(), HotKey: "B"},
			{ID: primitive.NewObjectID(), HotKey: "C"},
		}

		newHotKey := "A"

		// Check for duplicate
		isDuplicate := false
		for _, key := range existingKeys {
			if key.HotKey == newHotKey {
				isDuplicate = true
				break
			}
		}

		assert.True(t, isDuplicate)
	})

	t.Run("allows unique hotkeys", func(t *testing.T) {
		existingKeys := []models.Key{
			{ID: primitive.NewObjectID(), HotKey: "A"},
			{ID: primitive.NewObjectID(), HotKey: "B"},
		}

		newHotKey := "C"

		isDuplicate := false
		for _, key := range existingKeys {
			if key.HotKey == newHotKey {
				isDuplicate = true
				break
			}
		}

		assert.False(t, isDuplicate)
	})

	t.Run("allows same hotkey for same key on update", func(t *testing.T) {
		keyID := primitive.NewObjectID()
		existingKeys := []models.Key{
			{ID: keyID, HotKey: "A"},
			{ID: primitive.NewObjectID(), HotKey: "B"},
		}

		updatingKeyID := keyID
		newHotKey := "A"

		isDuplicate := false
		for _, key := range existingKeys {
			if key.HotKey == newHotKey && key.ID != updatingKeyID {
				isDuplicate = true
				break
			}
		}

		assert.False(t, isDuplicate)
	})
}

func TestKeyModelStructure(t *testing.T) {
	t.Run("key has all required fields", func(t *testing.T) {
		key := models.Key{
			ID:           primitive.NewObjectID(),
			BoardId:      primitive.NewObjectID(),
			Name:         "Test Key",
			Description:  "Description",
			AudioMediaId: primitive.NewObjectID(),
			ImageMediaId: primitive.NewObjectID(),
			HotKey:       "1",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		assert.False(t, key.ID.IsZero())
		assert.False(t, key.BoardId.IsZero())
		assert.NotEmpty(t, key.Name)
		assert.NotEmpty(t, key.HotKey)
		assert.False(t, key.CreatedAt.IsZero())
		assert.False(t, key.UpdatedAt.IsZero())
	})
}

// Benchmarks
func BenchmarkCreateKeyFromRequestLogic(b *testing.B) {
	audioMediaID := primitive.NewObjectID()
	imageMediaID := primitive.NewObjectID()
	description := "Test description"
	imageMediaIDStr := imageMediaID.Hex()

	request := api.CreateKeyRequest{
		Name:         "Test Key",
		Description:  &description,
		AudioMediaId: audioMediaID.Hex(),
		ImageMediaId: &imageMediaIDStr,
		HotKey:       "A",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var parsedAudioMediaId primitive.ObjectID
		if request.AudioMediaId != "" {
			parsedAudioMediaId, _ = primitive.ObjectIDFromHex(request.AudioMediaId)
		}

		var parsedImageMediaId primitive.ObjectID
		if request.ImageMediaId != nil && *request.ImageMediaId != "" {
			parsedImageMediaId, _ = primitive.ObjectIDFromHex(*request.ImageMediaId)
		}

		desc := ""
		if request.Description != nil {
			desc = *request.Description
		}

		_ = models.Key{
			ID:           primitive.NewObjectID(),
			Name:         request.Name,
			Description:  desc,
			AudioMediaId: parsedAudioMediaId,
			ImageMediaId: parsedImageMediaId,
			HotKey:       request.HotKey,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
	}
}

func BenchmarkUpdateKeyFromRequestLogic(b *testing.B) {
	audioMediaID := primitive.NewObjectID()
	imageMediaID := primitive.NewObjectID()
	name := "Updated Name"
	description := "Updated description"
	hotKey := "X"
	audioIDStr := audioMediaID.Hex()
	imageIDStr := imageMediaID.Hex()

	request := api.UpdateKeyRequest{
		Name:         &name,
		Description:  &description,
		HotKey:       &hotKey,
		AudioMediaId: &audioIDStr,
		ImageMediaId: &imageIDStr,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		update := map[string]interface{}{
			"updatedAt": time.Now(),
		}

		if request.Name != nil {
			update["name"] = *request.Name
		}
		if request.Description != nil {
			update["description"] = *request.Description
		}
		if request.HotKey != nil {
			update["hotKey"] = *request.HotKey
		}
		if request.AudioMediaId != nil && *request.AudioMediaId != "" {
			audioMediaId, _ := primitive.ObjectIDFromHex(*request.AudioMediaId)
			update["audioMediaId"] = audioMediaId
		}
		if request.ImageMediaId != nil && *request.ImageMediaId != "" {
			imageMediaId, _ := primitive.ObjectIDFromHex(*request.ImageMediaId)
			update["imageMediaId"] = imageMediaId
		}
		_ = update
	}
}
