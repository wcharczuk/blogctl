package async

// NewParallelQueue returns a new parallel queue worker.
func NewParallelQueue(numWorkers int, action func(interface{}) error) *ParallelQueue {
	return &ParallelQueue{
		action:     action,
		numWorkers: numWorkers,
		latch:      &Latch{},
		workers:    make(chan *Queue, numWorkers),
		work:       make(chan interface{}, DefaultQueueMaxWork),
	}
}

// ParallelQueue is a queude with multiple workers..
type ParallelQueue struct {
	latch      *Latch
	numWorkers int
	workers    chan *Queue
	action     func(interface{}) error
	errors     chan error
	work       chan interface{}
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

// WithErrors returns the error channel.
func (pq *ParallelQueue) WithErrors(errors chan error) *ParallelQueue {
	pq.errors = errors
	return pq
}

// Errors returns a channel to read action errors from.
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
	for x := 0; x < pq.numWorkers; x++ {
		worker := &Queue{
			latch:  &Latch{},
			errors: pq.errors,
			work:   make(chan interface{}),
		}
		worker.action = pq.AndReturn(worker, pq.action)
		worker.Start()
		pq.workers <- worker
	}
	go pq.Dispatch()
	<-pq.latch.NotifyStarted()
}

// Dispatch processes work items in a loop.
func (pq *ParallelQueue) Dispatch() {
	pq.latch.Started()
	var workItem interface{}
	var worker *Queue
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
func (pq *ParallelQueue) AndReturn(worker *Queue, action func(interface{}) error) func(interface{}) error {
	return func(workItem interface{}) error {
		defer func() {
			pq.workers <- worker
		}()
		return action(workItem)
	}
}

// Close stops the queue.
func (pq *ParallelQueue) Close() error {
	pq.latch.Stopping()
	<-pq.latch.NotifyStopped()

	for x := 0; x < pq.numWorkers; x++ {
		worker := <-pq.workers
		if err := worker.Close(); err != nil {
			return err
		}
	}
	return nil
}
