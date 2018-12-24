package stringutil

const (
	// Empty is the empty string
	Empty string = ""

	// RuneSpace is a single rune representing a space.
	RuneSpace rune = ' '

	// RuneNewline is a single rune representing a newline.
	RuneNewline rune = '\n'

	// LowerA is the ascii int value for 'a'
	LowerA uint = uint('a')
	// LowerZ is the ascii int value for 'z'
	LowerZ uint = uint('z')
)

var (
	// LowerDiff is the difference between lower Z and lower A
	LowerDiff = (LowerZ - LowerA)

	// LowerLetters is a runset of lowercase letters.
	LowerLetters = []rune("abcdefghijklmnopqrstuvwxyz")

	// UpperLetters is a runset of uppercase letters.
	UpperLetters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	// Letters is a runset of both lower and uppercase letters.
	Letters = append(LowerLetters, UpperLetters...)

	// Numbers is a runset of numeric characters.
	Numbers = []rune("0123456789")

	// LettersAndNumbers is a runset of letters and numeric characters.
	LettersAndNumbers = append(Letters, Numbers...)

	// Symbols is a runset of symbol characters.
	Symbols = []rune(`!@#$%^&*()_+-=[]{}\|:;`)

	// LettersNumbersAndSymbols is a runset of letters, numbers and symbols.
	LettersNumbersAndSymbols = append(LettersAndNumbers, Symbols...)
)
