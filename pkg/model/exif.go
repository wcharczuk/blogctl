package model

// Exif are known values for a subset of the full image exif data.
type Exif struct {
	CameraMake  string
	CameraModel string

	LensMake  string
	LensModel string

	FNumber      string
	ExposureTime string
	FocalLength  string

	ISOSpeedRatings string
}
