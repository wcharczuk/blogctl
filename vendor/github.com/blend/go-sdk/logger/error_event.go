package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// these are compile time assertions
var (
	_ Event            = &ErrorEvent{}
	_ EventHeadings    = &ErrorEvent{}
	_ EventLabels      = &ErrorEvent{}
	_ EventAnnotations = &ErrorEvent{}
)

// Errorf returns a new error event based on format and arguments.
func Errorf(flag Flag, format string, args ...interface{}) *ErrorEvent {
	return &ErrorEvent{
		EventMeta: NewEventMeta(flag),
		err:       fmt.Errorf(format, args...),
	}
}

// NewErrorEvent returns a new error event.
func NewErrorEvent(flag Flag, err error) *ErrorEvent {
	return &ErrorEvent{
		EventMeta: NewEventMeta(flag),
		err:       err,
	}
}

// NewErrorEventWithState returns a new error event with state.
func NewErrorEventWithState(flag Flag, err error, state interface{}) *ErrorEvent {
	return &ErrorEvent{
		EventMeta: NewEventMeta(flag),
		err:       err,
		state:     state,
	}
}

// NewErrorEventListener returns a new error event listener.
func NewErrorEventListener(listener func(*ErrorEvent)) Listener {
	return func(e Event) {
		if typed, isTyped := e.(*ErrorEvent); isTyped {
			listener(typed)
		}
	}
}

// ErrorEvent is an event that wraps an error.
type ErrorEvent struct {
	*EventMeta

	err   error
	state interface{}
}

// IsError indicates if we should write to the error writer or not.
func (e *ErrorEvent) IsError() bool {
	return true
}

// WithHeadings sets the headings.
func (e *ErrorEvent) WithHeadings(headings ...string) *ErrorEvent {
	e.headings = headings
	return e
}

// WithLabel sets a label on the event for later filtering.
func (e *ErrorEvent) WithLabel(key, value string) *ErrorEvent {
	e.AddLabelValue(key, value)
	return e
}

// WithAnnotation adds an annotation to the event.
func (e *ErrorEvent) WithAnnotation(key, value string) *ErrorEvent {
	e.AddAnnotationValue(key, value)
	return e
}

// WithTimestamp sets the event timestamp.
func (e *ErrorEvent) WithTimestamp(ts time.Time) *ErrorEvent {
	e.ts = ts
	return e
}

// WithFlag sets the event flag.
func (e *ErrorEvent) WithFlag(flag Flag) *ErrorEvent {
	e.flag = flag
	return e
}

// WithFlagTextColor sets the flag text color.
func (e *ErrorEvent) WithFlagTextColor(color AnsiColor) *ErrorEvent {
	e.flagTextColor = color
	return e
}

// WithErr sets the error.
func (e *ErrorEvent) WithErr(err error) *ErrorEvent {
	e.err = err
	return e
}

// Err returns the underlying error.
func (e *ErrorEvent) Err() error {
	return e.err
}

// WithState sets the state.
func (e *ErrorEvent) WithState(state interface{}) *ErrorEvent {
	e.state = state
	return e
}

// State returns underlying state, typically an http.Request.
func (e *ErrorEvent) State() interface{} {
	return e.state
}

// WriteText implements TextWritable.
func (e *ErrorEvent) WriteText(formatter TextFormatter, buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("%+v", e.err))
}

// WriteJSON implements JSONWritable.
func (e *ErrorEvent) WriteJSON() JSONObj {
	var errorJSON interface{}
	if _, ok := e.err.(json.Marshaler); ok {
		errorJSON = e.err
	} else {
		errorJSON = e.err.Error()
	}
	return JSONObj{
		JSONFieldErr: errorJSON,
	}
}
