package model

import (
	"github.com/wcharczuk/blogctl/pkg/config"
)

// ViewModel is the type passed to view rendering.
type ViewModel struct {
	Title  string
	Extra  map[string]interface{}
	Config config.Config
	Posts  []*Post
	Tags   []Tag

	Post Post
	Tag  Tag
}

// TitleOrDefault returns the title.
func (vm ViewModel) TitleOrDefault() string {
	if vm.Title != "" {
		return vm.Title
	}
	if !vm.Post.IsZero() {
		return vm.Post.TitleOrDefault()
	}
	return vm.Config.TitleOrDefault()
}
