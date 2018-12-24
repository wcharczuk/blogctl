package webutil

const (
	// HeaderAcceptEncoding is the "Accept-Encoding" header.
	// It indicates what types of encodings the request will accept responses as.
	// It typically enables or disables compressed (gzipped) responses.
	HeaderAcceptEncoding = "Accept-Encoding"

	// HeaderSetCookie is the header that sets cookies in a response.
	HeaderSetCookie = "Set-Cookie"

	// HeaderCookie is the request cookie header.
	HeaderCookie = "Cookie"

	// HeaderDate is the "Date" header.
	// It provides a timestamp the response was generated at.
	// It is typically used by client cache control to invalidate expired items.
	HeaderDate = "Date"

	// HeaderCacheControl is the "Cache-Control" header.
	// It indicates if and how clients should cache responses.
	// Typical values for this include "no-cache", "max-age", "min-fresh", and "max-stale" variants.
	HeaderCacheControl = "Cache-Control"

	// HeaderConnection is the "Connection" header.
	// It is used to indicate if the connection should remain open by the server
	// after the final response bytes are sent.
	// This allows the connection to be re-used, helping mitigate connection negotiation
	// penalites in making requests.
	HeaderConnection = "Connection"

	// HeaderContentEncoding is the "Content-Encoding" header.
	// It is used to indicate what the response encoding is.
	// Typical values are "gzip", "deflate", "compress", "br", and "identity" indicating no compression.
	HeaderContentEncoding = "Content-Encoding"

	// HeaderContentLength is the "Content-Length" header.
	// If provided, it specifies the size of the request or response.
	HeaderContentLength = "Content-Length"

	// HeaderContentType is the "Content-Type" header.
	// It specifies the MIME-type of the request or response.
	HeaderContentType = "Content-Type"

	// HeaderUserAgent is the user agent header.
	// It typically indicates what is making the request.
	HeaderUserAgent = "User-Agent"

	// HeaderServer is the "Server" header.
	// It is an informational header to tell the client what server software was used.
	HeaderServer = "Server"

	// HeaderVary is the "Vary" header.
	// It is used to indicate what fields should be used by the client as cache keys.
	HeaderVary = "Vary"

	// HeaderXServedBy is the "X-Served-By" header.
	// It is an informational header that indicates what software was used to generate the response.
	HeaderXServedBy = "X-Served-By"

	// HeaderXFrameOptions is the "X-Frame-Options" header.
	// It indicates if a browser is allowed to render the response in a <frame> element or not.
	HeaderXFrameOptions = "X-Frame-Options"

	// HeaderXXSSProtection is the "X-Xss-Protection" header.
	// It is a feature of internet explorer, and indicates if the browser should allow
	// requests across domain boundaries.
	HeaderXXSSProtection = "X-Xss-Protection"

	// HeaderXContentTypeOptions is the "X-Content-Type-Options" header.
	HeaderXContentTypeOptions = "X-Content-Type-Options"

	// HeaderStrictTransportSecurity is the hsts header.
	HeaderStrictTransportSecurity = "Strict-Transport-Security"

	// ContentTypeApplicationJSON is a content type for JSON responses.
	// We specify chartset=utf-8 so that clients know to use the UTF-8 string encoding.
	ContentTypeApplicationJSON = "application/json; charset=UTF-8"

	// ContentTypeHTML is a content type for html responses.
	// We specify chartset=utf-8 so that clients know to use the UTF-8 string encoding.
	ContentTypeHTML = "text/html; charset=utf-8"

	//ContentTypeXML is a content type for XML responses.
	// We specify chartset=utf-8 so that clients know to use the UTF-8 string encoding.
	ContentTypeXML = "text/xml; charset=utf-8"

	// ContentTypeText is a content type for text responses.
	// We specify chartset=utf-8 so that clients know to use the UTF-8 string encoding.
	ContentTypeText = "text/plain; charset=utf-8"

	// ConnectionKeepAlive is a value for the "Connection" header and
	// indicates the server should keep the tcp connection open
	// after the last byte of the response is sent.
	ConnectionKeepAlive = "keep-alive"

	// ContentEncodingIdentity is the identity (uncompressed) content encoding.
	ContentEncodingIdentity = "identity"
	// ContentEncodingGZIP is the gzip (compressed) content encoding.
	ContentEncodingGZIP = "gzip"
)
