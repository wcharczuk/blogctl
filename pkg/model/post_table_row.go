package model

import "time"

// PostTableRow is a row for table output.
type PostTableRow struct {
	Title    string
	Location string
	Posted   time.Time
	Tags     string
	Slug     string
	PostType string
}
