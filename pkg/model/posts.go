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
