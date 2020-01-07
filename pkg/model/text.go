package model

// Text is a text post.
type Text struct {
	SourcePath string `json:"sourcePath" yaml:"sourcePath"`
	Template   string `json:"template,omitempty" yaml:"template,omitempty"`
	Output     string `json:"output,omitempty" yaml:"output,omitempty"`
}

// IsZero returns if the post is set or not.
func (t Text) IsZero() bool {
	return t.Template == ""
}
