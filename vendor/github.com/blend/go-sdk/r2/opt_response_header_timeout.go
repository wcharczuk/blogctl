package r2

import (
	"net/http"
	"time"
)

// OptResponseHeaderTimeout sets the client transport ResponseHeaderTimeout.
func OptResponseHeaderTimeout(d time.Duration) Option {
	return func(r *Request) error {
		if r.Client.Transport == nil {
			r.Client.Transport = &http.Transport{}
		}
		if typed, ok := r.Client.Transport.(*http.Transport); ok {
			typed.ResponseHeaderTimeout = d
		}
		return nil
	}
}
