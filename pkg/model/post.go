package model

import (
	"fmt"
	"path/filepath"

	"github.com/wcharczuk/photoblog/pkg/constants"
	"github.com/wcharczuk/photoblog/pkg/stringutil"
)

// Post is a single post item.
type Post struct {
	Original string `json:"original" yaml:"original"`
	Image    Image  `json:"image" yaml:"image"`
	Meta     Meta   `json:"meta" yaml:"meta"`
	Previous *Post  `json:"-" yaml:"-"`
	Next     *Post  `json:"-" yaml:"-"`
}

// HasPrevious returns if there is a previous post.
func (p Post) HasPrevious() bool {
	return p.Previous != nil && !p.Previous.IsZero()
}

// HasNext returns if there is a next post.
func (p Post) HasNext() bool {
	return p.Next != nil && !p.Next.IsZero()
}

// IsZero returns if the post is set.
func (p Post) IsZero() bool {
	return p.Image.IsZero()
}

// TitleOrDefault returns the title for the post.
// It is coalesced from the meta.Title and the filename.
func (p Post) TitleOrDefault() string {
	if p.Meta.Title != "" {
		return p.Meta.Title
	}
	return filepath.Base(p.Original)
}

// Slug returns the fully qualified identifier for the post.
// It is in the form /Year/Month/Day/Slug
func (p Post) Slug() string {
	titleSlug := stringutil.Slugify(p.TitleOrDefault())
	return fmt.Sprintf("%d/%d/%d/%s", p.Meta.Posted.Year(), p.Meta.Posted.Month(), p.Meta.Posted.Day(), titleSlug)
}

// SlugIndex returns the fully qualified identifier for the post with the trailing index.html.
// It is in the form /Year/Month/Day/Slug/index.html
func (p Post) SlugIndex() string {
	titleSlug := stringutil.Slugify(p.TitleOrDefault())
	return fmt.Sprintf("%d/%d/%d/%s/%s", p.Meta.Posted.Year(), p.Meta.Posted.Month(), p.Meta.Posted.Day(), titleSlug, constants.FileIndex)
}

// SourceOriginal returns the fully qualified image source path.
func (p Post) SourceOriginal() string {
	return filepath.Join(p.Slug(), constants.ImageOriginal)
}

// SourceForSize returns the image source for a given image size in pixels.
func (p Post) SourceForSize(size int) string {
	return filepath.Join(p.Slug(), fmt.Sprintf(constants.ImageSizeFormat, size))
}

// SourceLarge returns the fully qualified image source path.
func (p Post) SourceLarge() string {
	return filepath.Join(p.Slug(), fmt.Sprintf(constants.ImageSizeFormat, constants.SizeLarge))
}

// SourceMedium returns the fully qualified image source path.
func (p Post) SourceMedium() string {
	return filepath.Join(p.Slug(), fmt.Sprintf(constants.ImageSizeFormat, constants.SizeMedium))
}

// SourceSmall returns the fully qualified image source path.
func (p Post) SourceSmall() string {
	return filepath.Join(p.Slug(), fmt.Sprintf(constants.ImageSizeFormat, constants.SizeSmall))
}
