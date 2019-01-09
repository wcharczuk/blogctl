package logger

import (
	"net/http"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
)

const (

	// DefaultListenerName is a default.
	DefaultListenerName = "default"

	// DefaultRecoverPanics is a default.
	DefaultRecoverPanics = true
)

var (
	// DefaultListenerWorkers is the default number of workers per listener.
	DefaultListenerWorkers = runtime.NumCPU()
)

// FatalExit creates a logger and calls `SyncFatalExit` on it.
func FatalExit(err error) {
	if err == nil {
		return
	}
	log := Sync()
	log.Enable(Fatal)
	log.SyncFatalExit(err)
}

// New returns a new logger with a given set of enabled flags, without a writer provisioned.
func New(flags ...Flag) *Logger {
	l := &Logger{
		recoverPanics: DefaultRecoverPanics,
		flags:         NewFlagSet(flags...),
	}
	l.writeWorker = NewWorker(l, l.Write, DefaultWriteQueueDepth)
	l.writeWorker.Start()
	return l
}

// NewFromConfig returns a new logger from a config.
func NewFromConfig(cfg *Config) *Logger {
	l := &Logger{
		recoverPanics:            DefaultRecoverPanics,
		flags:                    NewFlagSetFromValues(cfg.GetFlags()...),
		writeWorkerQueueDepth:    cfg.GetWriteQueueDepth(),
		listenerWorkerQueueDepth: cfg.GetListenerQueueDepth(),
	}
	l.writeWorker = NewWorker(l, l.Write, l.writeWorkerQueueDepth)
	l.writeWorker.Start()
	return l.
		WithHeading(cfg.GetHeading()).
		WithRecoverPanics(cfg.GetRecoverPanics()).
		WithHiddenFlags(NewFlagSetFromValues(cfg.GetHiddenFlags()...)).
		WithWriters(cfg.GetWriters()...)

}

// NewFromEnv returns a new agent with settings read from the environment,
// including the underlying writer.
func NewFromEnv() (*Logger, error) {
	cfg, err := NewConfigFromEnv()
	if err != nil {
		return nil, err
	}
	return NewFromConfig(cfg), nil
}

// MustNewFromEnv returns a new logger based on settings from the environment.
func MustNewFromEnv() *Logger {
	cfg, err := NewConfigFromEnv()
	if err != nil {
		panic(err)
	}
	return NewFromConfig(cfg)
}

// All returns a valid logger that fires any and all events, and includes a writer.
// It is effectively an alias to:
// 	New().WithFlags(NewFlagSetAll()).WithWriter(NewWriterFromEnv())
func All() *Logger {
	return New().WithFlags(NewFlagSetAll()).WithWriter(NewWriterFromEnv())
}

// None returns a valid agent that won't fire any events.
func None() *Logger {
	return &Logger{
		flags: NewFlagSetNone(),
	}
}

// NewText returns a new text logger.
// It is a shortcut for `New().WithWriter(NewTextWriterFromEnv())`
func NewText() *Logger {
	return New().WithFlagsFromEnv().WithWriter(NewTextWriterFromEnv())
}

// NewJSON returns a new json logger.
// It is a shortcut for `New().WithWriter(NewJSONWriterFromEnv())`
func NewJSON() *Logger {
	return New().WithFlagsFromEnv().WithWriter(NewJSONWriterFromEnv())
}

// Sync returns a valid agent that only processes events synchronously.
func Sync() *Logger {
	return &Logger{
		recoverPanics: DefaultRecoverPanics,
		flags:         NewFlagSetFromEnv(),
		writers:       []Writer{NewWriterFromEnv()},
	}
}

// Assertions
var (
	_ FullLogger = (*Logger)(nil)
)

// Logger is a handler for various logging events with descendent handlers.
type Logger struct {
	writers []Writer

	heading                  string
	writeWorkerQueueDepth    int
	listenerWorkerQueueDepth int

	state int32

	flagsLock sync.Mutex
	flags     *FlagSet

	hiddenFlagsLock sync.Mutex
	hiddenFlags     *FlagSet

	workersLock sync.Mutex
	workers     map[Flag]map[string]*Worker

	writeWorkerLock sync.Mutex
	writeWorker     *Worker

	recoverPanics bool
}

