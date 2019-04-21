package configutil

import "time"

var (
	_ DurationSource = (*Duration)(nil)
)

// Duration implements value provider.
type Duration time.Duration

// Duration returns the value for a constant.
func (dc Duration) Duration() (*time.Duration, error) {
	value := time.Duration(dc)
	return &value, nil
}
