package model

// Meta is extra data for a post.
type Meta struct {
	Title    string            `json:"title" yaml:"title"`
	Location string            `json:"location" yaml:"location"`
	Comments string            `json:"comments" yaml:"comments"`
	Tags     []string          `json:"tags" yaml:"tags"`
	Extra    map[string]string `json:"extra" yaml:"extra"`
}
