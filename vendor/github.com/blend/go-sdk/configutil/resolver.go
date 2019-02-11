package configutil

// Resolver is a type that can be resolved.
type Resolver interface {
	Resolve() error
}
