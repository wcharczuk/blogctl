package model

// TagPosts are posts associated with tags.
type TagPosts struct {
	Tag   string
	Posts []Post
}

// TagPostsByTag orders tag posts by tag.
type TagPostsByTag []TagPosts

// Len implements sorter.
func (o TagPostsByTag) Len() int {
	return len(o)
}

// Swap implements sorter.
func (o TagPostsByTag) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

// Less implements sorter.
func (o TagPostsByTag) Less(i, j int) bool {
	return o[i].Tag < o[j].Tag
}
