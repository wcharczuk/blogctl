package web

import (
	"compress/gzip"
	"net/http"
)

// NewCompressedResponseWriter returns a new gzipped response writer.
func NewCompressedResponseWriter(w http.ResponseWriter) *CompressedResponseWriter {
	return &CompressedResponseWriter{
		innerResponse: w,
	}
}

// CompressedResponseWriter is a response writer that compresses output.
type CompressedResponseWriter struct {
	gzipWriter    *gzip.Writer
	innerResponse http.ResponseWriter
	statusCode    int
	contentLength int
}

func (crw *CompressedResponseWriter) ensureCompressedStream() {
	if crw.gzipWriter == nil {
		crw.gzipWriter = gzip.NewWriter(crw.innerResponse)
	}
}

// Write writes the byes to the stream.
func (crw *CompressedResponseWriter) Write(b []byte) (int, error) {
	crw.ensureCompressedStream()
	_, err := crw.gzipWriter.Write(b)
	crw.contentLength += len(b)
	return len(b), err
}

// Header returns the headers for the response.
func (crw *CompressedResponseWriter) Header() http.Header {
	return crw.innerResponse.Header()
}

// WriteHeader writes a status code.
func (crw *CompressedResponseWriter) WriteHeader(code int) {
	crw.statusCode = code
	crw.innerResponse.WriteHeader(code)
}

// InnerResponse returns the backing http response.
func (crw *CompressedResponseWriter) InnerResponse() http.ResponseWriter {
	return crw.innerResponse
}

// StatusCode returns the status code for the request.
func (crw *CompressedResponseWriter) StatusCode() int {
	return crw.statusCode
}

// ContentLength returns the content length for the request.
func (crw *CompressedResponseWriter) ContentLength() int {
	return crw.contentLength
}

// Flush pushes any buffered data out to the response.
func (crw *CompressedResponseWriter) Flush() error {
	crw.ensureCompressedStream()
	return crw.gzipWriter.Flush()
}

// Close closes any underlying resources.
func (crw *CompressedResponseWriter) Close() error {
	if crw.gzipWriter != nil {
		err := crw.gzipWriter.Close()
		crw.gzipWriter = nil
		return err
	}
	return nil
}
