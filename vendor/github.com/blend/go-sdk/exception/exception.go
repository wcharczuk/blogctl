package exception

import (
	"encoding/json"
	"fmt"
	"io"
)

// New returns a new exception with a call stack.
func New(class interface{}) Exception {
	return newWithStartDepth(class, defaultNewStartDepth)
}

func newWithStartDepth(class interface{}, startDepth int) Exception {
	if class == nil {
		return nil
	}

	if typed, isTyped := class.(Exception); isTyped {
		return typed
	} else if err, isErr := class.(error); isErr {
		return &Ex{
			class: err,
			stack: callers(startDepth),
		}
	} else if str, isStr := class.(string); isStr {
		return &Ex{
			class: Class(str),
			stack: callers(startDepth),
		}
	}
	return &Ex{
		class: Class(fmt.Sprint(class)),
		stack: callers(startDepth),
	}
}

// Nest nests an arbitrary number of exceptions.
func Nest(err ...error) Exception {
	var ex Exception
	var last Exception
	var didSet bool

	for _, e := range err {
		if e != nil {
			var wrappedEx *Ex
			if typedEx, isTyped := e.(*Ex); !isTyped {
				wrappedEx = &Ex{
					class: e,
					stack: callers(defaultStartDepth),
				}
			} else {
				wrappedEx = typedEx
			}

			if wrappedEx != ex {
				if ex == nil {
					ex = wrappedEx
					last = wrappedEx
				} else {
					last.WithInner(wrappedEx)
					last = wrappedEx
				}
				didSet = true
			}
		}
	}
	if didSet {
		return ex
	}
	return nil
}

// Exception is an exception.
type Exception interface {
	error
	fmt.Formatter
	json.Marshaler

	WithClass(error) Exception
	Class() error
	WithMessage(string) Exception
	WithMessagef(string, ...interface{}) Exception
	Message() string
	WithInner(error) Exception
	Inner() error
	WithStack(StackTrace) Exception
	Stack() StackTrace

	Decompose() map[string]interface{}
}

// Ex is an error with a stack trace.
// It also can have an optional cause, it implements `Exception`
type Ex struct {
	// Class disambiguates between errors, it can be used to identify the type of the error.
	class error
	// Message adds further detail to the error, and shouldn't be used for disambiguation.
	message string
	// Inner holds the original error in cases where we're wrapping an error with a stack trace.
	inner error
	// Stack is the call stack frames used to create the stack output.
	stack StackTrace
}

// Format allows for conditional expansion in printf statements
// based on the token and flags used.
// 	%+v : class + message + stack
// 	%v, %c : class
// 	%m : message
// 	%t : stack
func (e *Ex) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if e.class != nil && len(e.class.Error()) > 0 {
				fmt.Fprintf(s, "%s", e.class)
				if len(e.message) > 0 {
					fmt.Fprintf(s, "\nmessage: %s", e.message)
				}
			} else if len(e.message) > 0 {
				io.WriteString(s, e.message)
			}
			e.stack.Format(s, verb)
		} else if s.Flag('-') {
			e.stack.Format(s, verb)
		} else {
			io.WriteString(s, e.class.Error())
			if len(e.message) > 0 {
				fmt.Fprintf(s, "\nmessage: %s", e.message)
			}
		}
		if e.inner != nil {
			if typed, ok := e.inner.(fmt.Formatter); ok {
				fmt.Fprint(s, "\ninner: ")
				typed.Format(s, verb)
			} else {
				fmt.Fprintf(s, "\ninner: %v", e.inner)
			}
		}
		return
	case 'c':
		io.WriteString(s, e.class.Error())
	case 'i':
		if e.inner != nil {
			io.WriteString(s, e.inner.Error())
		}
	case 'm':
		io.WriteString(s, e.message)
	case 'q':
		fmt.Fprintf(s, "%q", e.message)
	}
}

// Error implements the `error` interface
func (e *Ex) Error() string {
	return e.class.Error()
}

// Decompose breaks the exception down to be marshalled into an intermediate format.
func (e *Ex) Decompose() map[string]interface{} {
	values := map[string]interface{}{}
	values["Class"] = e.class.Error()
	values["Message"] = e.message
	if e.stack != nil {
		values["Stack"] = e.Stack().Strings()
	}
	if e.inner != nil {
		if typed, isTyped := e.inner.(Exception); isTyped {
			values["Inner"] = typed.Decompose()
		} else {
			values["Inner"] = e.inner.Error()
		}
	}
	return values
}

// MarshalJSON is a custom json marshaler.
func (e *Ex) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Decompose())
}

// WithClass sets the exception class and returns the exepction.
func (e *Ex) WithClass(class error) Exception {
	e.class = class
	return e
}

// Class returns the exception class.
// This error should be equatable, that is, you should be able to use it to test
// if an error is a similar class to another error.
func (e *Ex) Class() error {
	return e.class
}

// WithInner sets inner or causing exception.
func (e *Ex) WithInner(err error) Exception {
	e.inner = newWithStartDepth(err, defaultNewStartDepth)
	return e
}

// Inner returns an optional nested exception.
func (e *Ex) Inner() error {
	return e.inner
}

// WithMessage sets the exception message.
func (e *Ex) WithMessage(message string) Exception {
	e.message = message
	return e
}

// WithMessagef sets the message based on a format and args, and returns the exception.
func (e *Ex) WithMessagef(format string, args ...interface{}) Exception {
	e.message = fmt.Sprintf(format, args...)
	return e
}

// Message returns the exception descriptive message.
func (e *Ex) Message() string {
	return e.message
}

// WithStack sets the stack.
func (e *Ex) WithStack(stack StackTrace) Exception {
	e.stack = stack
	return e
}

// Stack returns the stack provider.
// This is typically the runtime []uintptr or []string if restored after the fact.
func (e *Ex) Stack() StackTrace {
	return e.stack
}
