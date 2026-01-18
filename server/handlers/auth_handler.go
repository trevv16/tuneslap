package handlers

import (
	"time"
	api "tuneslap/generated"
	"tuneslap/models"
	"tuneslap/repositories"
	"tuneslap/services"
	"tuneslap/validation"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// validateSignupRequest validates a SignupRequest from generated types
func validateSignupRequest(body *api.SignupRequest) validation.ValidationResult {
	var errors []validation.ValidationError

	if body.Name == nil || *body.Name == "" {
		errors = append(errors, validation.ValidationError{Field: "name", Message: "Name is required"})
	} else if len(*body.Name) < 3 || len(*body.Name) > 100 {
		errors = append(errors, validation.ValidationError{Field: "name", Message: "Name must be between 3 and 100 characters"})
	}

	if body.Email == nil || *body.Email == "" {
		errors = append(errors, validation.ValidationError{Field: "email", Message: "Email is required"})
	} else {
		// Use validator for email format
		validator := validator.New()
		if err := validator.Var(*body.Email, "email"); err != nil {
			errors = append(errors, validation.ValidationError{Field: "email", Message: "Email must be a valid email address"})
		}
	}

	if body.Password == nil || *body.Password == "" {
		errors = append(errors, validation.ValidationError{Field: "password", Message: "Password is required"})
	} else if len(*body.Password) < 8 || len(*body.Password) > 128 {
		errors = append(errors, validation.ValidationError{Field: "password", Message: "Password must be between 8 and 128 characters"})
	}

	return validation.ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// validateSigninRequest validates a SigninRequest from generated types
func validateSigninRequest(body *api.SigninRequest) validation.ValidationResult {
	var errors []validation.ValidationError

	if body.Email == nil || *body.Email == "" {
		errors = append(errors, validation.ValidationError{Field: "email", Message: "Email is required"})
	} else {
		validator := validator.New()
		if err := validator.Var(*body.Email, "email"); err != nil {
			errors = append(errors, validation.ValidationError{Field: "email", Message: "Email must be a valid email address"})
		}
	}

	if body.Password == nil || *body.Password == "" {
		errors = append(errors, validation.ValidationError{Field: "password", Message: "Password is required"})
	}

	return validation.ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// validateForgotRequest validates a ForgotRequest from generated types
func validateForgotRequest(body *api.ForgotRequest) validation.ValidationResult {
	var errors []validation.ValidationError

	if body.Email == nil || *body.Email == "" {
		errors = append(errors, validation.ValidationError{Field: "email", Message: "Email is required"})
	} else {
		validator := validator.New()
		if err := validator.Var(*body.Email, "email"); err != nil {
			errors = append(errors, validation.ValidationError{Field: "email", Message: "Email must be a valid email address"})
		}
	}

	return validation.ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// validateResetRequest validates a ResetRequest from generated types
func validateResetRequest(body *api.ResetRequest) validation.ValidationResult {
	var errors []validation.ValidationError

	if body.Token == nil || *body.Token == "" {
		errors = append(errors, validation.ValidationError{Field: "token", Message: "Token is required"})
	} else if len(*body.Token) < 32 || len(*body.Token) > 255 {
		errors = append(errors, validation.ValidationError{Field: "token", Message: "Token must be between 32 and 255 characters"})
	}

	if body.Password == nil || *body.Password == "" {
		errors = append(errors, validation.ValidationError{Field: "password", Message: "Password is required"})
	} else if len(*body.Password) < 8 || len(*body.Password) > 128 {
		errors = append(errors, validation.ValidationError{Field: "password", Message: "Password must be between 8 and 128 characters"})
	}

	return validation.ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// AuthHandler provides authentication-related operations
type AuthHandler struct {
	userRepo  *repositories.UserRepository
	validator *validation.UserValidator
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		userRepo:  repositories.NewUserRepository(),
		validator: validation.NewUserValidator(),
	}
}

func SignUpHandler(c *fiber.Ctx) error {
	handler := NewAuthHandler()

	// Parse request body using generated type
	var body api.SignupRequest
	if err := c.BodyParser(&body); err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate the request
	validationResult := validateSignupRequest(&body)
	if !validationResult.IsValid {
		return SendValidationErrorResponse(c, validationResult)
	}

	// Hash password
	hash, err := services.HashPassword(*body.Password)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Error hashing password", err)
	}

	// Create user object
	user := models.User{
		Name:         *body.Name,
		Email:        *body.Email,
		PasswordHash: hash,
	}

	// Create user using repository
	_, err = handler.userRepo.CreateUser(user)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Error creating user", err)
	}

	// Create response using generated type
	response := api.NewSignupResponse()
	success := true
	message := "User created successfully"
	response.Success = &success
	response.Message = &message

	return SendSuccessResponse(c, fiber.StatusCreated, "User created successfully", response)
}

func SignInHandler(c *fiber.Ctx) error {
	handler := NewAuthHandler()

	// Parse request body using generated type
	var body api.SigninRequest
	if err := c.BodyParser(&body); err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate the request
	validationResult := validateSigninRequest(&body)
	if !validationResult.IsValid {
		return SendValidationErrorResponse(c, validationResult)
	}

	// Fetch user by email using repository
	user, err := handler.userRepo.GetByEmail(*body.Email)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or password", err)
	}

	// Check password
	if err := services.CheckPasswordHash(*body.Password, user.PasswordHash); err != nil {
		return SendErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or password", err)
	}

	// Generate JWT token
	token, err := services.GenerateJWT(user.ID.Hex())
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Error generating token", err)
	}

	// Create response using generated type
	response := api.NewSigninResponse()
	success := true
	message := "User logged in successfully"
	response.Success = &success
	response.Message = &message

	// Create data with token
	signinData := api.NewSigninData()
	signinData.Token = &token
	response.Data = signinData

	return SendSuccessResponse(c, fiber.StatusOK, "User logged in successfully", response)
}

