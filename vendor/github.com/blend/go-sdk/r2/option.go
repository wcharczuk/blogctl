package r2

// Option is a modifier for a request.
type Option func(*Request) error
