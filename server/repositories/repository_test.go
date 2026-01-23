package repositories

import (
	"testing"
	"time"
	"tuneslap/models"
	"tuneslap/testutils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestNewRepository tests repository creation
func TestNewRepository(t *testing.T) {
	repo := NewRepository[models.User]("test_collection")

	assert.NotNil(t, repo)
	assert.Equal(t, "test_collection", repo.collectionName)
	assert.Equal(t, 10*time.Second, repo.timeout)
	assert.Equal(t, 5*time.Minute, repo.cacheTTL)
}

// TestNewArrayRepository tests array repository creation
func TestNewArrayRepository(t *testing.T) {
	repo := NewArrayRepository[models.Key]("boards")

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.Repository)
	assert.Equal(t, "boards", repo.collectionName)
}

// TestGetCacheKey tests cache key generation
func TestGetCacheKey(t *testing.T) {
	repo := NewRepository[models.User]("users")

	tests := []struct {
		name      string
		operation string
		params    []interface{}
		expected  string
	}{
		{
			name:      "simple operation",
			operation: "findOne",
			params:    nil,
			expected:  "users:findOne",
		},
		{
			name:      "operation with string param",
			operation: "findByEmail",
			params:    []interface{}{"test@example.com"},
			expected:  "users:findByEmail:test@example.com",
		},
		{
			name:      "operation with multiple params",
			operation: "findAll",
			params:    []interface{}{"filter1", 100},
			expected:  "users:findAll:filter1:100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.getCacheKey(tt.operation, tt.params...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Integration tests (require MongoDB)
func TestRepositoryIntegration(t *testing.T) {
	cleanup := testutils.SetupTestMongoDBWithSkip(t)
	defer cleanup()

	// Also setup Redis for cache tests
	redisCleanup := testutils.SetupTestRedisWithSkip(t)
	defer redisCleanup()

	// Create test repository
	repo := NewRepository[models.User]("test_users")

	t.Run("Create", func(t *testing.T) {
		user := models.User{
			ID:           primitive.NewObjectID(),
			Name:         "Test User",
			Email:        "create@test.com",
			PasswordHash: "hash123",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		created, err := repo.Create(user)
		assert.NoError(t, err)
		assert.NotEqual(t, primitive.NilObjectID, created.ID)
		assert.Equal(t, "Test User", created.Name)
	})

	t.Run("FindByID", func(t *testing.T) {
		// First create a user
		userID := primitive.NewObjectID()
		user := models.User{
			ID:           userID,
			Name:         "FindByID User",
			Email:        "findbyid@test.com",
			PasswordHash: "hash123",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		_, err := repo.Create(user)
		assert.NoError(t, err)

		// Find the user
		found, err := repo.FindByID(userID)
		assert.NoError(t, err)
		assert.Equal(t, "FindByID User", found.Name)
	})

	t.Run("FindOne", func(t *testing.T) {
		// First create a user
		email := "findone@test.com"
		user := models.User{
			ID:           primitive.NewObjectID(),
			Name:         "FindOne User",
			Email:        email,
			PasswordHash: "hash123",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		_, err := repo.Create(user)
		assert.NoError(t, err)

		// Find by email
		found, err := repo.FindOne(bson.M{"email": email})
		assert.NoError(t, err)
		assert.Equal(t, "FindOne User", found.Name)
	})

	t.Run("FindAll", func(t *testing.T) {
		// Create multiple users
		for i := 0; i < 3; i++ {
			user := models.User{
				ID:           primitive.NewObjectID(),
				Name:         "FindAll User",
				Email:        "findall" + string(rune('0'+i)) + "@test.com",
				PasswordHash: "hash123",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			_, err := repo.Create(user)
			assert.NoError(t, err)
		}

		// Find all with filter
		users, err := repo.FindAll(bson.M{"name": "FindAll User"})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 3)
	})

	t.Run("Update", func(t *testing.T) {
		// Create a user
		userID := primitive.NewObjectID()
		user := models.User{
			ID:           userID,
			Name:         "Original Name",
			Email:        "update@test.com",
			PasswordHash: "hash123",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		_, err := repo.Create(user)
		assert.NoError(t, err)

		// Update the user
		update := bson.M{
			"$set": bson.M{
				"name": "Updated Name",
			},
		}

		updated, err := repo.Update(userID, update)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Name", updated.Name)
	})

	t.Run("Delete", func(t *testing.T) {
		// Create a user
		userID := primitive.NewObjectID()
		user := models.User{
			ID:           userID,
			Name:         "Delete User",
			Email:        "delete@test.com",
			PasswordHash: "hash123",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		_, err := repo.Create(user)
		assert.NoError(t, err)

		// Delete the user
		err = repo.Delete(userID)
		assert.NoError(t, err)

		// Verify deletion
		_, err = repo.FindByID(userID)
		assert.Error(t, err)
	})

	t.Run("Count", func(t *testing.T) {
		// Clear and create specific test data
		testutils.ClearCollection("test_users")

		for i := 0; i < 5; i++ {
			user := models.User{
				ID:           primitive.NewObjectID(),
				Name:         "Count User",
				Email:        "count" + string(rune('0'+i)) + "@test.com",
				PasswordHash: "hash123",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			_, err := repo.Create(user)
			assert.NoError(t, err)
		}

		// Count all
		count, err := repo.Count(bson.M{"name": "Count User"})
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})
}

// TestArrayRepositoryIntegration tests array repository operations
func TestArrayRepositoryIntegration(t *testing.T) {
	cleanup := testutils.SetupTestMongoDBWithSkip(t)
	defer cleanup()

	redisCleanup := testutils.SetupTestRedisWithSkip(t)
	defer redisCleanup()

	// Create board repo for testing array operations
	boardRepo := NewRepository[models.Board]("test_boards")
	arrayRepo := NewArrayRepository[models.Key]("test_boards")

	// First create a board
	boardID := primitive.NewObjectID()
	board := models.Board{
		ID:            boardID,
		AuthorId:      primitive.NewObjectID(),
		Name:          "Test Board",
		Layout:        models.GridLayout,
		Keys:          []models.Key{},
		Collaborators: []models.Collaborator{},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err := boardRepo.Create(board)
	assert.NoError(t, err)

	t.Run("CreateInArray", func(t *testing.T) {
		key := models.Key{
			ID:        primitive.NewObjectID(),
			BoardId:   boardID,
			Name:      "Test Key",
			HotKey:    "A",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		created, err := arrayRepo.CreateInArray(boardID, "keys", key)
		assert.NoError(t, err)
		assert.Equal(t, "Test Key", created.Name)
	})

	t.Run("GetArrayElement", func(t *testing.T) {
		// First add a key
		keyID := primitive.NewObjectID()
		key := models.Key{
			ID:        keyID,
			BoardId:   boardID,
			Name:      "Get Test Key",
			HotKey:    "B",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err := arrayRepo.CreateInArray(boardID, "keys", key)
		assert.NoError(t, err)

		// Get the key
		found, err := arrayRepo.GetArrayElement(boardID, "keys", keyID)
		assert.NoError(t, err)
		assert.Equal(t, "Get Test Key", found.Name)
	})

	t.Run("UpdateInArray", func(t *testing.T) {
		// First add a key
		keyID := primitive.NewObjectID()
		key := models.Key{
			ID:        keyID,
			BoardId:   boardID,
			Name:      "Original Key",
			HotKey:    "C",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err := arrayRepo.CreateInArray(boardID, "keys", key)
		assert.NoError(t, err)

		// Update the key
		update := bson.M{
			"name": "Updated Key",
		}

		err = arrayRepo.UpdateInArray(boardID, "keys", keyID, update)
		assert.NoError(t, err)

		// Verify update
		found, err := arrayRepo.GetArrayElement(boardID, "keys", keyID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Key", found.Name)
	})

	t.Run("DeleteFromArray", func(t *testing.T) {
		// First add a key
		keyID := primitive.NewObjectID()
		key := models.Key{
			ID:        keyID,
			BoardId:   boardID,
			Name:      "Delete Key",
			HotKey:    "D",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err := arrayRepo.CreateInArray(boardID, "keys", key)
		assert.NoError(t, err)

		// Delete the key
		err = arrayRepo.DeleteFromArray(boardID, "keys", keyID)
		assert.NoError(t, err)

		// Verify deletion
		_, err = arrayRepo.GetArrayElement(boardID, "keys", keyID)
		assert.Error(t, err)
	})

	t.Run("GetAllArrayElements", func(t *testing.T) {
		// Create a new board for this test
		newBoardID := primitive.NewObjectID()
		newBoard := models.Board{
			ID:            newBoardID,
			AuthorId:      primitive.NewObjectID(),
			Name:          "Test Board 2",
			Layout:        models.GridLayout,
			Keys:          []models.Key{},
			Collaborators: []models.Collaborator{},
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		_, err := boardRepo.Create(newBoard)
		assert.NoError(t, err)

		// Add multiple keys
		for i := 0; i < 3; i++ {
			key := models.Key{
				ID:        primitive.NewObjectID(),
				BoardId:   newBoardID,
				Name:      "Key " + string(rune('1'+i)),
				HotKey:    string(rune('E' + i)),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			_, err := arrayRepo.CreateInArray(newBoardID, "keys", key)
			assert.NoError(t, err)
		}

		// Get all keys
		keys, err := arrayRepo.GetAllArrayElements(newBoardID, "keys")
		assert.NoError(t, err)
		assert.Len(t, keys, 3)
	})
}

// Benchmark tests
func BenchmarkRepositoryCreate(b *testing.B) {
	cleanup, err := testutils.SetupTestMongoDB()
	if err != nil {
		b.Skipf("Skipping benchmark - MongoDB not available: %v", err)
	}
	defer cleanup()

	repo := NewRepository[models.User]("benchmark_users")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := models.User{
			ID:           primitive.NewObjectID(),
			Name:         "Benchmark User",
			Email:        "bench" + string(rune('0'+i%10)) + "@test.com",
			PasswordHash: "hash123",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		repo.Create(user)
	}
}

func BenchmarkRepositoryFindByID(b *testing.B) {
	cleanup, err := testutils.SetupTestMongoDB()
	if err != nil {
		b.Skipf("Skipping benchmark - MongoDB not available: %v", err)
	}
	defer cleanup()

	repo := NewRepository[models.User]("benchmark_users")

	// Create a user to find
	userID := primitive.NewObjectID()
	user := models.User{
		ID:           userID,
		Name:         "Benchmark User",
		Email:        "benchfind@test.com",
		PasswordHash: "hash123",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	repo.Create(user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.FindByID(userID)
	}
}
