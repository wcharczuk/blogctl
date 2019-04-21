package r2

// OptResponseBodyInterceptor sets the response reader on the request.
// This should be used to do things that modify how we read the response.
func OptResponseBodyInterceptor(interceptor ReaderInterceptor) Option {
	return func(r *Request) error {
		r.ResponseBodyInterceptor = interceptor
		return nil
	}
}
