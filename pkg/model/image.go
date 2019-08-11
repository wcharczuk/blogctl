package model

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
