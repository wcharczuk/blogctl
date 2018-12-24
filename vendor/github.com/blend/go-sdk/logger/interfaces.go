package logger

import (
	"io"
	"time"
)

// Event is an interface representing methods necessary to trigger listeners.
type Event interface {
	Flag() Flag
	Timestamp() time.Time
}

// EventMetaProvider provides the full suite of event meta.
type EventMetaProvider interface {
	Event
	EventEntity
	EventHeadings
	EventLabels
	EventAnnotations
}

// Listener is a function that can be triggered by events.
type Listener func(e Event)

// EventEntity is a type that provides an entity value.
type EventEntity interface {
	SetEntity(string)
	Entity() string
}

// EventHeadings determines if we should add another output field, `event-headings` to output.
type EventHeadings interface {
	SetHeadings(...string)
	Headings() []string
}

// EventLabels is a type that provides labels.
type EventLabels interface {
	Labels() map[string]string
}

// EventAnnotations is a type that provides annotations.
type EventAnnotations interface {
	Annotations() map[string]string
}

// EventEnabled determines if we should allow an event to be triggered or not.
type EventEnabled interface {
	IsEnabled() bool
}

// EventWritable lets us disable implicit writing for some events.
type EventWritable interface {
	IsWritable() bool
}

// EventError determines if we should write the event to the error stream.
type EventError interface {
	IsError() bool
}

// Listenable is an interface.
type Listenable interface {
	Listen(Flag, string, Listener)
}

// Triggerable is an interface.
type Triggerable interface {
	Trigger(Event)
}

// SyncTriggerable is an interface.
type SyncTriggerable interface {
	SyncTrigger(Event)
}

// OutputReceiver is an interface
type OutputReceiver interface {
	Infof(string, ...Any)
	Sillyf(string, ...Any)
	Debugf(string, ...Any)
}

// SyncOutputReceiver is an interface
type SyncOutputReceiver interface {
	SyncInfof(string, ...Any)
	SyncSillyf(string, ...Any)
	SyncDebugf(string, ...Any)
}

// ErrorOutputReceiver is an interface
type ErrorOutputReceiver interface {
	Warningf(string, ...Any)
	Errorf(string, ...Any)
	Fatalf(string, ...Any)
}

// SyncErrorOutputReceiver is an interface
type SyncErrorOutputReceiver interface {
	SyncWarningf(string, ...Any)
	SyncErrorf(string, ...Any)
	SyncFatalf(string, ...Any)
}

// ErrorReceiver is an interface
type ErrorReceiver interface {
	Warning(error)
	Error(error)
	Fatal(error)
}

// SyncErrorReceiver is an interface
type SyncErrorReceiver interface {
	SyncWarning(error)
	SyncError(error)
	SyncFatal(error)
}

// SyncLogger is a logger that implements syncronous methods.
type SyncLogger interface {
	Listenable
	SyncTriggerable
	SyncOutputReceiver
	SyncErrorOutputReceiver
	SyncErrorReceiver
}

// AsyncLogger is a logger that implements async methods.
type AsyncLogger interface {
	Listenable
	Triggerable
	OutputReceiver
	ErrorOutputReceiver
	ErrorReceiver
}

// FullReceiver is every possible receiving / output interface.
type FullReceiver interface {
	SyncTriggerable
	SyncOutputReceiver
	SyncErrorOutputReceiver
	SyncErrorReceiver
	Triggerable
	OutputReceiver
	ErrorOutputReceiver
	ErrorReceiver
}

// FullLogger is every possible interface.
type FullLogger interface {
	Listenable
	SyncTriggerable
	SyncOutputReceiver
	SyncErrorOutputReceiver
	SyncErrorReceiver
	Triggerable
	OutputReceiver
	ErrorOutputReceiver
	ErrorReceiver
}

// Writer is a type that can consume events.
type Writer interface {
	Write(Event) error
	WriteError(Event) error
	Output() io.Writer
	ErrorOutput() io.Writer
	OutputFormat() OutputFormat
}

// --------------------------------------------------------------------------------
// testing helpers
// --------------------------------------------------------------------------------

// MarshalEvent marshals an object as a logger event.
func MarshalEvent(obj interface{}) (Event, bool) {
	typed, isTyped := obj.(Event)
	return typed, isTyped
}

// MarshalEventHeadings marshals an object as an event heading provider.
func MarshalEventHeadings(obj interface{}) (EventHeadings, bool) {
	typed, isTyped := obj.(EventHeadings)
	return typed, isTyped
}

// MarshalEventEnabled marshals an object as an event enabled provider.
func MarshalEventEnabled(obj interface{}) (EventEnabled, bool) {
	typed, isTyped := obj.(EventEnabled)
	return typed, isTyped
}

// MarshalEventWritable marshals an object as an event writable provider.
func MarshalEventWritable(obj interface{}) (EventWritable, bool) {
	typed, isTyped := obj.(EventWritable)
	return typed, isTyped
}

// MarshalEventMetaProvider marshals an object as an event meta provider.
func MarshalEventMetaProvider(obj interface{}) (EventMetaProvider, bool) {
	typed, isTyped := obj.(EventMetaProvider)
	return typed, isTyped
}
