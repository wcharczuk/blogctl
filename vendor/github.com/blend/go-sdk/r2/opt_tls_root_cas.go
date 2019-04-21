package r2

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
)

// OptTLSRootCAs sets the client tls root ca pool.
func OptTLSRootCAs(pool *x509.CertPool) Option {
	return func(r *Request) error {
		if r.Client == nil {
			r.Client = &http.Client{}
		}
		if r.Client.Transport == nil {
			r.Client.Transport = &http.Transport{}
		}
		if typed, ok := r.Client.Transport.(*http.Transport); ok {
			if typed.TLSClientConfig == nil {
				typed.TLSClientConfig = &tls.Config{}
			}
			typed.TLSClientConfig.RootCAs = pool
		}
		return nil
	}
}
