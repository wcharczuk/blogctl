package engine

import (
	"github.com/wcharczuk/photoblog/pkg/config"
	"github.com/wcharczuk/photoblog/pkg/model"
)

// HomeViewModel is the type passed to view rendering.
type HomeViewModel struct {
	Config config.Config
	Posts  []model.Post
}

// SingleViewModel is the type passed to view rendering.
type SingleViewModel struct {
	Config   config.Config
	Post     model.Post
	Previous model.Post
	Next     model.Post
}
