package configutil

import "time"

var (
	_ StringSource   = (*StringFunc)(nil)
	_ StringsSource  = (*StringsFunc)(nil)
	_ BoolSource     = (*BoolFunc)(nil)
	_ IntSource      = (*IntFunc)(nil)
	_ Float64Source  = (*Float64Func)(nil)
	_ DurationSource = (*DurationFunc)(nil)
)

// StringFunc is a value source from a function.
type StringFunc func() (*string, error)

// String returns an invocation of the function.
func (svf StringFunc) String() (*string, error) {
	return svf()
}

// StringsFunc is a value source from a function.
type StringsFunc func() ([]string, error)

// Strings returns an invocation of the function.
func (svf StringsFunc) Strings() ([]string, error) {
	return svf()
}

// BoolFunc is a bool value source.
// It can be used with configutil.SetBool
type BoolFunc func() (*bool, error)

// Bool returns an invocation of the function.
func (vf BoolFunc) Bool() (*bool, error) {
	return vf()
}

// IntFunc is an int value source from a commandline flag.
type IntFunc func() (*int, error)

// Int returns an invocation of the function.
func (vf IntFunc) Int() (*int, error) {
	return vf()
}

// Float64Func is a float value source from a commandline flag.
type Float64Func func() (*float64, error)

// Float64 returns an invocation of the function.
func (vf Float64Func) Float64() (*float64, error) {
	return vf()
}

// DurationFunc is a value source from a function.
type DurationFunc func() (*time.Duration, error)

// Duration returns an invocation of the function.
func (vf DurationFunc) Duration() (*time.Duration, error) {
	return vf()
}
