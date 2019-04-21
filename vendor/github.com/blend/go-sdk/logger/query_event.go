package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/blend/go-sdk/stringutil"
	"github.com/blend/go-sdk/timeutil"

	"github.com/blend/go-sdk/ansi"
)

// these are compile time assertions
var (
	_ Event          = (*QueryEvent)(nil)
	_ TextWritable   = (*QueryEvent)(nil)
	_ json.Marshaler = (*QueryEvent)(nil)
)

// NewQueryEvent creates a new query event.
func NewQueryEvent(body string, elapsed time.Duration, options ...EventMetaOption) *QueryEvent {
	return &QueryEvent{
		EventMeta: NewEventMeta(Query, options...),
		Body:      body,
		Elapsed:   elapsed,
	}
}

// NewQueryEventListener returns a new listener for spiffy events.
func NewQueryEventListener(listener func(context.Context, *QueryEvent)) Listener {
	return func(ctx context.Context, e Event) {
		if typed, isTyped := e.(*QueryEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// QueryEvent represents a database query.
type QueryEvent struct {
	*EventMeta

	Database   string
	Engine     string
	Username   string
	QueryLabel string
	Body       string
	Elapsed    time.Duration
	Err        error
}

// WriteText writes the event text to the output.
func (e QueryEvent) WriteText(tf TextFormatter, wr io.Writer) {
	io.WriteString(wr, "[")
	if len(e.Engine) > 0 {
		io.WriteString(wr, tf.Colorize(e.Engine, ansi.ColorLightWhite))
		io.WriteString(wr, Space)
	}
	if len(e.Username) > 0 {
		io.WriteString(wr, tf.Colorize(e.Username, ansi.ColorLightWhite))
		io.WriteString(wr, "@")
	}
	io.WriteString(wr, tf.Colorize(e.Database, ansi.ColorLightWhite))
	io.WriteString(wr, "]")

	if len(e.QueryLabel) > 0 {
		io.WriteString(wr, Space)
		io.WriteString(wr, fmt.Sprintf("[%s]", tf.Colorize(e.QueryLabel, ansi.ColorLightWhite)))
	}

	io.WriteString(wr, Space)
	io.WriteString(wr, e.Elapsed.String())

	if e.Err != nil {
		io.WriteString(wr, Space)
		io.WriteString(wr, tf.Colorize("failed", ansi.ColorRed))
	}

	if len(e.Body) > 0 {
		io.WriteString(wr, Space)
		io.WriteString(wr, stringutil.CompressSpace(e.Body))
	}
}

// MarshalJSON implements json.Marshaler.
func (e QueryEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(MergeDecomposed(e.EventMeta.Decompose(), map[string]interface{}{
		"engine":     e.Engine,
		"database":   e.Database,
		"username":   e.Username,
		"queryLabel": e.QueryLabel,
		"body":       e.Body,
		"err":        e.Err,
		"elapsed":    timeutil.Milliseconds(e.Elapsed),
	}))
}
