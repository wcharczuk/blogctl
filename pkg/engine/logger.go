package engine

// Logger is a responder for logging.
type Logger interface {
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
}
