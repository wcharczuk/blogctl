package logger

import (
	"bytes"
	"net/http"
	"time"
)

// these are compile time assertions
var (
	_ Event            = &HTTPResponseEvent{}
	_ EventHeadings    = &HTTPResponseEvent{}
	_ EventLabels      = &HTTPResponseEvent{}
	_ EventAnnotations = &HTTPResponseEvent{}
)

// NewHTTPResponseEvent is an event representing a response to an http request.
func NewHTTPResponseEvent(req *http.Request) *HTTPResponseEvent {
	return &HTTPResponseEvent{
		EventMeta: NewEventMeta(HTTPResponse),
		req:       req,
	}
}

// NewHTTPResponseEventListener returns a new web request event listener.
func NewHTTPResponseEventListener(listener func(*HTTPResponseEvent)) Listener {
	return func(e Event) {
		if typed, isTyped := e.(*HTTPResponseEvent); isTyped {
			listener(typed)
		}
	}
}

// HTTPResponseEvent is an event type for responses.
type HTTPResponseEvent struct {
	*EventMeta

	req   *http.Request
	route string

	contentLength   int
	contentType     string
	contentEncoding string

	statusCode int
	elapsed    time.Duration

	state map[interface{}]interface{}
}

// WithHeadings sets the headings.
func (e *HTTPResponseEvent) WithHeadings(headings ...string) *HTTPResponseEvent {
	e.headings = headings
	return e
}

// WithLabel sets a label on the event for later filtering.
func (e *HTTPResponseEvent) WithLabel(key, value string) *HTTPResponseEvent {
	e.AddLabelValue(key, value)
	return e
}

// WithAnnotation adds an annotation to the event.
func (e *HTTPResponseEvent) WithAnnotation(key, value string) *HTTPResponseEvent {
	e.AddAnnotationValue(key, value)
	return e
}

// WithFlag sets the event flag.
func (e *HTTPResponseEvent) WithFlag(flag Flag) *HTTPResponseEvent {
	e.flag = flag
	return e
}

// WithTimestamp sets the timestamp.
func (e *HTTPResponseEvent) WithTimestamp(ts time.Time) *HTTPResponseEvent {
	e.ts = ts
	return e
}

// WithRequest sets the request metadata.
func (e *HTTPResponseEvent) WithRequest(req *http.Request) *HTTPResponseEvent {
	e.req = req
	return e
}

// Request returns the request metadata.
func (e *HTTPResponseEvent) Request() *http.Request {
	return e.req
}

// WithRoute sets the mux route.
func (e *HTTPResponseEvent) WithRoute(route string) *HTTPResponseEvent {
	e.route = route
	return e
}

// Route is the mux route of the request.
func (e *HTTPResponseEvent) Route() string {
	return e.route
}

// WithStatusCode sets the status code.
func (e *HTTPResponseEvent) WithStatusCode(statusCode int) *HTTPResponseEvent {
	e.statusCode = statusCode
	return e
}

// StatusCode is the HTTP status code of the response.
func (e *HTTPResponseEvent) StatusCode() int {
	return e.statusCode
}

// WithElapsed sets the elapsed time.
func (e *HTTPResponseEvent) WithElapsed(elapsed time.Duration) *HTTPResponseEvent {
	e.elapsed = elapsed
	return e
}

// Elapsed returns the elapsed time.
func (e *HTTPResponseEvent) Elapsed() time.Duration {
	return e.elapsed
}

// WithContentLength sets the content length.
func (e *HTTPResponseEvent) WithContentLength(contentLength int) *HTTPResponseEvent {
	e.contentLength = contentLength
	return e
}

// ContentLength is the size of the response.
func (e *HTTPResponseEvent) ContentLength() int {
	return e.contentLength
}

// WithContentType sets the content type.
func (e *HTTPResponseEvent) WithContentType(contentType string) *HTTPResponseEvent {
	e.contentType = contentType
	return e
}

// ContentType is the type of the response.
func (e *HTTPResponseEvent) ContentType() string {
	return e.contentType
}

// WithContentEncoding sets the content encoding.
func (e *HTTPResponseEvent) WithContentEncoding(contentEncoding string) *HTTPResponseEvent {
	e.contentEncoding = contentEncoding
	return e
}

// ContentEncoding is the encoding of the response.
func (e *HTTPResponseEvent) ContentEncoding() string {
	return e.contentEncoding
}

// WithState sets the request state.
func (e *HTTPResponseEvent) WithState(state map[interface{}]interface{}) *HTTPResponseEvent {
	e.state = state
	return e
}

// State returns the state of the request.
func (e *HTTPResponseEvent) State() map[interface{}]interface{} {
	return e.state
}

// WriteText implements TextWritable.
func (e *HTTPResponseEvent) WriteText(formatter TextFormatter, buf *bytes.Buffer) {
	TextWriteHTTPResponse(formatter, buf, e.req, e.statusCode, e.contentLength, e.contentType, e.elapsed)
}

// WriteJSON implements JSONWritable.
func (e *HTTPResponseEvent) WriteJSON() JSONObj {
	return JSONWriteHTTPResponse(e.req, e.statusCode, e.contentLength, e.contentType, e.contentEncoding, e.elapsed)
}
