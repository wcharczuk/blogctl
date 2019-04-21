package r2

import "github.com/blend/go-sdk/webutil"

const (
	// MethodGet is a method.
	MethodGet = "GET"
	// MethodPost is a method.
	MethodPost = "POST"
	// MethodPut is a method.
	MethodPut = "PUT"
	// MethodPatch is a method.
	MethodPatch = "PATCH"
	// MethodDelete is a method.
	MethodDelete = "DELETE"
	// MethodOptions is a method.
	MethodOptions = "OPTIONS"
)

const (
	// HeaderConnection is a http header.
	HeaderConnection = "Connection"
	// HeaderContentType is a http header.
	HeaderContentType = "Content-Type"
)

const (
	// ConnectionKeepAlive is a connection header value.
	ConnectionKeepAlive = "keep-alive"
)

const (
	// ContentTypeApplicationJSON is a content type header value.
	ContentTypeApplicationJSON = webutil.ContentTypeApplicationJSON
	// ContentTypeApplicationXML is a content type header value.
	ContentTypeApplicationXML = webutil.ContentTypeApplicationXML
	// ContentTypeApplicationFormEncoded is a content type header value.
	ContentTypeApplicationFormEncoded = webutil.ContentTypeApplicationFormEncoded
	// ContentTypeApplicationOctetStream is a content type header value.
	ContentTypeApplicationOctetStream = webutil.ContentTypeApplicationOctetStream
)
