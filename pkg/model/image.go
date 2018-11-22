package model

// Image represents a posted image.
type Image struct {
	Width  int  `json:"width" yaml:"width"`
	Height int  `json:"height" yaml:"height"`
	Exif   Exif `json:"exif" yaml:"exif"`
}

// IsZero returns if the image has been processed or not.
func (i Image) IsZero() bool {
	return i.Width == 0 || i.Height == 0
}
