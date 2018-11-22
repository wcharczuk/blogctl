package config

import "github.com/wcharczuk/photoblog/pkg/constants"

// Layout details configuration options for the layout.
type Layout struct {
	// Post is the path to the single post page template
	Post string `json:"post" yaml:"post"`
	// Pages are extra pages you can add.
	Pages []string `json:"pages" yaml:"pages"`
	// partials is a list of paths to include as optional partials
	Partials []string `json:"partials" yaml:"partials"`
	// statics are paths to copy as is to the output /static folder
	Statics []string `json:"statics" yaml:"statics"`
}

// PostOrDefault returns the single post template or a default.
func (l Layout) PostOrDefault() string {
	if l.Post != "" {
		return l.Post
	}
	return constants.TemplatePost
}

// PagesOrDefault returns page file paths or defaults.
func (l Layout) PagesOrDefault() []string {
	if len(l.Pages) > 0 {
		return l.Pages
	}
	return []string{constants.OutputFileIndex}
}

// PartialsOrDefault returns partial file paths or defaults.
func (l Layout) PartialsOrDefault() []string {
	if len(l.Partials) > 0 {
		return l.Partials
	}
	return nil
}

// StaticsOrDefault returns static file paths or defaults.
func (l Layout) StaticsOrDefault() []string {
	if len(l.Statics) > 0 {
		return l.Statics
	}
	return []string{constants.DiscoveryPathStatic}
}
