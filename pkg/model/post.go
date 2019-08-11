package model

import (
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/wcharczuk/blogctl/pkg/constants"
)

// Post is a single post item.
type Post struct {
	Path string `json:"outputPath,omitempty" yaml:"outputPath,omitempty"`
	Slug string `json:"slug,omitempty" yaml:"slug,omitempty"`
	Meta Meta   `json:"meta" yaml:"meta"`

	TextPath   string         `json:"textPath,omitempty" yaml:"textPath,omitempty"`
	ImagePath  string         `json:"imagePath,omitempty" yaml:"imagePath,omitempty"`
	Image      Image          `json:"image,omitempty" yaml:"image,omitempty"`
	ImageSizes map[int]string `json:"imageSizes" yaml:"imageSizes,omitempty"`

	Template *template.Template `json:"-" yaml:"-"`
	Previous *Post              `json:"-" yaml:"-"`
	Next     *Post              `json:"-" yaml:"-"`
}

// IsZero returns if the post is set.
func (p Post) IsZero() bool {
	return p.ImagePath == "" && p.TextPath == ""
}

// IsImage returns if the post is an image post.
func (p Post) IsImage() bool {
	return p.ImagePath != ""
}

// IsText returns if the post is an text post.
func (p Post) IsText() bool {
	return p.TextPath != ""
}

// HasPrevious returns if there is a previous post.
func (p Post) HasPrevious() bool {
	return p.Previous != nil && !p.Previous.IsZero()
}

// HasNext returns if there is a next post.
func (p Post) HasNext() bool {
	return p.Next != nil && !p.Next.IsZero()
}

// TitleOrDefault returns the title for the post.
func (p Post) TitleOrDefault() string {
	return p.Meta.Title
}

// SlugIndex is a helper that returns the fully qualified path for the post's index.html.
// It is in the form /Year/Month/Day/Slug/index.html
func (p Post) SlugIndex() string {
	return filepath.Join(p.Slug, constants.FileIndex)
}

// ImageSourceOriginal returns the fully qualified image source path.
func (p Post) ImageSourceOriginal() string {
	return filepath.Join(p.Slug, constants.FileImageOriginal)
}

// ImageSourceForSize returns the image source for a given image size in pixels.
func (p Post) ImageSourceForSize(size int) string {
	return filepath.Join(p.Slug, fmt.Sprintf(constants.ImageSizeFormat, size))
}

// ImageSourceLarge returns the fully qualified image source path.
func (p Post) ImageSourceLarge() string {
	return p.ImageSourceForSize(constants.SizeLarge)
}

// ImageSourceMedium returns the fully qualified image source path.
func (p Post) ImageSourceMedium() string {
	return p.ImageSourceForSize(constants.SizeMedium)
}

// ImageSourceSmall returns the fully qualified image source path.
func (p Post) ImageSourceSmall() string {
	return p.ImageSourceForSize(constants.SizeSmall)
}
