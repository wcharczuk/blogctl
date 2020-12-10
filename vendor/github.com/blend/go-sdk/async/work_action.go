package async

import "context"

// WorkAction is an action handler for a queue.
type WorkAction func(context.Context, interface{}) error
