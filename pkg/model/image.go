package model

import "image"

// Image represents a posted image.
type Image struct {
	Width  int               `json:"width" yaml:"width"`
	Height int               `json:"height" yaml:"height"`
	Exif   Exif              `json:"exif" yaml:"exif"`
	Paths  map[string]string `json:"paths,omitempty" yaml:"paths,omitempty"`
}

// IsZero returns if the image has been processed or not.
func (i Image) IsZero() bool {
	return i.Width == 0 || i.Height == 0
}

// LongDimension returns the longest dimension of the image.
func (i Image) LongDimension() int {
	if i.Width > i.Height {
		return i.Width
	}
	return i.Height
}

// Ratio returns the ratio of the width to the height.
// As an example, for 3:2 images, the ratio is 1.5.
func (i Image) Ratio() float64 {
	return float64(i.Width) / float64(i.Height)
}

// Scaled returns the image dimensions scaled to a given long dimension.
func (i Image) Scaled(longDimension int) image.Rectangle {
	scaled := int(float64(longDimension) / i.Ratio())
	if i.Width > i.Height {
		return image.Rectangle{
			Max: image.Point{X: longDimension, Y: scaled},
		}
	}
	return image.Rectangle{
		Max: image.Point{X: scaled, Y: longDimension},
	}
}
