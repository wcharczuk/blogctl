package graceful

import (
	"os"
	"os/signal"
	"syscall"
)

// Graceful is a server that can start and shutdown.
type Graceful interface {
	Start() error
	Stop() error
	NotifyStarted() <-chan struct{}
	NotifyStopped() <-chan struct{}
}

// Shutdown starts an hosted process and responds to SIGINT and SIGTERM to shut the app down.
// It will return any errors returned by app.Start() that are not caused by shutting down the server.
func Shutdown(hosted Graceful) error {
	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, os.Interrupt, syscall.SIGTERM)
	return ShutdownBySignal(hosted, terminateSignal)
}

// ShutdownBySignal gracefully stops a hosted process based on an os signal channel.
func ShutdownBySignal(hosted Graceful, terminateSignal chan os.Signal) error {
	shutdown := make(chan struct{})
	shutdownAbort := make(chan struct{})
	shutdownComplete := make(chan struct{})
	server := make(chan struct{})
	errors := make(chan error, 2)

	go func() {
		if err := hosted.Start(); err != nil {
			errors <- err
		}
		close(server)
	}()

	go func() {
		select {
		case <-shutdown:
			if err := hosted.Stop(); err != nil {
				errors <- err
			}
			close(shutdownComplete)
			return
		case <-shutdownAbort:
			close(shutdownComplete)
			return
		}
	}()

	select {
	case <-terminateSignal: // if we've issued a shutdown, wait for the server to exit
		close(shutdown)
		<-shutdownComplete
		<-server
	case <-server: // if the server exited
		close(shutdownAbort) // quit the signal listener
		<-shutdownComplete
	}

	if len(errors) > 0 {
		return <-errors
	}
	return nil
}
