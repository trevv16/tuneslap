package validation

import (
	api "tuneslap/generated"
	"tuneslap/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BoardValidator provides validation functionality for board-related operations
type BoardValidator struct {
	*Validator
}

// NewBoardValidator creates a new board validator instance
func NewBoardValidator() *BoardValidator {
	return &BoardValidator{Validator: NewValidator()}
}

// Validate validates data for BaseHandler compatibility
func (v *BoardValidator) Validate(data interface{}) ValidationResult {
	switch req := data.(type) {
	case *api.CreateBoardRequest:
		return ValidateCreateBoardRequest(req)
	case *api.UpdateBoardRequest:
		return ValidateUpdateBoardRequest(req)
	default:
		return ValidationResult{
			IsValid: false,
			Errors: []ValidationError{
				{Field: "request", Message: "Invalid request type"},
			},
		}
	}
}

// ValidateCreateBoard validates board creation data
func (v *BoardValidator) ValidateCreateBoard(data *api.CreateBoardRequest) ValidationResult {
	return ValidateCreateBoardRequest(data)
}

// ValidateUpdateBoard validates board update data
func (v *BoardValidator) ValidateUpdateBoard(data *api.UpdateBoardRequest) ValidationResult {
	return ValidateUpdateBoardRequest(data)
}

// ValidateBoardOwnership validates if user owns the board
func (v *BoardValidator) ValidateBoardOwnership(board models.Board, authorId primitive.ObjectID) ValidationResult {
	if board.AuthorId != authorId {
		return ValidationResult{
			IsValid: false,
			Errors: []ValidationError{
				{Field: "authorId", Message: "You don't have permission to access this board"},
			},
		}
	}
	return ValidationResult{IsValid: true}
}

// ValidateCreateBoardRequest validates a CreateBoardRequest
func ValidateCreateBoardRequest(req *api.CreateBoardRequest) ValidationResult {
	var errors []ValidationError

	// Name: required, min=MinNameLength, max=MaxNameLength, alphanumspace, excludesall=<>{}[]()&|\
	if req.Name == "" {
		errors = append(errors, ValidationError{Field: "name", Message: "name is required"})
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

	// Layout: required, oneof=grid list
	if req.Layout == "" {
		errors = append(errors, ValidationError{Field: "layout", Message: "layout is required"})
	} else {
		if err := validateOneOf("layout", req.Layout, []string{"grid", "list"}); err != nil {
			errors = append(errors, *err)
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

// ValidateUpdateBoardRequest validates an UpdateBoardRequest
func ValidateUpdateBoardRequest(req *api.UpdateBoardRequest) ValidationResult {
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

	// Layout: oneof=grid list
	if req.Layout != nil && *req.Layout != "" {
		if err := validateOneOf("layout", *req.Layout, []string{"grid", "list"}); err != nil {
			errors = append(errors, *err)
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
