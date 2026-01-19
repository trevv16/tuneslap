package utils

import (
	"time"

	api "tuneslap/generated"
	"tuneslap/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Helper functions for ObjectID conversion
func objectIDToString(id primitive.ObjectID) *string {
	if id.IsZero() {
		return nil
	}
	s := id.Hex()
	return &s
}

func stringToObjectID(s *string) (primitive.ObjectID, error) {
	if s == nil || *s == "" {
		return primitive.NilObjectID, nil
	}
	return primitive.ObjectIDFromHex(*s)
}

// Helper functions for DateTime conversion
func dateTimeToTime(dt primitive.DateTime) *time.Time {
	t := dt.Time()
	if t.Year() < 1900 || t.Year() > 9999 {
		now := time.Now()
		return &now
	}
	return &t
}

func timeToDateTime(t *time.Time) primitive.DateTime {
	if t == nil {
		return primitive.NewDateTimeFromTime(time.Now())
	}
	return primitive.NewDateTimeFromTime(*t)
}

// Helper function for string pointer conversion
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func stringVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Helper function for float conversion
func float64ToFloat32(f float64) *float32 {
	val := float32(f)
	return &val
}

func float32ToFloat64(f *float32) float64 {
	if f == nil {
		return 0
	}
	return float64(*f)
}

// Helper function for int conversion
func intToInt32(i int) *int32 {
	val := int32(i)
	return &val
}

func int32ToInt(i *int32) int {
	if i == nil {
		return 0
	}
	return int(*i)
}

// Helper function for int array conversion
func intArrayToInt32Array(arr [2]int) []int32 {
	if len(arr) == 0 {
		return nil
	}
	result := make([]int32, len(arr))
	for i, v := range arr {
		result[i] = int32(v)
	}
	return result
}

func int32ArrayToIntArray(arr []int32) [2]int {
	if len(arr) < 2 {
		return [2]int{0, 0}
	}
	return [2]int{int(arr[0]), int(arr[1])}
}

func int32ArrayToInt4Array(arr []int32) [4]int {
	if len(arr) < 4 {
		return [4]int{0, 0, 0, 0}
	}
	return [4]int{int(arr[0]), int(arr[1]), int(arr[2]), int(arr[3])}
}

// User conversions
func ToUserResponse(user models.User) *api.UserResponse {
	return &api.UserResponse{
		Id:        objectIDToString(user.ID),
		Name:      stringPtr(user.Name),
		Email:     stringPtr(user.Email),
		ImageUrl:  stringPtr(user.ImageUrl),
		CreatedAt: dateTimeToTime(primitive.NewDateTimeFromTime(user.CreatedAt)),
		UpdatedAt: dateTimeToTime(primitive.NewDateTimeFromTime(user.UpdatedAt)),
	}
}

// Collaborator conversions
func ToCollaboratorResponse(collaborator models.Collaborator) *api.CollaboratorResponse {
	return &api.CollaboratorResponse{
		Id:        objectIDToString(collaborator.ID),
		UserId:    objectIDToString(collaborator.UserId),
		Email:     stringPtr(collaborator.Email),
		Name:      nil, // Name is not in internal Collaborator model
		Role:      stringPtr(collaborator.Role),
		ImageUrl:  nil, // ImageUrl is not in internal Collaborator model
		CreatedAt: dateTimeToTime(primitive.NewDateTimeFromTime(collaborator.CreatedAt)),
		UpdatedAt: dateTimeToTime(primitive.NewDateTimeFromTime(collaborator.UpdatedAt)),
	}
}

// Key conversions
func ToKeyResponse(key models.Key) *api.KeyResponse {
	return &api.KeyResponse{
		Id:           objectIDToString(key.ID),
		BoardId:      objectIDToString(key.BoardId),
		Name:         stringPtr(key.Name),
		Description:  stringPtr(key.Description),
		AudioMediaId: objectIDToString(key.AudioMediaId),
		AudioUrl:     stringPtr(key.AudioUrl),
		ImageMediaId: objectIDToString(key.ImageMediaId),
		ImageUrl:     stringPtr(key.ImageUrl),
		HotKey:       stringPtr(key.HotKey),
		CreatedAt:    dateTimeToTime(primitive.NewDateTimeFromTime(key.CreatedAt)),
		UpdatedAt:    dateTimeToTime(primitive.NewDateTimeFromTime(key.UpdatedAt)),
	}
}

// ProcessingParams conversions
func ToMediaProcessingParamsAudio(params models.AudioProcessingParams) *api.MediaProcessingParamsAudio {
	return &api.MediaProcessingParamsAudio{
		ContentType:   stringPtr(params.ContentType),
		TrimStart:     float64ToFloat32(params.TrimStart),
		TrimEnd:       float64ToFloat32(params.TrimEnd),
		Normalize:     &params.Normalize,
		FadeIn:        float64ToFloat32(params.FadeIn),
		FadeOut:       float64ToFloat32(params.FadeOut),
		Speed:         float64ToFloat32(params.Speed),
		Pitch:         float64ToFloat32(params.Pitch),
		OutputFormats: params.OutputFormats,
	}
}

func ToMediaProcessingParamsImage(params models.ImageProcessingParams) *api.MediaProcessingParamsImage {
	formatStr := ""
	if params.Format > 0 {
		// Format is stored as int in internal model, need to convert to string
		// This might need adjustment based on actual format values
		formatStr = string(rune(params.Format))
	}

	resizeTo := make([]int32, 2)
	resizeTo[0] = int32(params.ResizeTo[0])
	resizeTo[1] = int32(params.ResizeTo[1])

	crop := make([]int32, 4)
	crop[0] = int32(params.Crop[0])
	crop[1] = int32(params.Crop[1])
	crop[2] = int32(params.Crop[2])
	crop[3] = int32(params.Crop[3])

	applyFilters := []string{}
	if params.ApplyFilters != "" {
		applyFilters = []string{params.ApplyFilters}
	}

	return &api.MediaProcessingParamsImage{
		ResizeTo:     resizeTo,
		Format:       stringPtr(formatStr),
		Crop:         crop,
		AspectRatio:  stringPtr(params.AspectRatio),
		ApplyFilters: applyFilters,
	}
}

func ToMediaProcessingParams(params models.ProcessingParams) *api.MediaProcessingParams {
	result := &api.MediaProcessingParams{}

	if params.Audio != nil {
		result.Audio = ToMediaProcessingParamsAudio(*params.Audio)
	}
	if params.Image != nil {
		result.Image = ToMediaProcessingParamsImage(*params.Image)
	}

	return result
}

func ProcessingParamsFromAPI(params *api.MediaProcessingParams) models.ProcessingParams {
	result := models.ProcessingParams{}

	if params != nil {
		if params.Audio != nil {
			audio := models.AudioProcessingParams{
				ContentType:   stringVal(params.Audio.ContentType),
				TrimStart:     float32ToFloat64(params.Audio.TrimStart),
				TrimEnd:       float32ToFloat64(params.Audio.TrimEnd),
				Normalize:     params.Audio.Normalize != nil && *params.Audio.Normalize,
				FadeIn:        float32ToFloat64(params.Audio.FadeIn),
				FadeOut:       float32ToFloat64(params.Audio.FadeOut),
				Speed:         float32ToFloat64(params.Audio.Speed),
				Pitch:         float32ToFloat64(params.Audio.Pitch),
				OutputFormats: params.Audio.OutputFormats,
			}
			result.Audio = &audio
		}
		if params.Image != nil {
			crop := [4]int{0, 0, 0, 0}
			if len(params.Image.Crop) >= 4 {
				crop = int32ArrayToInt4Array(params.Image.Crop)
			}
			resizeTo := [2]int{0, 0}
			if len(params.Image.ResizeTo) >= 2 {
				resizeTo = int32ArrayToIntArray(params.Image.ResizeTo)
			}
			formatInt := 0
			if params.Image.Format != nil && *params.Image.Format != "" {
				// Convert format string to int if needed
				formatInt = int((*params.Image.Format)[0])
			}
			applyFilters := ""
			if len(params.Image.ApplyFilters) > 0 {
				applyFilters = params.Image.ApplyFilters[0]
			}
			image := models.ImageProcessingParams{
				ResizeTo:     resizeTo,
				Format:       formatInt,
				Crop:         crop,
				AspectRatio:  stringVal(params.Image.AspectRatio),
				ApplyFilters: applyFilters,
			}
			result.Image = &image
		}
	}

	return result
}

// ProcessingActivity conversions
func ToMediaProcessingActivity(activity models.ProcessingActivity) *api.MediaProcessingActivity {
	createdAt := activity.CreatedAt.Time()
	updatedAt := activity.UpdatedAt.Time()

	// Handle invalid dates
	if createdAt.Year() < 1900 || createdAt.Year() > 9999 {
		now := time.Now()
		createdAt = now
	}
	if updatedAt.Year() < 1900 || updatedAt.Year() > 9999 {
		now := time.Now()
		updatedAt = now
	}

	return &api.MediaProcessingActivity{
		Status:    stringPtr(activity.Status),
		Message:   stringPtr(activity.Message),
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}
}

func ProcessingActivityFromAPI(activity *api.MediaProcessingActivity) models.ProcessingActivity {
	result := models.ProcessingActivity{
		Status:  stringVal(activity.Status),
		Message: stringVal(activity.Message),
	}

	if activity.CreatedAt != nil {
		result.CreatedAt = primitive.NewDateTimeFromTime(*activity.CreatedAt)
	} else {
		result.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	}

	if activity.UpdatedAt != nil {
		result.UpdatedAt = primitive.NewDateTimeFromTime(*activity.UpdatedAt)
	} else {
		result.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	}

	return result
}

// Media conversions
func ToMediaResponse(media models.Media) *api.MediaResponse {
	// Convert ProcessingActivity with validation (preserve Media.ToJSONResponse logic)
	var validActivities []models.ProcessingActivity
	for _, activity := range media.ProcessingActivity {
		activityCreatedAt := activity.CreatedAt.Time()
		activityUpdatedAt := activity.UpdatedAt.Time()

		if activityCreatedAt.Year() >= 1900 && activityCreatedAt.Year() <= 9999 &&
			activityUpdatedAt.Year() >= 1900 && activityUpdatedAt.Year() <= 9999 {
			validActivities = append(validActivities, activity)
		}
	}

	// If no valid activities, create a default one
	if len(validActivities) == 0 {
		validActivities = []models.ProcessingActivity{
			{
				Status:    models.ProcessingStatusPending,
				Message:   "Queued for processing",
				CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
				UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
			},
		}
	}

	// Convert to API format
	apiActivities := make([]api.MediaProcessingActivity, len(validActivities))
	for i, activity := range validActivities {
		apiActivities[i] = *ToMediaProcessingActivity(activity)
	}

	// Handle dimensions
	dimensions := make([]int32, 2)
	if len(media.Dimensions) >= 2 {
		dimensions[0] = int32(media.Dimensions[0])
		dimensions[1] = int32(media.Dimensions[1])
	}

	// Handle duration
	var duration *float32
	if media.Duration > 0 {
		d := float32(media.Duration)
		duration = &d
	}

	createdAt := media.CreatedAt.Time()
	updatedAt := media.UpdatedAt.Time()

	// Handle invalid dates
	if createdAt.Year() < 1900 || createdAt.Year() > 9999 {
		now := time.Now()
		createdAt = now
	}
	if updatedAt.Year() < 1900 || updatedAt.Year() > 9999 {
		now := time.Now()
		updatedAt = now
	}

	return &api.MediaResponse{
		Id:                 objectIDToString(media.ID),
		AuthorId:           objectIDToString(media.AuthorId),
		MediaType:          stringPtr(media.MediaType),
		FileName:           stringPtr(media.FileName),
		Description:        stringPtr(media.Description),
		FileUrl:            stringPtr(media.FileUrl),
		ProcessedUrl:       stringPtr(media.ProcessedUrl),
		WaveformUrl:        stringPtr(media.WaveformUrl),
		ContentType:        stringPtr(media.ContentType),
		FileSize:           intToInt32(int(media.FileSize)),
		Status:             stringPtr(media.Status),
		ProcessingParams:   ToMediaProcessingParams(media.ProcessingParams),
		ProcessingActivity: apiActivities,
		Dimensions:         dimensions,
		Duration:           duration,
		CreatedAt:          &createdAt,
		UpdatedAt:          &updatedAt,
	}
}

// Board conversions
func ToBoardResponse(board models.Board) *api.BoardResponse {
	// Convert collaborators
	collaborators := make([]api.CollaboratorResponse, len(board.Collaborators))
	for i, collab := range board.Collaborators {
		collaborators[i] = *ToCollaboratorResponse(collab)
	}

	// Convert keys
	keys := make([]api.KeyResponse, len(board.Keys))
	for i, key := range board.Keys {
		keys[i] = *ToKeyResponse(key)
	}

	return &api.BoardResponse{
		Id:            objectIDToString(board.ID),
		AuthorId:      objectIDToString(board.AuthorId),
		Name:          stringPtr(board.Name),
		Description:   stringPtr(board.Description),
		Layout:        stringPtr(string(board.Layout)),
		ImageUrl:      stringPtr(board.ImageUrl),
		Collaborators: collaborators,
		Keys:          keys,
		CreatedAt:     dateTimeToTime(primitive.NewDateTimeFromTime(board.CreatedAt)),
		UpdatedAt:     dateTimeToTime(primitive.NewDateTimeFromTime(board.UpdatedAt)),
	}
}

// Request to Model conversions
func BoardFromCreateRequest(req *api.CreateBoardRequest, authorId primitive.ObjectID) models.Board {
	layout := models.GridLayout
	if req.Layout != nil {
		layoutStr := *req.Layout
		if layoutStr == "list" {
			layout = models.ListLayout
		}
	}

	return models.Board{
		ID:            primitive.NewObjectID(),
		AuthorId:      authorId,
		Name:          stringVal(req.Name),
		Description:   stringVal(req.Description),
		Layout:        layout,
		ImageUrl:      stringVal(req.ImageUrl),
		Collaborators: []models.Collaborator{}, // Initialize as empty array, not nil
		Keys:          []models.Key{},          // Initialize as empty array, not nil
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func BoardFromUpdateRequest(board models.Board, req *api.UpdateBoardRequest) models.Board {
	if req.Name != nil {
		board.Name = *req.Name
	}
	if req.Description != nil {
		board.Description = *req.Description
	}
	if req.Layout != nil {
		layoutStr := *req.Layout
		if layoutStr == "list" {
			board.Layout = models.ListLayout
		} else {
			board.Layout = models.GridLayout
		}
	}
	if req.ImageUrl != nil {
		board.ImageUrl = *req.ImageUrl
	}
	board.UpdatedAt = time.Now()
	return board
}

// Media request conversions
func MediaFromCreateRequest(req *api.CreateMediaRequest, authorId primitive.ObjectID) models.Media {
	status := models.ProcessingStatusPending
	if req.Status != nil && *req.Status != "" {
		status = *req.Status
	}

	var dimensions [2]int
	if len(req.Dimensions) == 2 {
		dimensions = [2]int{int(req.Dimensions[0]), int(req.Dimensions[1])}
	}

	var fileSize int64
	if req.FileSize != nil {
		fileSize = int64(*req.FileSize)
	}

	var duration float64
	if req.Duration != nil {
		duration = float64(*req.Duration)
	}

	return models.Media{
		ID:          primitive.NewObjectID(),
		AuthorId:    authorId,
		MediaType:   stringVal(req.MediaType),
		FileName:    stringVal(req.FileName),
		Description: stringVal(req.Description),
		FileUrl:     stringVal(req.FileUrl),
		ContentType: stringVal(req.ContentType),
		FileSize:    fileSize,
		Dimensions:  dimensions,
		Duration:    duration,
		Status:      status,
		ProcessingActivity: []models.ProcessingActivity{
			{
				Status:    models.ProcessingStatusPending,
				Message:   "Queued for processing",
				CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
				UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
			},
		},
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}
}

// Key request conversions
func KeyFromCreateRequest(req *api.CreateKeyRequest, boardId primitive.ObjectID) models.Key {
	var audioMediaId primitive.ObjectID
	if req.AudioMediaId != "" {
		audioMediaId, _ = primitive.ObjectIDFromHex(req.AudioMediaId)
	}

	var imageMediaId primitive.ObjectID
	if req.ImageMediaId != nil && *req.ImageMediaId != "" {
		imageMediaId, _ = primitive.ObjectIDFromHex(*req.ImageMediaId)
	}

	description := ""
	if req.Description != nil {
		description = *req.Description
	}

	return models.Key{
		ID:           primitive.NewObjectID(),
		BoardId:      boardId,
		Name:         req.Name,
		Description:  description,
		AudioMediaId: audioMediaId,
		ImageMediaId: imageMediaId,
		HotKey:       req.HotKey,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// Collaborator request conversions
func CollaboratorFromCreateRequest(req *api.CreateCollaboratorRequest) models.Collaborator {
	return models.Collaborator{
		ID:        primitive.NewObjectID(),
		Email:     stringVal(req.Email),
		Role:      stringVal(req.Role),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
