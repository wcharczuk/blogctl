package stringutil

import "strings"

// CSV produces a csv from a given set of values.
func CSV(values []string) string {
	return strings.Join(values, ",")
}
