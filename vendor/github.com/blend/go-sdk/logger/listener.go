package logger

import "context"

// Listener is a function that can be triggered by events.
type Listener func(context.Context, Event)
