package image

import (
	"fmt"
	"tuneslap/models"

	"github.com/h2non/bimg"
)

// DefaultBlurSigma is the default blur intensity when blur filter is applied
const DefaultBlurSigma = 5.0

// transform applies filters and crop operations to an image
func transform(input []byte, params models.ImageProcessingParams) ([]byte, error) {
	var err error

	// Apply crop if specified (all values must be positive)
	if isCropValid(params.Crop) {
		input, err = applyCrop(input, params.Crop)
		if err != nil {
			return nil, fmt.Errorf("crop error: %w", err)
		}
	}

	// Apply filters
	if params.ApplyFilters == "grayscale" {
		input, err = applyGrayscale(input)
		if err != nil {
			return nil, fmt.Errorf("grayscale error: %w", err)
		}
	} else if params.ApplyFilters == "blur" {
		input, err = applyBlur(input, DefaultBlurSigma)
		if err != nil {
			return nil, fmt.Errorf("blur error: %w", err)
		}
	}

	return input, nil
}

// isCropValid checks if crop parameters are valid (all values must be positive)
func isCropValid(crop [4]int) bool {
	// crop: [x, y, width, height]
	// width and height must be positive
	return crop[2] > 0 && crop[3] > 0
}

// applyCrop extracts a region from the image
func applyCrop(input []byte, crop [4]int) ([]byte, error) {
	// crop: [x, y, width, height]
	return bimg.NewImage(input).Extract(crop[0], crop[1], crop[2], crop[3])
}

// applyGrayscale converts the image to grayscale
func applyGrayscale(input []byte) ([]byte, error) {
	return bimg.NewImage(input).Colourspace(bimg.InterpretationBW)
}

// applyBlur applies a Gaussian blur to the image
func applyBlur(input []byte, sigma float64) ([]byte, error) {
	// bimg's GaussianBlur requires min/max sigma values
	// Using Process with blur option instead for better control
	options := bimg.Options{
		GaussianBlur: bimg.GaussianBlur{
			Sigma:   sigma,
			MinAmpl: 0.2,
		},
	}
	return bimg.NewImage(input).Process(options)
}
