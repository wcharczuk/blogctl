package configutil

// ConfigResolver is a type that can be resolved.
type ConfigResolver interface {
	Resolve() error
}
