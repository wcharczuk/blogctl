package logger

import (
	"net/http"
	"strconv"
)

// AnsiColor represents an ansi color code fragment.
type AnsiColor string

func (acc AnsiColor) escaped() string {
	return "\033[" + string(acc)
}

// Apply returns a string with the color code applied.
func (acc AnsiColor) Apply(text string) string {
	return acc.escaped() + text + ColorReset.escaped()
}

const (
	// RuneSpace is a single rune representing a space.
	RuneSpace rune = ' '

	// RuneNewline is a single rune representing a newline.
	RuneNewline rune = '\n'

	// ColorBlack is the posix escape code fragment for black.
	ColorBlack AnsiColor = "30m"

	// ColorRed is the posix escape code fragment for red.
	ColorRed AnsiColor = "31m"

	// ColorGreen is the posix escape code fragment for green.
	ColorGreen AnsiColor = "32m"

	// ColorYellow is the posix escape code fragment for yellow.
	ColorYellow AnsiColor = "33m"

	// ColorBlue is the posix escape code fragment for blue.
	ColorBlue AnsiColor = "34m"

	// ColorPurple is the posix escape code fragement for magenta (purple)
	ColorPurple AnsiColor = "35m"

	// ColorCyan is the posix escape code fragement for cyan.
	ColorCyan AnsiColor = "36m"

	// ColorWhite is the posix escape code fragment for white.
	ColorWhite AnsiColor = "37m"

	// ColorLightBlack is the posix escape code fragment for black.
	ColorLightBlack AnsiColor = "90m"

	// ColorLightRed is the posix escape code fragment for red.
	ColorLightRed AnsiColor = "91m"

	// ColorLightGreen is the posix escape code fragment for green.
	ColorLightGreen AnsiColor = "92m"

	// ColorLightYellow is the posix escape code fragment for yellow.
	ColorLightYellow AnsiColor = "93m"

	// ColorLightBlue is the posix escape code fragment for blue.
	ColorLightBlue AnsiColor = "94m"

	// ColorLightPurple is the posix escape code fragement for magenta (purple)
	ColorLightPurple AnsiColor = "95m"

	// ColorLightCyan is the posix escape code fragement for cyan.
	ColorLightCyan AnsiColor = "96m"

	// ColorLightWhite is the posix escape code fragment for white.
	ColorLightWhite AnsiColor = "97m"

	// ColorGray is an alias to ColorLightBlack to preserve backwards compatibility.
	ColorGray AnsiColor = ColorLightBlack

	// ColorReset is the posix escape code fragment to reset all formatting.
	ColorReset AnsiColor = "0m"
)

var (
	// DefaultFlagTextColors is the default color for each known flag.
	DefaultFlagTextColors = map[Flag]AnsiColor{
		Info:    ColorLightWhite,
		Silly:   ColorLightBlack,
		Debug:   ColorLightYellow,
		Warning: ColorLightYellow,
		Error:   ColorRed,
		Fatal:   ColorRed,
	}

	// DefaultFlagTextColor is the default flag color.
	DefaultFlagTextColor = ColorLightWhite
)

// GetFlagTextColor returns the color for a flag.
func GetFlagTextColor(flag Flag) AnsiColor {
	if color, hasColor := DefaultFlagTextColors[flag]; hasColor {
		return color
	}
	return DefaultFlagTextColor
}

// ColorizeByStatusCode returns a value colored by an http status code.
func ColorizeByStatusCode(statusCode int, value string) string {
	if statusCode >= http.StatusOK && statusCode < 300 { //the http 2xx range is ok
		return ColorGreen.Apply(value)
	} else if statusCode == http.StatusInternalServerError {
		return ColorRed.Apply(value)
	}
	return ColorYellow.Apply(value)
}

// ColorizeStatusCode colorizes a status code.
func ColorizeStatusCode(statusCode int) string {
	return ColorizeByStatusCode(statusCode, strconv.Itoa(statusCode))
}