func ForgotPasswordHandler(c *fiber.Ctx) error {
	handler := NewAuthHandler()

	// Parse request body using generated type
	var body api.ForgotRequest
	if err := c.BodyParser(&body); err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate the request
	validationResult := validateForgotRequest(&body)
	if !validationResult.IsValid {
		return SendValidationErrorResponse(c, validationResult)
	}

	// Fetch user by email using repository
	user, err := handler.userRepo.GetByEmail(*body.Email)
	if err != nil {
		// Don't reveal if user exists or not for security
		response := api.NewForgotResponse()
		success := true
		message := "If the email exists, a reset link has been sent"
		response.Success = &success
		response.Message = &message
		return SendSuccessResponse(c, fiber.StatusOK, "If the email exists, a reset link has been sent", response)
	}

	// Generate reset token
	token := uuid.New().String()
	expiry := time.Now().Add(30 * time.Minute)

	// Update user with reset token using repository
	_, err = handler.userRepo.UpdateUser(user.ID, &models.UpdateUserRequest{
		ResetToken:     &token,
		ResetExpiresAt: &expiry,
	})
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Error updating user", err)
	}

	// TODO: Send email with reset link/token here

	response := api.NewForgotResponse()
	success := true
	message := "If the email exists, a reset link has been sent"
	response.Success = &success
	response.Message = &message
	return SendSuccessResponse(c, fiber.StatusOK, "If the email exists, a reset link has been sent", response)
}

func ResetPasswordHandler(c *fiber.Ctx) error {
	handler := NewAuthHandler()

	// Parse request body using generated type
	var body api.ResetRequest
	if err := c.BodyParser(&body); err != nil {
		return SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate the request
	validationResult := validateResetRequest(&body)
	if !validationResult.IsValid {
		return SendValidationErrorResponse(c, validationResult)
	}

	// Fetch user by reset token using repository
	user, err := handler.userRepo.GetByResetToken(*body.Token)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token", err)
	}

	// Check if token is expired
	if user.ResetExpiresAt == nil || time.Now().After(*user.ResetExpiresAt) {
		return SendErrorResponse(c, fiber.StatusUnauthorized, "Token has expired", nil)
	}

	// Hash new password
	hash, err := services.HashPassword(*body.Password)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Error hashing password", err)
	}

	// Update user with new password and clear reset token using repository
	_, err = handler.userRepo.UpdateUser(user.ID, &models.UpdateUserRequest{
		ResetToken:     nil,
		ResetExpiresAt: nil,
		PasswordHash:   hash,
	})
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Error updating user", err)
	}

	// TODO: Send email notification that password has been reset

	response := api.NewResetResponse()
	success := true
	message := "Password reset successful"
	response.Success = &success
	response.Message = &message
	return SendSuccessResponse(c, fiber.StatusOK, "Password reset successful", response)
}
