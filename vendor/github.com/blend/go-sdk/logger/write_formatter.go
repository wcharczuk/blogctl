package logger

import (
	"context"
	"io"
)

// WriteFormatter is a formatter for writing events to output writers.
type WriteFormatter interface {
	WriteFormat(context.Context, io.Writer, Event) error
}
