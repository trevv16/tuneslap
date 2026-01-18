package audio

import (
	"io"
	"os"
	"strings"
)

func GetProcessedFileName(fileName string) string {
	return strings.Join([]string{"processed", fileName}, "-")
}

func GetAudioFileDataFromPath(path string) ([]byte, error) {
	// open the file
	fileData, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fileData.Close()

	return io.ReadAll(fileData)
}
