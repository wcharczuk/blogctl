package async

import "context"

// AutoflushAction is an action called by an autoflush buffer.
type AutoflushAction func(context.Context, []interface{}) error
