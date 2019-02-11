package configutil

import "time"

var (
	_ ValueSource         = (*Const)(nil)
	_ BoolValueSource     = (*BoolConst)(nil)
	_ IntValueSource      = (*IntConst)(nil)
	_ FloatValueSource    = (*FloatConst)(nil)
	_ DurationValueSource = (*DurationConst)(nil)
)

// Const implements value provider.
type Const string

// Value returns the value for a constant.
func (cv Const) Value() (*string, error) {
	value := string(cv)
	return &value, nil
}

// BoolConst implements value provider.
type BoolConst bool

// BoolValue returns the value for a constant.
func (bc BoolConst) BoolValue() (*bool, error) {
	value := bool(bc)
	return &value, nil
}

// IntConst implements value provider.
type IntConst int

// IntValue returns the value for a constant.
func (ic IntConst) IntValue() (*int, error) {
	value := int(ic)
	return &value, nil
}

// FloatConst implements value provider.
type FloatConst float64

// FloatValue returns the value for a constant.
func (ic FloatConst) FloatValue() (*float64, error) {
	value := float64(ic)
	return &value, nil
}

// DurationConst implements value provider.
type DurationConst time.Duration

// DurationValue returns the value for a constant.
func (dc DurationConst) DurationValue() (*time.Duration, error) {
	value := time.Duration(dc)
	return &value, nil
}
