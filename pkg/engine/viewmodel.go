package engine

import (
	"github.com/wcharczuk/blogctl/pkg/config"
	"github.com/wcharczuk/blogctl/pkg/model"
)

// TitleOrDefaultProvider is a type that provides a title.
type TitleOrDefaultProvider interface {
	TitleOrDefault() string
}

// ViewModel is the type passed to view rendering.
type ViewModel struct {
	Config    config.Config
	Posts     []model.Post
	Tags      []model.TagPosts
	Post      model.Post
	PostIndex int
	Tag       model.TagPosts
	TagIndex  int
}

// TitleOrDefault returns the title.
func (vm ViewModel) TitleOrDefault() string {
	if !vm.Post.IsZero() {
		return vm.Post.TitleOrDefault()
	}
	return vm.Config.TitleOrDefault()
}
