package logger

import (
	"bytes"
	"net/http"
	"time"
)

// these are compile time assertions
var (
	_ Event            = &HTTPRequestEvent{}
	_ EventHeadings    = &HTTPRequestEvent{}
	_ EventLabels      = &HTTPRequestEvent{}
	_ EventAnnotations = &HTTPRequestEvent{}
)

// NewHTTPRequestEvent creates a new web request event.
func NewHTTPRequestEvent(req *http.Request) *HTTPRequestEvent {
	return &HTTPRequestEvent{
		EventMeta: NewEventMeta(HTTPRequest),
		req:       req,
	}
}

// NewHTTPRequestEventListener returns a new web request event listener.
func NewHTTPRequestEventListener(listener func(*HTTPRequestEvent)) Listener {
	return func(e Event) {
		if typed, isTyped := e.(*HTTPRequestEvent); isTyped {
			listener(typed)
		}
	}
}

// HTTPRequestEvent is an event type for http responses.
type HTTPRequestEvent struct {
	*EventMeta
	req   *http.Request
	route string
	state map[interface{}]interface{}
}

// WithHeadings sets the headings.
func (e *HTTPRequestEvent) WithHeadings(headings ...string) *HTTPRequestEvent {
	e.headings = headings
	return e
}

// WithLabel sets a label on the event for later filtering.
func (e *HTTPRequestEvent) WithLabel(key, value string) *HTTPRequestEvent {
	e.AddLabelValue(key, value)
	return e
}

// WithAnnotation adds an annotation to the event.
func (e *HTTPRequestEvent) WithAnnotation(key, value string) *HTTPRequestEvent {
	e.AddAnnotationValue(key, value)
	return e
}

// WithFlag sets the event flag.
func (e *HTTPRequestEvent) WithFlag(flag Flag) *HTTPRequestEvent {
	e.flag = flag
	return e
}

// WithTimestamp sets the timestamp.
func (e *HTTPRequestEvent) WithTimestamp(ts time.Time) *HTTPRequestEvent {
	e.ts = ts
	return e
}

// WithRequest sets the request metadata.
func (e *HTTPRequestEvent) WithRequest(req *http.Request) *HTTPRequestEvent {
	e.req = req
	return e
}

// Request returns the request metadata.
func (e *HTTPRequestEvent) Request() *http.Request {
	return e.req
}

// WithRoute sets the mux route.
func (e *HTTPRequestEvent) WithRoute(route string) *HTTPRequestEvent {
	e.route = route
	return e
}

// Route is the mux route of the request.
func (e *HTTPRequestEvent) Route() string {
	return e.route
}

// WithState sets the request state.
func (e *HTTPRequestEvent) WithState(state map[interface{}]interface{}) *HTTPRequestEvent {
	e.state = state
	return e
}

// State returns the state of the request.
func (e *HTTPRequestEvent) State() map[interface{}]interface{} {
	return e.state
}

// WriteText implements TextWritable.
func (e *HTTPRequestEvent) WriteText(formatter TextFormatter, buf *bytes.Buffer) {
	TextWriteHTTPRequest(formatter, buf, e.req)
}

// WriteJSON implements JSONWritable.
func (e *HTTPRequestEvent) WriteJSON() JSONObj {
	return JSONWriteHTTPRequest(e.req)
}
