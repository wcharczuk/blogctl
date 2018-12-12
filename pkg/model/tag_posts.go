package model

// TagPosts are posts associated with tags.
type TagPosts struct {
	Tag   string
	Posts []Post
}
