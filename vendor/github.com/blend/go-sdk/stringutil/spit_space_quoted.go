package stringutil

import "unicode"

// SplitSpaceQuoted splits a corpus on space but treats quoted strings
// i.e. within `"` as being atomic chunks.
func SplitSpaceQuoted(text string) (output []string) {
	if len(text) == 0 {
		return
	}
	var state int
	var word []rune
	var opened rune
	for _, r := range text {
		switch state {
		case 0: // word
			if isQuote(r) {
				opened = r
				word = append(word, r)
				state = 2
			} else if unicode.IsSpace(r) {
				if len(word) > 0 {
					output = append(output, string(word))
					word = nil
				}
				state = 1
			} else {
				word = append(word, r)
			}
		case 1: // we've seen a space
			if !unicode.IsSpace(r) {
				if isQuote(r) {
					opened = r
					state = 2
				} else {
					state = 0
				}
				word = append(word, r)
			}
		case 2: // quoted section
			if matchesQuote(opened, r) {
				state = 0
			}
			word = append(word, r)
		}
	}

	if len(word) > 0 {
		output = append(output, string(word))
	}
	return
}

func isQuote(r rune) bool {
	return r == '"' || r == '\'' || r == '“' || r == '”' || r == '`'
}

func matchesQuote(a, b rune) bool {
	if a == '“' && b == '”' {
		return true
	}
	if a == '”' && b == '“' {
		return true
	}
	return a == b
}
