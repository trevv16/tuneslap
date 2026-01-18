package validation

import (
	api "tuneslap/generated"
	"tuneslap/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CollaboratorValidator provides validation functionality for collaborator-related operations
type CollaboratorValidator struct {
	*Validator
}

// NewCollaboratorValidator creates a new collaborator validator instance
func NewCollaboratorValidator() *CollaboratorValidator {
	return &CollaboratorValidator{Validator: NewValidator()}
}

// Validate validates data for ArrayHandler compatibility
func (v *CollaboratorValidator) Validate(data interface{}) ValidationResult {
	switch req := data.(type) {
	case *api.CreateCollaboratorRequest:
		return ValidateCreateCollaboratorRequest(req)
	case *api.UpdateCollaboratorRequest:
		return ValidateUpdateCollaboratorRequest(req)
	default:
		return ValidationResult{
			IsValid: false,
			Errors: []ValidationError{
				{Field: "request", Message: "Invalid request type"},
			},
		}
	}
}

// ValidateCreateCollaborator validates collaborator creation data
func (v *CollaboratorValidator) ValidateCreateCollaborator(data *api.CreateCollaboratorRequest) ValidationResult {
	return ValidateCreateCollaboratorRequest(data)
}

// ValidateUpdateCollaborator validates collaborator update data
func (v *CollaboratorValidator) ValidateUpdateCollaborator(data *api.UpdateCollaboratorRequest) ValidationResult {
	return ValidateUpdateCollaboratorRequest(data)
}

// ValidateCollaboratorUniqueness validates that user is not already a collaborator
func (v *CollaboratorValidator) ValidateCollaboratorUniqueness(userId primitive.ObjectID, collaborators []models.Collaborator) ValidationResult {
	for _, collaborator := range collaborators {
		if collaborator.UserId == userId {
			return ValidationResult{
				IsValid: false,
				Errors: []ValidationError{
					{Field: "userId", Message: "User is already a collaborator on this board"},
				},
			}
		}
	}
	return ValidationResult{IsValid: true}
}

// ValidateCreateCollaboratorRequest validates a CreateCollaboratorRequest
func ValidateCreateCollaboratorRequest(req *api.CreateCollaboratorRequest) ValidationResult {
	var errors []ValidationError

	// Email: required, email
	if err := validateRequiredString("email", req.Email); err != nil {
		errors = append(errors, *err)
	} else if req.Email != nil {
		if err := validateEmailHelper("email", *req.Email); err != nil {
			errors = append(errors, *err)
		}
	}

	// Role: required, oneof=owner editor viewer
	if err := validateRequiredString("role", req.Role); err != nil {
		errors = append(errors, *err)
	} else if req.Role != nil {
		if err := validateOneOf("role", *req.Role, []string{"owner", "editor", "viewer"}); err != nil {
			errors = append(errors, *err)
		}
	}

	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// ValidateUpdateCollaboratorRequest validates an UpdateCollaboratorRequest
func ValidateUpdateCollaboratorRequest(req *api.UpdateCollaboratorRequest) ValidationResult {
	var errors []ValidationError

	// Role: required, oneof=owner editor viewer
	if err := validateRequiredString("role", req.Role); err != nil {
		errors = append(errors, *err)
	} else if req.Role != nil {
		if err := validateOneOf("role", *req.Role, []string{"owner", "editor", "viewer"}); err != nil {
			errors = append(errors, *err)
		}
	}

	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}
