package repositories

import (
	"testing"
	"time"
	"tuneslap/models"
	"tuneslap/testutils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewUserRepository(t *testing.T) {
	repo := NewUserRepository()
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.Repository)
	assert.NotNil(t, repo.validator)
}

func TestUserRepository_GetValidator(t *testing.T) {
	repo := NewUserRepository()
	validator := repo.GetValidator()
	assert.NotNil(t, validator)
}

// Integration tests requiring MongoDB
func TestUserRepositoryIntegration(t *testing.T) {
	cleanup := testutils.SetupTestMongoDBWithSkip(t)
	defer cleanup()

	redisCleanup := testutils.SetupTestRedisWithSkip(t)
	defer redisCleanup()

	repo := NewUserRepository()

	t.Run("CreateUser", func(t *testing.T) {
		user := models.User{
			Name:         "Test User",
			Email:        "createuser@test.com",
			PasswordHash: "hashedpassword123",
		}

		created, err := repo.CreateUser(user)
		assert.NoError(t, err)
		assert.NotEqual(t, primitive.NilObjectID, created.ID)
		assert.Equal(t, "Test User", created.Name)
		assert.Equal(t, "createuser@test.com", created.Email)
		assert.NotZero(t, created.CreatedAt)
		assert.NotZero(t, created.UpdatedAt)
	})

	t.Run("GetByEmail", func(t *testing.T) {
		// Create a user first
		user := models.User{
			Name:         "Email Test User",
			Email:        "getbyemail@test.com",
			PasswordHash: "hashedpassword123",
		}

		_, err := repo.CreateUser(user)
		assert.NoError(t, err)

		// Get by email
		found, err := repo.GetByEmail("getbyemail@test.com")
		assert.NoError(t, err)
		assert.Equal(t, "Email Test User", found.Name)
		assert.Equal(t, "getbyemail@test.com", found.Email)
	})

	t.Run("GetByEmail_NotFound", func(t *testing.T) {
		_, err := repo.GetByEmail("nonexistent@test.com")
		assert.Error(t, err)
	})

	t.Run("GetMe", func(t *testing.T) {
		// Create a user first
		user := models.User{
			Name:         "GetMe Test User",
			Email:        "getme@test.com",
			PasswordHash: "hashedpassword123",
		}

		created, err := repo.CreateUser(user)
		assert.NoError(t, err)

		// Get by ID (GetMe)
		found, err := repo.GetMe(created.ID)
		assert.NoError(t, err)
		assert.Equal(t, "GetMe Test User", found.Name)
	})

	t.Run("GetByResetToken", func(t *testing.T) {
		// Create a user with reset token
		user := models.User{
			Name:         "Reset Token User",
			Email:        "resettoken@test.com",
			PasswordHash: "hashedpassword123",
		}

		created, err := repo.CreateUser(user)
		assert.NoError(t, err)

		// Update with reset token (keeping the name)
		token := "test-reset-token-123456789012345678901234"
		expiry := time.Now().Add(30 * time.Minute)
		_, err = repo.UpdateUser(created.ID, &models.UpdateUserRequest{
			Name:           "Reset Token User", // Keep the original name
			ResetToken:     &token,
			ResetExpiresAt: &expiry,
		})
		assert.NoError(t, err)

		// Get by reset token
		found, err := repo.GetByResetToken(token)
		assert.NoError(t, err)
		assert.Equal(t, "Reset Token User", found.Name)
		assert.NotNil(t, found.ResetToken)
		assert.Equal(t, token, *found.ResetToken)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		// Create a user first
		user := models.User{
			Name:         "Update Test User",
			Email:        "updateuser@test.com",
			PasswordHash: "hashedpassword123",
		}

		created, err := repo.CreateUser(user)
		assert.NoError(t, err)

		// Update user
		newName := "Updated Name"
		newImageUrl := "https://example.com/avatar.png"
		_, err = repo.UpdateUser(created.ID, &models.UpdateUserRequest{
			Name:     newName,
			ImageUrl: newImageUrl,
		})
		assert.NoError(t, err)

		// Verify update
		found, err := repo.GetMe(created.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Name", found.Name)
		assert.Equal(t, "https://example.com/avatar.png", found.ImageUrl)
	})
}

// Benchmark tests
func BenchmarkUserRepository_CreateUser(b *testing.B) {
	cleanup, err := testutils.SetupTestMongoDB()
	if err != nil {
		b.Skipf("Skipping benchmark - MongoDB not available: %v", err)
	}
	defer cleanup()

	repo := NewUserRepository()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := models.User{
			Name:         "Benchmark User",
			Email:        "bench" + string(rune('0'+i%10)) + "@test.com",
			PasswordHash: "hashedpassword123",
		}
		repo.CreateUser(user)
	}
}

func BenchmarkUserRepository_GetByEmail(b *testing.B) {
	cleanup, err := testutils.SetupTestMongoDB()
	if err != nil {
		b.Skipf("Skipping benchmark - MongoDB not available: %v", err)
	}
	defer cleanup()

	repo := NewUserRepository()

	// Create a user to find
	user := models.User{
		Name:         "Benchmark User",
		Email:        "benchemail@test.com",
		PasswordHash: "hashedpassword123",
	}
	repo.CreateUser(user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.GetByEmail("benchemail@test.com")
	}
}
