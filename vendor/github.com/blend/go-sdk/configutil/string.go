package configutil

var (
	_ StringSource = (*String)(nil)
)

// String implements value provider.
type String string

// StringValue returns the value for a constant.
func (s String) String() (*string, error) {
	value := string(s)
	if value == "" {
		return nil, nil
	}
	return &value, nil
}
