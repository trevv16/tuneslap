package testutils

import (
	"context"
	"fmt"
	"os"
	"time"
	"tuneslap/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	TestDBName     = "tuneslap_test"
	TestMongoURI   = "mongodb://localhost:27017"
	DefaultTimeout = 10 * time.Second
)

var testMongoClient *mongo.Client

// SetupTestMongoDB initializes a test MongoDB connection
// Returns a cleanup function to be called after tests
func SetupTestMongoDB() (func(), error) {
	// Save original environment
	originalURI := os.Getenv("MONGODB_URI")
	originalDB := os.Getenv("DATABASE")

	// Set test environment
	os.Setenv("MONGODB_URI", TestMongoURI)
	os.Setenv("DATABASE", TestDBName)

	// Start MongoDB
	err := database.StartMongoDB()
	if err != nil {
		// Restore environment on error
		os.Setenv("MONGODB_URI", originalURI)
		os.Setenv("DATABASE", originalDB)
		return nil, fmt.Errorf("failed to start test MongoDB: %w", err)
	}

	// Return cleanup function
	cleanup := func() {
		// Clear test collections before closing
		ClearTestCollections()
		database.CloseMongoDB()
		// Restore original environment
		os.Setenv("MONGODB_URI", originalURI)
		os.Setenv("DATABASE", originalDB)
	}

	return cleanup, nil
}

// SetupTestMongoDBWithSkip sets up MongoDB and skips the test if unavailable
func SetupTestMongoDBWithSkip(t TestingT) func() {
	cleanup, err := SetupTestMongoDB()
	if err != nil {
		t.Skipf("Skipping test - MongoDB not available: %v", err)
		return func() {}
	}
	return cleanup
}

// ClearTestCollections clears all test collections
func ClearTestCollections() error {
	collections := []string{"users", "boards", "media"}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	for _, collName := range collections {
		coll := database.GetCollection(collName)
		if coll != nil {
			_, err := coll.DeleteMany(ctx, bson.M{})
			if err != nil {
				return fmt.Errorf("failed to clear collection %s: %w", collName, err)
			}
		}
	}

	return nil
}

// ClearCollection clears a specific collection
func ClearCollection(collectionName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	coll := database.GetCollection(collectionName)
	if coll == nil {
		return fmt.Errorf("collection %s not found", collectionName)
	}

	_, err := coll.DeleteMany(ctx, bson.M{})
	return err
}

// InsertTestData inserts test documents into a collection
func InsertTestData(collectionName string, documents []interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	coll := database.GetCollection(collectionName)
	if coll == nil {
		return fmt.Errorf("collection %s not found", collectionName)
	}

	_, err := coll.InsertMany(ctx, documents)
	return err
}

// InsertOne inserts a single document into a collection
func InsertOne(collectionName string, document interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	coll := database.GetCollection(collectionName)
	if coll == nil {
		return fmt.Errorf("collection %s not found", collectionName)
	}

	_, err := coll.InsertOne(ctx, document)
	return err
}

// GetTestCollection returns a test collection
func GetTestCollection(name string) *mongo.Collection {
	return database.GetCollection(name)
}

// CountDocuments counts documents in a collection matching a filter
func CountDocuments(collectionName string, filter bson.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	coll := database.GetCollection(collectionName)
	if coll == nil {
		return 0, fmt.Errorf("collection %s not found", collectionName)
	}

	return coll.CountDocuments(ctx, filter)
}

// FindOne finds a single document in a collection
func FindOne(collectionName string, filter bson.M, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	coll := database.GetCollection(collectionName)
	if coll == nil {
		return fmt.Errorf("collection %s not found", collectionName)
	}

	return coll.FindOne(ctx, filter).Decode(result)
}

// CreateTestIndexes creates indexes for test collections
func CreateTestIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	// Users collection indexes
	usersCol := database.GetCollection("users")
	if usersCol != nil {
		_, err := usersCol.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		})
		if err != nil {
			return fmt.Errorf("failed to create users email index: %w", err)
		}
	}

	// Boards collection indexes
	boardsCol := database.GetCollection("boards")
	if boardsCol != nil {
		_, err := boardsCol.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bson.D{{Key: "authorId", Value: 1}},
		})
		if err != nil {
			return fmt.Errorf("failed to create boards authorId index: %w", err)
		}
	}

	// Media collection indexes
	mediaCol := database.GetCollection("media")
	if mediaCol != nil {
		_, err := mediaCol.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bson.D{{Key: "authorId", Value: 1}},
		})
		if err != nil {
			return fmt.Errorf("failed to create media authorId index: %w", err)
		}
	}

	return nil
}
