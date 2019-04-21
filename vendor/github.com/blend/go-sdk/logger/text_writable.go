package logger

import "io"

// TextWritable is an event that can be written.
type TextWritable interface {
	WriteText(TextFormatter, io.Writer)
}
