package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"
	"tuneslap/config"
	api "tuneslap/generated"
	"tuneslap/models"
	"tuneslap/repositories"
	"tuneslap/services/storage"
	"tuneslap/tasks"
	"tuneslap/utils"
	"tuneslap/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MediaHandler provides media-specific operations using the base handler
type MediaHandler struct {
	*BaseHandler[models.Media, api.CreateMediaRequest, api.UpdateMediaRequest]
	mediaRepo *repositories.MediaRepository
}

// NewMediaHandler creates a new media handler
func NewMediaHandler() *MediaHandler {
	mediaRepo := repositories.NewMediaRepository()
	mediaValidator := validation.NewMediaValidator()

	return &MediaHandler{
		BaseHandler: NewBaseHandler[models.Media, api.CreateMediaRequest, api.UpdateMediaRequest](
			mediaRepo.Repository,
			mediaValidator,
		),
		mediaRepo: mediaRepo,
	}
}

// MediaResponseMapper maps media model to response format
func (h *MediaHandler) MediaResponseMapper(media models.Media) interface{} {
	return utils.ToMediaResponse(media)
}

// deriveContentTypeFromFileName derives content type from file extension
func deriveContentTypeFromFileName(fileName *string) string {
	if fileName == nil || *fileName == "" {
		return "application/octet-stream"
	}

	lowerName := strings.ToLower(*fileName)

	// Audio types
	if strings.HasSuffix(lowerName, ".mp3") {
		return "audio/mpeg"
	}
	if strings.HasSuffix(lowerName, ".wav") {
		return "audio/x-wav"
	}
	if strings.HasSuffix(lowerName, ".webm") {
		return "audio/webm"
	}
	if strings.HasSuffix(lowerName, ".ogg") {
		return "audio/ogg"
	}
	if strings.HasSuffix(lowerName, ".aac") {
		return "audio/aac"
	}

	// Image types
	if strings.HasSuffix(lowerName, ".jpg") || strings.HasSuffix(lowerName, ".jpeg") {
		return "image/jpeg"
	}
	if strings.HasSuffix(lowerName, ".png") {
		return "image/png"
	}
	if strings.HasSuffix(lowerName, ".gif") {
		return "image/gif"
	}
	if strings.HasSuffix(lowerName, ".webp") {
		return "image/webp"
	}
	if strings.HasSuffix(lowerName, ".svg") {
		return "image/svg+xml"
	}

	// Default
	return "application/octet-stream"
}

// HandleGetAllMedia handles GET all media with filtering
func (h *MediaHandler) HandleGetAllMedia(c *fiber.Ctx) error {
	return h.HandleGetAll(
		c,
		func(authorId primitive.ObjectID) bson.M {
			filter := bson.M{"authorId": authorId}

			// Add mediaType filter if provided
			if mediaType := c.Query("mediaType"); mediaType != "" {
				if mediaType == "image" || mediaType == "audio" {
					filter["mediaType"] = mediaType
				}
			}

			// Add contentType filter if provided
			if contentType := c.Query("contentType"); contentType != "" {
				validContentTypes := []string{
					"image/jpeg", "image/png", "image/gif", "image/webp",
					"audio/mpeg", "audio/mp3", "audio/mp4", "audio/wav", "audio/ogg", "audio/webm",
				}
				for _, validType := range validContentTypes {
					if contentType == validType {
						filter["contentType"] = contentType
						break
					}
				}
			}

			return filter
		},
		h.MediaResponseMapper,
		nil,
		"media", // dataFieldName to match OpenAPI spec
	)
}

// HandleGetMediaById handles GET media by ID
func (h *MediaHandler) HandleGetMediaById(c *fiber.Ctx) error {
	authorId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}

	return h.HandleGetById(c, "mediaId", authorId, h.MediaResponseMapper, nil)
}

