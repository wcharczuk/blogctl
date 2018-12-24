package async

import (
	"time"

	"github.com/blend/go-sdk/exception"
)

// NewInterval returns a new worker that runs an action on an interval.
func NewInterval(action func() error, interval time.Duration) *Interval {
	return &Interval{
		interval: interval,
		action:   action,
		latch:    &Latch{},
	}
}

// Interval is a managed goroutine that does things.
type Interval struct {
	delay    time.Duration
	interval time.Duration
	action   func() error
	latch    *Latch
	errors   chan error
}

// WithDelay sets a start delay time.
func (i *Interval) WithDelay(d time.Duration) *Interval {
	i.delay = d
	return i
}

// Delay returns the start delay.
func (i *Interval) Delay() time.Duration {
	return i.delay
}

// WithInterval sets the inteval. It must be set before `.Start()` is called.
func (i *Interval) WithInterval(d time.Duration) *Interval {
	i.interval = d
	return i
}

// Interval returns the interval for the ticker.
func (i Interval) Interval() time.Duration {
	return i.interval
}

// IsRunning returns if the worker is running.
func (i *Interval) IsRunning() bool {
	return i.latch.IsRunning()
}

// Latch returns the inteval worker latch.
func (i *Interval) Latch() *Latch {
	return i.latch
}

// WithAction sets the interval action.
func (i *Interval) WithAction(action func() error) *Interval {
	i.action = action
	return i
}

// Action returns the interval action.
func (i *Interval) Action() func() error {
	return i.action
}

// WithErrors returns the error channel.
func (i *Interval) WithErrors(errors chan error) *Interval {
	i.errors = errors
	return i
}

// Errors returns a channel to read action errors from.
func (i *Interval) Errors() chan error {
	return i.errors
}

// NotifyStarted returns the notify started signal.
func (i *Interval) NotifyStarted() <-chan struct{} {
	return i.latch.NotifyStarted()
}

// NotifyStopped returns the notify stopped signal.
func (i *Interval) NotifyStopped() <-chan struct{} {
	return i.latch.NotifyStopped()
}

// Start starts the worker.
func (i *Interval) Start() error {
	if !i.latch.CanStart() {
		return exception.New(ErrCannotStart)
	}

	i.latch.Starting()
	go func() {
		i.latch.Started()

		if i.delay > 0 {
			time.Sleep(i.delay)
		}

		tick := time.Tick(i.interval)
		var err error
		for {
			select {
			case <-tick:
				err = i.action()
				if err != nil && i.errors != nil {
					i.errors <- err
				}
			case <-i.latch.NotifyStopping():
				i.latch.Stopped()
				return
			}
		}
	}()
	<-i.latch.NotifyStarted()
	return nil
}

// Stop stops the worker.
func (i *Interval) Stop() error {
	if !i.latch.CanStop() {
		return exception.New(ErrCannotStop)
	}
	i.latch.Stopping()
	<-i.latch.NotifyStopped()
	return nil
}
