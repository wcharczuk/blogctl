package async

import "context"

// WorkerFinalizer is an action handler for a queue.
type WorkerFinalizer func(context.Context, *Worker) error
