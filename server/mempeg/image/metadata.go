package image

import "github.com/h2non/bimg"

type ImageMetadata struct {
	Dimensions  [2]int
	ContentType string
	FileSize    int64
}

func getMetadata(input []byte) (ImageMetadata, error) {
	metadata, err := bimg.Metadata(input)
	if err != nil {
		return ImageMetadata{}, err
	}

	// convert to ImageMetadata
	imageMetadata := ImageMetadata{
		Dimensions:  [2]int{metadata.Size.Width, metadata.Size.Height},
		ContentType: metadata.Type,
		FileSize:    int64(len(input)),
	}

	return imageMetadata, nil
}
