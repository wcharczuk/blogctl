package r2

import (
	"fmt"
	"net/url"
)

// OptPort sets a custom port for the request url.
func OptPort(port int32) Option {
	return func(r *Request) error {
		if r.URL == nil {
			r.URL = &url.URL{}
		}
		r.URL.Host = fmt.Sprintf("%s:%d", r.URL.Hostname(), port)
		return nil
	}
}
