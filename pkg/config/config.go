package config

import "github.com/wcharczuk/photoblog/pkg/constants"

// These are set by ldflags.
var (
	Version = ""
	GitRef  = ""
)

// Config is the blog config
type Config struct {
	Title  string `json:"title" yaml:"title"`
	Images string `json:"images" yaml:"images"`
	Output string `json:"output" yaml:"output"`
	Layout Layout `json:"layout" yaml:"layout"`
	Extra  Extra  `json:"extra" yaml:"extra"`
}

// TitleOrDefault returns the title or a default.
func (c Config) TitleOrDefault() string {
	if c.Title != "" {
		return c.Title
	}
	return "Unset"
}

// ImagesOrDefault returns the images path.
func (c Config) ImagesOrDefault() string {
	if c.Images != "" {
		return c.Images
	}
	return constants.ImagesPath
}

// OutputOrDefault returns the output path.
func (c Config) OutputOrDefault() string {
	if c.Output != "" {
		return c.Output
	}
	return constants.OutputPath
}
