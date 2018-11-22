package config

import "github.com/wcharczuk/photoblog/pkg/constants"

// Layout details configuration options for the layout.
type Layout struct {
	// index is the path to the index.html
	Index string `json:"index" yaml:"index"`
	// home is the path to the list or home page template
	Home string `json:"home" yaml:"home"`
	// single is the path to the single post page template
	Single string `json:"single" yaml:"single"`
	// partials is a list of paths to include as optional partials
	Partials []string `json:"partials" yaml:"partials"`
	// statics are paths to copy as is to the output /static folder
	Statics []string `json:"statics" yaml:"statics"`
}

// GetIndex returns the index template or a default.
func (l Layout) GetIndex() string {
	if l.Index != "" {
		return l.Index
	}
	return constants.TemplateIndex
}

// GetHome returns the home template or a default.
func (l Layout) GetHome() string {
	if l.Home != "" {
		return l.Home
	}
	return constants.TemplateHome
}

// GetSingle returns the single post template or a default.
func (l Layout) GetSingle() string {
	if l.Single != "" {
		return l.Single
	}
	return constants.TemplateSingle
}

// GetPartials returns partial file paths or defaults.
func (l Layout) GetPartials() []string {
	if len(l.Partials) > 0 {
		return l.Partials
	}
	return []string{constants.DiscoveryPathPartials}
}

// GetStatics returns static file paths or defaults.
func (l Layout) GetStatics() []string {
	if len(l.Statics) > 0 {
		return l.Statics
	}
	return []string{constants.DiscoveryPathStatic}
}
