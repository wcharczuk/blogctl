package configutil

import "time"

var (
	_ ValueSource         = (*ValueFunc)(nil)
	_ BoolValueSource     = (*BoolValueFunc)(nil)
	_ IntValueSource      = (*IntValueFunc)(nil)
	_ DurationValueSource = (*DurationValueFunc)(nil)
)

// ValueFunc is a value source from a function.
type ValueFunc func() (*string, error)

// Value returns an invocation of the function.
func (vf ValueFunc) Value() (*string, error) {
	return vf()
}

// BoolValueFunc is a bool value source from a commandline flag.
type BoolValueFunc func() (*bool, error)

// BoolValue returns an invocation of the function.
func (vf BoolValueFunc) BoolValue() (*bool, error) {
	return vf()
}

// IntValueFunc is an int value source from a commandline flag.
type IntValueFunc func() (*int, error)

// IntValue returns an invocation of the function.
func (vf IntValueFunc) IntValue() (*int, error) {
	return vf()
}

// FloatValueFunc is a float value source from a commandline flag.
type FloatValueFunc func() (*float64, error)

// FloatValue returns an invocation of the function.
func (vf FloatValueFunc) FloatValue() (*float64, error) {
	return vf()
}

// DurationValueFunc is a value source from a function.
type DurationValueFunc func() (*time.Duration, error)

// DurationValue returns an invocation of the function.
func (vf DurationValueFunc) DurationValue() (*time.Duration, error) {
	return vf()
}
