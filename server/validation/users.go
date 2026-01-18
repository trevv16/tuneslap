package validation

import (
	"strings"
	api "tuneslap/generated"
	"tuneslap/models"
)

// UserValidator provides validation functionality for user-related operations
type UserValidator struct {
	*Validator
}

// NewUserValidator creates a new user validator instance
func NewUserValidator() *UserValidator {
	return &UserValidator{Validator: NewValidator()}
}

// ValidateCreateUser validates user creation data
func (v *UserValidator) ValidateCreateUser(user models.User) ValidationResult {
	return v.Validate(user)
}

// Validate validates data for BaseHandler compatibility
func (v *UserValidator) Validate(data interface{}) ValidationResult {
	switch req := data.(type) {
	case *api.UpdateMeRequest:
		return ValidateUpdateMeRequest(req)
	default:
		return ValidationResult{
			IsValid: false,
			Errors: []ValidationError{
				{Field: "request", Message: "Invalid request type"},
			},
		}
	}
}

// ValidateUpdateMe validates user profile update data
func (v *UserValidator) ValidateUpdateMe(data *api.UpdateMeRequest) ValidationResult {
	return ValidateUpdateMeRequest(data)
}

// ValidateUpdateUser validates user update data
func (v *UserValidator) ValidateUpdateUser(data *models.UpdateUserRequest) ValidationResult {
	return v.Validate(data)
}

// ValidateEmail validates email format and uniqueness
func (v *UserValidator) ValidateEmail(email string, existingEmails []string) ValidationResult {
	// Check email format
	if err := v.validate.Var(email, "required,email"); err != nil {
		return ValidationResult{
			IsValid: false,
			Errors: []ValidationError{
				{Field: "email", Message: "Email must be a valid email address"},
			},
		}
	}

	// Check for duplicates
	for _, existing := range existingEmails {
		if strings.EqualFold(email, existing) {
			return ValidationResult{
				IsValid: false,
				Errors: []ValidationError{
					{Field: "email", Message: "Email already exists"},
				},
			}
		}
	}

	return ValidationResult{IsValid: true}
}

// ValidateUpdateMeRequest validates an UpdateMeRequest
func ValidateUpdateMeRequest(req *api.UpdateMeRequest) ValidationResult {
	var errors []ValidationError

	// Name: omitempty, min=MinNameLength, max=MaxNameLength, alphanumspace, excludesall=<>{}[]()&|\
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

	// ImageUrl: omitempty, url, excludesall=<>{}[]()&|\
	if err := validateOptionalString("imageUrl", req.ImageUrl, 0, 0, "url", "excludesall"); err != nil {
		errors = append(errors, *err)
	}

	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}
