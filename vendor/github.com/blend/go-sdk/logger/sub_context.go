package logger

var (
	// This is a compile time assertion `SubContext` implements `FullReceiver`.
	_ FullReceiver = &SubContext{}
)

// SubContext is a sub-reference to a logger with a specific heading and set of default labels for messages.
// It implements the full logger suite but forwards them up to the parent logger.
type SubContext struct {
	log         *Logger
	headings    []string
	labels      map[string]string
	annotations map[string]string
}

// Logger returns the underlying logger.
func (sc *SubContext) Logger() *Logger {
	return sc.log
}

// SubContext returns a further sub-context with a given heading.
func (sc *SubContext) SubContext(heading string) *SubContext {
	return &SubContext{
		headings:    append(sc.headings, heading),
		labels:      sc.labels,
		annotations: sc.annotations,
		log:         sc.log,
	}
}

// Headings returns the headings.
func (sc *SubContext) Headings() []string {
	return sc.headings
}

// WithLabel adds a label.
func (sc *SubContext) WithLabel(key, value string) *SubContext {
	if sc.labels == nil {
		sc.labels = map[string]string{}
	}
	sc.labels[key] = value
	return sc
}

// WithLabels sets the labels.
func (sc *SubContext) WithLabels(labels map[string]string) *SubContext {
	sc.labels = labels
	return sc
}

// Labels returns the sub-context labels.
func (sc *SubContext) Labels() map[string]string {
	return sc.labels
}

// WithAnnotations sets the annotations.
func (sc *SubContext) WithAnnotations(annotations map[string]string) *SubContext {
	sc.annotations = annotations
	return sc
}

// WithAnnotation adds an annotation.
func (sc *SubContext) WithAnnotation(key, value string) *SubContext {
	if sc.annotations == nil {
		sc.annotations = map[string]string{}
	}
	sc.annotations[key] = value
	return sc
}

// Annotations returns the sub-context annotations.
func (sc *SubContext) Annotations() map[string]string {
	return sc.annotations
}

// Sillyf writes a message.
func (sc *SubContext) Sillyf(format string, args ...Any) {
	msg := Messagef(Silly, format, args...)
	sc.injectHeadings(msg)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.Trigger(msg)
}

// SyncSillyf synchronously writes a message.
func (sc *SubContext) SyncSillyf(format string, args ...Any) {
	msg := Messagef(Silly, format, args...)
	sc.injectHeadings(msg)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.SyncTrigger(msg)
}

// Infof writes a message.
func (sc *SubContext) Infof(format string, args ...Any) {
	msg := Messagef(Info, format, args...)
	sc.injectHeadings(msg)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.Trigger(msg)
}

// SyncInfof synchronously writes a message.
func (sc *SubContext) SyncInfof(format string, args ...Any) {
	msg := Messagef(Info, format, args...)
	sc.injectHeadings(msg)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.SyncTrigger(msg)
}

// Debugf writes a message.
func (sc *SubContext) Debugf(format string, args ...Any) {
	msg := Messagef(Debug, format, args...)
	sc.injectHeadings(msg)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.Trigger(msg)
}

// SyncDebugf synchronously writes a message.
func (sc *SubContext) SyncDebugf(format string, args ...Any) {
	msg := Messagef(Debug, format, args...)
	sc.injectHeadings(msg)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.SyncTrigger(msg)
}

// Warningf writes an error message.
func (sc *SubContext) Warningf(format string, args ...Any) {
	msg := Errorf(Warning, format, args...)
	sc.injectHeadings(msg)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.Trigger(msg)
}

// Warning writes an error message.
func (sc *SubContext) Warning(err error) {
	msg := NewErrorEvent(Warning, err)
	sc.injectHeadings(msg)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.Trigger(msg)
}

// SyncWarningf synchronously writes an error message.
func (sc *SubContext) SyncWarningf(format string, args ...Any) {
	msg := Errorf(Warning, format, args...)
	sc.injectHeadings(msg)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.SyncTrigger(msg)
}

// SyncWarning writes a message.
func (sc *SubContext) SyncWarning(err error) {
	msg := NewErrorEvent(Warning, err)
	sc.injectHeadings(msg)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.SyncTrigger(msg)
}

// Errorf writes an error  message.
func (sc *SubContext) Errorf(format string, args ...Any) {
	msg := Errorf(Error, format, args...)
	sc.injectHeadings(msg)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.Trigger(msg)
}

// Error writes an error message.
func (sc *SubContext) Error(err error) {
	msg := NewErrorEvent(Error, err)
	sc.injectHeadings(msg)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.Trigger(msg)
}

// SyncErrorf synchronously writes an error message.
func (sc *SubContext) SyncErrorf(format string, args ...Any) {
	msg := Errorf(Error, format, args...)
	sc.injectHeadings(msg)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.SyncTrigger(msg)
}

// SyncError writes an error message.
func (sc *SubContext) SyncError(err error) {
	msg := NewErrorEvent(Error, err)
	sc.injectHeadings(msg)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.SyncTrigger(msg)
}

// Fatalf writes an error  message.
func (sc *SubContext) Fatalf(format string, args ...Any) {
	msg := Errorf(Fatal, format, args...)
	sc.injectHeadings(msg)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.Trigger(msg)
}

// Fatal writes an error message.
func (sc *SubContext) Fatal(err error) {
	msg := NewErrorEvent(Fatal, err)
	sc.injectHeadings(msg)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.Trigger(msg)
}

// SyncFatalf synchronously writes an error message.
func (sc *SubContext) SyncFatalf(format string, args ...Any) {
	msg := Errorf(Fatal, format, args...)
	sc.injectHeadings(msg)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.SyncTrigger(msg)
}

// SyncFatal writes an error message.
func (sc *SubContext) SyncFatal(err error) {
	msg := NewErrorEvent(Fatal, err)
	sc.injectHeadings(msg)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.SyncTrigger(msg)
}

// Trigger triggers listeners asynchronously.
func (sc *SubContext) Trigger(e Event) {
	sc.injectHeadings(e)
	sc.injectLabels(e)
	sc.injectAnnotations(e)
	sc.log.trigger(true, e)
}

// SyncTrigger triggers event listeners synchronously.
func (sc *SubContext) SyncTrigger(e Event) {
	sc.injectHeadings(e)
	sc.injectLabels(e)
	sc.injectAnnotations(e)
	sc.log.trigger(false, e)
}

// injectHeadings injects the sub-context's headings into an event if it supports headings.
func (sc *SubContext) injectHeadings(e Event) {
	if len(sc.headings) == 0 {
		return
	}
	if typed, isTyped := e.(EventHeadings); isTyped {
		typed.SetHeadings(append(sc.headings, typed.Headings()...)...)
	}
}

// injectHeadings injects the sub-context's labels into an event if it supports labels.
func (sc *SubContext) injectLabels(e Event) {
	if sc.labels == nil {
		return
	}
	if typed, isTyped := e.(EventLabels); isTyped {
		for key, value := range sc.labels {
			typed.Labels()[key] = value
		}
	}
}

// injectHeadings injects the sub-context's annotations into an event if it supports annotations.
func (sc *SubContext) injectAnnotations(e Event) {
	if sc.annotations == nil {
		return
	}
	if typed, isTyped := e.(EventAnnotations); isTyped {
		for key, value := range sc.annotations {
			typed.Annotations()[key] = value
		}
	}
}
