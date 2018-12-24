package async

import (
	"sync"
	"sync/atomic"
)

const (
	// LatchStopped is a latch lifecycle state.
	LatchStopped int32 = 0
	// LatchStarting is a latch lifecycle state.
	LatchStarting int32 = 1
	// LatchRunning is a latch lifecycle state.
	LatchRunning int32 = 2
	// LatchStopping is a latch lifecycle state.
	LatchStopping int32 = 3
)

// NewLatch returns a new latch.
func NewLatch() *Latch {
	return &Latch{
		starting: make(chan struct{}),
		started:  make(chan struct{}),
		stopping: make(chan struct{}),
		stopped:  make(chan struct{}),
	}
}

// Latch is a helper to coordinate goroutine lifecycles.
// The lifecycle is generally as follows.
// 0 - stopped / idle
// 1 - starting
// 2 - running
// 3 - stopping
// goto 0
// Each state includes a transition notification, i.e. `Starting()` triggers `NotifyStarting`
type Latch struct {
	sync.Mutex
	state int32

	starting chan struct{}
	started  chan struct{}
	stopping chan struct{}
	stopped  chan struct{}
}

// Reset resets the latch.
func (l *Latch) Reset() {
	l.Lock()
	l.state = LatchStopped
	l.starting = make(chan struct{})
	l.started = make(chan struct{})
	l.stopping = make(chan struct{})
	l.stopped = make(chan struct{})
	l.Unlock()
}

// CanStart returns if the latch can start.
func (l *Latch) CanStart() bool {
	return atomic.LoadInt32(&l.state) == LatchStopped
}

// CanStop returns if the latch can stop.
func (l *Latch) CanStop() bool {
	return atomic.LoadInt32(&l.state) == LatchRunning
}

// IsStopping returns if the latch is waiting to finish stopping.
func (l *Latch) IsStopping() bool {
	return atomic.LoadInt32(&l.state) == LatchStopping
}

// IsStopped returns if the latch is stopped.
func (l *Latch) IsStopped() (isStopped bool) {
	return atomic.LoadInt32(&l.state) == LatchStopped
}

// IsStarting indicates the latch is waiting to be scheduled.
func (l *Latch) IsStarting() bool {
	return atomic.LoadInt32(&l.state) == LatchStarting
}

// IsRunning indicates we can signal to stop.
func (l *Latch) IsRunning() bool {
	return atomic.LoadInt32(&l.state) == LatchRunning
}

// NotifyStarting returns the started signal.
// It is used to coordinate the transition from stopped -> starting.
func (l *Latch) NotifyStarting() (notifyStarting <-chan struct{}) {
	l.Lock()
	notifyStarting = l.starting
	l.Unlock()
	return
}

// NotifyStarted returns the started signal.
// It is used to coordinate the transition from starting -> started.
func (l *Latch) NotifyStarted() (notifyStarted <-chan struct{}) {
	l.Lock()
	notifyStarted = l.started
	l.Unlock()
	return
}

// NotifyStopping returns the should stop signal.
// It is used to trigger the transition from running -> stopping -> stopped.
func (l *Latch) NotifyStopping() (notifyStopping <-chan struct{}) {
	l.Lock()
	notifyStopping = l.stopping
	l.Unlock()
	return
}

// NotifyStopped returns the stopped signal.
// It is used to coordinate the transition from stopping -> stopped.
func (l *Latch) NotifyStopped() (notifyStopped <-chan struct{}) {
	l.Lock()
	notifyStopped = l.stopped
	l.Unlock()
	return
}

// Starting signals the latch is starting.
// This is typically done before you kick off a goroutine.
func (l *Latch) Starting() {
	if l.IsStarting() {
		return
	}
	l.Lock()
	atomic.StoreInt32(&l.state, LatchStarting)
	if l.starting != nil {
		close(l.starting)
	}
	l.started = make(chan struct{})
	l.Unlock()
}

// Started signals that the latch is started and has entered
// the `IsRunning` state.
func (l *Latch) Started() {
	if l.IsRunning() {
		return
	}
	l.Lock()
	atomic.StoreInt32(&l.state, LatchRunning)
	if l.started != nil {
		close(l.started)
	}
	l.stopping = make(chan struct{})
	l.Unlock()
}

// Stopping signals the latch to stop.
// It could also be thought of as `SignalStopping`.
func (l *Latch) Stopping() {
	if l.IsStopping() {
		return
	}
	l.Lock()
	atomic.StoreInt32(&l.state, LatchStopping)
	if l.stopping != nil {
		close(l.stopping)
	}
	l.stopped = make(chan struct{})
	l.Unlock()
}

// Stopped signals the latch has stopped.
func (l *Latch) Stopped() {
	if l.IsStopped() {
		return
	}
	l.Lock()
	atomic.StoreInt32(&l.state, LatchStopped)
	if l.stopped != nil {
		close(l.stopped)
	}
	l.starting = make(chan struct{})
	l.Unlock()
}
