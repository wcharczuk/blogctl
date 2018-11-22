package model

import (
	"fmt"
	"path/filepath"
	"time"
)

// Post is a single post item.
type Post struct {
	ImagePath string    `json:"path" yaml:"path"`
	Posted    time.Time `json:"posted" yaml:"posted"`
	Image     Image     `json:"image" yaml:"image"`
	Meta      Meta      `json:"meta" yaml:"meta"`
}

// Title returns the title for the post.
// It is coalesced from the meta.Title and the filename.
func (p Post) Title() string {
	if p.Meta.Title != "" {
		return p.Meta.Title
	}
	return filepath.Base(p.ImagePath)
}

// Slug returns the fully qualified identifier for the post.
// It is in the form /Year/Month/Day/Slug
func (p Post) Slug() string {
	return fmt.Sprintf("%d/%d/%d/%s", p.Posted.Year(), p.Posted.Month(), p.Posted.Day(), p.Title())
}
