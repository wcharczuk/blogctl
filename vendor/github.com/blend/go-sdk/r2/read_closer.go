package r2

import "io"

// NewReadCloser creates a new read closer from a given body and interceptor.
func NewReadCloser(body io.ReadCloser, interceptor ReaderInterceptor) *ReadCloser {
	return &ReadCloser{
		Closer: body,
		Reader: interceptor(body),
	}
}

// ReadCloser allows you to split a ReaderCloser up into separate components.
// This lets use apply `io.LimitReader` and the like to response bodies, but preserve
// the original Close functionality.
type ReadCloser struct {
	Closer io.Closer
	Reader io.Reader
}

// Close calls the closer.
func (rc *ReadCloser) Close() error {
	return rc.Closer.Close()
}

// Read calls the reader.
func (rc *ReadCloser) Read(contents []byte) (int, error) {
	return rc.Reader.Read(contents)
}
