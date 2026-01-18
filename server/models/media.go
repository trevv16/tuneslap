package models

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ProcessingStatusPending    = "pending"
	ProcessingStatusProcessing = "processing"
	ProcessingStatusDone       = "done"
	ProcessingStatusError      = "error"
)

type ProcessingActivity struct {
	Status    string             `json:"status" bson:"status" validate:"required,oneof=pending processing done error"`
	Message   string             `json:"message" bson:"message"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

// MarshalJSON implements custom JSON marshaling to handle primitive.DateTime
func (pa ProcessingActivity) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Status    string    `json:"status"`
		Message   string    `json:"message"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	// Convert primitive.DateTime to time.Time, handling invalid dates
	createdAt := pa.CreatedAt.Time()
	updatedAt := pa.UpdatedAt.Time()

	// If times are invalid, use current time
	if createdAt.Year() < 1900 || createdAt.Year() > 9999 {
		createdAt = time.Now()
	}
	if updatedAt.Year() < 1900 || updatedAt.Year() > 9999 {
		updatedAt = time.Now()
	}

	aux := Alias{
		Status:    pa.Status,
		Message:   pa.Message,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return json.Marshal(aux)
}

type AudioProcessingParams struct {
	ContentType   string   `json:"contentType,omitempty"` // "audio/webm", "audio/mp3"
	TrimStart     float64  `json:"trimStart,omitempty"`   // seconds
	TrimEnd       float64  `json:"trimEnd,omitempty"`     // seconds
	Normalize     bool     `json:"normalize,omitempty"`
	FadeIn        float64  `json:"fadeIn,omitempty"`        // seconds
	FadeOut       float64  `json:"fadeOut,omitempty"`       // seconds
	Speed         float64  `json:"speed,omitempty"`         // multiplier (e.g., 1.25x)
	Pitch         float64  `json:"pitch,omitempty"`         // semitone ratio
	OutputFormats []string `json:"outputFormats,omitempty"` // e.g., ["webm", "mp3"]
}

type ImageProcessingParams struct {
	ResizeTo     [2]int `json:"resizeTo,omitempty"`     // [width, height]
	Format       int    `json:"format,omitempty"`       // "webp", "png"
	Crop         [4]int `json:"crop,omitempty"`         // [x, y, width, height]
	AspectRatio  string `json:"aspectRatio,omitempty"`  // "16:9", "1:1", etc.
	ApplyFilters string `json:"applyFilters,omitempty"` // ["grayscale", "blur", etc.]
}

type ProcessingParams struct {
	Audio *AudioProcessingParams `json:"audio,omitempty"`
	Image *ImageProcessingParams `json:"image,omitempty"`
}

type Media struct {
	ID                 primitive.ObjectID   `json:"id" bson:"_id" validate:"required"`
	AuthorId           primitive.ObjectID   `json:"authorId" bson:"authorId" validate:"required"`
	MediaType          string               `json:"mediaType" bson:"mediaType" validate:"required,oneof=image audio"`
	FileName           string               `json:"fileName" bson:"fileName" validate:"required,min=3,max=100,alphanumspace,excludesall=<>{}[]()&|\\"`
	Description        string               `json:"description,omitempty" bson:"description,omitempty" validate:"omitempty,max=1000"`
	FileUrl            string               `json:"fileUrl" bson:"fileUrl" validate:"required,url,excludesall=<>{}[]()&|\\"`
	ProcessedUrl       string               `json:"processedUrl" bson:"processedUrl" validate:"required,url,excludesall=<>{}[]()&|\\"`
	WaveformUrl        string               `json:"waveformUrl" bson:"waveformUrl" validate:"required,url,excludesall=<>{}[]()&|\\"`
	ContentType        string               `json:"contentType" bson:"contentType" validate:"required,min=3,max=100,alphanumspace,excludesall=<>{}[]()&|\\"`
	FileSize           int64                `json:"fileSize" bson:"fileSize" validate:"required,min=1,max=1000000000"`
	Status             string               `json:"status" bson:"status" validate:"required,oneof=pending processing done error"`
	ProcessingParams   ProcessingParams     `json:"processingParams,omitempty" bson:"processingParams,omitempty"`
	ProcessingActivity []ProcessingActivity `json:"processingActivity,omitempty" bson:"processingActivity,omitempty"`
	Dimensions         [2]int               `json:"dimensions,omitempty" bson:"dimensions,omitempty"`
	Duration           float64              `json:"duration,omitempty" bson:"duration,omitempty"`
	CreatedAt          primitive.DateTime   `json:"createdAt" bson:"createdAt"`
	UpdatedAt          primitive.DateTime   `json:"updatedAt" bson:"updatedAt"`
}

// MarshalJSON implements custom JSON marshaling to handle primitive.DateTime
func (m Media) MarshalJSON() ([]byte, error) {
	type Alias struct {
		ID                 primitive.ObjectID   `json:"id"`
		AuthorId           primitive.ObjectID   `json:"authorId"`
		MediaType          string               `json:"mediaType"`
		FileName           string               `json:"fileName"`
		Description        string               `json:"description,omitempty"`
		FileUrl            string               `json:"fileUrl"`
		ProcessedUrl       string               `json:"processedUrl"`
		WaveformUrl        string               `json:"waveformUrl"`
		ContentType        string               `json:"contentType"`
		FileSize           int64                `json:"fileSize"`
		Status             string               `json:"status"`
		ProcessingParams   ProcessingParams     `json:"processingParams,omitempty"`
		ProcessingActivity []ProcessingActivity `json:"processingActivity,omitempty"`
		Dimensions         [2]int               `json:"dimensions,omitempty"`
		Duration           float64              `json:"duration,omitempty"`
		CreatedAt          time.Time            `json:"createdAt"`
		UpdatedAt          time.Time            `json:"updatedAt"`
	}

	// Convert primitive.DateTime to time.Time, handling invalid dates
	createdAt := m.CreatedAt.Time()
	updatedAt := m.UpdatedAt.Time()

	// If times are invalid, use current time
	if createdAt.Year() < 1900 || createdAt.Year() > 9999 {
		createdAt = time.Now()
	}
	if updatedAt.Year() < 1900 || updatedAt.Year() > 9999 {
		updatedAt = time.Now()
	}

	// Filter out invalid ProcessingActivity entries
	var validActivities []ProcessingActivity
	for _, activity := range m.ProcessingActivity {
		// Check if the activity has valid dates
		activityCreatedAt := activity.CreatedAt.Time()
		activityUpdatedAt := activity.UpdatedAt.Time()

		if activityCreatedAt.Year() >= 1900 && activityCreatedAt.Year() <= 9999 &&
			activityUpdatedAt.Year() >= 1900 && activityUpdatedAt.Year() <= 9999 {
			validActivities = append(validActivities, activity)
		}
	}

	// If no valid activities, create a default one
	if len(validActivities) == 0 {
		validActivities = []ProcessingActivity{
			{
				Status:    ProcessingStatusPending,
				Message:   "Processing started",
				CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
				UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
			},
		}
	}

	aux := Alias{
		ID:                 m.ID,
		AuthorId:           m.AuthorId,
		MediaType:          m.MediaType,
		FileName:           m.FileName,
		Description:        m.Description,
		FileUrl:            m.FileUrl,
		ProcessedUrl:       m.ProcessedUrl,
		WaveformUrl:        m.WaveformUrl,
		ContentType:        m.ContentType,
		FileSize:           m.FileSize,
		Status:             m.Status,
		ProcessingParams:   m.ProcessingParams,
		ProcessingActivity: validActivities,
		Dimensions:         m.Dimensions,
		Duration:           m.Duration,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
	}

	return json.Marshal(aux)
}

type MediaStats struct {
	ImageCount       int   `json:"imageCount"`
	AudioCount       int   `json:"audioCount"`
	UsedStorage      int64 `json:"usedStorage"`
	AvailableStorage int64 `json:"availableStorage"`
}

type MediaUrl struct {
	ID  primitive.ObjectID `json:"id" bson:"_id"`
	Url string             `json:"url" bson:"url"`
}

// JSONResponse represents the JSON response format for Media
type MediaJSONResponse struct {
	ID                 primitive.ObjectID   `json:"id"`
	AuthorId           primitive.ObjectID   `json:"authorId"`
	MediaType          string               `json:"mediaType"`
	FileName           string               `json:"fileName"`
	Description        string               `json:"description,omitempty"`
	FileUrl            string               `json:"fileUrl"`
	ProcessedUrl       string               `json:"processedUrl"`
	WaveformUrl        string               `json:"waveformUrl"`
	ContentType        string               `json:"contentType"`
	FileSize           int64                `json:"fileSize"`
	Status             string               `json:"status"`
	ProcessingParams   ProcessingParams     `json:"processingParams,omitempty"`
	ProcessingActivity []ProcessingActivity `json:"processingActivity,omitempty"`
	Dimensions         [2]int               `json:"dimensions,omitempty"`
	Duration           float64              `json:"duration,omitempty"`
	CreatedAt          time.Time            `json:"createdAt"`
	UpdatedAt          time.Time            `json:"updatedAt"`
}

// ToJSONResponse converts Media to JSON response format
func (m Media) ToJSONResponse() MediaJSONResponse {
	// Convert primitive.DateTime to time.Time, handling invalid dates
	createdAt := m.CreatedAt.Time()
	updatedAt := m.UpdatedAt.Time()

	// If times are invalid, use current time
	if createdAt.Year() < 1900 || createdAt.Year() > 9999 {
		createdAt = time.Now()
	}
	if updatedAt.Year() < 1900 || updatedAt.Year() > 9999 {
		updatedAt = time.Now()
	}

	// Filter out invalid ProcessingActivity entries
	var validActivities []ProcessingActivity
	for _, activity := range m.ProcessingActivity {
		// Check if the activity has valid dates
		activityCreatedAt := activity.CreatedAt.Time()
		activityUpdatedAt := activity.UpdatedAt.Time()

		if activityCreatedAt.Year() >= 1900 && activityCreatedAt.Year() <= 9999 &&
			activityUpdatedAt.Year() >= 1900 && activityUpdatedAt.Year() <= 9999 {
			validActivities = append(validActivities, activity)
		}
	}

	// If no valid activities, create a default one
	if len(validActivities) == 0 {
		validActivities = []ProcessingActivity{
			{
				Status:    ProcessingStatusPending,
				Message:   "Processing started",
				CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
				UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
			},
		}
	}

	return MediaJSONResponse{
		ID:                 m.ID,
		AuthorId:           m.AuthorId,
		MediaType:          m.MediaType,
		FileName:           m.FileName,
		Description:        m.Description,
		FileUrl:            m.FileUrl,
		ProcessedUrl:       m.ProcessedUrl,
		WaveformUrl:        m.WaveformUrl,
		ContentType:        m.ContentType,
		FileSize:           m.FileSize,
		Status:             m.Status,
		ProcessingParams:   m.ProcessingParams,
		ProcessingActivity: validActivities,
		Dimensions:         m.Dimensions,
		Duration:           m.Duration,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
	}
}
