package stringutil

import (
	"math/rand"
	"time"
)

var (
	provider = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// Random returns a random selection of runes from the set.
func Random(runeset []rune, length int) string {
	runes := make([]rune, length)
	for index := range runes {
		runes[index] = runeset[provider.Intn(len(runeset))]
	}
	return string(runes)
}

// CombineRunsets combines given runsets into a single runset.
func CombineRunsets(runesets ...[]rune) []rune {
	output := []rune{}
	for _, set := range runesets {
		output = append(output, set...)
	}
	return output
}
