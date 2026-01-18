package image

import (
	"tuneslap/models"

	"github.com/h2non/bimg"
)

func transform(input []byte, params models.ImageProcessingParams) ([]byte, error) {
	// img := bimg.NewImage(input)

	if params.ApplyFilters == "grayscale" {
		input, _ = bimg.NewImage(input).Colourspace(bimg.InterpretationBW)
	}
	// else if params.ApplyFilters == "blur" {
	// 		input, _ = bimg.NewImage(input).Blur(params.BlurRadius)
	// }

	return input, nil
}
