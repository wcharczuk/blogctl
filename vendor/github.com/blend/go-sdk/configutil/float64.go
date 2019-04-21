package configutil

var (
	_ Float64Source = (*Float64)(nil)
)

// Float64 implements value provider.
type Float64 float64

// Float64 returns the value for a constant.
func (f Float64) Float64() (*float64, error) {
	value := float64(f)
	return &value, nil
}
