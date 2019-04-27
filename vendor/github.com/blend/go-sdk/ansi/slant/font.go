package slant

// Font represents a single font.
type Font struct {
	// Height of one char
	Height int
	// Baseline is the height of letters not including descenders.
	Baseline int
	// Width of the widest char
	Width int
	// Hardblank symbol is the non-smushable space character.
	Hardblank string
	// A string for each line of the char
	Letters [][]string
}
