package configutil

// ConfigOptions are options built for reading configs.
type ConfigOptions struct {
	Resolver func(interface{}) error
	Paths    []string
}
