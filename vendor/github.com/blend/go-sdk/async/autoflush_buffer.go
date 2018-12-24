package async

import (
	"sync"
	"time"

	"github.com/blend/go-sdk/collections"
	"github.com/blend/go-sdk/exception"
)

// NewAutoflushBuffer creates a new autoflush buffer.
func NewAutoflushBuffer(maxLen int, interval time.Duration) *AutoflushBuffer {
	return &AutoflushBuffer{
		maxLen:       maxLen,
		interval:     interval,
		flushOnAbort: true,
		contents:     collections.NewRingBufferWithCapacity(maxLen),
		latch:        &Latch{},
	}
}

// AutoflushBuffer is a backing store that operates either on a fixed length flush or a fixed interval flush.
// A handler should be provided but without one the buffer will just clear.
// Adds that would cause fixed length flushes do not block on the flush handler.
type AutoflushBuffer struct {
	maxLen   int
	interval time.Duration

	contents     *collections.RingBuffer
	contentsLock sync.Mutex

	flushOnAbort bool
	handler      func(obj []interface{})

	latch *Latch
}

// WithFlushOnAbort sets if we should flush on aborts or not.
// This defaults to true.
func (ab *AutoflushBuffer) WithFlushOnAbort(should bool) *AutoflushBuffer {
	ab.flushOnAbort = should
	return ab
}

// ShouldFlushOnAbort returns if the buffer will do one final flush on abort.
func (ab *AutoflushBuffer) ShouldFlushOnAbort() bool {
	return ab.flushOnAbort
}

// Interval returns the flush interval.
func (ab *AutoflushBuffer) Interval() time.Duration {
	return ab.interval
}

// MaxLen returns the maximum buffer length before a flush is triggered.
func (ab *AutoflushBuffer) MaxLen() int {
	return ab.maxLen
}

// WithFlushHandler sets the buffer flush handler and returns a reference to the buffer.
func (ab *AutoflushBuffer) WithFlushHandler(handler func(objs []interface{})) *AutoflushBuffer {
	ab.handler = handler
	return ab
}

// NotifyStarted returns the started signal.
func (ab *AutoflushBuffer) NotifyStarted() <-chan struct{} {
	return ab.latch.NotifyStarted()
}

// NotifyStopped returns the started stopped.
func (ab *AutoflushBuffer) NotifyStopped() <-chan struct{} {
	return ab.latch.NotifyStopped()
}

// Start starts the buffer flusher.
func (ab *AutoflushBuffer) Start() error {
	if !ab.latch.CanStart() {
		return exception.New(ErrCannotStart)
	}
	ab.latch.Starting()
	go func() {
		ab.latch.Started()
		ab.runLoop()
	}()
	<-ab.latch.NotifyStarted()
	return nil
}

// Stop stops the buffer flusher.
func (ab *AutoflushBuffer) Stop() error {
	if !ab.latch.CanStop() {
		return exception.New(ErrCannotStop)
	}
	ab.latch.Stopping()
	<-ab.latch.NotifyStopped()
	return nil
}

// Add adds a new object to the buffer, blocking if it triggers a flush.
// If the buffer is full, it will call the flush handler on a separate goroutine.
func (ab *AutoflushBuffer) Add(obj interface{}) {
	ab.contentsLock.Lock()
	defer ab.contentsLock.Unlock()

	ab.contents.Enqueue(obj)
	if ab.contents.Len() >= ab.maxLen {
		ab.flushUnsafeAsync()
	}
}

// AddMany adds many objects to the buffer at once.
func (ab *AutoflushBuffer) AddMany(objs ...interface{}) {
	ab.contentsLock.Lock()
	defer ab.contentsLock.Unlock()

	for _, obj := range objs {
		ab.contents.Enqueue(obj)
		if ab.contents.Len() >= ab.maxLen {
			ab.flushUnsafeAsync()
		}
	}
}

// Flush clears the buffer, if a handler is provided it is passed the contents of the buffer.
// This call is synchronous, in that it will call the flush handler on the same goroutine.
func (ab *AutoflushBuffer) Flush() {
	ab.contentsLock.Lock()
	defer ab.contentsLock.Unlock()
	ab.flushUnsafe()
}

// FlushAsync clears the buffer, if a handler is provided it is passed the contents of the buffer.
// This call is asynchronous, in that it will call the flush handler on its own goroutine.
func (ab *AutoflushBuffer) FlushAsync() {
	ab.contentsLock.Lock()
	defer ab.contentsLock.Unlock()
	ab.flushUnsafeAsync()
}

// flushUnsafeAsync flushes the buffer without acquiring any locks.
func (ab *AutoflushBuffer) flushUnsafeAsync() {
	if ab.handler != nil {
		if ab.contents.Len() > 0 {
			contents := ab.contents.Drain()
			go ab.handler(contents)
		}
	} else {
		ab.contents.Clear()
	}
}

// flushUnsafeAsync flushes the buffer without acquiring any locks.
func (ab *AutoflushBuffer) flushUnsafe() {
	if ab.handler != nil {
		if ab.contents.Len() > 0 {
			ab.handler(ab.contents.Drain())
		}
	} else {
		ab.contents.Clear()
	}
}

func (ab *AutoflushBuffer) runLoop() {
	ticker := time.Tick(ab.interval)
	for {
		select {
		case <-ticker:
			ab.FlushAsync()
		case <-ab.latch.NotifyStopping():
			if ab.flushOnAbort {
				ab.Flush()
			}
			ab.latch.Stopped()
			return
		}
	}
}
