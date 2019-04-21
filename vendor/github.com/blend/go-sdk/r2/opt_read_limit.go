package r2

import "io"

// OptReadLimit applies an `io.LimitReader` to a request's response.
func OptReadLimit(byteCount int64) Option {
	return OptResponseBodyInterceptor(func(r io.Reader) io.Reader {
		return io.LimitReader(r, byteCount)
	})
}
