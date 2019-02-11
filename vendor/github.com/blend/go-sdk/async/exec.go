package async

import "github.com/blend/go-sdk/exception"

// Exec runs an action and passes any errors to the returned errors channel.
func Exec(action func() error) chan error {
	errors := make(chan error, 1)
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
	return errors
}