// CanStart returns if the latch can start.
func (l *Logger) CanStart() bool {
	return atomic.LoadInt32(&l.state) == LoggerStopped
}

// CanStop returns if the latch can stop.
func (l *Logger) CanStop() bool {
	return atomic.LoadInt32(&l.state) == LoggerStarted
}

// WithHeading returns the logger heading.
func (l *Logger) WithHeading(heading string) *Logger {
	l.heading = heading
	return l
}

// Heading returns the logger heading.
func (l *Logger) Heading() string {
	return l.heading
}

// WithWriteWorkerQueueDepth sets the worker queue depth.
func (l *Logger) WithWriteWorkerQueueDepth(queueDepth int) *Logger {
	l.writeWorkerQueueDepth = queueDepth
	return l
}

// WriteWorkerQueueDepth returns the worker queue depth.
func (l *Logger) WriteWorkerQueueDepth() int {
	return l.writeWorkerQueueDepth
}

// WithListenerWorkerQueueDepth sets the worker queue depth.
func (l *Logger) WithListenerWorkerQueueDepth(queueDepth int) *Logger {
	l.listenerWorkerQueueDepth = queueDepth
	return l
}

// ListenerWorkerQueueDepth returns the worker queue depth.
func (l *Logger) ListenerWorkerQueueDepth() int {
	return l.listenerWorkerQueueDepth
}

// Writers returns the output writers for events.
func (l *Logger) Writers() []Writer {
	return l.writers
}

// WithWriters sets the logger writers, overwriting any existing writers.
func (l *Logger) WithWriters(writers ...Writer) *Logger {
	l.writers = writers
	return l
}

// WithWriter adds a logger writer.
func (l *Logger) WithWriter(writer Writer) *Logger {
	l.writers = append(l.writers, writer)
	return l
}

// RecoversPanics returns if we should recover panics in logger listeners.
func (l *Logger) RecoversPanics() bool {
	return l.recoverPanics
}

// WithRecoverPanics sets the recoverPanics sets if the logger should trap panics in event handlers.
func (l *Logger) WithRecoverPanics(value bool) *Logger {
	l.recoverPanics = value
	return l
}

// Flags returns the logger flag set.
func (l *Logger) Flags() *FlagSet {
	return l.flags
}

// WithFlags sets the logger flag set.
func (l *Logger) WithFlags(flags *FlagSet) *Logger {
	l.flagsLock.Lock()
	defer l.flagsLock.Unlock()
	l.flags = flags
	return l
}

// WithHiddenFlags sets the hidden flag set.
// These flags mark events as to be omitted from output.
func (l *Logger) WithHiddenFlags(flags *FlagSet) *Logger {
	l.hiddenFlagsLock.Lock()
	defer l.hiddenFlagsLock.Unlock()
	l.hiddenFlags = flags
	return l
}

// WithFlagsFromEnv adds flags from the environment.
func (l *Logger) WithFlagsFromEnv() *Logger {
	l.flagsLock.Lock()
	defer l.flagsLock.Unlock()

	l.hiddenFlagsLock.Lock()
	defer l.hiddenFlagsLock.Unlock()

	if l.flags != nil {
		l.flags.CoalesceWith(NewFlagSetFromEnv())
	} else {
		l.flags = NewFlagSetFromEnv()
	}

	if l.hiddenFlags != nil {
		l.hiddenFlags.CoalesceWith(NewHiddenFlagSetFromEnv())
	} else {
		l.hiddenFlags = NewHiddenFlagSetFromEnv()
	}
	return l
}

// WithEnabled flips the bit flag for a given set of events.
func (l *Logger) WithEnabled(flags ...Flag) *Logger {
	l.Enable(flags...)
	return l
}

// Enable flips the bit flag for a given set of events.
func (l *Logger) Enable(flags ...Flag) {
	l.flagsLock.Lock()
	defer l.flagsLock.Unlock()

	if l.flags != nil {
		for _, flag := range flags {
			l.flags.Enable(flag)
		}
	} else {
		l.flags = NewFlagSet(flags...)
	}
}

// WithDisabled flips the bit flag for a given set of events.
func (l *Logger) WithDisabled(flags ...Flag) *Logger {
	l.Disable(flags...)
	return l
}

