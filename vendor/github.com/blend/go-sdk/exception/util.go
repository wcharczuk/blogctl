package exception

// ErrClass returns the exception class or the error message.
// This depends on if the err is itself an exception or not.
func ErrClass(err error) string {
	if err == nil {
		return ""
	}
	if ex := As(err); ex != nil && ex.Class() != nil {
		return ex.Class().Error()
	}
	return err.Error()
}

// ErrMessage returns the exception message.
// This depends on if the err is itself an exception or not.
// If it is not an exception, this will return empty string.
func ErrMessage(err error) string {
	if err == nil {
		return ""
	}
	if ex := As(err); ex != nil && ex.Class() != nil {
		return ex.Message()
	}
	return ""
}

// Is is a helper function that returns if an error is an exception.
func Is(err, cause error) bool {
	if err == nil || cause == nil {
		return false
	}
	if typed, isTyped := err.(Exception); isTyped && typed.Class() != nil {
		return (typed.Class() == cause) || (typed.Class().Error() == cause.Error())
	}
	return (err == cause) || (err.Error() == cause.Error())
}

// Inner returns an inner error if the error is an exception.
func Inner(err error) error {
	if typed := As(err); typed != nil {
		return typed.Inner()
	}
	return nil
}

// As is a helper method that returns an error as an exception.
func As(err error) Exception {
	if typed, typedOk := err.(Exception); typedOk {
		return typed
	}
	return nil
}
