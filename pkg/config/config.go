package config

import "github.com/wcharczuk/photoblog/pkg/constants"

// These are set by ldflags.
var (
	Version = ""
	GitRef  = ""
)

// Config is the blog config
type Config struct {
	Name   string `json:"name" yaml:"name"`
	Images string `json:"images" yaml:"images"`
	Output string `json:"output" yaml:"output"`
	Layout Layout `json:"layout" yaml:"layout"`
	Extra  Extra  `json:"extra" yaml:"extra"`
}

// GetImages returns the images path.
func (c Config) GetImages() string {
	if c.Images != "" {
		return c.Images
	}
	return constants.ImagesPath
}

// GetOutput returns the output path.
func (c Config) GetOutput() string {
	if c.Output != "" {
		return c.Output
	}
	return constants.OutputPath
}
