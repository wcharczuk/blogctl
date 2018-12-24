package webutil

import (
	"net/http"
	"strings"
)

// GetProto gets the request proto.
// X-FORWARDED-PROTO is checked first, then the original request proto is used.
func GetProto(r *http.Request) string {
	if r == nil {
		return ""
	}
	tryHeader := func(key string) (string, bool) {
		if headerVal := r.Header.Get(key); len(headerVal) > 0 {
			if !strings.ContainsRune(headerVal, ',') {
				return headerVal, true
			}
			return strings.SplitN(headerVal, ",", 2)[0], true
		}
		return "", false
	}

	for _, header := range []string{"X-FORWARDED-PROTO"} {
		if headerVal, ok := tryHeader(header); ok {
			return headerVal
		}
	}

	return r.Proto
}
