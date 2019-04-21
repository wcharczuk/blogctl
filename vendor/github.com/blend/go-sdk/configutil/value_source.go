package configutil

import "time"

// StringSource is a type that can return a value.
type StringSource interface {
	// String should return a string if the source has a given value.
	// It should return nil if the value is not present.
	// It should return an error if there was a problem fetching the value.
	String() (*string, error)
}

// StringsSource is a type that can return a value.
type StringsSource interface {
	// Strings should return a string array if the source has a given value.
	// It should return nil if the value is not present.
	// It should return an error if there was a problem fetching the value.
	Strings() ([]string, error)
}

// BoolSource is a type that can return a value.
type BoolSource interface {
	// Bool should return a bool if the source has a given value.
	// It should return nil if the value is not found.
	// It should return an error if there was a problem fetching the value.
	Bool() (*bool, error)
}

// IntSource is a type that can return a value.
type IntSource interface {
	// Int should return a int if the source has a given value.
	// It should return nil if the value is not found.
	// It should return an error if there was a problem fetching the value.
	Int() (*int, error)
}

// Float64Source is a type that can return a value.
type Float64Source interface {
	// Float should return a float64 if the source has a given value.
	// It should return nil if the value is not found.
	// It should return an error if there was a problem fetching the value.
	Float64() (*float64, error)
}

// DurationSource is a type that can return a time.Duration value.
type DurationSource interface {
	// Duration should return a time.Duration if the source has a given value.
	// It should return nil if the value is not present.
	// It should return an error if there was a problem fetching the value.
	Duration() (*time.Duration, error)
}
