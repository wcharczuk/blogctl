package model

// Tags orders tag posts by tag.
type Tags []Tag

// Len implements sorter.
func (o Tags) Len() int {
	return len(o)
}

// Swap implements sorter.
func (o Tags) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

// Less implements sorter.
func (o Tags) Less(i, j int) bool {
	return o[i].Tag < o[j].Tag
}
