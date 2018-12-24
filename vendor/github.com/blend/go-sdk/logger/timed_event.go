package logger

import (
	"bytes"
	"fmt"
	"time"
)

// these are compile time assertions
var (
	_ Event            = &TimedEvent{}
	_ EventHeadings    = &TimedEvent{}
	_ EventLabels      = &TimedEvent{}
	_ EventAnnotations = &TimedEvent{}
)

// Timedf returns a timed message event.
func Timedf(flag Flag, elapsed time.Duration, format string, args ...Any) *TimedEvent {
	return &TimedEvent{
		EventMeta: NewEventMeta(flag),
		message:   fmt.Sprintf(format, args...),
		elapsed:   elapsed,
	}
}

// NewTimedEventListener returns a new timed event listener.
func NewTimedEventListener(listener func(e *TimedEvent)) Listener {
	return func(e Event) {
		if typed, isTyped := e.(*TimedEvent); isTyped {
			listener(typed)
		}
	}
}

// TimedEvent is a message event with an elapsed time.
type TimedEvent struct {
	*EventMeta

	message string
	elapsed time.Duration
}

// WithHeadings sets the headings.
func (e *TimedEvent) WithHeadings(headings ...string) *TimedEvent {
	e.headings = headings
	return e
}

// WithLabel sets a label on the event for later filtering.
func (e *TimedEvent) WithLabel(key, value string) *TimedEvent {
	e.AddLabelValue(key, value)
	return e
}

// WithAnnotation adds an annotation to the event.
func (e *TimedEvent) WithAnnotation(key, value string) *TimedEvent {
	e.AddAnnotationValue(key, value)
	return e
}

// WithFlag sets the timed message flag.
func (e *TimedEvent) WithFlag(flag Flag) *TimedEvent {
	e.flag = flag
	return e
}

// WithTimestamp sets the message timestamp.
func (e *TimedEvent) WithTimestamp(ts time.Time) *TimedEvent {
	e.ts = ts
	return e
}

// WithMessage sets the message.
func (e *TimedEvent) WithMessage(message string) *TimedEvent {
	e.message = message
	return e
}

// Message returns the string message.
func (e TimedEvent) Message() string {
	return e.message
}

// WithElapsed sets the elapsed time.
func (e *TimedEvent) WithElapsed(elapsed time.Duration) *TimedEvent {
	e.elapsed = elapsed
	return e
}

// Elapsed returns the elapsed time.
func (e TimedEvent) Elapsed() time.Duration {
	return e.elapsed
}

// String implements fmt.Stringer
func (e TimedEvent) String() string {
	return fmt.Sprintf("%s (%v)", e.message, e.elapsed)
}

// WriteText implements TextWritable.
func (e TimedEvent) WriteText(tf TextFormatter, buf *bytes.Buffer) {
	buf.WriteString(e.String())
}

// WriteJSON implements JSONWritable.
func (e TimedEvent) WriteJSON() JSONObj {
	return JSONObj{
		JSONFieldMessage: e.message,
		JSONFieldElapsed: Milliseconds(e.elapsed),
	}
}
