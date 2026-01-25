package image

import (
	"tuneslap/models"

	"github.com/h2non/bimg"
)

var defaultNormalizeOptions = bimg.Options{
	Width:         500,
	Height:        500,
	Quality:       90,
	Type:          bimg.WEBP,
	Compression:   6,
	StripMetadata: true, // Should always be true
}

// NormalizeDefault processes an image with default options (500x500, WEBP, quality 90)
func NormalizeDefault(input []byte) ([]byte, error) {
	return bimg.NewImage(input).Process(defaultNormalizeOptions)
}

// stringToImageType converts a format string to bimg.ImageType
func stringToImageType(format string) bimg.ImageType {
	switch format {
	case "jpeg", "jpg":
		return bimg.JPEG
	case "png":
		return bimg.PNG
	case "webp":
		return bimg.WEBP
	case "gif":
		return bimg.GIF
	case "svg":
		return bimg.SVG
	case "tiff":
		return bimg.TIFF
	default:
		return bimg.WEBP // Default to WEBP
	}
}

// Normalize processes an image with custom dimensions from params
// Falls back to default dimensions if params.ResizeTo is not set
func Normalize(input []byte, params models.ImageProcessingParams) ([]byte, error) {
	options := bimg.Options{
		Width:         defaultNormalizeOptions.Width,
		Height:        defaultNormalizeOptions.Height,
		Quality:       defaultNormalizeOptions.Quality,
		Type:          bimg.WEBP,
		Compression:   defaultNormalizeOptions.Compression,
		StripMetadata: true, // Should always be true
	}

	// Use custom dimensions if provided
	if params.ResizeTo[0] > 0 && params.ResizeTo[1] > 0 {
		options.Width = params.ResizeTo[0]
		options.Height = params.ResizeTo[1]
	}

	// Support custom output format if specified
	if params.Format != "" {
		options.Type = stringToImageType(params.Format)
	}

	return bimg.NewImage(input).Process(options)
}

// GetDefaultNormalizeOptions returns a copy of the default normalize options
// Useful for testing and configuration inspection
func GetDefaultNormalizeOptions() bimg.Options {
	return defaultNormalizeOptions
}
