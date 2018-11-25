package model

import (
	"fmt"
	"path/filepath"

	"github.com/wcharczuk/photoblog/pkg/constants"
	"github.com/wcharczuk/photoblog/pkg/stringutil"
)

// CreateSlug creates a slug for a post.
func CreateSlug(p Post) string {
	titleSlug := stringutil.Slugify(p.TitleOrDefault())
	return fmt.Sprintf("%d/%d/%d/%s", p.Meta.Posted.Year(), p.Meta.Posted.Month(), p.Meta.Posted.Day(), titleSlug)
}

// Post is a single post item.
type Post struct {
	OriginalPath string `json:"originalPath" yaml:"originalPath"`
	OutputPath   string `json:"outputPath" yaml:"outputPath"`
	Slug         string `json:"slug" yaml:"slug"`
	Image        Image  `json:"image" yaml:"image"`
	Meta         Meta   `json:"meta" yaml:"meta"`

	Index    int   `json:"index" yaml:"index"`
	Previous *Post `json:"-" yaml:"-"`
	Next     *Post `json:"-" yaml:"-"`
}

// IsZero returns if the post is set.
func (p Post) IsZero() bool {
	return p.OriginalPath == ""
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
// It is coalesced from the meta.Title and the filename.
func (p Post) TitleOrDefault() string {
	if p.Meta.Title != "" {
		return p.Meta.Title
	}
	return filepath.Base(p.OriginalPath)
}

// SlugIndex is a helper that returns the fully qualified path for the post's index.html.
// It is in the form /Year/Month/Day/Slug/index.html
func (p Post) SlugIndex() string {
	return filepath.Join(p.Slug, constants.FileIndex)
}

// ImageSourceOriginal returns the fully qualified image source path.
func (p Post) ImageSourceOriginal() string {
	return filepath.Join(p.Slug, constants.ImageOriginal)
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