// Disable flips the bit flag for a given set of events.
func (l *Logger) Disable(flags ...Flag) {
	l.flagsLock.Lock()
	defer l.flagsLock.Unlock()
	for _, flag := range flags {
		l.flags.Disable(flag)
	}
}

// IsEnabled asserts if a flag value is set or not.
func (l *Logger) IsEnabled(flag Flag) (enabled bool) {
	l.flagsLock.Lock()
	if l.flags == nil {
		enabled = false
		l.flagsLock.Unlock()
		return
	}
	enabled = l.flags.IsEnabled(flag)
	l.flagsLock.Unlock()
	return
}

// Hide disallows automatic logging for each event emitted under the provided list of flags.
func (l *Logger) Hide(flags ...Flag) {
	l.hiddenFlagsLock.Lock()
	defer l.hiddenFlagsLock.Unlock()

	if l.hiddenFlags != nil {
		for _, flag := range flags {
			l.hiddenFlags.Enable(flag)
		}
	} else {
		l.hiddenFlags = NewFlagSet(flags...)
	}
}

// IsHidden asserts if a flag is hidden or not.
func (l *Logger) IsHidden(flag Flag) bool {
	if l.hiddenFlags == nil {
		return false
	}
	return l.hiddenFlags.IsEnabled(flag)
}

// WithHidden hides a set of flags and returns logger
func (l *Logger) WithHidden(flags ...Flag) *Logger {
	l.Hide(flags...)
	return l
}

// Show allows automatic logging for each event emitted under the provided list of flags.
func (l *Logger) Show(flags ...Flag) {
	l.hiddenFlagsLock.Lock()
	defer l.hiddenFlagsLock.Unlock()
	for _, flag := range flags {
		l.hiddenFlags.Disable(flag)
	}
}

// HasListeners returns if there are registered listener for an event.
func (l *Logger) HasListeners(flag Flag) bool {
	l.workersLock.Lock()
	defer l.workersLock.Unlock()

	if l.workers == nil {
		return false
	}
	workers, hasWorkers := l.workers[flag]
	if !hasWorkers {
		return false
	}
	return len(workers) > 0
}

// HasListener returns if a specific listener is registerd for a flag.
func (l *Logger) HasListener(flag Flag, listenerName string) bool {
	l.workersLock.Lock()
	defer l.workersLock.Unlock()

	if l.workers == nil {
		return false
	}
	workers, hasWorkers := l.workers[flag]
	if !hasWorkers {
		return false
	}
	_, hasWorker := workers[listenerName]
	return hasWorker
}

// SubContext returns a sub context for the logger.
// It provides a more limited api surface area but lets you
// decoarate events with context specific headings, labels and annotations.
func (l *Logger) SubContext(heading string) *SubContext {
	return &SubContext{
		log:      l,
		headings: []string{heading},
	}
}

// Listen adds a listener for a given flag.
func (l *Logger) Listen(flag Flag, listenerName string, listener Listener) {
	l.workersLock.Lock()
	defer l.workersLock.Unlock()

	if l.workers == nil {
		l.workers = map[Flag]map[string]*Worker{}
	}

	w := NewWorker(l, listener, l.listenerWorkerQueueDepth)
	if listeners, hasListeners := l.workers[flag]; hasListeners {
		listeners[listenerName] = w
	} else {
		l.workers[flag] = map[string]*Worker{
			listenerName: w,
		}
	}
	w.Start()
}

// RemoveListeners clears *all* listeners for a Flag.
func (l *Logger) RemoveListeners(flag Flag) {
	l.workersLock.Lock()
	defer l.workersLock.Unlock()

	if l.workers == nil {
		return
	}

	listeners, hasListeners := l.workers[flag]
	if !hasListeners {
		return
	}

	for _, w := range listeners {
		w.Close()
	}

	delete(l.workers, flag)
}

// RemoveListener clears a specific listener for a Flag.
func (l *Logger) RemoveListener(flag Flag, listenerName string) {
	l.workersLock.Lock()
	defer l.workersLock.Unlock()

	if l.workers == nil {
		return
	}

	listeners, hasListeners := l.workers[flag]
	if !hasListeners {
		return
	}

	worker, hasWorker := listeners[listenerName]
	if !hasWorker {
		return
	}

	worker.Close()
	delete(listeners, listenerName)

	if len(listeners) == 0 {
		delete(l.workers, flag)
	}
}

