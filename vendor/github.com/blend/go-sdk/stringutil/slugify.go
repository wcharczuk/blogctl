package stringutil

import (
	"unicode"
)

// Slugify replaces non-letter or digit runes with '-'.
func Slugify(v string) string {
	runes := []rune(v)
	var c rune
	for index := range runes {
		c = runes[index]
		if !(unicode.IsLetter(c) || unicode.IsDigit(c)) {
			runes[index] = '-'
		}
	}
	return string(runes)
}
