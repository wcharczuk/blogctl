package model

import "time"

// PostTableRow is a row for table output.
type PostTableRow struct {
	Title    string
	Location string
	Posted   time.Time
	Slug     string
	PostType string
}
