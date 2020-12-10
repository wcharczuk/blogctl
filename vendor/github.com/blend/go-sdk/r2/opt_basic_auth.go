package r2

import "github.com/blend/go-sdk/webutil"

// OptBasicAuth is an option that sets the http basic auth.
func OptBasicAuth(username, password string) Option {
	return RequestOption(webutil.OptBasicAuth(username, password))
}
