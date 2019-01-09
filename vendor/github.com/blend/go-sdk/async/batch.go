package async

import (
	"context"
	"runtime"
)

// NewBatch creates a new batch processor.
func NewBatch(action QueueAction) *Batch {
	return &Batch{
		latch:      &Latch{},
		action:     action,
		numWorkers: runtime.NumCPU(),
	}
}

// Batch is a batch of work executed by a fixed count of workers.
type Batch struct {
	latch      *Latch
	numWorkers int
	action     QueueAction
	work       chan interface{}
	errors     chan error
	workers    chan *Worker
}

// WithWork sets the work channel.
func (b *Batch) WithWork(work chan interface{}) *Batch {
	b.work = work
	return b
}

// Work returns the work channel.
func (b *Batch) Work() chan interface{} {
	return b.work
}

// Add adds an item.
func (b *Batch) Add(item interface{}) {
	b.work <- item
}

// WithNumWorkers sets the number of workers.
// It defaults to `runtime.NumCPU()`
func (b *Batch) WithNumWorkers(numWorkers int) *Batch {
	b.numWorkers = numWorkers
	return b
}

// NumWorkers returns the number of worker route
func (b *Batch) NumWorkers() int {
	return b.numWorkers
}

// Latch returns the worker latch.
func (b *Batch) Latch() *Latch {
	return b.latch
}

// WithErrors sets the error channel.
func (b *Batch) WithErrors(errors chan error) *Batch {
	b.errors = errors
	return b
}

// Errors returns a channel to read action errors from.
func (b *Batch) Errors() chan error {
	return b.errors
}

// Process exeuctes the action for all the work items.
func (b *Batch) Process() {
	b.ProcessContext(context.Background())
}

// ProcessContext exeuctes the action for all the work items.
func (b *Batch) ProcessContext(ctx context.Context) {
	// initialize the workers
	b.workers = make(chan *Worker, b.numWorkers)
	for x := 0; x < b.numWorkers; x++ {
		worker := &Worker{
			latch:  &Latch{},
			work:   make(chan interface{}),
			errors: b.errors,
		}
		worker.action = b.andReturn(worker, b.action)
		worker.Start()
		b.workers <- worker
	}

	defer func() {
		for x := 0; x < b.numWorkers; x++ {
			worker := <-b.workers
			worker.Stop()
		}
	}()

	numWorkItems := len(b.work)
	var worker *Worker
	var workItem interface{}
	for x := 0; x < numWorkItems; x++ {
		workItem = <-b.work
		select {
		case worker = <-b.workers:
			worker.Enqueue(workItem)
		case <-ctx.Done():
			b.latch.Stopped()
			return
		case <-b.latch.NotifyStopping():
			b.latch.Stopped()
			return
		}
	}
}

// Abort aborts the work in progress.
func (b *Batch) Abort() {
	b.latch.Stopping()
	<-b.latch.NotifyStopped()
}

// AndReturn creates an action handler that returns a given worker to the worker queue.
// It wraps any action provided to the queue.
func (b *Batch) andReturn(worker *Worker, action QueueAction) QueueAction {
	return func(ctx context.Context, workItem interface{}) error {
		defer func() {
			b.workers <- worker
		}()
		return action(ctx, workItem)
	}
}
