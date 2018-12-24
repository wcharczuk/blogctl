package web

import (
	"net/http"
)

// Handler is the most basic route handler.
type Handler func(http.ResponseWriter, *http.Request, *Route, RouteParameters)

// WrapHandler wraps an http.Handler as a Handler.
func WrapHandler(handler http.Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request, _ *Route, _ RouteParameters) {
		handler.ServeHTTP(w, r)
	}
}

// PanicHandler is a handler for panics that also takes an error.
type PanicHandler func(http.ResponseWriter, *http.Request, interface{})

// Route is an entry in the route tree.
type Route struct {
	Handler
	Method string
	Path   string
	Params []string
}

// String returns the path.
func (r Route) String() string { return r.Path }

// StringWithMethod returns a string representation of the route.
// Namely: Method_Path
func (r Route) StringWithMethod() string {
	return r.Method + "_" + r.Path
}
