package repositories

import (
	"time"
	"tuneslap/config"
	api "tuneslap/generated"
	"tuneslap/models"
	"tuneslap/utils"
	"tuneslap/validation"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MediaRepository struct {
	*Repository[models.Media]
	validator *validation.MediaValidator
}

func NewMediaRepository() *MediaRepository {
	return &MediaRepository{
		Repository: NewRepository[models.Media]("media"),
		validator:  validation.NewMediaValidator(),
	}
}

// GetValidator returns the media validator
func (r *MediaRepository) GetValidator() *validation.MediaValidator {
	return r.validator
}

// GetByAuthor retrieves all media for a specific author
func (r *MediaRepository) GetByAuthor(authorId primitive.ObjectID) ([]models.Media, error) {
	return r.FindAll(bson.M{"authorId": authorId})
}

// GetById retrieves media by ID for a specific author
func (r *MediaRepository) GetById(id primitive.ObjectID, authorId primitive.ObjectID) (models.Media, error) {
	return r.FindOne(bson.M{
		"_id":      id,
		"authorId": authorId,
	})
}

// GetByIdUnscoped retrieves media by ID without author filter
func (r *MediaRepository) GetByIdUnscoped(id primitive.ObjectID) (models.Media, error) {
	return r.FindOne(bson.M{
		"_id": id,
	})
}

// CreateMedia creates a new media
func (r *MediaRepository) CreateMedia(media *api.CreateMediaRequest, authorId primitive.ObjectID) (models.Media, error) {
	newMedia := utils.MediaFromCreateRequest(media, authorId)
	return r.Create(newMedia)
}

// UpdateMedia updates media by ID for a specific author
func (r *MediaRepository) UpdateMedia(id primitive.ObjectID, authorId primitive.ObjectID, updateData *api.UpdateMediaRequest) (models.Media, error) {
	update := bson.M{
		"$set": bson.M{
			"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	// Only update fields that are provided
	if updateData.Description != nil && *updateData.Description != "" {
		update["$set"].(bson.M)["description"] = updateData.GetDescription()
	}

	// First update the document
	_, err := r.Update(id, update)
	if err != nil {
		return models.Media{}, err
	}

	// Then fetch the updated document with author filter
	return r.FindOne(bson.M{
		"_id":      id,
		"authorId": authorId,
	})
}

// UpdateMediaUnscoped updates media by ID without author filter
// Only updates fields that have non-zero values to avoid overwriting existing data
func (r *MediaRepository) UpdateMediaUnscoped(id primitive.ObjectID, updateData *models.Media) (models.Media, error) {
	setFields := bson.M{
		"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
	}

	// Only set fields that have values
	if updateData.FileName != "" {
		setFields["fileName"] = updateData.FileName
	}
	if updateData.ProcessedUrl != "" {
		setFields["processedUrl"] = updateData.ProcessedUrl
	}
	if updateData.WaveformUrl != "" {
		setFields["waveformUrl"] = updateData.WaveformUrl
	}
	if updateData.ContentType != "" {
		setFields["contentType"] = updateData.ContentType
	}
	if updateData.FileSize != 0 {
		setFields["fileSize"] = updateData.FileSize
	}
	if updateData.Dimensions[0] != 0 || updateData.Dimensions[1] != 0 {
		setFields["dimensions"] = updateData.Dimensions
	}
	if updateData.Duration != 0 {
		setFields["duration"] = updateData.Duration
	}
	if updateData.Status != "" {
		setFields["status"] = updateData.Status
	}
	if updateData.ProcessingParams.Audio != nil || updateData.ProcessingParams.Image != nil {
		setFields["processingParams"] = updateData.ProcessingParams
	}
	if len(updateData.ProcessingActivity) > 0 {
		setFields["processingActivity"] = updateData.ProcessingActivity
	}

	update := bson.M{"$set": setFields}

	// First update the document
	_, err := r.Update(id, update)
	if err != nil {
		return models.Media{}, err
	}

	// Then fetch the updated document
	return r.FindOne(bson.M{
		"_id": id,
	})
}

// DeleteMedia deletes media by ID for a specific author
func (r *MediaRepository) DeleteMedia(id primitive.ObjectID, authorId primitive.ObjectID) error {
	// We need to check if the media exists and belongs to the author first
	_, err := r.FindOne(bson.M{
		"_id":      id,
		"authorId": authorId,
	})
	if err != nil {
		return err
	}

	return r.Delete(id)
}

// AggregateMedia performs aggregation on media
func (r *MediaRepository) AggregateMedia(pipeline interface{}, authorId primitive.ObjectID) ([]models.Media, int64, error) {
	return r.AggregateWithCount(pipeline, bson.M{"authorId": authorId})
}

// GetUrlsForMedia retrieves URLs for multiple media IDs
func (r *MediaRepository) GetUrlsForMedia(mediaIds []primitive.ObjectID) (map[primitive.ObjectID]string, error) {
	media, err := r.FindAll(bson.M{
		"_id": bson.M{"$in": mediaIds},
	})
	if err != nil {
		return nil, err
	}

	urls := make(map[primitive.ObjectID]string)
	for _, m := range media {
		urls[m.ID] = m.FileUrl
		if m.ProcessedUrl != "" {
			urls[m.ID] = m.ProcessedUrl
		} else {
			urls[m.ID] = m.FileUrl
		}
	}

	return urls, nil
}

// GetMyMediaStats calculates media statistics for a user
func (r *MediaRepository) GetMyMediaStats(authorId primitive.ObjectID) (models.MediaStats, error) {
	// Get all media for the user
	media, err := r.FindAll(bson.M{"authorId": authorId})
	if err != nil {
		return models.MediaStats{}, err
	}

	// Ensure media is not nil (defensive check)
	if media == nil {
		media = []models.Media{}
	}

	// Calculate stats from media
	imageCount := 0
	audioCount := 0
	totalStorage := int64(0)

	for _, m := range media {
		switch m.MediaType {
		case "image":
			imageCount++
		case "audio":
			audioCount++
		}
		totalStorage += m.FileSize
	}

	// Get configured max storage from environment variable
	maxStorage := config.GetMaxStorageBytes()
	var availableStorage int64

	if maxStorage == -1 {
		// Unlimited storage (not set or invalid)
		// Use -1 as sentinel value to indicate unlimited
		availableStorage = -1
	} else {
		// Calculate available storage: maxStorage - totalStorage
		availableStorage = maxStorage - totalStorage
		// Ensure availableStorage never goes below 0
		if availableStorage < 0 {
			availableStorage = 0
		}
	}

	return models.MediaStats{
		ImageCount:       imageCount,
		AudioCount:       audioCount,
		UsedStorage:      totalStorage,
		AvailableStorage: availableStorage,
	}, nil
}
