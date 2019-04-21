package collections

// Error is an error string.
type Error string

// Error implements error.
func (e Error) Error() string { return string(e) }

// Labels is a loose type alias to map[string]string
type Labels = map[string]string

// Vars is a loose type alias to map[string]string
type Vars = map[string]interface{}

// Any is a loose type alias to interface{}.
type Any = interface{}
