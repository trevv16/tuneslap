package image

import (
	"strings"

	"github.com/h2non/bimg"
)

func GetProcessedFileName(fileName string) string {
	return strings.Join([]string{"processed", fileName}, "-")
}

func GetContentTypeFromFormat(format int) string {
	switch format {
	case int(bimg.WEBP):
		return "image/webp"
	case int(bimg.PNG):
		return "image/png"
	default:
		return "image/png"
	}
}
