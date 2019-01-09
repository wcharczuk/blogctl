package async

import (
	"context"
	"runtime"
)

// NewParallelQueue returns a new parallel queue worker.
func NewParallelQueue(action QueueAction) *ParallelQueue {
	return &ParallelQueue{
		latch:      &Latch{},
		work:       make(chan interface{}, DefaultQueueMaxWork),
		action:     action,
		numWorkers: runtime.NumCPU(),
	}
}

// ParallelQueue is a queude with multiple workers..
type ParallelQueue struct {
	latch      *Latch
	numWorkers int
	action     QueueAction
	workers    chan *Worker
	work       chan interface{}
	errors     chan error
}

// WithNumWorkers sets the number of workers.
// It defaults to `runtime.NumCPU()`
func (pq *ParallelQueue) WithNumWorkers(numWorkers int) *ParallelQueue {
	pq.numWorkers = numWorkers
	return pq
}

// NumWorkers returns the number of worker route
func (pq *ParallelQueue) NumWorkers() int {
	return pq.numWorkers
}

// WithWork sets the work channel.
func (pq *ParallelQueue) WithWork(work chan interface{}) *ParallelQueue {
	pq.work = work
	return pq
}

// Work returns the work channel.
func (pq *ParallelQueue) Work() chan interface{} {
	return pq.work
}

// Latch returns the worker latch.
func (pq *ParallelQueue) Latch() *Latch {
	return pq.latch
}

// WithErrors sets the error channel.
func (pq *ParallelQueue) WithErrors(errors chan error) *ParallelQueue {
	pq.errors = errors
	return pq
}

// Errors returns a channel to read action errors from.
// You must provide it with `WithErrors`.
func (pq *ParallelQueue) Errors() chan error {
	return pq.errors
}

// Enqueue adds an item to the work queue.
func (pq *ParallelQueue) Enqueue(obj interface{}) {
	pq.work <- obj
}

// Start starts the worker.
func (pq *ParallelQueue) Start() {
	pq.latch.Starting()
	if pq.workers == nil {
		pq.initializeWorkers()
	}
	pq.startWorkers()
	go pq.dispatch()
	<-pq.latch.NotifyStarted()
}

// Close stops the queue.
// Any work left in the queue will be discarded.
func (pq *ParallelQueue) Close() error {
	pq.stopWorkers()
	pq.latch.Stopping()
	<-pq.latch.NotifyStopped()
	return nil
}

// helpers

// StartWorkers starts all workers.
func (pq *ParallelQueue) startWorkers() {
	for x := 0; x < pq.numWorkers; x++ {
		worker := <-pq.workers
		worker.Start()
		pq.workers <- worker
	}
}

// StopWorkers closes all workers.
func (pq *ParallelQueue) stopWorkers() {
	for x := 0; x < pq.numWorkers; x++ {
		worker := <-pq.workers
		worker.Stop()
		pq.workers <- worker
	}
}

// InitializeWorkers initializes the workers.
func (pq *ParallelQueue) initializeWorkers() {
	pq.workers = make(chan *Worker, pq.numWorkers)
	for x := 0; x < pq.numWorkers; x++ {
		worker := &Worker{
			latch:  &Latch{},
			errors: pq.errors,
			work:   make(chan interface{}),
		}
		worker.action = pq.andReturn(worker, pq.action)
		pq.workers <- worker
	}
}

// Dispatch processes work items in a loop.
func (pq *ParallelQueue) dispatch() {
	pq.latch.Started()
	var workItem interface{}
	var worker *Worker
	for {
		select {
		case workItem = <-pq.work:
			select {
			case worker = <-pq.workers:
				worker.work <- workItem
			case <-pq.latch.NotifyStopping():
				pq.latch.Stopped()
				return
			}
		case <-pq.latch.NotifyStopping():
			pq.latch.Stopped()
			return
		}
	}
}

// AndReturn creates an action handler that returns a given worker to the worker queue.
// It wraps any action provided to the queue.
func (pq *ParallelQueue) andReturn(worker *Worker, action QueueAction) QueueAction {
	return func(ctx context.Context, workItem interface{}) error {
		defer func() {
			pq.workers <- worker
		}()
		return action(ctx, workItem)
	}
}
