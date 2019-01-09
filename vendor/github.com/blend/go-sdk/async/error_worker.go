package async

import (
	"context"

	"github.com/blend/go-sdk/exception"
)

// NewErrorWorker creates a new error worker.
func NewErrorWorker(action ErrorAction) *ErrorWorker {
	return &ErrorWorker{
		latch:  &Latch{},
		action: action,
		work:   make(chan error),
	}
}

// ErrorWorker is a worker that is pushed work as errors over a channel.
type ErrorWorker struct {
	latch    *Latch
	action   ErrorAction
	work     chan error
	fallback func(error)
}

// Latch returns the worker latch.
func (ew *ErrorWorker) Latch() *Latch {
	return ew.latch
}

// WithWork sets the work channel.
// It allows you to override the default (non-buffered) channel with
// a buffer of your chosing.
func (ew *ErrorWorker) WithWork(work chan error) *ErrorWorker {
	ew.work = work
	return ew
}

// Work returns the work channel.
func (ew *ErrorWorker) Work() chan error {
	return ew.work
}

// WithFallback sets the fallback collector.
func (ew *ErrorWorker) WithFallback(action func(error)) *ErrorWorker {
	ew.fallback = action
	return ew
}

// Enqueue adds an item to the error work queue.
func (ew *ErrorWorker) Enqueue(obj error) {
	ew.work <- obj
}

// Start starts the worker.
func (ew *ErrorWorker) Start() {
	ew.StartContext(context.Background())
}

// StartContext starts the worker with a given context.
func (ew *ErrorWorker) StartContext(ctx context.Context) {
	ew.latch.Starting()
	go ew.Dispatch(ctx)
	<-ew.latch.NotifyStarted()
}

// Dispatch starts the listen loop for work.
func (ew *ErrorWorker) Dispatch(ctx context.Context) {
	ew.latch.Started()
	var workItem error
	for {
		select {
		case workItem = <-ew.work:
			ew.Execute(ctx, workItem)
		case <-ew.latch.NotifyStopping():
			ew.latch.Stopped()
			return
		}
	}
}

// Execute invokes the action and recovers panics.
func (ew *ErrorWorker) Execute(ctx context.Context, workItem error) {
	defer func() {
		if r := recover(); r != nil {
			if ew.fallback != nil {
				ew.fallback(exception.New(r))
			}
		}
	}()
	if err := ew.action(ctx, workItem); err != nil {
		if ew.fallback != nil {
			ew.fallback(err)
		}
	}
}

// Stop stop the worker.
// The work left in the queue will remain.
func (ew *ErrorWorker) Stop() {
	ew.latch.Stopping()
	<-ew.latch.NotifyStopped()
}

// Drain stops the worker and synchronously finishes work.
func (ew *ErrorWorker) Drain() {
	ew.DrainContext(context.Background())
}

// DrainContext stops the worker and synchronously drains the the remaining work
// with a given context.
func (ew *ErrorWorker) DrainContext(ctx context.Context) {
	ew.latch.Stopping()
	<-ew.latch.NotifyStopped()
	remaining := len(ew.work)
	stopped := make(chan struct{})
	go func() {
		defer func() {
			close(stopped)
		}()
		for x := 0; x < remaining; x++ {
			ew.Execute(ctx, <-ew.work)
		}
	}()
	<-stopped
}

// Close stops the worker.
func (ew *ErrorWorker) Close() error {
	ew.latch.Stopping()
	<-ew.latch.NotifyStopped()
	ew.work = nil
	return nil
}
