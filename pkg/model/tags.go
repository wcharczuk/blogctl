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

// TableRows returns the table rows for the given slice of tags.
func (o Tags) TableRows() []TagTableRow {
	output := make([]TagTableRow, len(o))
	for index := range o {
		output[index] = o[index].TableRow()
	}
	return output
}
