package logger

import (
	"bytes"
	"time"
)

// these are compile time assertions
var (
	_ Event            = &RPCEvent{}
	_ EventHeadings    = &RPCEvent{}
	_ EventLabels      = &RPCEvent{}
	_ EventAnnotations = &RPCEvent{}
)

// NewRPCEvent creates a new rpc event.
func NewRPCEvent(method string, elapsed time.Duration) *RPCEvent {
	return &RPCEvent{
		EventMeta: NewEventMeta(RPC),
		method:    method,
		elapsed:   elapsed,
	}
}

// NewRPCEventListener returns a new web request event listener.
func NewRPCEventListener(listener func(*RPCEvent)) Listener {
	return func(e Event) {
		if typed, isTyped := e.(*RPCEvent); isTyped {
			listener(typed)
		}
	}
}

// RPCEvent is an event type for rpc
type RPCEvent struct {
	*EventMeta
	engine      string
	peer        string
	method      string
	userAgent   string
	authority   string
	contentType string
	elapsed     time.Duration
	err         error
}

// WithHeadings sets the headings.
func (e *RPCEvent) WithHeadings(headings ...string) *RPCEvent {
	e.headings = headings
	return e
}

// WithLabel sets a label on the event for later filtering.
func (e *RPCEvent) WithLabel(key, value string) *RPCEvent {
	e.AddLabelValue(key, value)
	return e
}

// WithAnnotation adds an annotation to the event.
func (e *RPCEvent) WithAnnotation(key, value string) *RPCEvent {
	e.AddAnnotationValue(key, value)
	return e
}

// WithFlag sets the event flag.
func (e *RPCEvent) WithFlag(flag Flag) *RPCEvent {
	e.flag = flag
	return e
}

// WithTimestamp sets the timestamp.
func (e *RPCEvent) WithTimestamp(ts time.Time) *RPCEvent {
	e.ts = ts
	return e
}

// WithEngine sets the engine.
func (e *RPCEvent) WithEngine(engine string) *RPCEvent {
	e.engine = engine
	return e
}

// Engine returns the engine.
func (e RPCEvent) Engine() string {
	return e.engine
}

// WithPeer sets the peer.
func (e *RPCEvent) WithPeer(peer string) *RPCEvent {
	e.peer = peer
	return e
}

// Peer returns the peer.
func (e RPCEvent) Peer() string {
	return e.peer
}

// WithMethod sets the method.
func (e *RPCEvent) WithMethod(method string) *RPCEvent {
	e.method = method
	return e
}

// Method returns the method.
func (e RPCEvent) Method() string {
	return e.method
}

// WithAuthority sets the authority.
func (e *RPCEvent) WithAuthority(authority string) *RPCEvent {
	e.authority = authority
	return e
}

// Authority returns the authority.
func (e *RPCEvent) Authority() string {
	return e.authority
}

// WithUserAgent sets the user agent.
func (e *RPCEvent) WithUserAgent(userAgent string) *RPCEvent {
	e.userAgent = userAgent
	return e
}

// UserAgent returns the user agent.
func (e *RPCEvent) UserAgent() string {
	return e.userAgent
}

// WithContentType sets the content type.
func (e *RPCEvent) WithContentType(contentType string) *RPCEvent {
	e.contentType = contentType
	return e
}

// ContentType is the type of the response.
func (e *RPCEvent) ContentType() string {
	return e.contentType
}

// WithElapsed sets the elapsed time.
func (e *RPCEvent) WithElapsed(elapsed time.Duration) *RPCEvent {
	e.elapsed = elapsed
	return e
}

// Elapsed returns the elapsed time.
func (e *RPCEvent) Elapsed() time.Duration {
	return e.elapsed
}

// WithErr sets the error on the event.
func (e *RPCEvent) WithErr(err error) *RPCEvent {
	e.err = err
	return e
}

// Err returns the event err (if any).
func (e RPCEvent) Err() error {
	return e.err
}

// WriteText implements TextWritable.
func (e *RPCEvent) WriteText(tf TextFormatter, buf *bytes.Buffer) {

	if e.engine != "" {
		buf.WriteString("[")
		buf.WriteString(tf.Colorize(e.engine, ColorLightWhite))
		buf.WriteString("]")
	}
	if e.method != "" {
		if e.engine != "" {
			buf.WriteRune(RuneSpace)
		}
		buf.WriteString(tf.Colorize(e.method, ColorBlue))
	}
	if e.peer != "" {
		buf.WriteRune(RuneSpace)
		buf.WriteString(e.peer)
	}
	if e.authority != "" {
		buf.WriteRune(RuneSpace)
		buf.WriteString(e.authority)
	}
	if e.userAgent != "" {
		buf.WriteRune(RuneSpace)
		buf.WriteString(e.userAgent)
	}
	if e.contentType != "" {
		buf.WriteRune(RuneSpace)
		buf.WriteString(e.contentType)
	}

	buf.WriteRune(RuneSpace)
	buf.WriteString(e.elapsed.String())

	if e.err != nil {
		buf.WriteRune(RuneSpace)
		buf.WriteString(tf.Colorize("failed", ColorRed))
	}
}

// WriteJSON implements JSONWritable.
func (e *RPCEvent) WriteJSON() JSONObj {
	return JSONObj{
		"engine":         e.engine,
		"peer":           e.peer,
		"method":         e.method,
		"authority":      e.authority,
		"userAgent":      e.userAgent,
		"contentType":    e.contentType,
		JSONFieldElapsed: Milliseconds(e.elapsed),
		JSONFieldErr:     e.err,
	}
}
