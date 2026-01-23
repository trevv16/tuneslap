package repositories

import (
	"testing"

	api "tuneslap/generated"
	"tuneslap/models"
	"tuneslap/testutils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewMediaRepository(t *testing.T) {
	repo := NewMediaRepository()
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.Repository)
	assert.NotNil(t, repo.validator)
}

func TestMediaRepository_GetValidator(t *testing.T) {
	repo := NewMediaRepository()
	validator := repo.GetValidator()
	assert.NotNil(t, validator)
}

// Integration tests requiring MongoDB
func TestMediaRepositoryIntegration(t *testing.T) {
	cleanup := testutils.SetupTestMongoDBWithSkip(t)
	defer cleanup()

	redisCleanup := testutils.SetupTestRedisWithSkip(t)
	defer redisCleanup()

	repo := NewMediaRepository()
	authorID := primitive.NewObjectID()

	t.Run("CreateMedia", func(t *testing.T) {
		mediaType := "audio"
		fileName := "test.mp3"
		fileUrl := "https://example.com/test.mp3"
		contentType := "audio/mpeg"
		fileSize := int32(1024)

		req := &api.CreateMediaRequest{
			MediaType:   &mediaType,
			FileName:    &fileName,
			FileUrl:     &fileUrl,
			ContentType: &contentType,
			FileSize:    &fileSize,
		}

		created, err := repo.CreateMedia(req, authorID)
		assert.NoError(t, err)
		assert.NotEqual(t, primitive.NilObjectID, created.ID)
		assert.Equal(t, "audio", created.MediaType)
		assert.Equal(t, "test.mp3", created.FileName)
		assert.Equal(t, models.ProcessingStatusPending, created.Status)
	})

	t.Run("GetByAuthor", func(t *testing.T) {
		// Create some media
		for i := 0; i < 3; i++ {
			mediaType := "audio"
			fileName := "author_test_" + string(rune('1'+i)) + ".mp3"
			fileUrl := "https://example.com/" + fileName
			contentType := "audio/mpeg"
			fileSize := int32(1024)

			req := &api.CreateMediaRequest{
				MediaType:   &mediaType,
				FileName:    &fileName,
				FileUrl:     &fileUrl,
				ContentType: &contentType,
				FileSize:    &fileSize,
			}
			_, err := repo.CreateMedia(req, authorID)
			assert.NoError(t, err)
		}

		// Get by author
		media, err := repo.GetByAuthor(authorID)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(media), 3)
	})

	t.Run("GetById", func(t *testing.T) {
		// Create media
		mediaType := "image"
		fileName := "getbyid.png"
		fileUrl := "https://example.com/getbyid.png"
		contentType := "image/png"
		fileSize := int32(2048)

		req := &api.CreateMediaRequest{
			MediaType:   &mediaType,
			FileName:    &fileName,
			FileUrl:     &fileUrl,
			ContentType: &contentType,
			FileSize:    &fileSize,
		}

		created, err := repo.CreateMedia(req, authorID)
		assert.NoError(t, err)

		// Get by ID
		found, err := repo.GetById(created.ID, authorID)
		assert.NoError(t, err)
		assert.Equal(t, "getbyid.png", found.FileName)
		assert.Equal(t, "image", found.MediaType)
	})

	t.Run("GetByIdUnscoped", func(t *testing.T) {
		// Create media
		mediaType := "audio"
		fileName := "unscoped.mp3"
		fileUrl := "https://example.com/unscoped.mp3"
		contentType := "audio/mpeg"
		fileSize := int32(1024)

		req := &api.CreateMediaRequest{
			MediaType:   &mediaType,
			FileName:    &fileName,
			FileUrl:     &fileUrl,
			ContentType: &contentType,
			FileSize:    &fileSize,
		}

		created, err := repo.CreateMedia(req, authorID)
		assert.NoError(t, err)

		// Get unscoped (different author can access)
		found, err := repo.GetByIdUnscoped(created.ID)
		assert.NoError(t, err)
		assert.Equal(t, "unscoped.mp3", found.FileName)
	})

	t.Run("UpdateMedia", func(t *testing.T) {
		// Create media
		mediaType := "audio"
		fileName := "update.mp3"
		fileUrl := "https://example.com/update.mp3"
		contentType := "audio/mpeg"
		fileSize := int32(1024)

		req := &api.CreateMediaRequest{
			MediaType:   &mediaType,
			FileName:    &fileName,
			FileUrl:     &fileUrl,
			ContentType: &contentType,
			FileSize:    &fileSize,
		}

		created, err := repo.CreateMedia(req, authorID)
		assert.NoError(t, err)

		// Update media
		newDesc := "Updated description"
		updateReq := &api.UpdateMediaRequest{
			Description: &newDesc,
		}

		updated, err := repo.UpdateMedia(created.ID, authorID, updateReq)
		assert.NoError(t, err)
		assert.Equal(t, "Updated description", updated.Description)
	})

	t.Run("UpdateMediaUnscoped", func(t *testing.T) {
		// Create media
		mediaType := "audio"
		fileName := "unscoped_update.mp3"
		fileUrl := "https://example.com/unscoped_update.mp3"
		contentType := "audio/mpeg"
		fileSize := int32(1024)

		req := &api.CreateMediaRequest{
			MediaType:   &mediaType,
			FileName:    &fileName,
			FileUrl:     &fileUrl,
			ContentType: &contentType,
			FileSize:    &fileSize,
		}

		created, err := repo.CreateMedia(req, authorID)
		assert.NoError(t, err)

		// Update unscoped
		updateData := &models.Media{
			Status:       models.ProcessingStatusDone,
			ProcessedUrl: "https://example.com/processed.webm",
			Duration:     120.5,
		}

		updated, err := repo.UpdateMediaUnscoped(created.ID, updateData)
		assert.NoError(t, err)
		assert.Equal(t, models.ProcessingStatusDone, updated.Status)
		assert.Equal(t, "https://example.com/processed.webm", updated.ProcessedUrl)
	})

	t.Run("DeleteMedia", func(t *testing.T) {
		// Create media
		mediaType := "audio"
		fileName := "delete.mp3"
		fileUrl := "https://example.com/delete.mp3"
		contentType := "audio/mpeg"
		fileSize := int32(1024)

		req := &api.CreateMediaRequest{
			MediaType:   &mediaType,
			FileName:    &fileName,
			FileUrl:     &fileUrl,
			ContentType: &contentType,
			FileSize:    &fileSize,
		}

		created, err := repo.CreateMedia(req, authorID)
		assert.NoError(t, err)

		// Delete media
		err = repo.DeleteMedia(created.ID, authorID)
		assert.NoError(t, err)

		// Verify deletion
		_, err = repo.GetById(created.ID, authorID)
		assert.Error(t, err)
	})

	t.Run("GetUrlsForMedia", func(t *testing.T) {
		// Create multiple media items
		var mediaIds []primitive.ObjectID
		for i := 0; i < 3; i++ {
			mediaType := "audio"
			fileName := "url_test_" + string(rune('1'+i)) + ".mp3"
			fileUrl := "https://example.com/" + fileName
			contentType := "audio/mpeg"
			fileSize := int32(1024)

			req := &api.CreateMediaRequest{
				MediaType:   &mediaType,
				FileName:    &fileName,
				FileUrl:     &fileUrl,
				ContentType: &contentType,
				FileSize:    &fileSize,
			}

			created, err := repo.CreateMedia(req, authorID)
			assert.NoError(t, err)
			mediaIds = append(mediaIds, created.ID)
		}

		// Get URLs
		urls, err := repo.GetUrlsForMedia(mediaIds)
		assert.NoError(t, err)
		assert.Len(t, urls, 3)

		for _, id := range mediaIds {
			_, exists := urls[id]
			assert.True(t, exists)
		}
	})

	t.Run("GetMyMediaStats", func(t *testing.T) {
		// Create a new author for clean stats
		statsAuthorID := primitive.NewObjectID()

		// Create some audio media
		for i := 0; i < 2; i++ {
			mediaType := "audio"
			fileName := "stats_audio_" + string(rune('1'+i)) + ".mp3"
			fileUrl := "https://example.com/" + fileName
			contentType := "audio/mpeg"
			fileSize := int32(1000)

			req := &api.CreateMediaRequest{
				MediaType:   &mediaType,
				FileName:    &fileName,
				FileUrl:     &fileUrl,
				ContentType: &contentType,
				FileSize:    &fileSize,
			}
			_, err := repo.CreateMedia(req, statsAuthorID)
			assert.NoError(t, err)
		}

		// Create some image media
		for i := 0; i < 3; i++ {
			mediaType := "image"
			fileName := "stats_image_" + string(rune('1'+i)) + ".png"
			fileUrl := "https://example.com/" + fileName
			contentType := "image/png"
			fileSize := int32(500)

			req := &api.CreateMediaRequest{
				MediaType:   &mediaType,
				FileName:    &fileName,
				FileUrl:     &fileUrl,
				ContentType: &contentType,
				FileSize:    &fileSize,
			}
			_, err := repo.CreateMedia(req, statsAuthorID)
			assert.NoError(t, err)
		}

		// Get stats
		stats, err := repo.GetMyMediaStats(statsAuthorID)
		assert.NoError(t, err)
		assert.Equal(t, 2, stats.AudioCount)
		assert.Equal(t, 3, stats.ImageCount)
		assert.Equal(t, int64(3500), stats.UsedStorage) // 2*1000 + 3*500
	})
}

