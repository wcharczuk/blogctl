package engine

import (
	"github.com/wcharczuk/photoblog/pkg/config"
	"github.com/wcharczuk/photoblog/pkg/model"
)

// TitleOrDefaultProvider is a type that provides a title.
type TitleOrDefaultProvider interface {
	TitleOrDefault() string
}

// ViewModel is the type passed to view rendering.
type ViewModel struct {
	Config config.Config
	Posts  []model.Post
	Post   model.Post
}

// TitleOrDefault returns the title.
func (vm ViewModel) TitleOrDefault() string {
	if !vm.Post.IsZero() {
		return vm.Post.TitleOrDefault()
	}
	return vm.Config.TitleOrDefault()
}
