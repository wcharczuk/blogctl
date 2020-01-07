package model

import "time"

// Meta is extra data for a post.
type Meta struct {
	Posted   time.Time         `json:"posted" yaml:"posted"`
	Title    string            `json:"title" yaml:"title"`
	Location string            `json:"location,omitempty" yaml:"location,omitempty"`
	Comments string            `json:"comments,omitempty" yaml:"comments,omitempty"`
	Tags     []string          `json:"tags,omitempty" yaml:"tags,omitempty"`
	Extra    map[string]string `json:"extra,omitempty" yaml:"extra,omitempty"`
}
