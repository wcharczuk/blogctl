package r2

import (
	"crypto/tls"
	"net/http"
)

// OptTLSClientCert adds a client cert and key to the request.
func OptTLSClientCert(cert, key []byte) Option {
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
			cert, err := tls.X509KeyPair(cert, key)
			if err != nil {
				return err
			}
			typed.TLSClientConfig.Certificates = append(typed.TLSClientConfig.Certificates, cert)
		}
		return nil
	}
}
