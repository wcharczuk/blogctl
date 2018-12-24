package async

import "github.com/blend/go-sdk/exception"

// Errors
var (
	ErrCannotStart exception.Class = "cannot start; already started"
	ErrCannotStop  exception.Class = "cannot stop; already stopped"
)
