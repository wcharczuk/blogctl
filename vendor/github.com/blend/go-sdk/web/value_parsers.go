package web

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// AsBool parses a value as an bool.
// If the input error is set it short circuits.
func AsBool(value string, inputErr error) (output bool, err error) {
	if inputErr != nil {
		err = inputErr
		return
	}
	switch strings.ToLower(value) {
	case "1", "true", "yes":
		output = true
	case "0", "false", "no":
		output = false
	default:
		err = fmt.Errorf("invalid boolean value")
	}
	return
}

// AsInt parses a value as an int.
// If the input error is set it short circuits.
func AsInt(value string, inputErr error) (output int, err error) {
	if inputErr != nil {
		err = inputErr
		return
	}
	output, err = strconv.Atoi(value)
	return
}

// AsInt64 parses a value as an int64.
// If the input error is set it short circuits.
func AsInt64(value string, inputErr error) (output int64, err error) {
	if inputErr != nil {
		err = inputErr
		return
	}
	output, err = strconv.ParseInt(value, 10, 64)
	return
}

// AsFloat64 parses a value as an float64.
// If the input error is set it short circuits.
func AsFloat64(value string, inputErr error) (output float64, err error) {
	if inputErr != nil {
		err = inputErr
		return
	}
	output, err = strconv.ParseFloat(value, 64)
	return
}

// AsDuration parses a value as an time.Duration.
// If the input error is set it short circuits.
func AsDuration(value string, inputErr error) (output time.Duration, err error) {
	if inputErr != nil {
		err = inputErr
		return
	}
	output, err = time.ParseDuration(value)
	return
}

// AsString just returns the string directly from a value error pair.
func AsString(value string, _ error) string {
	return value
}

// AsCSV just returns the string directly from a value error pair.
func AsCSV(value string, err error) ([]string, error) {
	if err != nil {
		return nil, err
	}
	return strings.Split(value, ","), nil
}
