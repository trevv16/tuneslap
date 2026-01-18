package validation

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Helper functions for pointer-based validation

// validateRequiredString checks if a string pointer is not nil and not empty
func validateRequiredString(field string, value *string) *ValidationError {
	if value == nil || *value == "" {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s is required", field),
		}
	}
	return nil
}

// validateOptionalString validates an optional string field
func validateOptionalString(field string, value *string, min, max int, rules ...string) *ValidationError {
	if value == nil || *value == "" {
		return nil // Optional field, skip validation if nil or empty
	}

	val := *value
	var errors []*ValidationError

	// Check length constraints
	if min > 0 && len(val) < min {
		errors = append(errors, &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must be at least %d characters long", field, min),
		})
	}
	if max > 0 && len(val) > max {
		errors = append(errors, &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must be at most %d characters long", field, max),
		})
	}

	// Check validation rules
	for _, rule := range rules {
		switch rule {
		case "alphanumspace":
			if err := validateAlphaNumSpaceHelper(field, val); err != nil {
				errors = append(errors, err)
			}
		case "excludesall":
			if err := validateExcludesAllHelper(field, val, "<>{}[]()&|\\"); err != nil {
				errors = append(errors, err)
			}
		case "url":
			if err := validateURLHelper(field, val); err != nil {
				errors = append(errors, err)
			}
		case "email":
			if err := validateEmailHelper(field, val); err != nil {
				errors = append(errors, err)
			}
		}
	}

	if len(errors) > 0 {
		return errors[0] // Return first error
	}
	return nil
}

// validateStringLength validates string length
func validateStringLength(field, value string, min, max int) *ValidationError {
	if min > 0 && len(value) < min {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must be at least %d characters long", field, min),
		}
	}
	if max > 0 && len(value) > max {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must be at most %d characters long", field, max),
		}
	}
	return nil
}

// validateAlphaNumSpaceHelper validates alphanumeric and space characters
func validateAlphaNumSpaceHelper(field, value string) *ValidationError {
	for _, char := range value {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == ' ') {
			return &ValidationError{
				Field:   field,
				Message: fmt.Sprintf("%s can only contain letters, numbers, and spaces", field),
			}
		}
	}
	return nil
}

// validateExcludesAllHelper validates that string doesn't contain excluded characters
func validateExcludesAllHelper(field, value, excluded string) *ValidationError {
	for _, char := range excluded {
		if strings.ContainsRune(value, char) {
			return &ValidationError{
				Field:   field,
				Message: fmt.Sprintf("%s cannot contain special characters: %s", field, excluded),
			}
		}
	}
	return nil
}

// validateURLHelper validates URL format
func validateURLHelper(field, value string) *ValidationError {
	_, err := url.ParseRequestURI(value)
	if err != nil {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must be a valid URL", field),
		}
	}
	return nil
}

// validateEmailHelper validates email format
func validateEmailHelper(field, value string) *ValidationError {
	validator := validator.New()
	if err := validator.Var(value, "email"); err != nil {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must be a valid email address", field),
		}
	}
	return nil
}

// validateOneOf validates that value is one of the allowed options
func validateOneOf(field, value string, options []string) *ValidationError {
	for _, option := range options {
		if value == option {
			return nil
		}
	}
	return &ValidationError{
		Field:   field,
		Message: fmt.Sprintf("%s must be one of: %s", field, strings.Join(options, " ")),
	}
}

// validateIntRange validates integer range
func validateIntRange(field string, value *int32, min, max int64) *ValidationError {
	if value == nil {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s is required", field),
		}
	}
	val := int64(*value)
	if min > 0 && val < min {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must be at least %d", field, min),
		}
	}
	if max > 0 && val > max {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must be at most %d", field, max),
		}
	}
	return nil
}

// validateNonNegativeFloat validates that an optional float is >= 0
func validateNonNegativeFloat(field string, value *float32) *ValidationError {
	if value == nil {
		return nil // Optional field
	}
	if *value < 0 {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must be non-negative", field),
		}
	}
	return nil
}

// validateDimensions validates dimensions array
func validateDimensions(field string, value []int32) *ValidationError {
	if value != nil && len(value) != 2 {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must be an array of exactly 2 elements", field),
		}
	}
	return nil
}