// Trigger fires the listeners for a given event asynchronously.
// The invocations will be queued in a work queue and processed by a fixed worker count.
// There are no order guarantees on when these events will be processed.
// This call will not block on the event listeners.
func (l *Logger) Trigger(e Event) {
	l.trigger(true, e)
}

// SyncTrigger fires the listeners for a given event synchronously.
// The invocations will be triggered immediately, blocking the call.
func (l *Logger) SyncTrigger(e Event) {
	l.trigger(false, e)
}

func (l *Logger) trigger(async bool, e Event) {
	if !async && l.recoverPanics {
		defer func() {
			if r := recover(); r != nil {
				l.Write(Errorf(Fatal, "%+v", r))
			}
		}()
	}

	if async && !l.isStarted() {
		return
	}

	if typed, isTyped := e.(EventEnabled); isTyped && !typed.IsEnabled() {
		return
	}

	flag := e.Flag()
	if l.IsEnabled(flag) {
		if l.heading != "" {
			if typed, isTyped := e.(EventHeadings); isTyped {
				if len(typed.Headings()) > 0 {
					typed.SetHeadings(append([]string{l.heading}, typed.Headings()...)...)
				} else {
					typed.SetHeadings(l.heading)
				}
			}
		}

		var workers map[string]*Worker
		l.workersLock.Lock()
		if l.workers != nil {
			if flagWorkers, hasWorkers := l.workers[flag]; hasWorkers {
				workers = flagWorkers
			}
		}
		l.workersLock.Unlock()

		for _, worker := range workers {
			if async {
				worker.Work <- e
			} else {
				worker.Listener(e)
			}
		}

		// check if the flag is globally hidden from output.
		if l.IsHidden(flag) {
			return
		}

		// check if the event controls if it should be written or not.
		if typed, isTyped := e.(EventWritable); isTyped && !typed.IsWritable() {
			return
		}

		if async {
			l.writeWorker.Work <- e
		} else {
			l.Write(e)
		}
	}
}

// --------------------------------------------------------------------------------
// Builtin Flag Handlers (infof, debugf etc.)
// --------------------------------------------------------------------------------

// Sillyf logs an incredibly verbose message to the output stream.
func (l *Logger) Sillyf(format string, args ...interface{}) {
	l.trigger(true, Messagef(Silly, format, args...))
}

// SyncSillyf logs an incredibly verbose message to the output stream synchronously.
func (l *Logger) SyncSillyf(format string, args ...interface{}) {
	l.trigger(false, Messagef(Silly, format, args...))
}

// Infof logs an informational message to the output stream.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.trigger(true, Messagef(Info, format, args...))
}

// SyncInfof logs an informational message to the output stream synchronously.
func (l *Logger) SyncInfof(format string, args ...interface{}) {
	l.trigger(false, Messagef(Info, format, args...))
}

// Debugf logs a debug message to the output stream.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.trigger(true, Messagef(Debug, format, args...))
}

// SyncDebugf logs an debug message to the output stream synchronously.
func (l *Logger) SyncDebugf(format string, args ...interface{}) {
	l.trigger(false, Messagef(Debug, format, args...))
}

// Warningf logs a debug message to the output stream.
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.trigger(false, Errorf(Warning, format, args...))
}

// SyncWarningf logs an warning message to the output stream synchronously.
func (l *Logger) SyncWarningf(format string, args ...interface{}) {
	l.trigger(false, Errorf(Warning, format, args...))
}

// Warning logs a warning error to std err.
func (l *Logger) Warning(err error) error {
	l.trigger(true, NewErrorEvent(Warning, err))
	return err
}

// SyncWarning synchronously logs a warning to std err.
func (l *Logger) SyncWarning(err error) error {
	l.trigger(false, NewErrorEvent(Warning, err))
	return err
}

// WarningWithReq logs a warning error to std err with a request.
func (l *Logger) WarningWithReq(err error, req *http.Request) error {
	l.trigger(true, NewErrorEventWithState(Warning, err, req))
	return err
}

