package r2

import (
	"net/http"
	"time"
)

// OptTLSHandshakeTimeout sets the client transport TLSHandshakeTimeout.
func OptTLSHandshakeTimeout(d time.Duration) Option {
	return func(r *Request) error {
		if r.Client == nil {
			r.Client = &http.Client{}
		}
		if r.Client.Transport == nil {
			r.Client.Transport = &http.Transport{}
		}
		if typed, ok := r.Client.Transport.(*http.Transport); ok {
			typed.TLSHandshakeTimeout = d
		}
		return nil
	}
}
