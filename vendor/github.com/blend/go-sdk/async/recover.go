package async

import (
	"sync"

	"github.com/blend/go-sdk/exception"
)

// Recover runs an action and passes any errors to the given errors channel.
func Recover(errors chan error, action func() error) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errors <- exception.New(r)
			}
		}()

		if err := action(); err != nil {
			errors <- err
		}
	}()
}

// RecoverGroup runs a recovery against a specific wait group with an error collector.
func RecoverGroup(wg *sync.WaitGroup, errors chan error, action func() error) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errors <- exception.New(r)
			}
			wg.Done()
		}()

		if err := action(); err != nil {
			errors <- err
		}
	}()
}
