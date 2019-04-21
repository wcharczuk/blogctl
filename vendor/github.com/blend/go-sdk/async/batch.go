package async

import (
	"context"
	"runtime"
)

// NewBatch creates a new batch processor.
// Batch processes are a known quantity of work that needs to be processed in parallel.
func NewBatch(action WorkAction, work chan interface{}, options ...BatchOption) *Batch {
	b := Batch{
		Action:      action,
		Work:        work,
		Parallelism: runtime.NumCPU(),
	}
	for _, option := range options {
		option(&b)
	}
	return &b
}

// BatchOption is an option for the batch worker.
type BatchOption func(*Batch)

// OptBatchErrors sets the batch worker error return channel.
func OptBatchErrors(errors chan error) BatchOption {
	return func(i *Batch) {
		i.Errors = errors
	}
}

// OptBatchParallelism sets the batch worker parallelism, or the number of workers to create.
func OptBatchParallelism(parallelism int) BatchOption {
	return func(i *Batch) {
		i.Parallelism = parallelism
	}
}

// Batch is a batch of work executed by a fixed count of workers.
type Batch struct {
	Action      WorkAction
	Parallelism int
	Work        chan interface{}
	Errors      chan error
}

// Process executes the action for all the work items.
func (b *Batch) Process(ctx context.Context) {
	// initialize the workers
	workers := make(chan *Worker, b.Parallelism)

	returnWorker := func(ctx context.Context, worker *Worker) error {
		workers <- worker
		return nil
	}

	for x := 0; x < b.Parallelism; x++ {
		worker := NewWorker(b.Action)
		worker.Errors = b.Errors
		worker.Finalizer = returnWorker
		go worker.Start()
		<-worker.NotifyStarted()
		workers <- worker
	}

	defer func() {
		for x := 0; x < b.Parallelism; x++ {
			worker := <-workers
			worker.Stop()
		}
	}()

	numWorkItems := len(b.Work)
	var worker *Worker
	var workItem interface{}
	for x := 0; x < numWorkItems; x++ {
		workItem = <-b.Work
		select {
		case worker = <-workers:
			worker.Enqueue(workItem)
		case <-ctx.Done():
			return
		}
	}
}