// Errorf writes an event to the log and triggers event listeners.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.trigger(true, Errorf(Error, format, args...))
}

// SyncErrorf synchronously triggers a error.
func (l *Logger) SyncErrorf(format string, args ...interface{}) {
	l.trigger(false, Errorf(Error, format, args...))
}

// Error logs an error to std err.
func (l *Logger) Error(err error) error {
	l.trigger(true, NewErrorEvent(Error, err))
	return err
}

// SyncError synchronously logs an error to std err.
func (l *Logger) SyncError(err error) error {
	l.trigger(false, NewErrorEvent(Error, err))
	return err
}

// ErrorWithReq logs an error to std err with a request.
func (l *Logger) ErrorWithReq(err error, req *http.Request) error {
	l.trigger(true, NewErrorEventWithState(Error, err, req))
	return err
}

// Fatalf writes an event to the log and triggers event listeners.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.trigger(true, Errorf(Fatal, format, args...))
}

// SyncFatalf synchronously triggers a fatal.
func (l *Logger) SyncFatalf(format string, args ...interface{}) {
	l.trigger(false, Errorf(Fatal, format, args...))
}

// Fatal logs the result of a panic to std err.
func (l *Logger) Fatal(err error) error {
	l.trigger(true, NewErrorEvent(Fatal, err))
	return err
}

// SyncFatal synchronously logs a fatal to std err.
func (l *Logger) SyncFatal(err error) error {
	l.trigger(false, NewErrorEvent(Fatal, err))
	return err
}

// SyncFatalExit logs the result of a fatal error to std err and calls `exit(1)`
func (l *Logger) SyncFatalExit(err error) {
	l.SyncFatal(err)
	l.Drain()
	os.Exit(1)
}

// Write writes an event synchronously to the writer either as a normal even or as an error.
func (l *Logger) Write(e Event) {
	ll := len(l.writers)
	if typed, isTyped := e.(EventError); isTyped && typed.IsError() {
		for index := 0; index < ll; index++ {
			l.writers[index].WriteError(e)
		}
		return
	}
	for index := 0; index < ll; index++ {
		l.writers[index].Write(e)
	}
}

// --------------------------------------------------------------------------------
// finalizers
// --------------------------------------------------------------------------------

// Close releases shared resources for the agent.
func (l *Logger) Close() (err error) {
	l.flagsLock.Lock()
	defer l.flagsLock.Unlock()

	if l.flags != nil {
		l.flags.SetNone()
	}

	l.setStopping()

	l.workersLock.Lock()
	defer l.workersLock.Unlock()

	for _, workers := range l.workers {
		for _, worker := range workers {
			worker.Close()
		}
	}

	for key := range l.workers {
		delete(l.workers, key)
	}
	l.workers = nil

	l.writeWorkerLock.Lock()
	defer l.writeWorkerLock.Unlock()

	l.writeWorker.Close()
	l.writeWorker = nil

	l.setStopped()

	return nil
}

// Drain waits for the agent to finish its queue of events before closing.
func (l *Logger) Drain() error {
	l.workersLock.Lock()
	defer l.workersLock.Unlock()

	l.setStopping()

	for _, workers := range l.workers {
		for _, worker := range workers {
			worker.Drain()
		}
	}

	l.writeWorkerLock.Lock()
	defer l.writeWorkerLock.Unlock()

	if l.writeWorker != nil {
		l.writeWorker.Drain()
	}

	l.setStarted()

	return nil
}

func (l *Logger) isStarted() bool {
	return atomic.LoadInt32(&l.state) == LoggerStarted
}

func (l *Logger) isStopping() bool {
	return atomic.LoadInt32(&l.state) == LoggerStopping
}

func (l *Logger) isStopped() bool {
	return atomic.LoadInt32(&l.state) == LoggerStopped
}

func (l *Logger) setStarted() {
	atomic.StoreInt32(&l.state, LoggerStarted)
}

func (l *Logger) setStopping() {
	atomic.StoreInt32(&l.state, LoggerStopping)
}

func (l *Logger) setStopped() {
	atomic.StoreInt32(&l.state, LoggerStopped)
}
