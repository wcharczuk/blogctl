package async

import "context"

// ErrorAction is an action handler for a queue.
type ErrorAction func(context.Context, error) error
