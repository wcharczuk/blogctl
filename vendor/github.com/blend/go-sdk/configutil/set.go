package configutil

import "time"

// Set coalesces a given list of sources into a variable.
func Set(destination *string, sources ...ValueSource) error {
	var value *string
	var err error
	for _, source := range sources {
		value, err = source.Value()
		if err != nil {
			return err
		}
		if value != nil {
			*destination = *value
			return nil
		}
	}
	return nil
}

// SetBool coalesces a given list of sources into a variable.
func SetBool(destination *bool, sources ...BoolValueSource) error {
	var value *bool
	var err error
	for _, source := range sources {
		value, err = source.BoolValue()
		if err != nil {
			return err
		}
		if value != nil {
			*destination = *value
			return nil
		}
	}
	return nil
}

// SetInt coalesces a given list of sources into a variable.
func SetInt(destination *int, sources ...IntValueSource) error {
	var value *int
	var err error
	for _, source := range sources {
		value, err = source.IntValue()
		if err != nil {
			return err
		}
		if value != nil {
			*destination = *value
			return nil
		}
	}
	return nil
}

// SetFloat coalesces a given list of sources into a variable.
func SetFloat(destination *float64, sources ...FloatValueSource) error {
	var value *float64
	var err error
	for _, source := range sources {
		value, err = source.FloatValue()
		if err != nil {
			return err
		}
		if value != nil {
			*destination = *value
			return nil
		}
	}
	return nil
}

// SetDuration coalesces a given list of sources into a variable.
func SetDuration(destination *time.Duration, sources ...DurationValueSource) error {
	var value *time.Duration
	var err error
	for _, source := range sources {
		value, err = source.DurationValue()
		if err != nil {
			return err
		}
		if value != nil {
			*destination = *value
			return nil
		}
	}
	return nil
}
