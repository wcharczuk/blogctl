package configutil

import (
	"time"

	"github.com/blend/go-sdk/env"
)

var (
	_ ValueSource         = (*Env)(nil)
	_ BoolValueSource     = (*Env)(nil)
	_ IntValueSource      = (*Env)(nil)
	_ FloatValueSource    = (*Env)(nil)
	_ DurationValueSource = (*Env)(nil)
)

// Env is a value provider where the string represent the variable name.
type Env string

// Value returns a given environment value.
func (e Env) Value() (*string, error) {
	if env.Env().Has(string(e)) {
		value := env.Env().String(string(e))
		return &value, nil
	}
	return nil, nil
}

// BoolValue returns a given environment value.
func (e Env) BoolValue() (*bool, error) {
	if env.Env().Has(string(e)) {
		value := env.Env().Bool(string(e))
		return &value, nil
	}
	return nil, nil
}

// IntValue returns a given environment value.
func (e Env) IntValue() (*int, error) {
	if env.Env().Has(string(e)) {
		value, err := env.Env().Int(string(e))
		if err != nil {
			return nil, err
		}
		return &value, nil
	}
	return nil, nil
}

// FloatValue returns a given environment value.
func (e Env) FloatValue() (*float64, error) {
	if env.Env().Has(string(e)) {
		value, err := env.Env().Float64(string(e))
		if err != nil {
			return nil, err
		}
		return &value, nil
	}
	return nil, nil
}

// DurationValue returns a given environment value.
func (e Env) DurationValue() (*time.Duration, error) {
	if env.Env().Has(string(e)) {
		value, err := env.Env().Duration(string(e))
		if err != nil {
			return nil, err
		}
		return &value, nil
	}
	return nil, nil
}
