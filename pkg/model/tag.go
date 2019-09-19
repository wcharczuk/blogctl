package model

// Tag are posts associated with tags.
type Tag struct {
	Tag   string
	Posts []*Post
}

// TableRow returns the ansi table row form of the tag.
func (t Tag) TableRow() TagTableRow {
	return TagTableRow{
		Tag:   t.Tag,
		Posts: len(t.Posts),
	}
}
