/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stringutil

import "strings"

// SplitCSV splits a corpus by the `,`, dropping leading or trailing whitespace unless quoted.
func SplitCSV(text string) (output []string) {
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
				state = 1
			} else if isCSVDelim(r) {
				output = append(output, strings.TrimSpace(string(word)))
				word = nil
			} else {
				word = append(word, r)
			}
		case 1: // we're in a quoted section
			if matchesQuote(opened, r) {
				state = 0
				continue
			}
			word = append(word, r)
		}
	}

	if len(word) > 0 {
		output = append(output, strings.TrimSpace(string(word)))
	}
	return
}

func isCSVDelim(r rune) bool {
	return r == rune(',')
}
