package stringutil

import (
	"net/url"
	"strings"
)

// Slugify replaces whitespace with '-' and url escapes.
func Slugify(v string) string {
	v = strings.Replace(v, " ", "-", -1)
	v = strings.Replace(v, "\t", "-", -1)
	v = strings.Replace(v, "\n", "-", -1)
	return url.PathEscape(v)
}
