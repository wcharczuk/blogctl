package r2

// Defaults is a helper to create requests with a common set of base options.
type Defaults []Option

// Add adds new options to the default set.
func (d Defaults) Add(options ...Option) Defaults {
	d = append(d, options...)
	return d
}

// ConcatWith concats the options with a given set of new options for output.
func (d Defaults) ConcatWith(options ...Option) []Option {
	return []Option(append(d, options...))
}
