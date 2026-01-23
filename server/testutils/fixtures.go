package testutils

import (
	"time"
	"tuneslap/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Test user fixtures
var (
	TestUserID       = primitive.NewObjectID()
	TestUser2ID      = primitive.NewObjectID()
	TestBoardID      = primitive.NewObjectID()
	TestBoard2ID     = primitive.NewObjectID()
	TestMediaID      = primitive.NewObjectID()
	TestMedia2ID     = primitive.NewObjectID()
	TestKeyID        = primitive.NewObjectID()
	TestKey2ID       = primitive.NewObjectID()
	TestCollabID     = primitive.NewObjectID()
	TestCollab2ID    = primitive.NewObjectID()
)

// CreateTestUser creates a test user with the given parameters
func CreateTestUser(id primitive.ObjectID, name, email, passwordHash string) models.User {
	now := time.Now()
	return models.User{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// CreateDefaultTestUser creates a default test user
func CreateDefaultTestUser() models.User {
	return CreateTestUser(
		TestUserID,
		"Test User",
		"test@example.com",
		"$2a$10$hashedpassword", // bcrypt hash placeholder
	)
}

// CreateSecondTestUser creates a second test user
func CreateSecondTestUser() models.User {
	return CreateTestUser(
		TestUser2ID,
		"Test User 2",
		"test2@example.com",
		"$2a$10$hashedpassword2",
	)
}

// CreateTestBoard creates a test board with the given parameters
func CreateTestBoard(id, authorID primitive.ObjectID, name, description string, layout models.LayoutType) models.Board {
	now := time.Now()
	return models.Board{
		ID:            id,
		AuthorId:      authorID,
		Name:          name,
		Description:   description,
		Layout:        layout,
		Collaborators: []models.Collaborator{},
		Keys:          []models.Key{},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// CreateDefaultTestBoard creates a default test board
func CreateDefaultTestBoard() models.Board {
	return CreateTestBoard(
		TestBoardID,
		TestUserID,
		"Test Board",
		"A test board for testing",
		models.GridLayout,
	)
}

// CreateSecondTestBoard creates a second test board
func CreateSecondTestBoard() models.Board {
	return CreateTestBoard(
		TestBoard2ID,
		TestUserID,
		"Test Board 2",
		"Another test board",
		models.ListLayout,
	)
}

// CreateTestMedia creates test media with the given parameters
func CreateTestMedia(id, authorID primitive.ObjectID, mediaType, fileName, contentType string, fileSize int64) models.Media {
	now := primitive.NewDateTimeFromTime(time.Now())
	return models.Media{
		ID:          id,
		AuthorId:    authorID,
		MediaType:   mediaType,
		FileName:    fileName,
		Description: "Test media description",
		FileUrl:     "https://storage.example.com/test/" + fileName,
		ContentType: contentType,
		FileSize:    fileSize,
		Status:      models.ProcessingStatusPending,
		ProcessingActivity: []models.ProcessingActivity{
			{
				Status:    models.ProcessingStatusPending,
				Message:   "Queued for processing",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// CreateDefaultTestAudioMedia creates default test audio media
func CreateDefaultTestAudioMedia() models.Media {
	return CreateTestMedia(
		TestMediaID,
		TestUserID,
		"audio",
		"test-audio.mp3",
		"audio/mpeg",
		1024000,
	)
}

// CreateDefaultTestImageMedia creates default test image media
func CreateDefaultTestImageMedia() models.Media {
	return CreateTestMedia(
		TestMedia2ID,
		TestUserID,
		"image",
		"test-image.png",
		"image/png",
		512000,
	)
}

// CreateTestKey creates a test key with the given parameters
func CreateTestKey(id, boardID, audioMediaID, imageMediaID primitive.ObjectID, name, description, hotKey string) models.Key {
	now := time.Now()
	return models.Key{
		ID:           id,
		BoardId:      boardID,
		Name:         name,
		Description:  description,
		AudioMediaId: audioMediaID,
		ImageMediaId: imageMediaID,
		HotKey:       hotKey,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// CreateDefaultTestKey creates a default test key
func CreateDefaultTestKey() models.Key {
	return CreateTestKey(
		TestKeyID,
		TestBoardID,
		TestMediaID,
		TestMedia2ID,
		"Test Key",
		"A test soundboard key",
		"A",
	)
}

// CreateSecondTestKey creates a second test key
func CreateSecondTestKey() models.Key {
	return CreateTestKey(
		TestKey2ID,
		TestBoardID,
		TestMediaID,
		primitive.NilObjectID,
		"Test Key 2",
		"Another test key",
		"B",
	)
}

// CreateTestCollaborator creates a test collaborator with the given parameters
func CreateTestCollaborator(id, userID primitive.ObjectID, email, role string) models.Collaborator {
	now := time.Now()
	return models.Collaborator{
		ID:        id,
		UserId:    userID,
		Email:     email,
		Role:      role,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// CreateDefaultTestCollaborator creates a default test collaborator
func CreateDefaultTestCollaborator() models.Collaborator {
	return CreateTestCollaborator(
		TestCollabID,
		TestUser2ID,
		"collab@example.com",
		"editor",
	)
}

// CreateViewerCollaborator creates a collaborator with viewer role
func CreateViewerCollaborator() models.Collaborator {
	return CreateTestCollaborator(
		TestCollab2ID,
		TestUser2ID,
		"viewer@example.com",
		"viewer",
	)
}

// CreateTestBoardWithCollaborators creates a board with collaborators
func CreateTestBoardWithCollaborators(collaborators []models.Collaborator) models.Board {
	board := CreateDefaultTestBoard()
	board.Collaborators = collaborators
	return board
}

// CreateTestBoardWithKeys creates a board with keys
func CreateTestBoardWithKeys(keys []models.Key) models.Board {
	board := CreateDefaultTestBoard()
	board.Keys = keys
	return board
}

// CreateTestBoardFull creates a board with both collaborators and keys
func CreateTestBoardFull() models.Board {
	board := CreateDefaultTestBoard()
	board.Collaborators = []models.Collaborator{CreateDefaultTestCollaborator()}
	board.Keys = []models.Key{CreateDefaultTestKey(), CreateSecondTestKey()}
	return board
}

// CreateProcessedTestMedia creates media in processed state
func CreateProcessedTestMedia() models.Media {
	media := CreateDefaultTestAudioMedia()
	media.Status = models.ProcessingStatusDone
	media.ProcessedUrl = "https://storage.example.com/processed/test-audio.webm"
	media.Duration = 5.5
	now := primitive.NewDateTimeFromTime(time.Now())
	media.ProcessingActivity = append(media.ProcessingActivity, models.ProcessingActivity{
		Status:    models.ProcessingStatusDone,
		Message:   "Processing completed",
		CreatedAt: now,
		UpdatedAt: now,
	})
	return media
}

// CreateUserWithResetToken creates a user with reset token set
func CreateUserWithResetToken() models.User {
	user := CreateDefaultTestUser()
	token := "test-reset-token-12345678901234567890"
	expiry := time.Now().Add(30 * time.Minute)
	user.ResetToken = &token
	user.ResetExpiresAt = &expiry
	return user
}

// TestSignupRequest represents a signup request for testing
type TestSignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateValidSignupRequest creates a valid signup request
func CreateValidSignupRequest() TestSignupRequest {
	return TestSignupRequest{
		Name:     "New User",
		Email:    "newuser@example.com",
		Password: "securepassword123",
	}
}

// TestSigninRequest represents a signin request for testing
type TestSigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateValidSigninRequest creates a valid signin request
func CreateValidSigninRequest() TestSigninRequest {
	return TestSigninRequest{
		Email:    "test@example.com",
		Password: "testpassword123",
	}
}

// TestBoardRequest represents a board creation/update request for testing
type TestBoardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Layout      string `json:"layout,omitempty"`
	ImageUrl    string `json:"imageUrl,omitempty"`
}

// CreateValidBoardRequest creates a valid board request
func CreateValidBoardRequest() TestBoardRequest {
	return TestBoardRequest{
		Name:        "New Board",
		Description: "A new test board",
		Layout:      "grid",
	}
}

// TestMediaRequest represents a media creation request for testing
type TestMediaRequest struct {
	MediaType   string `json:"mediaType"`
	FileName    string `json:"fileName"`
	Description string `json:"description,omitempty"`
	FileUrl     string `json:"fileUrl"`
	ContentType string `json:"contentType"`
	FileSize    int64  `json:"fileSize"`
}

// CreateValidMediaRequest creates a valid media request
func CreateValidMediaRequest() TestMediaRequest {
	return TestMediaRequest{
		MediaType:   "audio",
		FileName:    "new-audio.mp3",
		Description: "A new audio file",
		FileUrl:     "https://storage.example.com/uploads/new-audio.mp3",
		ContentType: "audio/mpeg",
		FileSize:    2048000,
	}
}

// TestKeyRequest represents a key creation request for testing
type TestKeyRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	AudioMediaId string `json:"audioMediaId"`
	ImageMediaId string `json:"imageMediaId,omitempty"`
	HotKey       string `json:"hotKey"`
}

// CreateValidKeyRequest creates a valid key request
func CreateValidKeyRequest() TestKeyRequest {
	return TestKeyRequest{
		Name:         "New Key",
		Description:  "A new soundboard key",
		AudioMediaId: TestMediaID.Hex(),
		HotKey:       "C",
	}
}

// TestCollaboratorRequest represents a collaborator creation request for testing
type TestCollaboratorRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// CreateValidCollaboratorRequest creates a valid collaborator request
func CreateValidCollaboratorRequest() TestCollaboratorRequest {
	return TestCollaboratorRequest{
		Email: "newcollab@example.com",
		Role:  "editor",
	}
}
