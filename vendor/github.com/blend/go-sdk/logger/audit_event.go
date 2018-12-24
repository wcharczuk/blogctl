package logger

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

// these are compile time assertions
var (
	_ Event            = &AuditEvent{}
	_ EventHeadings    = &AuditEvent{}
	_ EventLabels      = &AuditEvent{}
	_ EventAnnotations = &AuditEvent{}
)

// NewAuditEvent returns a new audit event.
func NewAuditEvent(principal, verb string) *AuditEvent {
	return &AuditEvent{
		EventMeta: NewEventMeta(Audit),
		principal: principal,
		verb:      verb,
	}
}

// NewAuditEventListener returns a new audit event listener.
func NewAuditEventListener(listener func(me *AuditEvent)) Listener {
	return func(e Event) {
		if typed, isTyped := e.(*AuditEvent); isTyped {
			listener(typed)
		}
	}
}

// AuditEvent is a common type of event detailing a business action by a subject.
type AuditEvent struct {
	*EventMeta

	context       string
	principal     string
	verb          string
	noun          string
	subject       string
	property      string
	remoteAddress string
	userAgent     string
	extra         map[string]string
}

// WithHeadings sets the headings.
func (e *AuditEvent) WithHeadings(headings ...string) *AuditEvent {
	e.headings = headings
	return e
}

// WithLabel sets a label on the event for later filtering.
func (e *AuditEvent) WithLabel(key, value string) *AuditEvent {
	e.AddLabelValue(key, value)
	return e
}

// WithAnnotation adds an annotation to the event.
func (e *AuditEvent) WithAnnotation(key, value string) *AuditEvent {
	e.AddAnnotationValue(key, value)
	return e
}

// WithFlag sets the audit event flag
func (e *AuditEvent) WithFlag(f Flag) *AuditEvent {
	e.flag = f
	return e
}

// WithTimestamp sets the message timestamp.
func (e *AuditEvent) WithTimestamp(ts time.Time) *AuditEvent {
	e.ts = ts
	return e
}

// WithContext sets the context.
func (e *AuditEvent) WithContext(context string) *AuditEvent {
	e.context = context
	return e
}

// Context returns the audit context.
func (e *AuditEvent) Context() string {
	return e.context
}

// WithPrincipal sets the principal.
func (e *AuditEvent) WithPrincipal(principal string) *AuditEvent {
	e.principal = principal
	return e
}

// Principal returns the principal.
func (e AuditEvent) Principal() string {
	return e.principal
}

// WithVerb sets the verb.
func (e *AuditEvent) WithVerb(verb string) *AuditEvent {
	e.verb = verb
	return e
}

// Verb returns the verb.
func (e AuditEvent) Verb() string {
	return e.verb
}

// WithNoun sets the noun.
func (e *AuditEvent) WithNoun(noun string) *AuditEvent {
	e.noun = noun
	return e
}

// Noun returns the noun.
func (e AuditEvent) Noun() string {
	return e.noun
}

// WithSubject sets the subject.
func (e *AuditEvent) WithSubject(subject string) *AuditEvent {
	e.subject = subject
	return e
}

// Subject returns the subject.
func (e AuditEvent) Subject() string {
	return e.subject
}

// WithProperty sets the property.
func (e *AuditEvent) WithProperty(property string) *AuditEvent {
	e.property = property
	return e
}

// Property returns the property.
func (e AuditEvent) Property() string {
	return e.property
}

// WithRemoteAddress sets the remote address.
func (e *AuditEvent) WithRemoteAddress(remoteAddr string) *AuditEvent {
	e.remoteAddress = remoteAddr
	return e
}

// RemoteAddress returns the remote address.
func (e AuditEvent) RemoteAddress() string {
	return e.remoteAddress
}

// WithUserAgent sets the user agent.
func (e *AuditEvent) WithUserAgent(userAgent string) *AuditEvent {
	e.userAgent = userAgent
	return e
}

// UserAgent returns the user agent.
func (e AuditEvent) UserAgent() string {
	return e.userAgent
}

// WithExtra sets the extra info.
func (e *AuditEvent) WithExtra(extra map[string]string) *AuditEvent {
	e.extra = extra
	return e
}

// Extra returns the extra information.
func (e AuditEvent) Extra() map[string]string {
	return e.extra
}

// WriteText implements TextWritable.
func (e AuditEvent) WriteText(formatter TextFormatter, buf *bytes.Buffer) {
	if len(e.context) > 0 {
		buf.WriteString(formatter.Colorize("Context:", ColorGray))
		buf.WriteString(e.context)
		buf.WriteRune(RuneSpace)
	}
	if len(e.principal) > 0 {
		buf.WriteString(formatter.Colorize("Principal:", ColorGray))
		buf.WriteString(e.principal)
		buf.WriteRune(RuneSpace)
	}
	if len(e.verb) > 0 {
		buf.WriteString(formatter.Colorize("Verb:", ColorGray))
		buf.WriteString(e.verb)
		buf.WriteRune(RuneSpace)
	}
	if len(e.noun) > 0 {
		buf.WriteString(formatter.Colorize("Noun:", ColorGray))
		buf.WriteString(e.noun)
		buf.WriteRune(RuneSpace)
	}
	if len(e.subject) > 0 {
		buf.WriteString(formatter.Colorize("Subject:", ColorGray))
		buf.WriteString(e.subject)
		buf.WriteRune(RuneSpace)
	}
	if len(e.property) > 0 {
		buf.WriteString(formatter.Colorize("Property:", ColorGray))
		buf.WriteString(e.property)
		buf.WriteRune(RuneSpace)
	}
	if len(e.remoteAddress) > 0 {
		buf.WriteString(formatter.Colorize("Remote Addr:", ColorGray))
		buf.WriteString(e.remoteAddress)
		buf.WriteRune(RuneSpace)
	}
	if len(e.userAgent) > 0 {
		buf.WriteString(formatter.Colorize("UA:", ColorGray))
		buf.WriteString(e.userAgent)
		buf.WriteRune(RuneSpace)
	}
	if len(e.extra) > 0 {
		var values []string
		for key, value := range e.extra {
			values = append(values, fmt.Sprintf("%s%s", formatter.Colorize(key+":", ColorGray), value))
		}
		buf.WriteString(strings.Join(values, " "))
	}
}

// WriteJSON implements JSONWritable.
func (e AuditEvent) WriteJSON() JSONObj {
	return JSONObj{
		"context":    e.context,
		"principal":  e.principal,
		"verb":       e.verb,
		"noun":       e.noun,
		"subject":    e.subject,
		"property":   e.property,
		"remoteAddr": e.remoteAddress,
		"ua":         e.userAgent,
		"extra":      e.extra,
	}
}