// HandleCreateMedia handles CREATE media with storage validation
func (h *MediaHandler) HandleCreateMedia(c *fiber.Ctx) error {
	authorId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}

	return h.HandleCreate(c, authorId, func(req api.CreateMediaRequest, authorId primitive.ObjectID) (models.Media, error) {
		// Derive contentType from fileName if not provided
		if req.ContentType == nil || *req.ContentType == "" {
			contentType := deriveContentTypeFromFileName(req.FileName)
			req.ContentType = &contentType
		}

		// Check storage limit if configured
		mediaStats, err := h.mediaRepo.GetMyMediaStats(authorId)
		if err != nil {
			return models.Media{}, err
		}

		// Demo mode limits
		if config.IsDemoMode() {
			// Check file size limit (10MB in demo mode)
			if req.FileSize != nil && int64(*req.FileSize) > config.DemoMaxFileSize {
				return models.Media{}, fmt.Errorf("file size (%d bytes) exceeds demo mode limit (%d bytes)", *req.FileSize, config.DemoMaxFileSize)
			}
			// Check media count limit (5 uploads in demo mode)
			totalCount := mediaStats.ImageCount + mediaStats.AudioCount
			if totalCount >= config.DemoMaxMediaCount {
				return models.Media{}, fmt.Errorf("demo mode limit reached: maximum %d uploads allowed", config.DemoMaxMediaCount)
			}
		}

		// If storage limit is set (availableStorage is not -1 for unlimited), validate
		if mediaStats.AvailableStorage != -1 {
			if req.FileSize != nil && int64(*req.FileSize) > mediaStats.AvailableStorage {
				return models.Media{}, fmt.Errorf("file size (%d bytes) exceeds available storage (%d bytes)", *req.FileSize, mediaStats.AvailableStorage)
			}
		}

		// Create the media using repository with authorId
		createdMedia, err := h.mediaRepo.CreateMedia(&req, authorId)
		if err != nil {
			// Clean up storage if database insert fails
			storage.DeleteMedia(createdMedia)
			return models.Media{}, err
		}

		return createdMedia, nil
	}, h.MediaResponseMapper)
}

// HandleUpdateMedia handles UPDATE media
func (h *MediaHandler) HandleUpdateMedia(c *fiber.Ctx) error {
	authorId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}

	return h.HandleUpdate(c, "mediaId", authorId, func(id primitive.ObjectID, authorId primitive.ObjectID, req api.UpdateMediaRequest) (models.Media, error) {
		return h.mediaRepo.UpdateMedia(id, authorId, &req)
	}, h.MediaResponseMapper)
}

// HandleDeleteMedia handles DELETE media with S3 cleanup
func (h *MediaHandler) HandleDeleteMedia(c *fiber.Ctx) error {
	authorId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}

	mediaId, err := GetValidObjectId(c, "mediaId")
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid media ID", err)
	}

	// Get the media first to clean up S3 using repository
	media, err := h.mediaRepo.GetById(mediaId, authorId)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusNotFound, "Media not found", err)
	}

	// Delete from storage first
	err = storage.DeleteMedia(media)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete media from storage", err)
	}

	// Then delete from database
	return h.HandleDelete(c, "mediaId", authorId)
}

// HandleGetMyMediaStats handles GET media stats (specialized operation)
func (h *MediaHandler) HandleGetMyMediaStats(c *fiber.Ctx) error {
	authorId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}

	mediaStats, err := h.mediaRepo.GetMyMediaStats(authorId)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Error getting media stats", err)
	}

	return SendSuccessResponse(c, fiber.StatusOK, "Media stats retrieved", mediaStats)
}

// HandleProcessMedia handles media processing (specialized operation)
func (h *MediaHandler) HandleProcessMedia(c *fiber.Ctx) error {
	mediaId, err := GetValidObjectId(c, "mediaId")
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid media ID", err)
	}

	authorId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}

	// Get the media first to check if it exists using repository
	_, err = h.mediaRepo.GetById(mediaId, authorId)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusNotFound, "Media not found", err)
	}

	// Parse request body using generated type
	var body api.MediaProcessingParams
	if err := c.BodyParser(&body); err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid processing params", err)
	}

	// Validate audio processing params if present
	if body.Audio != nil {
		validationResult := validation.ValidateMediaProcessingParamsAudio(body.Audio)
		if !validationResult.IsValid {
			return SendValidationErrorResponse(c, validationResult)
		}
	}

	// Convert to model processing params using converter
	processingParams := utils.ProcessingParamsFromAPI(&body)

	// Create new media task
	task, err := tasks.NewMediaProcessTask(mediaId, authorId, processingParams)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create processing task", err)
	}

	// Get the asynq client
	client := tasks.GetClient()

	// Enqueue task with options
	info, err := client.Enqueue(task, asynq.Queue("media"), asynq.MaxRetry(5), asynq.Timeout(60*time.Second))
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to enqueue task", err)
	}

	fmt.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	return SendSuccessResponse(c, fiber.StatusCreated, "Processing job started", nil)
}

