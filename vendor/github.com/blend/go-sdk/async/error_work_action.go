package async

import "context"

// ErrorWorkAction is an action for an error queue.
type ErrorWorkAction func(context.Context, error)
