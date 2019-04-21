package logger

import (
	"net/http"
	"strconv"

	"github.com/blend/go-sdk/ansi"
)

var (
	// DefaultFlagTextColors is the default color for each known flag.
	DefaultFlagTextColors = map[string]ansi.Color{
		Info:    ansi.ColorLightWhite,
		Debug:   ansi.ColorLightYellow,
		Warning: ansi.ColorLightYellow,
		Error:   ansi.ColorRed,
		Fatal:   ansi.ColorRed,
	}

	// DefaultFlagTextColor is the default flag color.
	DefaultFlagTextColor = ansi.ColorLightWhite
)

// FlagTextColor returns the color for a flag.
func FlagTextColor(flag string) ansi.Color {
	if color, hasColor := DefaultFlagTextColors[flag]; hasColor {
		return color
	}
	return DefaultFlagTextColor
}

// ColorizeByStatusCode returns a value colored by an http status code.
func ColorizeByStatusCode(statusCode int, value string) string {
	if statusCode >= http.StatusOK && statusCode < 300 { //the http 2xx range is ok
		return ansi.ColorGreen.Apply(value)
	} else if statusCode == http.StatusInternalServerError {
		return ansi.ColorRed.Apply(value)
	}
	return ansi.ColorYellow.Apply(value)
}

// ColorizeByStatusCodeWithFormatter returns a value colored by an http status code with a given formatter.
func ColorizeByStatusCodeWithFormatter(tf TextFormatter, statusCode int, value string) string {
	if statusCode >= http.StatusOK && statusCode < 300 { //the http 2xx range is ok
		return tf.Colorize(value, ansi.ColorGreen)
	} else if statusCode == http.StatusInternalServerError {
		return tf.Colorize(value, ansi.ColorRed)
	}
	return tf.Colorize(value, ansi.ColorYellow)
}

// ColorizeStatusCode colorizes a status code.
func ColorizeStatusCode(statusCode int) string {
	return ColorizeByStatusCode(statusCode, strconv.Itoa(statusCode))
}

// ColorizeStatusCodeWithFormatter colorizes a status code with a given formatter.
func ColorizeStatusCodeWithFormatter(tf TextFormatter, statusCode int) string {
	return ColorizeByStatusCodeWithFormatter(tf, statusCode, strconv.Itoa(statusCode))
}
