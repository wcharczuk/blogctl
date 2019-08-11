package engine

import "fmt"

// Logger is a responder for logging.
type Logger interface {
	Info(args ...interface{})
	Debug(args ...interface{})
	Warning(error) error
	Error(error) error
}

// MaybeInfof writes an info message if the logger is set.
func MaybeInfof(log Logger, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Info(fmt.Sprintf(format, args...))
}

// MaybeDebugf writes a debug message if the logger is set.
func MaybeDebugf(log Logger, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Debug(fmt.Sprintf(format, args...))
}
