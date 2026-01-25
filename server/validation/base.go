package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validation constants
const (
	// MaxFileSize is the maximum allowed file size in bytes (1GB)
	MaxFileSize int64 = 1_000_000_000

	// String length limits
	MinNameLength             = 3
	MaxNameLength             = 100
	MinFileNameLength         = 3
	MaxFileNameLength         = 100
	MaxDescriptionLength      = 1000
	MaxShortDescriptionLength = 500
	MinContentTypeLength      = 3
	MaxContentTypeLength      = 100
	HotKeyLength              = 1
)

// ValidationError represents a validation error with field and message
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationResult represents the result of validation
type ValidationResult struct {
	IsValid bool              `json:"isValid"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

// Validator provides validation functionality for repositories
type Validator struct {
	validate *validator.Validate
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	v := validator.New()

	// Register custom validations
	v.RegisterValidation("alphanumspace", validateAlphaNumSpace)
	v.RegisterValidation("excludesall", validateExcludesAll)
	v.RegisterValidation("filename", validateFileName)
	v.RegisterValidation("contenttype", validateContentType)

	return &Validator{validate: v}
}

// Validate validates a struct and returns validation result
func (v *Validator) Validate(data interface{}) ValidationResult {
	err := v.validate.Struct(data)
	if err == nil {
		return ValidationResult{IsValid: true}
	}

	var errors []ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		tag := err.Tag()
		param := err.Param()

		message := v.getErrorMessage(field, tag, param)
		errors = append(errors, ValidationError{
			Field:   field,
			Message: message,
		})
	}

	return ValidationResult{
		IsValid: false,
		Errors:  errors,
	}
}

// getErrorMessage returns a user-friendly error message
func (v *Validator) getErrorMessage(field, tag, param string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, param)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	case "alphanumspace":
		return fmt.Sprintf("%s can only contain letters, numbers, and spaces", field)
	case "excludesall":
		return fmt.Sprintf("%s cannot contain special characters: %s", field, param)
	case "filename":
		return fmt.Sprintf("%s can only contain letters, numbers, spaces, hyphens, underscores, and dots", field)
	case "contenttype":
		return fmt.Sprintf("%s must be a valid MIME type", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// Custom validation functions
func validateAlphaNumSpace(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	for _, char := range value {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == ' ') {
			return false
		}
	}
	return true
}

func validateExcludesAll(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	excluded := fl.Param()
	for _, char := range excluded {
		if strings.ContainsRune(value, char) {
			return false
		}
	}
	return true
}

// validateFileName validates filename characters (allows alphanumeric, spaces, hyphens, underscores, dots)
func validateFileName(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	for _, char := range value {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == ' ' ||
			char == '-' ||
			char == '_' ||
			char == '.') {
			return false
		}
	}
	return true
}

// validateContentType validates MIME type format (allows alphanumeric, slashes, hyphens, plus signs, dots)
func validateContentType(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	for _, char := range value {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '/' ||
			char == '-' ||
			char == '+' ||
			char == '.' ||
			char == '*') {
			return false
		}
	}
	return true
}

// validateFileNameHelper validates filename characters (allows alphanumeric, spaces, hyphens, underscores, dots)
func validateFileNameHelper(field, value string) *ValidationError {
	for _, char := range value {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == ' ' ||
			char == '-' ||
			char == '_' ||
			char == '.') {
			return &ValidationError{
				Field:   field,
				Message: fmt.Sprintf("%s can only contain letters, numbers, spaces, hyphens, underscores, and dots", field),
			}
		}
	}
	return nil
}

// validateContentTypeHelper validates MIME type format (allows alphanumeric, slashes, hyphens, plus signs)
func validateContentTypeHelper(field, value string) *ValidationError {
	for _, char := range value {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '/' ||
			char == '-' ||
			char == '+' ||
			char == '.' ||
			char == '*') {
			return &ValidationError{
				Field:   field,
				Message: fmt.Sprintf("%s must be a valid MIME type (e.g., audio/x-wav, image/jpeg)", field),
			}
		}
	}
	return nil
}

// NewValidationError creates a new validation error from validation errors
func NewValidationError(errors []ValidationError) error {
	if len(errors) == 0 {
		return nil
	}

	var messages []string
	for _, err := range errors {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}

	return fmt.Errorf("validation failed: %s", strings.Join(messages, "; "))
}
