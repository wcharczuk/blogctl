package logger

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/blend/go-sdk/webutil"
)

// FormatFileSize returns a string representation of a file size in bytes.
func FormatFileSize(sizeBytes int) string {
	if sizeBytes >= 1<<30 {
		return strconv.Itoa(sizeBytes/Gigabyte) + "gb"
	} else if sizeBytes >= 1<<20 {
		return strconv.Itoa(sizeBytes/Megabyte) + "mb"
	} else if sizeBytes >= 1<<10 {
		return strconv.Itoa(sizeBytes/Kilobyte) + "kb"
	}
	return strconv.Itoa(sizeBytes)
}

// TextWriteHTTPRequest is a helper method to write request start events to a writer.
func TextWriteHTTPRequest(tf TextFormatter, buf *bytes.Buffer, req *http.Request) {
	if ip := webutil.GetRemoteAddr(req); len(ip) > 0 {
		buf.WriteString(ip)
		buf.WriteRune(RuneSpace)
	}
	buf.WriteString(tf.Colorize(req.Method, ColorBlue))
	buf.WriteRune(RuneSpace)
	buf.WriteString(req.URL.Path)
}

// TextWriteHTTPResponse is a helper method to write request complete events to a writer.
func TextWriteHTTPResponse(tf TextFormatter, buf *bytes.Buffer, req *http.Request, statusCode, contentLength int, contentType string, elapsed time.Duration) {
	buf.WriteString(webutil.GetRemoteAddr(req))
	buf.WriteRune(RuneSpace)
	buf.WriteString(tf.Colorize(req.Method, ColorBlue))
	buf.WriteRune(RuneSpace)
	buf.WriteString(req.URL.Path)
	buf.WriteRune(RuneSpace)
	buf.WriteString(tf.ColorizeByStatusCode(statusCode, strconv.Itoa(statusCode)))
	buf.WriteRune(RuneSpace)
	buf.WriteString(elapsed.String())
	if len(contentType) > 0 {
		buf.WriteRune(RuneSpace)
		buf.WriteString(contentType)
	}
	buf.WriteRune(RuneSpace)
	buf.WriteString(FormatFileSize(contentLength))
}

// JSONWriteHTTPRequest marshals a request start as json.
func JSONWriteHTTPRequest(req *http.Request) JSONObj {
	return JSONObj{
		"ip":   webutil.GetRemoteAddr(req),
		"verb": req.Method,
		"path": req.URL.Path,
		"host": req.Host,
	}
}

// JSONWriteHTTPResponse marshals a request as json.
func JSONWriteHTTPResponse(req *http.Request, statusCode, contentLength int, contentType, contentEncoding string, elapsed time.Duration) JSONObj {
	return JSONObj{
		"ip":              webutil.GetRemoteAddr(req),
		"verb":            req.Method,
		"path":            req.URL.Path,
		"host":            req.Host,
		"contentLength":   contentLength,
		"contentType":     contentType,
		"contentEncoding": contentEncoding,
		"statusCode":      statusCode,
		JSONFieldElapsed:  Milliseconds(elapsed),
	}
}

// CompressWhitespace compresses whitespace characters into single spaces.
// It trims leading and trailing whitespace as well.
func CompressWhitespace(text string) (output string) {
	if len(text) == 0 {
		return
	}

	var state int
	for _, r := range text {
		switch state {
		case 0: // non-whitespace
			if unicode.IsSpace(r) {
				state = 1
			} else {
				output = output + string(r)
			}
		case 1: // whitespace
			if !unicode.IsSpace(r) {
				output = output + " " + string(r)
				state = 0
			}
		}
	}

	output = strings.TrimSpace(output)
	return
}
