package configutil

// AnyError returns the first non-nil error.
func AnyError(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}
