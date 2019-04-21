package configutil

var (
	_ StringsSource = (*Strings)(nil)
)

// Strings implements a value provider.
type Strings []string

// Strings returns the value for a constant.
func (s Strings) Strings() ([]string, error) {
	return []string(s), nil
}
