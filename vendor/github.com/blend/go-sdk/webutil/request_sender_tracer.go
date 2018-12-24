package webutil

import "net/http"

// RequestTracer is a tracer for requests.
type RequestTracer interface {
	Start(*http.Request) RequestTraceFinisher
}

// RequestTraceFinisher is a finisher for request traces.
type RequestTraceFinisher interface {
	Finish(*http.Request, *http.Response, error)
}
