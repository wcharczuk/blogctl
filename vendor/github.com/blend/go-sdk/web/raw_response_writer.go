package web

import (
	"net/http"
)

// NewRawResponseWriter creates a new uncompressed response writer.
func NewRawResponseWriter(w http.ResponseWriter) *RawResponseWriter {
	return &RawResponseWriter{
		innerResponse: w,
	}
}

// RawResponseWriter  a better response writer
type RawResponseWriter struct {
	innerResponse http.ResponseWriter
	contentLength int
	statusCode    int
}

// Write writes the data to the response.
func (rw *RawResponseWriter) Write(b []byte) (int, error) {
	written, err := rw.innerResponse.Write(b)
	rw.contentLength += written
	return written, err
}

// Header accesses the response header collection.
func (rw *RawResponseWriter) Header() http.Header {
	return rw.innerResponse.Header()
}

// WriteHeader is actually a terrible name and this writes the status code.
func (rw *RawResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.innerResponse.WriteHeader(code)
}

// InnerResponse returns the backing writer.
func (rw *RawResponseWriter) InnerResponse() http.ResponseWriter {
	return rw.innerResponse
}

// StatusCode returns the status code.
func (rw *RawResponseWriter) StatusCode() int {
	return rw.statusCode
}

// ContentLength returns the content length
func (rw *RawResponseWriter) ContentLength() int {
	return rw.contentLength
}

// Close disposes of the response writer.
func (rw *RawResponseWriter) Close() error {
	return nil
}
