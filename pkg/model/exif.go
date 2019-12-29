package model

import "time"

// Exif are known values for a subset of the full image exif data.
type Exif struct {
	CaptureDate     time.Time `json:"captureDate" yaml:"captureDate"`
	CameraMake      string    `json:"cameraMake" yaml:"cameraMake"`
	CameraModel     string    `json:"cameraModel" yaml:"cameraModel"`
	LensModel       string    `json:"lensModel" yaml:"lensModel"`
	FNumber         string    `json:"fNumber" yaml:"fNumber"`
	ExposureTime    string    `json:"exposureTime" yaml:"exposureTime"`
	FocalLength     string    `json:"focalLength" yaml:"focalLength"`
	ISOSpeedRatings string    `json:"isoSpeedRatings" yaml:"isoSpeedRatings"`
}
