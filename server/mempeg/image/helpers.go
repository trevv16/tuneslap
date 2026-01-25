package image

import (
	"strings"
)

func GetProcessedFileName(fileName string) string {
	return strings.Join([]string{"processed", fileName}, "-")
}

func GetContentTypeFromFormat(format string) string {
	switch format {
	case "webp":
		return "image/webp"
	case "png":
		return "image/png"
	case "jpeg", "jpg":
		return "image/jpeg"
	case "gif":
		return "image/gif"
	case "svg":
		return "image/svg+xml"
	case "tiff":
		return "image/tiff"
	default:
		return "image/webp"
	}
}
