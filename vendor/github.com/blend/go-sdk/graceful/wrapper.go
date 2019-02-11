package graceful

// New returns a graceful wrapper.
// Start and stop should not block.
func New(start, stop func() error) Graceful {
	return &gracefulWrapper{
		start:   start,
		stop:    stop,
		started: make(chan struct{}),
		stopped: make(chan struct{}),
	}
}

// GracefulWrapper wraps a set of functions to start and stop a process
// and returns a graceful handler.
type gracefulWrapper struct {
	start func() error
	stop  func() error

	started chan struct{}
	stopped chan struct{}
}

// Start calls the start handler.
func (gw *gracefulWrapper) Start() error {
	if err := gw.start(); err != nil {
		return err
	}
	close(gw.started)
	<-gw.stopped
	return nil
}

// Stop stops the process.
func (gw *gracefulWrapper) Stop() error {
	if err := gw.stop(); err != nil {
		return err
	}
	close(gw.stopped)
	return nil
}

// NotifyStarted returns the started channel.
func (gw *gracefulWrapper) NotifyStarted() <-chan struct{} {
	return gw.started
}

// NotifyStopped returns the stopped channel.
func (gw *gracefulWrapper) NotifyStopped() <-chan struct{} {
	return gw.stopped
}
