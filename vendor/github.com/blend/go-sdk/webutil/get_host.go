package webutil

import (
	"net/http"
	"strings"
)

// GetHost returns the request host, omiting the port if specified.
func GetHost(r *http.Request) string {
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

	for _, header := range []string{"X-FORWARDED-HOST"} {
		if headerVal, ok := tryHeader(header); ok {
			return headerVal
		}
	}
	if r.URL != nil && len(r.URL.Host) > 0 {
		return r.URL.Host
	}
	if strings.Contains(r.Host, ":") {
		return strings.SplitN(r.Host, ":", 2)[0]
	}
	return r.Host
}
