package audio

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"tuneslap/models"
	"tuneslap/services/storage"
)

func DownloadAudio(mediaClient storage.ObjectStorage, media models.Media) (string, error) {
	// Step 1: Get original file key
	originalFileKey := storage.GetMediaKey(media.AuthorId.Hex(), media.MediaType, media.FileName)

	// Step 2: Download original audio
	downloadedFilePath := filepath.Join(os.TempDir(), originalFileKey)
	err := mediaClient.DownloadFile(context.Background(), originalFileKey, downloadedFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to download original audio: %w", err)
	}

	return downloadedFilePath, nil
}
