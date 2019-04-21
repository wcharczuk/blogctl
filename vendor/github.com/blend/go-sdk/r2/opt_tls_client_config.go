package r2

import (
	"crypto/tls"
	"net/http"
)

// OptTLSClientConfig sets the tls config for the request.
// It will create a client, and a transport if unset.
func OptTLSClientConfig(cfg *tls.Config) Option {
	return func(r *Request) error {
		if r.Client == nil {
			r.Client = &http.Client{}
		}
		if r.Client.Transport == nil {
			r.Client.Transport = &http.Transport{}
		}
		if typed, ok := r.Client.Transport.(*http.Transport); ok {
			typed.TLSClientConfig = cfg
		}
		return nil
	}
}
