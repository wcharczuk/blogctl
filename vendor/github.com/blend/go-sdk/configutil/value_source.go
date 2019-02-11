package configutil

import "time"

// ValueSource is a type that can return a value.
type ValueSource interface {
	// Value should return a string if the source has a given value.
	// It should return nil if the value is not present.
	// It should return an error if there was a problem fetching the value.
	Value() (*string, error)
}

// BoolValueSource is a type that can return a value.
type BoolValueSource interface {
	// BoolValue should return a bool if the source has a given value.
	// It should return nil if the value is not found.
	// It should return an error if there was a problem fetching the value.
	BoolValue() (*bool, error)
}

// IntValueSource is a type that can return a value.
type IntValueSource interface {
	// IntValue should return a int if the source has a given value.
	// It should return nil if the value is not found.
	// It should return an error if there was a problem fetching the value.
	IntValue() (*int, error)
}

// FloatValueSource is a type that can return a value.
type FloatValueSource interface {
	// FloatValue should return a float64 if the source has a given value.
	// It should return nil if the value is not found.
	// It should return an error if there was a problem fetching the value.
	FloatValue() (*float64, error)
}

// DurationValueSource is a type that can return a time.Duration value.
type DurationValueSource interface {
	// DurationValue should return a time.Duration if the source has a given value.
	// It should return nil if the value is not present.
	// It should return an error if there was a problem fetching the value.
	DurationValue() (*time.Duration, error)
}
