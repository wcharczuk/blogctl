package logger

import "time"

// these are compile time assertions
var (
	_ Event            = &EventMeta{}
	_ EventHeadings    = &EventMeta{}
	_ EventLabels      = &EventMeta{}
	_ EventAnnotations = &EventMeta{}
)

// NewEventMeta returns a new event meta.
func NewEventMeta(flag Flag) *EventMeta {
	return &EventMeta{
		flag:        flag,
		ts:          time.Now().UTC(),
		labels:      make(Labels),
		annotations: make(Annotations),
	}
}

// EventMeta is the metadata common to events.
type EventMeta struct {
	flag          Flag
	flagTextColor AnsiColor
	ts            time.Time
	headings      []string
	entity        string
	labels        Labels
	annotations   Annotations
}

// Headings returns the event meta headings.
func (em *EventMeta) Headings() []string { return em.headings }

// SetHeadings sets the event meta headings.
func (em *EventMeta) SetHeadings(headings ...string) { em.headings = headings }

// Flag returnst the event meta flag.
func (em *EventMeta) Flag() Flag { return em.flag }

// SetFlag sets the flag.
func (em *EventMeta) SetFlag(flag Flag) { em.flag = flag }

// FlagTextColor returns a custom color for the flag.
func (em *EventMeta) FlagTextColor() AnsiColor { return em.flagTextColor }

// SetFlagTextColor sets the flag text color.
func (em *EventMeta) SetFlagTextColor(color AnsiColor) { em.flagTextColor = color }

// Timestamp returnst the event meta timestamp.
func (em *EventMeta) Timestamp() time.Time { return em.ts }

// SetTimestamp sets the timestamp.
func (em *EventMeta) SetTimestamp(ts time.Time) { em.ts = ts }

// AddLabelValue adds a label value
func (em *EventMeta) AddLabelValue(key, value string) { em.labels[key] = value }

// SetLabels sets the labels collection.
func (em *EventMeta) SetLabels(labels Labels) { em.labels = labels }

// Labels returns the event labels.
func (em *EventMeta) Labels() Labels { return em.labels }

// AddAnnotationValue adds an annotation value
func (em *EventMeta) AddAnnotationValue(key, value string) { em.annotations[key] = value }

// SetAnnotations sets the annotations collection.
func (em *EventMeta) SetAnnotations(annotations Annotations) { em.annotations = annotations }

// Annotations returns the event annotations.
func (em *EventMeta) Annotations() Annotations { return em.annotations }

// SetEntity sets the entity value.
func (em *EventMeta) SetEntity(value string) { em.entity = value }

// Entity returns an entity value.
func (em *EventMeta) Entity() string { return em.entity }
