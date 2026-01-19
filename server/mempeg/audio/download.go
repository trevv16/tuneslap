package audio

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"tuneslap/models"
	"tuneslap/services/storage"
)

// DirPermissions is the default permission for created directories (rwxr-xr-x)
const DirPermissions = 0755

func DownloadAudio(mediaClient storage.ObjectStorage, media models.Media) (string, error) {
	// Step 1: Get original file key
	originalFileKey := storage.GetMediaKey(media.AuthorId.Hex(), media.MediaType, media.FileName)

	// Step 2: Create destination path and ensure directory exists
	downloadedFilePath := filepath.Join(os.TempDir(), originalFileKey)
	downloadDir := filepath.Dir(downloadedFilePath)
	if err := os.MkdirAll(downloadDir, DirPermissions); err != nil {
		return "", fmt.Errorf("failed to create download directory: %w", err)
	}

	// Step 3: Download original audio
	err := mediaClient.DownloadFile(context.Background(), originalFileKey, downloadedFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to download original audio: %w", err)
	}

	return downloadedFilePath, nil
}
