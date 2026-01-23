package seed

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	"tuneslap/config"
	"tuneslap/database"
	"tuneslap/models"
	"tuneslap/services/storage"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// EnsureDemoData seeds the demo user, board, and keys if in demo mode
func EnsureDemoData() error {
	if !config.IsDemoMode() {
		log.Println("[Seed] Not in demo mode, skipping demo data seeding")
		return nil
	}

	log.Println("[Seed] Demo mode enabled, ensuring demo data exists...")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Seed demo user
	if err := seedDemoUser(ctx); err != nil {
		return err
	}

	// Seed demo media (audio files to MinIO)
	if err := seedDemoMedia(ctx); err != nil {
		return err
	}

	// Seed demo board with keys
	if err := seedDemoBoard(ctx); err != nil {
		return err
	}

	log.Println("[Seed] Demo data seeding completed successfully")
	return nil
}

func seedDemoUser(ctx context.Context) error {
	usersCollection := database.GetCollection("users")

	// Check if demo user exists
	var existingUser models.User
	err := usersCollection.FindOne(ctx, bson.M{"_id": config.DemoUserID}).Decode(&existingUser)
	if err == nil {
		log.Println("[Seed] Demo user already exists, skipping...")
		return nil
	}

	// Hash password (even though demo user won't login)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(config.DemoUserPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create demo user
	demoUser := models.User{
		ID:           config.DemoUserID,
		Name:         config.DemoUserName,
		Email:        config.DemoUserEmail,
		PasswordHash: string(hashedPassword),
		ImageUrl:     "",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	opts := options.Update().SetUpsert(true)
	_, err = usersCollection.UpdateOne(
		ctx,
		bson.M{"_id": config.DemoUserID},
		bson.M{"$setOnInsert": demoUser},
		opts,
	)
	if err != nil {
		log.Printf("[Seed] Error creating demo user: %v", err)
		return err
	}

	log.Println("[Seed] Demo user created successfully")
	return nil
}

// seedDemoMedia uploads audio files to MinIO and creates Media records
func seedDemoMedia(ctx context.Context) error {
	mediaCollection := database.GetCollection("media")
	audioFiles := config.GetDemoAudioFiles()

	// Get media storage
	mediaStorage, err := storage.GetMediaStorage()
	if err != nil {
		log.Printf("[Seed] Warning: Could not get media storage, skipping media upload: %v", err)
		// Don't fail - demo can still work with fallback URLs
		return nil
	}

	now := primitive.NewDateTimeFromTime(time.Now())
	uploadedCount := 0

	for _, af := range audioFiles {
		// Try to find the audio file first
		audioPath := findDemoAudioFile(af.LocalPath)
		if audioPath == "" {
			log.Printf("[Seed] Warning: Demo audio file not found: %s (skipping upload)", af.FileName)
			continue
		}

		// Get file info
		fileInfo, err := os.Stat(audioPath)
		if err != nil {
			log.Printf("[Seed] Warning: Could not stat demo audio file %s: %v", af.FileName, err)
			continue
		}

		// Check if media record already exists
		var existingMedia models.Media
		err = mediaCollection.FindOne(ctx, bson.M{"_id": af.MediaID}).Decode(&existingMedia)
		needsUpload := false
		if err != nil {
			// Media record doesn't exist, need to upload
			needsUpload = true
		} else {
			// Media record exists, verify the file actually exists in storage
			objectKey := storage.GetMediaKey(config.DemoUserID.Hex(), "audio", af.FileName)
			fileExists, checkErr := mediaStorage.FileExists(ctx, objectKey)
			if checkErr != nil {
				log.Printf("[Seed] Warning: Could not verify if file exists for %s: %v, will attempt re-upload", af.FileName, checkErr)
				needsUpload = true
			} else if !fileExists {
				// File doesn't exist in storage, even though record exists
				log.Printf("[Seed] Demo media %s record exists but file not found in storage, re-uploading...", af.FileName)
				needsUpload = true
			} else if existingMedia.FileUrl == "" || existingMedia.ProcessedUrl == "" {
				// File exists but URLs are missing, update the record
				log.Printf("[Seed] Demo media %s file exists but URLs missing, updating record...", af.FileName)
				fileURL := mediaStorage.GetFileURL(objectKey)
				if fileURL != "" {
					opts := options.Update().SetUpsert(true)
					_, updateErr := mediaCollection.UpdateOne(
						ctx,
						bson.M{"_id": af.MediaID},
						bson.M{"$set": bson.M{
							"fileUrl":      fileURL,
							"processedUrl": fileURL,
							"updatedAt":    now,
						}},
						opts,
					)
					if updateErr != nil {
						log.Printf("[Seed] Warning: Could not update URLs for %s: %v", af.FileName, updateErr)
					} else {
						log.Printf("[Seed] Updated file URLs for %s", af.FileName)
					}
				}
				continue
			} else {
				// File exists and URLs are set, all good
				log.Printf("[Seed] Demo media %s already exists with file in storage, skipping...", af.FileName)
				continue
			}
		}

		if !needsUpload {
			continue
		}

		// Upload to MinIO
		objectKey := storage.GetMediaKey(config.DemoUserID.Hex(), "audio", af.FileName)
		uploadReq := storage.UploadFileRequest{
			ObjectName:  objectKey,
			FilePath:    audioPath,
			ContentType: "audio/mpeg",
			Metadata: map[string]string{
				"demo": "true",
			},
		}

		if err := mediaStorage.UploadFile(ctx, uploadReq); err != nil {
			log.Printf("[Seed] Warning: Could not upload demo audio file %s: %v", af.FileName, err)
			continue
		}

		// Get the URL for the uploaded file
		fileURL := mediaStorage.GetFileURL(objectKey)
		if fileURL == "" {
			log.Printf("[Seed] Warning: Got empty file URL for %s after upload, skipping record creation", af.FileName)
			continue
		}

		// Validate file size is greater than 0
		if fileInfo.Size() == 0 {
			log.Printf("[Seed] Warning: Demo audio file %s has zero size, skipping record creation", af.FileName)
			continue
		}

		// Only create or update media record after successful upload and validation
		media := models.Media{
			ID:           af.MediaID,
			AuthorId:     config.DemoUserID,
			MediaType:    "audio",
			FileName:     af.FileName,
			Description:  fmt.Sprintf("Demo sound: %s", af.FileName),
			FileUrl:      fileURL,
			ProcessedUrl: fileURL, // Use same URL since it's already processed
			WaveformUrl:  "",      // No waveform for demo files
			ContentType:  "audio/mpeg",
			FileSize:     fileInfo.Size(),
			Status:       models.ProcessingStatusDone,
			Duration:     5.0, // Approximate duration for demo files
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		opts := options.Update().SetUpsert(true)
		_, err = mediaCollection.UpdateOne(
			ctx,
			bson.M{"_id": af.MediaID},
			bson.M{"$set": media},
			opts,
		)
		if err != nil {
			log.Printf("[Seed] Error: Could not create/update media record for %s: %v", af.FileName, err)
			// Don't continue here - the file was uploaded but record creation failed
			// This is a more serious error, but we'll log it and continue with other files
			continue
		}

		uploadedCount++
		log.Printf("[Seed] Uploaded demo audio: %s -> %s", af.FileName, fileURL)
	}

	log.Printf("[Seed] Demo media seeding completed: %d/%d files uploaded", uploadedCount, len(audioFiles))
	return nil
}

// findDemoAudioFile tries to find the audio file in server/assets/demo/
func findDemoAudioFile(localPath string) string {
	// All demo audio files should be in server/assets/demo/
	// Try relative paths for local dev and absolute for Docker
	paths := []string{
		localPath,                        // Direct path (assets/demo/filename.mp3)
		filepath.Join(".", localPath),    // Relative to server/
		filepath.Join("/app", localPath), // Docker container path
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return ""
}

func seedDemoBoard(ctx context.Context) error {
	boardsCollection := database.GetCollection("boards")
	mediaCollection := database.GetCollection("media")

	// Check if demo board exists
	var existingBoard models.Board
	err := boardsCollection.FindOne(ctx, bson.M{"_id": config.DemoBoardID}).Decode(&existingBoard)
	if err == nil {
		log.Println("[Seed] Demo board already exists, skipping...")
		return nil
	}

	// Build keys from demo data
	demoKeys := config.GetDemoKeys()
	keys := make([]models.Key, len(demoKeys))
	now := time.Now()

	for i, dk := range demoKeys {
		// Try to get the audio URL from the media record if it exists
		audioURL := dk.AudioURL // Fallback to default URL
		var audioMedia models.Media
		if err := mediaCollection.FindOne(ctx, bson.M{"_id": dk.AudioMediaID}).Decode(&audioMedia); err == nil {
			// Use the processed URL from the media record
			if audioMedia.ProcessedUrl != "" {
				audioURL = audioMedia.ProcessedUrl
			} else if audioMedia.FileUrl != "" {
				audioURL = audioMedia.FileUrl
			}
		}

		keys[i] = models.Key{
			ID:           dk.ID,
			BoardId:      config.DemoBoardID,
			Name:         dk.Name,
			Description:  "Demo sound: " + dk.Name,
			HotKey:       dk.HotKey,
			AudioMediaId: dk.AudioMediaID,
			AudioUrl:     audioURL,
			ImageMediaId: dk.ImageMediaID,
			ImageUrl:     dk.ImageURL, // Keep using external image URLs for now
			CreatedAt:    now,
			UpdatedAt:    now,
		}
	}

	// Create demo board
	demoBoard := models.Board{
		ID:            config.DemoBoardID,
		AuthorId:      config.DemoUserID,
		Name:          config.DemoBoardName,
		Description:   "Try out TuneSlap with this demo soundboard!",
		Layout:        models.GridLayout,
		ImageUrl:      "",
		Collaborators: []models.Collaborator{},
		Keys:          keys,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	opts := options.Update().SetUpsert(true)
	_, err = boardsCollection.UpdateOne(
		ctx,
		bson.M{"_id": config.DemoBoardID},
		bson.M{"$setOnInsert": demoBoard},
		opts,
	)
	if err != nil {
		log.Printf("[Seed] Error creating demo board: %v", err)
		return err
	}

	log.Printf("[Seed] Demo board created successfully with %d keys", len(keys))
	return nil
}
