package configutil

var (
	_ BoolSource = (*BoolValue)(nil)
)

// Bool returns a BoolValue for a given value.
func Bool(value *bool) *BoolValue {
	if value == nil {
		return nil
	}
	typed := BoolValue(*value)
	return &typed
}

// BoolValue implements value provider.
type BoolValue bool

// Bool returns the value for a constant.
func (b *BoolValue) Bool() (*bool, error) {
	if b == nil {
		return nil, nil
	}
	value := *b
	typed := bool(value)
	return &typed, nil
}
