package tasks

import (
	"context"
	"fmt"
	"log"
	"tuneslap/config"
	"tuneslap/database"
	"tuneslap/models"
	"tuneslap/seed"
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

	// Check for context cancellation before starting
	select {
	case <-ctx.Done():
		return fmt.Errorf("cleanup task cancelled before start: %w", ctx.Err())
	default:
	}

	// Get collections
	usersCollection := database.GetCollection("users")
	boardsCollection := database.GetCollection("boards")
	mediaCollection := database.GetCollection("media")

	// Find all non-demo users (critical operation)
	cursor, err := usersCollection.Find(ctx, bson.M{
		"email": bson.M{"$ne": config.DemoUserEmail},
	})
	if err != nil {
		log.Printf("[Cleanup] Critical error: failed to find non-demo users: %v", err)
		return fmt.Errorf("critical failure: failed to find non-demo users: %w", err)
	}
	defer cursor.Close(ctx)

	// Check context before decoding
	select {
	case <-ctx.Done():
		return fmt.Errorf("cleanup task cancelled during user query: %w", ctx.Err())
	default:
	}

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		log.Printf("[Cleanup] Critical error: failed to decode users: %v", err)
		return fmt.Errorf("critical failure: failed to decode users: %w", err)
	}

	log.Printf("[Cleanup] Found %d non-demo users to clean up", len(users))

	// Get storage client for media deletion (non-critical if fails)
	mediaStorage, err := storage.GetMediaStorage()
	if err != nil {
		log.Printf("[Cleanup] Warning: Could not get media storage, media files will not be deleted: %v", err)
	}

	// Track statistics for summary
	stats := struct {
		usersProcessed      int
		usersDeleted        int
		usersFailed         int
		mediaDeleted        int64
		mediaFailed         int
		boardsDeleted       int64
		boardsFailed        int
		collaboratorsFailed int
	}{
		usersProcessed: len(users),
	}

	// For each non-demo user, delete their data
	for i, user := range users {
		// Check context before processing each user
		select {
		case <-ctx.Done():
			log.Printf("[Cleanup] Task cancelled after processing %d/%d users", i, len(users))
			return fmt.Errorf("cleanup task cancelled during processing: %w", ctx.Err())
		default:
		}

		log.Printf("[Cleanup] Cleaning up data for user: %s (%s)", user.Email, user.ID.Hex())

		userDeleted := false

		// 1. Delete user's media files from storage and database
		mediaCursor, err := mediaCollection.Find(ctx, bson.M{"authorId": user.ID})
		if err != nil {
			log.Printf("[Cleanup] Error finding media for user %s: %v", user.ID.Hex(), err)
		} else {
			var mediaItems []models.Media
			if err := mediaCursor.All(ctx, &mediaItems); err != nil {
				log.Printf("[Cleanup] Error decoding media for user %s: %v", user.ID.Hex(), err)
			} else {
				// Delete each media file from storage (non-critical failures)
				for _, media := range mediaItems {
					if mediaStorage != nil {
						mediaKey := storage.GetMediaKey(media.AuthorId.Hex(), media.MediaType, media.FileName)
						if err := mediaStorage.DeleteFile(ctx, mediaKey); err != nil {
							log.Printf("[Cleanup] Warning: Failed to delete media file %s: %v", mediaKey, err)
							stats.mediaFailed++
						}
					}
				}
			}
			mediaCursor.Close(ctx)
		}

		// Delete media records from database
		deleteResult, err := mediaCollection.DeleteMany(ctx, bson.M{"authorId": user.ID})
		if err != nil {
			log.Printf("[Cleanup] Error deleting media records for user %s: %v", user.ID.Hex(), err)
			stats.mediaFailed++
		} else {
			stats.mediaDeleted += deleteResult.DeletedCount
			log.Printf("[Cleanup] Deleted %d media records for user %s", deleteResult.DeletedCount, user.ID.Hex())
		}

		// 2. Delete user's boards (keys and collaborators are embedded)
		deleteResult, err = boardsCollection.DeleteMany(ctx, bson.M{"authorId": user.ID})
		if err != nil {
			log.Printf("[Cleanup] Error deleting boards for user %s: %v", user.ID.Hex(), err)
			stats.boardsFailed++
		} else {
			stats.boardsDeleted += deleteResult.DeletedCount
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
			stats.collaboratorsFailed++
		}

		// 4. Delete the user (critical for this user, but non-critical for overall task)
		_, err = usersCollection.DeleteOne(ctx, bson.M{"_id": user.ID})
		if err != nil {
			log.Printf("[Cleanup] Error deleting user %s: %v", user.ID.Hex(), err)
			stats.usersFailed++
		} else {
			stats.usersDeleted++
			userDeleted = true
			log.Printf("[Cleanup] Deleted user: %s", user.Email)
		}

		// Log user completion status
		if !userDeleted {
			log.Printf("[Cleanup] Warning: User %s cleanup completed with some failures", user.Email)
		}
	}

	// Check context before re-seeding
	select {
	case <-ctx.Done():
		return fmt.Errorf("cleanup task cancelled before re-seeding: %w", ctx.Err())
	default:
	}

	// Log summary statistics
	log.Printf("[Cleanup] Cleanup summary: processed=%d, deleted=%d, failed=%d, media_deleted=%d, media_failed=%d, boards_deleted=%d, boards_failed=%d, collaborators_failed=%d",
		stats.usersProcessed, stats.usersDeleted, stats.usersFailed,
		stats.mediaDeleted, stats.mediaFailed, stats.boardsDeleted, stats.boardsFailed, stats.collaboratorsFailed)

	// Re-seed demo data immediately after cleanup (critical operation)
	log.Println("[Cleanup] Re-seeding demo data after cleanup...")
	if err := seed.EnsureDemoData(); err != nil {
		log.Printf("[Cleanup] Critical error: Failed to re-seed demo data after cleanup: %v", err)
		return fmt.Errorf("critical failure: failed to re-seed demo data: %w", err)
	}

	log.Println("[Cleanup] Demo data re-seeded successfully")
	log.Printf("[Cleanup] Demo cleanup task completed successfully. Cleaned up %d users.", stats.usersDeleted)
	return nil
}

// RegisterCleanupTasks registers cleanup task handlers to the provided mux
func RegisterCleanupTasks(mux *asynq.ServeMux) {
	log.Println("[Worker] Registering cleanup:demo task handler")
	mux.HandleFunc(TypeDemoCleanup, HandleDemoCleanupTask)
}
