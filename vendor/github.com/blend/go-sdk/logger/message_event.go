package logger

import (
	"bytes"
	"fmt"
	"time"
)

// these are compile time assertions
var (
	_ Event            = &MessageEvent{}
	_ EventHeadings    = &MessageEvent{}
	_ EventLabels      = &MessageEvent{}
	_ EventAnnotations = &MessageEvent{}
)

// Messagef returns a new Message Event.
func Messagef(flag Flag, format string, args ...interface{}) *MessageEvent {
	return &MessageEvent{
		EventMeta: NewEventMeta(flag),
		message:   fmt.Sprintf(format, args...),
	}
}

// NewMessageEventListener returns a new message event listener.
func NewMessageEventListener(listener func(*MessageEvent)) Listener {
	return func(e Event) {
		if typed, isTyped := e.(*MessageEvent); isTyped {
			listener(typed)
		}
	}
}

// MessageEvent is a common type of message.
type MessageEvent struct {
	*EventMeta
	message string
}

// WithHeadings sets the headings.
func (e *MessageEvent) WithHeadings(headings ...string) *MessageEvent {
	e.headings = headings
	return e
}

// WithLabel sets a label on the event for later filtering.
func (e *MessageEvent) WithLabel(key, value string) *MessageEvent {
	e.AddLabelValue(key, value)
	return e
}

// WithAnnotation adds an annotation to the event.
func (e *MessageEvent) WithAnnotation(key, value string) *MessageEvent {
	e.AddAnnotationValue(key, value)
	return e
}

// WithFlag sets the message flag.
func (e *MessageEvent) WithFlag(flag Flag) *MessageEvent {
	e.flag = flag
	return e
}

// WithFlagTextColor sets the message flag text color.
func (e *MessageEvent) WithFlagTextColor(color AnsiColor) *MessageEvent {
	e.flagTextColor = color
	return e
}

// WithTimestamp sets the message timestamp.
func (e *MessageEvent) WithTimestamp(ts time.Time) *MessageEvent {
	e.ts = ts
	return e
}

// WithMessage sets the message.
func (e *MessageEvent) WithMessage(message string) *MessageEvent {
	e.message = message
	return e
}

// Message returns the message.
func (e *MessageEvent) Message() string {
	return e.message
}

// WriteText implements TextWritable.
func (e *MessageEvent) WriteText(formatter TextFormatter, buf *bytes.Buffer) {
	buf.WriteString(e.message)
}

// WriteJSON implements JSONWriteable.
func (e *MessageEvent) WriteJSON() JSONObj {
	return JSONObj{
		JSONFieldMessage: e.message,
	}
}

// String returns the message event body.
func (e *MessageEvent) String() string {
	return e.message
}
