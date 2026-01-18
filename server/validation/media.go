package validation

import (
	api "tuneslap/generated"
	"tuneslap/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MediaValidator provides validation functionality for media-related operations
type MediaValidator struct {
	*Validator
}

// NewMediaValidator creates a new media validator instance
func NewMediaValidator() *MediaValidator {
	return &MediaValidator{Validator: NewValidator()}
}

// Validate validates data for BaseHandler compatibility
func (v *MediaValidator) Validate(data interface{}) ValidationResult {
	switch req := data.(type) {
	case *api.CreateMediaRequest:
		return ValidateCreateMediaRequest(req)
	case *api.UpdateMediaRequest:
		return ValidateUpdateMediaRequest(req)
	default:
		return ValidationResult{
			IsValid: false,
			Errors: []ValidationError{
				{Field: "request", Message: "Invalid request type"},
			},
		}
	}
}

// ValidateCreateMedia validates media creation data
func (v *MediaValidator) ValidateCreateMedia(data *api.CreateMediaRequest) ValidationResult {
	return ValidateCreateMediaRequest(data)
}

// ValidateUpdateMedia validates media update data
func (v *MediaValidator) ValidateUpdateMedia(data *api.UpdateMediaRequest) ValidationResult {
	return ValidateUpdateMediaRequest(data)
}

// ValidateMediaOwnership validates if user owns the media
func (v *MediaValidator) ValidateMediaOwnership(media models.Media, authorId primitive.ObjectID) ValidationResult {
	if media.AuthorId != authorId {
		return ValidationResult{
			IsValid: false,
			Errors: []ValidationError{
				{Field: "authorId", Message: "You don't have permission to access this media"},
			},
		}
	}
	return ValidationResult{IsValid: true}
}

// ValidateFileSize validates if file size is within acceptable limits
func (v *MediaValidator) ValidateFileSize(fileSize int64, maxSize int64) ValidationResult {
	if fileSize > maxSize {
		return ValidationResult{
			IsValid: false,
			Errors: []ValidationError{
				{Field: "fileSize", Message: "File size exceeds maximum allowed size"},
			},
		}
	}
	return ValidationResult{IsValid: true}
}

// ValidateCreateMediaRequest validates a CreateMediaRequest
func ValidateCreateMediaRequest(req *api.CreateMediaRequest) ValidationResult {
	var errors []ValidationError

	// MediaType: required, oneof=image audio
	if err := validateRequiredString("mediaType", req.MediaType); err != nil {
		errors = append(errors, *err)
	} else if req.MediaType != nil {
		if err := validateOneOf("mediaType", *req.MediaType, []string{"image", "audio"}); err != nil {
			errors = append(errors, *err)
		}
	}

	// FileName: required, min=MinFileNameLength, max=MaxFileNameLength, allows alphanumeric, spaces, hyphens, underscores, and dots
	if err := validateRequiredString("fileName", req.FileName); err != nil {
		errors = append(errors, *err)
	} else if req.FileName != nil {
		if err := validateStringLength("fileName", *req.FileName, MinFileNameLength, MaxFileNameLength); err != nil {
			errors = append(errors, *err)
		} else {
			// Allow alphanumeric, spaces, hyphens, underscores, and dots for filenames
			if err := validateFileNameHelper("fileName", *req.FileName); err != nil {
				errors = append(errors, *err)
			}
			if err := validateExcludesAllHelper("fileName", *req.FileName, "<>{}[]()&|\\"); err != nil {
				errors = append(errors, *err)
			}
		}
	}

	// Description: omitempty, max=MaxDescriptionLength
	if req.Description != nil && *req.Description != "" {
		if err := validateStringLength("description", *req.Description, 0, MaxDescriptionLength); err != nil {
			errors = append(errors, *err)
		}
	}

	// FileUrl: required, url, excludesall=<>{}[]()&|\
	if err := validateRequiredString("fileUrl", req.FileUrl); err != nil {
		errors = append(errors, *err)
	} else if req.FileUrl != nil {
		if err := validateURLHelper("fileUrl", *req.FileUrl); err != nil {
			errors = append(errors, *err)
		}
		if err := validateExcludesAllHelper("fileUrl", *req.FileUrl, "<>{}[]()&|\\"); err != nil {
			errors = append(errors, *err)
		}
	}

	// ContentType: optional, but if provided, min=MinContentTypeLength, max=MaxContentTypeLength, allows MIME type format (with slashes)
	// If not provided, it will be derived from file extension in the handler
	if req.ContentType != nil && *req.ContentType != "" {
		if err := validateStringLength("contentType", *req.ContentType, MinContentTypeLength, MaxContentTypeLength); err != nil {
			errors = append(errors, *err)
		} else {
			// Allow MIME type format (e.g., "audio/x-wav", "image/jpeg")
			if err := validateContentTypeHelper("contentType", *req.ContentType); err != nil {
				errors = append(errors, *err)
			}
			if err := validateExcludesAllHelper("contentType", *req.ContentType, "<>{}[]()&|\\"); err != nil {
				errors = append(errors, *err)
			}
		}
	}

	// FileSize: required, min=1, max=MaxFileSize
	if err := validateIntRange("fileSize", req.FileSize, 1, MaxFileSize); err != nil {
		errors = append(errors, *err)
	}

	// Status: omitempty, oneof=pending processing done error
	if req.Status != nil && *req.Status != "" {
		if err := validateOneOf("status", *req.Status, []string{"pending", "processing", "done", "error"}); err != nil {
			errors = append(errors, *err)
		}
	}

	// Dimensions: optional, must be [width, height] with exactly 2 elements
	if err := validateDimensions("dimensions", req.Dimensions); err != nil {
		errors = append(errors, *err)
	}

	// Duration: optional, must be non-negative
	if err := validateNonNegativeFloat("duration", req.Duration); err != nil {
		errors = append(errors, *err)
	}

	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// ValidateUpdateMediaRequest validates an UpdateMediaRequest
func ValidateUpdateMediaRequest(req *api.UpdateMediaRequest) ValidationResult {
	var errors []ValidationError

	// Description: omitempty, max=MaxDescriptionLength
	if req.Description != nil && *req.Description != "" {
		if err := validateStringLength("description", *req.Description, 0, MaxDescriptionLength); err != nil {
			errors = append(errors, *err)
		}
	}

	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// ValidateMediaProcessingParamsAudio validates audio processing parameters
func ValidateMediaProcessingParamsAudio(params *api.MediaProcessingParamsAudio) ValidationResult {
	if params == nil {
		return ValidationResult{IsValid: true}
	}

	var errors []ValidationError

	// TrimStart: optional, must be >= 0
	if err := validateNonNegativeFloat("trimStart", params.TrimStart); err != nil {
		errors = append(errors, *err)
	}

	// TrimEnd: optional, must be >= 0
	if err := validateNonNegativeFloat("trimEnd", params.TrimEnd); err != nil {
		errors = append(errors, *err)
	}

	// FadeIn: optional, must be >= 0
	if err := validateNonNegativeFloat("fadeIn", params.FadeIn); err != nil {
		errors = append(errors, *err)
	}

	// FadeOut: optional, must be >= 0
	if err := validateNonNegativeFloat("fadeOut", params.FadeOut); err != nil {
		errors = append(errors, *err)
	}

	// Speed: optional, must be > 0
	if params.Speed != nil && *params.Speed <= 0 {
		errors = append(errors, ValidationError{
			Field:   "speed",
			Message: "speed must be greater than 0",
		})
	}

	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}
