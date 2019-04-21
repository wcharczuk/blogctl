package async

import "context"

// ErrorWorkerFinalizer is an action handler for a queue.
type ErrorWorkerFinalizer func(context.Context, *ErrorWorker)
