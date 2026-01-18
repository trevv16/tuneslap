package image

import (
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

func normalizeDefault(input []byte) ([]byte, error) {
	return bimg.NewImage(input).Process(defaultNormalizeOptions)
}

// func normalize(input []byte, params models.ImageProcessingParams) ([]byte, error) {
// 	options := bimg.Options{
// 		Width:         defaultNormalizeOptions.Width,
// 		Height:        defaultNormalizeOptions.Height,
// 		Quality:       defaultNormalizeOptions.Quality,
// 		Type:          bimg.WEBP,
// 		Compression:   defaultNormalizeOptions.Compression,
// 		StripMetadata: true, // Should always be true
// 	}

// 	if params.ResizeTo[0] > 0 && params.ResizeTo[1] > 0 {
// 		options.Width = params.ResizeTo[0]
// 		options.Height = params.ResizeTo[1]
// 	}

// 	return bimg.NewImage(input).Process(options)
// }
