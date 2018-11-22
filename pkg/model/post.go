package model

import (
	"fmt"
	"path/filepath"
)

// Post is a single post item.
type Post struct {
	Original string `json:"original" yaml:"original"`
	File     string `json:"file" yaml:"file"`
	Image    Image  `json:"image" yaml:"image"`
	Meta     Meta   `json:"meta" yaml:"meta"`
}

// IsZero returns if the post is set.
func (p Post) IsZero() bool {
	return p.File == "" || p.Image.IsZero()
}

// TitleOrDefault returns the title for the post.
// It is coalesced from the meta.Title and the filename.
func (p Post) TitleOrDefault() string {
	if p.Meta.Title != "" {
		return p.Meta.Title
	}
	return filepath.Base(p.File)
}

// Slug returns the fully qualified identifier for the post.
// It is in the form /Year/Month/Day/Slug
func (p Post) Slug() string {
	return fmt.Sprintf("%d/%d/%d/%s", p.Meta.Posted.Year(), p.Meta.Posted.Month(), p.Meta.Posted.Day(), p.TitleOrDefault())
}

// Source returns the fully qualified image source path.
func (p Post) Source() string {
	return filepath.Join(p.Slug(), p.File)
}
