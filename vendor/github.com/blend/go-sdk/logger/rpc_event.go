package logger

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/timeutil"
)

// these are compile time assertions
var (
	_ Event = (*RPCEvent)(nil)
)

// NewRPCEvent creates a new rpc event.
func NewRPCEvent(method string, elapsed time.Duration) *RPCEvent {
	return &RPCEvent{
		EventMeta: NewEventMeta(RPC),
		Method:    method,
		Elapsed:   elapsed,
	}
}

// NewRPCEventListener returns a new web request event listener.
func NewRPCEventListener(listener func(context.Context, *RPCEvent)) Listener {
	return func(ctx context.Context, e Event) {
		if typed, isTyped := e.(*RPCEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// RPCEvent is an event type for rpc
type RPCEvent struct {
	*EventMeta
	Engine      string
	Peer        string
	Method      string
	UserAgent   string
	Authority   string
	ContentType string
	Elapsed     time.Duration
	Err         error
}

// WriteText implements TextWritable.
func (e RPCEvent) WriteText(tf TextFormatter, wr io.Writer) {

	if e.Engine != "" {
		io.WriteString(wr, "[")
		io.WriteString(wr, tf.Colorize(e.Engine, ansi.ColorLightWhite))
		io.WriteString(wr, "]")
	}
	if e.Method != "" {
		if e.Engine != "" {
			io.WriteString(wr, Space)
		}
		io.WriteString(wr, tf.Colorize(e.Method, ansi.ColorBlue))
	}
	if e.Peer != "" {
		io.WriteString(wr, Space)
		io.WriteString(wr, e.Peer)
	}
	if e.Authority != "" {
		io.WriteString(wr, Space)
		io.WriteString(wr, e.Authority)
	}
	if e.UserAgent != "" {
		io.WriteString(wr, Space)
		io.WriteString(wr, e.UserAgent)
	}
	if e.ContentType != "" {
		io.WriteString(wr, Space)
		io.WriteString(wr, e.ContentType)
	}

	io.WriteString(wr, Space)
	io.WriteString(wr, e.Elapsed.String())

	if e.Err != nil {
		io.WriteString(wr, Space)
		io.WriteString(wr, tf.Colorize("failed", ansi.ColorRed))
	}
}

// MarshalJSON implements json.Marshaler.
func (e RPCEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(MergeDecomposed(e.EventMeta.Decompose(), map[string]interface{}{
		"engine":      e.Engine,
		"peer":        e.Peer,
		"method":      e.Method,
		"userAgent":   e.UserAgent,
		"authority":   e.Authority,
		"contentType": e.ContentType,
		"elapsed":     timeutil.Milliseconds(e.Elapsed),
		"err":         e.Err,
	}))
}
