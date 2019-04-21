package async

// RunToError runs a given set of functions until one of them panics or errors.
// It is useful when you need to start multiple servers and exit if any of them crashes.
func RunToError(fns ...func() error) error {
	panicChan := make(chan interface{}, 1)
	errChan := make(chan error, 1)
	for _, fn := range fns {
		go func(fn func() error) {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			errChan <- fn()
		}(fn)
	}
	select {
	case p := <-panicChan:
		panic(p)
	case err := <-errChan:
		return err
	}
}
