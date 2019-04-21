package stringutil

import (
	"bytes"
	"unicode"
)

var nonTitleWords = map[string]bool{
	"and":     true,
	"the":     true,
	"a":       true,
	"an":      true,
	"but":     true,
	"or":      true,
	"on":      true,
	"in":      true,
	"with":    true,
	"for":     true,
	"either":  true,
	"neither": true,
	"nor":     true,
}

// Title returns a string in title case.
func Title(corpus string) string {
	output := bytes.NewBuffer(nil)
	runes := []rune(corpus)

	haveSeenLetter := false
	var r rune
	for x := 0; x < len(runes); x++ {
		r = runes[x]

		if unicode.IsLetter(r) {
			if !haveSeenLetter {
				output.WriteRune(unicode.ToUpper(r))
				haveSeenLetter = true
			} else {
				output.WriteRune(unicode.ToLower(r))
			}
		} else {
			output.WriteRune(r)
			haveSeenLetter = false
		}
	}
	return output.String()
}
