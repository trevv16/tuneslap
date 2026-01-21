package tasks

import (
	"context"
	"fmt"
	"log"
	"tuneslap/config"
	"tuneslap/database"
	"tuneslap/models"
	"tuneslap/services/storage"

	"github.com/hibiken/asynq"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	TypeDemoCleanup = "cleanup:demo"
)

// NewDemoCleanupTask creates a new demo cleanup task
func NewDemoCleanupTask() (*asynq.Task, error) {
	return asynq.NewTask(TypeDemoCleanup, nil), nil
}

// HandleDemoCleanupTask deletes all non-demo data from the database
// This includes all users except the demo user, and all their associated data
func HandleDemoCleanupTask(ctx context.Context, task *asynq.Task) error {
	log.Println("[Cleanup] Starting demo cleanup task...")

	// Get collections
	usersCollection := database.GetCollection("users")
	boardsCollection := database.GetCollection("boards")
	mediaCollection := database.GetCollection("media")

	// Find all non-demo users
	cursor, err := usersCollection.Find(ctx, bson.M{
		"email": bson.M{"$ne": config.DemoUserEmail},
	})
	if err != nil {
		log.Printf("[Cleanup] Error finding non-demo users: %v", err)
		return fmt.Errorf("failed to find non-demo users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		log.Printf("[Cleanup] Error decoding users: %v", err)
		return fmt.Errorf("failed to decode users: %w", err)
	}

	log.Printf("[Cleanup] Found %d non-demo users to clean up", len(users))

	// Get storage client for media deletion
	mediaStorage, err := storage.GetMediaStorage()
	if err != nil {
		log.Printf("[Cleanup] Warning: Could not get media storage, media files will not be deleted: %v", err)
	}

	// For each non-demo user, delete their data
	for _, user := range users {
		log.Printf("[Cleanup] Cleaning up data for user: %s (%s)", user.Email, user.ID.Hex())

		// 1. Delete user's media files from storage and database
		mediaCursor, err := mediaCollection.Find(ctx, bson.M{"authorId": user.ID})
		if err != nil {
			log.Printf("[Cleanup] Error finding media for user %s: %v", user.ID.Hex(), err)
		} else {
			var mediaItems []models.Media
			if err := mediaCursor.All(ctx, &mediaItems); err != nil {
				log.Printf("[Cleanup] Error decoding media for user %s: %v", user.ID.Hex(), err)
			} else {
				// Delete each media file from storage
				for _, media := range mediaItems {
					if mediaStorage != nil {
						mediaKey := storage.GetMediaKey(media.AuthorId.Hex(), media.MediaType, media.FileName)
						if err := mediaStorage.DeleteFile(ctx, mediaKey); err != nil {
							log.Printf("[Cleanup] Warning: Failed to delete media file %s: %v", mediaKey, err)
						}
					}
				}
			}
			mediaCursor.Close(ctx)
		}

		// Delete media records from database
		deleteResult, err := mediaCollection.DeleteMany(ctx, bson.M{"authorId": user.ID})
		if err != nil {
			log.Printf("[Cleanup] Error deleting media for user %s: %v", user.ID.Hex(), err)
		} else {
			log.Printf("[Cleanup] Deleted %d media records for user %s", deleteResult.DeletedCount, user.ID.Hex())
		}

		// 2. Delete user's boards (keys and collaborators are embedded)
		deleteResult, err = boardsCollection.DeleteMany(ctx, bson.M{"authorId": user.ID})
		if err != nil {
			log.Printf("[Cleanup] Error deleting boards for user %s: %v", user.ID.Hex(), err)
		} else {
			log.Printf("[Cleanup] Deleted %d boards for user %s", deleteResult.DeletedCount, user.ID.Hex())
		}

		// 3. Remove user from any collaborator lists on other boards
		_, err = boardsCollection.UpdateMany(ctx,
			bson.M{},
			bson.M{
				"$pull": bson.M{
					"collaborators": bson.M{"userId": user.ID},
				},
			},
		)
		if err != nil {
			log.Printf("[Cleanup] Error removing user %s from collaborator lists: %v", user.ID.Hex(), err)
		}

		// 4. Delete the user
		_, err = usersCollection.DeleteOne(ctx, bson.M{"_id": user.ID})
		if err != nil {
			log.Printf("[Cleanup] Error deleting user %s: %v", user.ID.Hex(), err)
		} else {
			log.Printf("[Cleanup] Deleted user: %s", user.Email)
		}
	}

	// Note: Demo data re-seeding happens on next app startup
	// We don't call it here to avoid import cycle with app package

	log.Printf("[Cleanup] Demo cleanup task completed. Cleaned up %d users.", len(users))
	return nil
}

// RegisterCleanupTasks registers cleanup task handlers to the provided mux
func RegisterCleanupTasks(mux *asynq.ServeMux) {
	log.Println("[Worker] Registering cleanup:demo task handler")
	mux.HandleFunc(TypeDemoCleanup, HandleDemoCleanupTask)
}
