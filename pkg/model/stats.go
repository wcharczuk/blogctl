package model

import (
	"strconv"
	"time"
)

// Stats are stats about the blog.
type Stats struct {
	NumPosts      int
	NumTags       int
	NumImagePosts int
	NumTextPosts  int

	Earliest time.Time
	Latest   time.Time
}

// TableData returns the stats as ansi table data.
func (s Stats) TableData() (columns []string, rows [][]string) {
	columns = []string{
		"number of posts",
		"number of tags",
		"number of image posts",
		"number of text posts",
		"earliest post",
		"latest post",
	}
	rows = [][]string{{
		strconv.Itoa(s.NumPosts),
		strconv.Itoa(s.NumTags),
		strconv.Itoa(s.NumImagePosts),
		strconv.Itoa(s.NumTextPosts),
		s.Earliest.Format(time.RFC3339),
		s.Latest.Format(time.RFC3339),
	},
	}
	return
}
