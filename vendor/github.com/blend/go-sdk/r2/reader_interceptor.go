package r2

import "io"

// ReaderInterceptor is a handler for request.ReadResponse
type ReaderInterceptor func(io.Reader) io.Reader
