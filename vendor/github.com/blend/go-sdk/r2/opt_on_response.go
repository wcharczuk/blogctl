package r2

import (
	"net/http"
	"time"
)

// OnResponseListener is an on response listener.
type OnResponseListener func(*http.Request, *http.Response, time.Time, error) error

// OptOnResponse adds an on response listener.
// If an OnResponse listener has already been addded, it will be merged with the existing listener.
func OptOnResponse(listener OnResponseListener) Option {
	return func(r *Request) error {
		r.OnResponse = append(r.OnResponse, listener)
		return nil
	}
}
