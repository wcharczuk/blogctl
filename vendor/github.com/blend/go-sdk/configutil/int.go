package configutil

var (
	_ IntSource = (*Int)(nil)
)

// Int implements value provider.
type Int int

// Int returns the value for a constant.
func (i Int) Int() (*int, error) {
	value := int(i)
	return &value, nil
}
