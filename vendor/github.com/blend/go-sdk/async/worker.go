package async

import (
	"context"

	"github.com/blend/go-sdk/exception"
)

// NewWorker creates a new worker.
func NewWorker(action QueueAction) *Worker {
	return &Worker{
		latch:  &Latch{},
		action: action,
		work:   make(chan interface{}),
	}
}

// Worker is a worker that is pushed work over a channel.
type Worker struct {
	latch  *Latch
	action QueueAction
	errors chan error
	work   chan interface{}
}

// Latch returns the worker latch.
func (w *Worker) Latch() *Latch {
	return w.latch
}

// WithWork sets the work channel.
// It allows you to override the default (non-buffered) channel with
// a buffer of your chosing.
func (w *Worker) WithWork(work chan interface{}) *Worker {
	w.work = work
	return w
}

// Work returns the work channel.
func (w *Worker) Work() chan interface{} {
	return w.work
}

// WithErrors returns the error channel.
func (w *Worker) WithErrors(errors chan error) *Worker {
	w.errors = errors
	return w
}

// Errors returns a channel to read action errors from.
func (w *Worker) Errors() chan error {
	return w.errors
}

// Enqueue adds an item to the work queue.
func (w *Worker) Enqueue(obj interface{}) {
	w.work <- obj
}

// Start starts the worker.
func (w *Worker) Start() {
	w.StartContext(context.Background())
}

// StartContext starts the worker with a given context.
func (w *Worker) StartContext(ctx context.Context) {
	w.latch.Starting()
	go w.Dispatch(ctx)
	<-w.latch.NotifyStarted()
}

// Dispatch starts the listen loop for work.
func (w *Worker) Dispatch(ctx context.Context) {
	w.latch.Started()
	var workItem interface{}
	for {
		select {
		case workItem = <-w.work:
			w.Execute(ctx, workItem)
		case <-w.latch.NotifyStopping():
			w.latch.Stopped()
			return
		}
	}
}

// Execute invokes the action and recovers panics.
func (w *Worker) Execute(ctx context.Context, workItem interface{}) {
	defer func() {
		if r := recover(); r != nil {
			if w.errors != nil {
				w.errors <- exception.New(r)
			}
		}
	}()
	if err := w.action(ctx, workItem); err != nil {
		if w.errors != nil {
			w.errors <- exception.New(err)
		}
	}
}

// Stop stop the worker.
// The work left in the queue will remain.
func (w *Worker) Stop() {
	w.latch.Stopping()
	<-w.latch.NotifyStopped()
}

// Drain stops the worker and synchronously finishes work.
func (w *Worker) Drain() {
	w.DrainContext(context.Background())
}

// DrainContext stops the worker and synchronously drains the the remaining work
// with a given context.
func (w *Worker) DrainContext(ctx context.Context) {
	w.latch.Stopping()
	<-w.latch.NotifyStopped()
	remaining := len(w.work)
	stopped := make(chan struct{})
	go func() {
		defer func() {
			close(stopped)
		}()
		for x := 0; x < remaining; x++ {
			w.Execute(ctx, <-w.work)
		}
	}()
	<-stopped
}

// Close stops the worker.
func (w *Worker) Close() error {
	w.latch.Stopping()
	<-w.latch.NotifyStopped()
	w.work = nil
	return nil
}
