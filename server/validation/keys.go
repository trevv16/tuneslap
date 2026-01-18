package validation

import (
	api "tuneslap/generated"
	"tuneslap/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// KeyValidator provides validation functionality for key-related operations
type KeyValidator struct {
	*Validator
}

// NewKeyValidator creates a new key validator instance
func NewKeyValidator() *KeyValidator {
	return &KeyValidator{Validator: NewValidator()}
}

// Validate validates data for ArrayHandler compatibility
func (v *KeyValidator) Validate(data interface{}) ValidationResult {
	switch req := data.(type) {
	case *api.CreateKeyRequest:
		return ValidateCreateKeyRequest(req)
	case *api.UpdateKeyRequest:
		return ValidateUpdateKeyRequest(req)
	default:
		return ValidationResult{
			IsValid: false,
			Errors: []ValidationError{
				{Field: "request", Message: "Invalid request type"},
			},
		}
	}
}

// ValidateCreateKey validates key creation data
func (v *KeyValidator) ValidateCreateKey(data *api.CreateKeyRequest) ValidationResult {
	return ValidateCreateKeyRequest(data)
}

// ValidateUpdateKey validates key update data
func (v *KeyValidator) ValidateUpdateKey(data *api.UpdateKeyRequest) ValidationResult {
	return ValidateUpdateKeyRequest(data)
}

// ValidateHotKeyUniqueness validates that hot key is unique within a board
func (v *KeyValidator) ValidateHotKeyUniqueness(hotKey string, boardKeys []models.Key, excludeId primitive.ObjectID) ValidationResult {
	for _, key := range boardKeys {
		if key.ID != excludeId && key.HotKey == hotKey {
			return ValidationResult{
				IsValid: false,
				Errors: []ValidationError{
					{Field: "hotKey", Message: "Hot key must be unique within the board"},
				},
			}
		}
	}
	return ValidationResult{IsValid: true}
}

// ValidateCreateKeyRequest validates a CreateKeyRequest
func ValidateCreateKeyRequest(req *api.CreateKeyRequest) ValidationResult {
	var errors []ValidationError

	// BoardId: required (non-pointer string, check not empty)
	if req.BoardId == "" {
		errors = append(errors, ValidationError{
			Field:   "boardId",
			Message: "boardId is required",
		})
	}

	// Name: required, min=MinNameLength, max=MaxNameLength, alphanumspace, excludesall=<>{}[]()&|\
	if req.Name == "" {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "name is required",
		})
	} else {
		if err := validateStringLength("name", req.Name, MinNameLength, MaxNameLength); err != nil {
			errors = append(errors, *err)
		} else {
			if err := validateAlphaNumSpaceHelper("name", req.Name); err != nil {
				errors = append(errors, *err)
			}
			if err := validateExcludesAllHelper("name", req.Name, "<>{}[]()&|\\"); err != nil {
				errors = append(errors, *err)
			}
		}
	}

	// Description: omitempty, max=MaxShortDescriptionLength, excludesall=<>{}[]()&|\
	if err := validateOptionalString("description", req.Description, 0, MaxShortDescriptionLength, "excludesall"); err != nil {
		errors = append(errors, *err)
	}

	// AudioMediaId: required (non-pointer string, check not empty)
	if req.AudioMediaId == "" {
		errors = append(errors, ValidationError{
			Field:   "audioMediaId",
			Message: "audioMediaId is required",
		})
	}

	// ImageMediaId: omitempty (pointer string)
	// No validation needed for optional field

	// HotKey: exactly HotKeyLength character
	if req.HotKey == "" {
		errors = append(errors, ValidationError{
			Field:   "hotKey",
			Message: "hotKey is required",
		})
	} else if len(req.HotKey) != HotKeyLength {
		errors = append(errors, ValidationError{
			Field:   "hotKey",
			Message: "hotKey must be exactly 1 character",
		})
	}

	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// ValidateUpdateKeyRequest validates an UpdateKeyRequest
func ValidateUpdateKeyRequest(req *api.UpdateKeyRequest) ValidationResult {
	var errors []ValidationError

	// Name: min=MinNameLength, max=MaxNameLength, alphanumspace, excludesall=<>{}[]()&|\
	if req.Name != nil && *req.Name != "" {
		if err := validateStringLength("name", *req.Name, MinNameLength, MaxNameLength); err != nil {
			errors = append(errors, *err)
		} else {
			if err := validateAlphaNumSpaceHelper("name", *req.Name); err != nil {
				errors = append(errors, *err)
			}
			if err := validateExcludesAllHelper("name", *req.Name, "<>{}[]()&|\\"); err != nil {
				errors = append(errors, *err)
			}
		}
	}

	// Description: omitempty, max=MaxShortDescriptionLength, excludesall=<>{}[]()&|\
	if err := validateOptionalString("description", req.Description, 0, MaxShortDescriptionLength, "excludesall"); err != nil {
		errors = append(errors, *err)
	}

	// HotKey: exactly HotKeyLength character
	if req.HotKey != nil && *req.HotKey != "" {
		if len(*req.HotKey) != HotKeyLength {
			errors = append(errors, ValidationError{
				Field:   "hotKey",
				Message: "hotKey must be exactly 1 character",
			})
		}
	}

	// AudioMediaId: omitempty
	// ImageMediaId: omitempty
	// No validation needed for optional fields

	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}
