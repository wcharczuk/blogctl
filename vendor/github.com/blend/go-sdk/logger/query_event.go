package logger

import (
	"bytes"
	"fmt"
	"time"
)

// these are compile time assertions
var (
	_ Event            = &QueryEvent{}
	_ EventHeadings    = &QueryEvent{}
	_ EventLabels      = &QueryEvent{}
	_ EventAnnotations = &QueryEvent{}
)

// NewQueryEvent creates a new query event.
func NewQueryEvent(body string, elapsed time.Duration) *QueryEvent {
	return &QueryEvent{
		EventMeta: NewEventMeta(Query),
		body:      body,
		elapsed:   elapsed,
	}
}

// NewQueryEventListener returns a new listener for spiffy events.
func NewQueryEventListener(listener func(e *QueryEvent)) Listener {
	return func(e Event) {
		if typed, isTyped := e.(*QueryEvent); isTyped {
			listener(typed)
		}
	}
}

// QueryEvent represents a database query.
type QueryEvent struct {
	*EventMeta

	database   string
	engine     string
	username   string
	queryLabel string
	body       string
	elapsed    time.Duration
	err        error
}

// WithHeadings sets the headings.
func (e *QueryEvent) WithHeadings(headings ...string) *QueryEvent {
	e.headings = headings
	return e
}

// WithLabel sets a label on the event for later filtering.
func (e *QueryEvent) WithLabel(key, value string) *QueryEvent {
	e.AddLabelValue(key, value)
	return e
}

// WithAnnotation adds an annotation to the event.
func (e *QueryEvent) WithAnnotation(key, value string) *QueryEvent {
	e.AddAnnotationValue(key, value)
	return e
}

// WithFlag sets the flag.
func (e *QueryEvent) WithFlag(flag Flag) *QueryEvent {
	e.flag = flag
	return e
}

// WithTimestamp sets the timestamp.
func (e *QueryEvent) WithTimestamp(ts time.Time) *QueryEvent {
	e.ts = ts
	return e
}

// WithUsername sets the engine.
func (e *QueryEvent) WithUsername(username string) *QueryEvent {
	e.username = username
	return e
}

// Username returns the username.
func (e QueryEvent) Username() string {
	return e.username
}

// WithEngine sets the engine.
func (e *QueryEvent) WithEngine(engine string) *QueryEvent {
	e.engine = engine
	return e
}

// Engine returns the engine.
func (e QueryEvent) Engine() string {
	return e.engine
}

// WithDatabase sets the database.
func (e *QueryEvent) WithDatabase(db string) *QueryEvent {
	e.database = db
	return e
}

// Database returns the event database.
func (e QueryEvent) Database() string {
	return e.database
}

// WithQueryLabel sets the query label.
func (e *QueryEvent) WithQueryLabel(queryLabel string) *QueryEvent {
	e.queryLabel = queryLabel
	return e
}

// QueryLabel returns the query label.
func (e QueryEvent) QueryLabel() string {
	return e.queryLabel
}

// WithBody sets the body.
func (e *QueryEvent) WithBody(body string) *QueryEvent {
	e.body = body
	return e
}

// Body returns the query body.
func (e QueryEvent) Body() string {
	return e.body
}

// WithElapsed sets the elapsed time.
func (e *QueryEvent) WithElapsed(elapsed time.Duration) *QueryEvent {
	e.elapsed = elapsed
	return e
}

// Elapsed returns the elapsed time.
func (e QueryEvent) Elapsed() time.Duration {
	return e.elapsed
}

// WithErr sets the error on the event.
func (e *QueryEvent) WithErr(err error) *QueryEvent {
	e.err = err
	return e
}

// Err returns the event err (if any).
func (e QueryEvent) Err() error {
	return e.err
}

// WriteText writes the event text to the output.
func (e QueryEvent) WriteText(tf TextFormatter, buf *bytes.Buffer) {
	buf.WriteString("[")
	if len(e.engine) > 0 {
		buf.WriteString(tf.Colorize(e.engine, ColorLightWhite))
		buf.WriteRune(RuneSpace)
	}
	if len(e.username) > 0 {
		buf.WriteString(tf.Colorize(e.username, ColorLightWhite))
		buf.WriteRune('@')
	}
	buf.WriteString(tf.Colorize(e.database, ColorLightWhite))
	buf.WriteString("]")

	if len(e.queryLabel) > 0 {
		buf.WriteRune(RuneSpace)
		buf.WriteString(fmt.Sprintf("[%s]", tf.Colorize(e.queryLabel, ColorLightWhite)))
	}

	buf.WriteRune(RuneSpace)
	buf.WriteString(e.elapsed.String())

	if e.err != nil {
		buf.WriteRune(RuneSpace)
		buf.WriteString(tf.Colorize("failed", ColorRed))
	}

	if len(e.body) > 0 {
		buf.WriteRune(RuneSpace)
		buf.WriteString(CompressWhitespace(e.body))
	}
}

// WriteJSON implements JSONWritable.
func (e QueryEvent) WriteJSON() JSONObj {
	return JSONObj{
		"engine":         e.engine,
		"database":       e.database,
		"username":       e.username,
		"queryLabel":     e.queryLabel,
		"body":           e.body,
		JSONFieldErr:     e.err,
		JSONFieldElapsed: Milliseconds(e.elapsed),
	}
}