// Benchmark tests
func BenchmarkMediaRepository_CreateMedia(b *testing.B) {
	cleanup, err := testutils.SetupTestMongoDB()
	if err != nil {
		b.Skipf("Skipping benchmark - MongoDB not available: %v", err)
	}
	defer cleanup()

	repo := NewMediaRepository()
	authorID := primitive.NewObjectID()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mediaType := "audio"
		fileName := "benchmark.mp3"
		fileUrl := "https://example.com/benchmark.mp3"
		contentType := "audio/mpeg"
		fileSize := int32(1024)

		req := &api.CreateMediaRequest{
			MediaType:   &mediaType,
			FileName:    &fileName,
			FileUrl:     &fileUrl,
			ContentType: &contentType,
			FileSize:    &fileSize,
		}
		repo.CreateMedia(req, authorID)
	}
}

func BenchmarkMediaRepository_GetMyMediaStats(b *testing.B) {
	cleanup, err := testutils.SetupTestMongoDB()
	if err != nil {
		b.Skipf("Skipping benchmark - MongoDB not available: %v", err)
	}
	defer cleanup()

	repo := NewMediaRepository()
	authorID := primitive.NewObjectID()

	// Create some media for stats
	for i := 0; i < 20; i++ {
		var mediaType string
		if i%2 == 0 {
			mediaType = "audio"
		} else {
			mediaType = "image"
		}
		fileName := "bench_" + string(rune('0'+i%10)) + ".mp3"
		fileUrl := "https://example.com/" + fileName
		contentType := "audio/mpeg"
		fileSize := int32(1024)

		req := &api.CreateMediaRequest{
			MediaType:   &mediaType,
			FileName:    &fileName,
			FileUrl:     &fileUrl,
			ContentType: &contentType,
			FileSize:    &fileSize,
		}
		repo.CreateMedia(req, authorID)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.GetMyMediaStats(authorID)
	}
}