// HandleGenerateUploadURL handles generation of signed upload URLs
func (h *MediaHandler) HandleGenerateUploadURL(c *fiber.Ctx) error {
	authorId, err := GetAuthorId(c)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid author ID", err)
	}

	// Parse request body
	requestBody := new(struct {
		FileName    string `json:"fileName" validate:"required,min=3,max=100"`
		ContentType string `json:"contentType" validate:"required"`
		FileSize    int64  `json:"fileSize" validate:"required,min=1,max=1000000000"`
	})
	if err := c.BodyParser(requestBody); err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate the request body using base validator (struct tag validation)
	baseValidator := validation.NewValidator()
	validationResult := baseValidator.Validate(requestBody)
	if !validationResult.IsValid {
		return SendValidationErrorResponse(c, validationResult)
	}

	// Check storage limit if configured
	mediaStats, err := h.mediaRepo.GetMyMediaStats(authorId)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Error getting media stats", err)
	}

	// Demo mode limits
	if config.IsDemoMode() {
		// Check file size limit (10MB in demo mode)
		if requestBody.FileSize > config.DemoMaxFileSize {
			return SendErrorResponse(c, fiber.StatusBadRequest, fmt.Sprintf("File size (%d bytes) exceeds demo mode limit (%d bytes)", requestBody.FileSize, config.DemoMaxFileSize), nil)
		}
		// Check media count limit (5 uploads in demo mode)
		totalCount := mediaStats.ImageCount + mediaStats.AudioCount
		if totalCount >= config.DemoMaxMediaCount {
			return SendErrorResponse(c, fiber.StatusBadRequest, fmt.Sprintf("Demo mode limit reached: maximum %d uploads allowed", config.DemoMaxMediaCount), nil)
		}
	}

	// If storage limit is set (availableStorage is not -1 for unlimited), validate
	if mediaStats.AvailableStorage != -1 {
		if requestBody.FileSize > mediaStats.AvailableStorage {
			return SendErrorResponse(c, fiber.StatusBadRequest, fmt.Sprintf("File size (%d bytes) exceeds available storage (%d bytes)", requestBody.FileSize, mediaStats.AvailableStorage), nil)
		}
	}

	// Determine media type from content type
	mediaType := "audio"
	if requestBody.ContentType != "" && requestBody.ContentType[:5] == "image" {
		mediaType = "image"
	}

	// Create a temporary media object to generate the key
	tempMedia := models.Media{
		AuthorId:  authorId,
		MediaType: mediaType,
		FileName:  requestBody.FileName,
	}

	// Generate object key
	objectKey := storage.GetMediaKey(tempMedia.AuthorId.Hex(), tempMedia.MediaType, tempMedia.FileName)

	// Get storage client
	storageClient, err := storage.GetUserUploadsStorage()
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to initialize storage client", err)
	}

	// Generate signed URL (valid for 15 minutes)
	signedURL, err := storageClient.GenerateSignedUploadURL(
		context.Background(),
		objectKey,
		requestBody.ContentType,
		15*time.Minute,
	)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate upload URL", err)
	}

	// Get the actual file URL that will be used to access the file after upload
	fileURL := storageClient.GetFileURL(objectKey)

	// Return the signed URL, object key, bucket name, and the actual file URL
	response := map[string]interface{}{
		"signedUrl":  signedURL,
		"objectKey":  objectKey,
		"bucketName": storageClient.GetBucketName(),
		"fileUrl":    fileURL, // Return the actual file URL to store
	}

	return SendSuccessResponse(c, fiber.StatusOK, "Upload URL generated", response)
}
