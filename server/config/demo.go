package config

import (
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Demo constants - single source of truth for demo data
const (
	DemoUserEmail    = "demo@tuneslap.com"
	DemoUserName     = "Demo User"
	DemoUserPassword = "demo123456" // Not a real password since demo user cannot login
	DemoBoardName    = "Demo Board"
)

// E2E Test User constants - separate from demo user for test isolation
// These credentials are used by Playwright E2E tests
const (
	E2ETestUserEmail    = "e2e-test@tuneslap.test"
	E2ETestUserName     = "E2E Test User"
	E2ETestUserPassword = "e2e-test-password-123"
	E2ETestBoardName    = "E2E Test Board"
)

// Demo IDs - fixed ObjectIDs for predictable seeding
var (
	DemoUserID  = mustObjectID("000000000000000000000001")
	DemoBoardID = mustObjectID("000000000000000000000002")

	// E2E Test IDs - fixed ObjectIDs for predictable E2E testing
	E2ETestUserID  = mustObjectID("000000000000000000000099")
	E2ETestBoardID = mustObjectID("000000000000000000000098")

	DemoKeyIDs = []primitive.ObjectID{
		mustObjectID("000000000000000000000101"),
		mustObjectID("000000000000000000000102"),
		mustObjectID("000000000000000000000103"),
		mustObjectID("000000000000000000000104"),
		mustObjectID("000000000000000000000105"),
		mustObjectID("000000000000000000000106"),
		mustObjectID("000000000000000000000107"),
		mustObjectID("000000000000000000000108"),
	}
	// Audio media IDs for demo keys
	DemoAudioMediaIDs = []primitive.ObjectID{
		mustObjectID("000000000000000000000201"), // applause
		mustObjectID("000000000000000000000202"), // drum-roll
		mustObjectID("000000000000000000000203"), // laughter
		mustObjectID("000000000000000000000204"), // air-horn
		mustObjectID("000000000000000000000205"), // whoosh
		mustObjectID("000000000000000000000206"), // bell-ding
		mustObjectID("000000000000000000000207"), // boing
		mustObjectID("000000000000000000000208"), // tada
	}
	// Image media IDs for demo keys
	DemoImageMediaIDs = []primitive.ObjectID{
		mustObjectID("000000000000000000000301"), // applause image
		mustObjectID("000000000000000000000302"), // drum-roll image
		mustObjectID("000000000000000000000303"), // laughter image
		mustObjectID("000000000000000000000304"), // air-horn image
		mustObjectID("000000000000000000000305"), // whoosh image
		mustObjectID("000000000000000000000306"), // bell-ding image
		mustObjectID("000000000000000000000307"), // boing image
		mustObjectID("000000000000000000000308"), // tada image
	}
)

// DemoKey represents a key in the demo board
type DemoKey struct {
	ID           primitive.ObjectID
	Name         string
	HotKey       string
	AudioURL     string
	ImageURL     string
	AudioMediaID primitive.ObjectID
	ImageMediaID primitive.ObjectID
}

// DemoAudioFile represents an audio file for demo seeding
type DemoAudioFile struct {
	MediaID   primitive.ObjectID
	FileName  string
	LocalPath string // Path in server/assets/demo/ or frontend/public/demo/audio/
}

// GetDemoKeys returns the 8 demo keys with their data
// Audio files should be placed in server/assets/demo/ or frontend/public/demo/audio/
// Each audio file should be 3-10 seconds long for best demo experience
func GetDemoKeys() []DemoKey {
	return []DemoKey{
		{
			ID:           DemoKeyIDs[0],
			Name:         "Applause",
			HotKey:       "1",
			AudioURL:     "/demo/audio/applause.mp3",
			ImageURL:     "https://images.unsplash.com/photo-1540575467063-178a50c2df87?w=400",
			AudioMediaID: DemoAudioMediaIDs[0],
			ImageMediaID: DemoImageMediaIDs[0],
		},
		{
			ID:           DemoKeyIDs[1],
			Name:         "Drum Roll",
			HotKey:       "2",
			AudioURL:     "/demo/audio/drum-roll.mp3",
			ImageURL:     "https://images.unsplash.com/photo-1519892300165-cb5542fb47c7?w=400",
			AudioMediaID: DemoAudioMediaIDs[1],
			ImageMediaID: DemoImageMediaIDs[1],
		},
		{
			ID:           DemoKeyIDs[2],
			Name:         "Laughter",
			HotKey:       "3",
			AudioURL:     "/demo/audio/laughter.mp3",
			ImageURL:     "https://images.unsplash.com/photo-1543610892-0b1f7e6d8ac1?w=400",
			AudioMediaID: DemoAudioMediaIDs[2],
			ImageMediaID: DemoImageMediaIDs[2],
		},
		{
			ID:           DemoKeyIDs[3],
			Name:         "Air Horn",
			HotKey:       "4",
			AudioURL:     "/demo/audio/air-horn.mp3",
			ImageURL:     "https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=400",
			AudioMediaID: DemoAudioMediaIDs[3],
			ImageMediaID: DemoImageMediaIDs[3],
		},
		{
			ID:           DemoKeyIDs[4],
			Name:         "Whoosh",
			HotKey:       "5",
			AudioURL:     "/demo/audio/whoosh.mp3",
			ImageURL:     "https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=400",
			AudioMediaID: DemoAudioMediaIDs[4],
			ImageMediaID: DemoImageMediaIDs[4],
		},
		{
			ID:           DemoKeyIDs[5],
			Name:         "Bell Ding",
			HotKey:       "6",
			AudioURL:     "/demo/audio/bell-ding.mp3",
			ImageURL:     "https://images.unsplash.com/photo-1513836279014-a89f7a76ae86?w=400",
			AudioMediaID: DemoAudioMediaIDs[5],
			ImageMediaID: DemoImageMediaIDs[5],
		},
		{
			ID:           DemoKeyIDs[6],
			Name:         "Boing",
			HotKey:       "7",
			AudioURL:     "/demo/audio/boing.mp3",
			ImageURL:     "https://images.unsplash.com/photo-1518640467707-6811f4a6ab73?w=400",
			AudioMediaID: DemoAudioMediaIDs[6],
			ImageMediaID: DemoImageMediaIDs[6],
		},
		{
			ID:           DemoKeyIDs[7],
			Name:         "Ta-Da",
			HotKey:       "8",
			AudioURL:     "/demo/audio/tada.mp3",
			ImageURL:     "https://images.unsplash.com/photo-1492684223066-81342ee5ff30?w=400",
			AudioMediaID: DemoAudioMediaIDs[7],
			ImageMediaID: DemoImageMediaIDs[7],
		},
	}
}

// GetDemoAudioFiles returns the list of audio files for demo seeding
func GetDemoAudioFiles() []DemoAudioFile {
	return []DemoAudioFile{
		{MediaID: DemoAudioMediaIDs[0], FileName: "applause.mp3", LocalPath: "assets/demo/applause.mp3"},
		{MediaID: DemoAudioMediaIDs[1], FileName: "drum-roll.mp3", LocalPath: "assets/demo/drum-roll.mp3"},
		{MediaID: DemoAudioMediaIDs[2], FileName: "laughter.mp3", LocalPath: "assets/demo/laughter.mp3"},
		{MediaID: DemoAudioMediaIDs[3], FileName: "air-horn.mp3", LocalPath: "assets/demo/air-horn.mp3"},
		{MediaID: DemoAudioMediaIDs[4], FileName: "whoosh.mp3", LocalPath: "assets/demo/whoosh.mp3"},
		{MediaID: DemoAudioMediaIDs[5], FileName: "bell-ding.mp3", LocalPath: "assets/demo/bell-ding.mp3"},
		{MediaID: DemoAudioMediaIDs[6], FileName: "boing.mp3", LocalPath: "assets/demo/boing.mp3"},
		{MediaID: DemoAudioMediaIDs[7], FileName: "tada.mp3", LocalPath: "assets/demo/tada.mp3"},
	}
}

func mustObjectID(hex string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		panic("invalid ObjectID: " + hex)
	}
	return id
}

// IsDemoUser checks if a user ID is the demo user
func IsDemoUser(userID primitive.ObjectID) bool {
	return userID == DemoUserID
}

// IsDemoBoard checks if a board ID is the demo board
func IsDemoBoard(boardID primitive.ObjectID) bool {
	return boardID == DemoBoardID
}

// IsE2ETestMode checks if E2E test mode is enabled via environment variable
func IsE2ETestMode() bool {
	return os.Getenv("E2E_TEST_MODE") == "true"
}

// IsE2ETestUser checks if a user ID is the E2E test user
func IsE2ETestUser(userID primitive.ObjectID) bool {
	return userID == E2ETestUserID
}

// IsE2ETestBoard checks if a board ID is the E2E test board
func IsE2ETestBoard(boardID primitive.ObjectID) bool {
	return boardID == E2ETestBoardID
}
