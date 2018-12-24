package async

import (
	"sync"

	"github.com/blend/go-sdk/exception"
)

const (
	// DefaultQueueMaxWork is the maximum number of work items before queueing blocks.
	DefaultQueueMaxWork = 1 << 10
)

// NewQueue returns a new queue worker.
func NewQueue(action func(interface{}) error) *Queue {
	return &Queue{
		action: action,
		latch:  &Latch{},
		work:   make(chan interface{}, DefaultQueueMaxWork),
	}
}

// Queue is a worker that is pushed work over a channel.
type Queue struct {
	latch  *Latch
	action func(interface{}) error
	errors chan error
	work   chan interface{}
}

// WithWork sets the work channel.
func (q *Queue) WithWork(work chan interface{}) *Queue {
	q.work = work
	return q
}

// Work returns the work channel.
func (q *Queue) Work() chan interface{} {
	return q.work
}

// Latch returns the worker latch.
func (q *Queue) Latch() *Latch {
	return q.latch
}

// WithErrors returns the error channel.
func (q *Queue) WithErrors(errors chan error) *Queue {
	q.errors = errors
	return q
}

// Errors returns a channel to read action errors from.
func (q *Queue) Errors() chan error {
	return q.errors
}

// Enqueue adds an item to the work queue.
func (q *Queue) Enqueue(obj interface{}) {
	q.work <- obj
}

// Start starts the worker.
func (q *Queue) Start() {
	q.latch.Starting()
	go func() {
		q.latch.Started()
		var err error
		var workItem interface{}
		for {
			select {
			case workItem = <-q.work:
				err = q.action(workItem)
				if err != nil && q.errors != nil {
					q.errors <- err
				}
			case <-q.latch.NotifyStopping():
				q.latch.Stopped()
				return
			}
		}
	}()
	<-q.latch.NotifyStarted()
}

// SafeAction invokes the action and recovers panics.
func (q *Queue) SafeAction(workItem interface{}) {
	defer func() {
		if r := recover(); r != nil {
			if q.errors != nil {
				q.errors <- exception.New(r)
			}
		}
	}()
	if err := q.action(workItem); err != nil {
		if q.errors != nil {
			q.errors <- exception.New(err)
		}
	}
}

// Close stops the queue.
func (q *Queue) Close() error {
	q.latch.Stopping()
	<-q.latch.NotifyStopped()

	remaining := len(q.work)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for x := 0; x < remaining; x++ {
			q.SafeAction(<-q.work)
		}
	}()
	wg.Wait()
	return nil
}
