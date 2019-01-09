package model

// Posts is a list of posts.
type Posts []Post

// First returns the first post in the list.
// It returns an empty post if the list is empty.
func (p Posts) First() (output Post) {
	if len(p) > 0 {
		output = p[0]
	}
	return
}

// Len implements sorter.
func (p Posts) Len() int {
	return len(p)
}

// Swap implements sorter.
func (p Posts) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Less implements sorter.
func (p Posts) Less(i, j int) bool {
	return p[i].Meta.Posted.After(p[j].Meta.Posted)
}
