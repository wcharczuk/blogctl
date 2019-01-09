package async

import "context"

// QueueAction is an action handler for a queue.
type QueueAction func(context.Context, interface{}) error
