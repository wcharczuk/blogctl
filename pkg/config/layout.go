package config

import "github.com/wcharczuk/photoblog/pkg/constants"

// Layout details configuration options for the layout.
type Layout struct {
	Post     string `json:"post" yaml:"post"`
	Pages    string `json:"pages" yaml:"pages"`
	Partials string `json:"partials" yaml:"partials"`
	Static   string `json:"static" yaml:"static"`
}

// PostOrDefault returns the single post template or a default.
func (l Layout) PostOrDefault() string {
	if l.Post != "" {
		return l.Post
	}
	return constants.TemplatePost
}

// PagesOrDefault returns page file paths or defaults.
func (l Layout) PagesOrDefault() string {
	if l.Pages != "" {
		return l.Pages
	}
	return constants.DiscoveryPathPages
}

// PartialsOrDefault returns partial file paths or defaults.
func (l Layout) PartialsOrDefault() string {
	if l.Partials != "" {
		return l.Partials
	}
	return constants.DiscoveryPathPartials
}

// StaticOrDefault returns static file paths or defaults.
func (l Layout) StaticOrDefault() string {
	if l.Static != "" {
		return l.Static
	}
	return constants.DiscoveryPathStatic
}
