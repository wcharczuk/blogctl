package stringutil

import (
	"strings"
)

// StripQuotes strips leading or trailing quotes.
func StripQuotes(v string) string {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "\"")
	v = strings.TrimSuffix(v, "\"")
	return v
}
