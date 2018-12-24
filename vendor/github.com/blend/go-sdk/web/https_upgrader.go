package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/blend/go-sdk/logger"
)

// NewHTTPSUpgrader returns a new HTTPSUpgrader which redirects HTTP to HTTPS
func NewHTTPSUpgrader() *HTTPSUpgrader {
	return &HTTPSUpgrader{}
}

// NewHTTPSUpgraderFromEnv returns a new https upgrader from enviroment variables.
func NewHTTPSUpgraderFromEnv() *HTTPSUpgrader {
	return NewHTTPSUpgraderFromConfig(NewHTTPSUpgraderConfigFromEnv())
}

// NewHTTPSUpgraderFromConfig creates a new https upgrader from a config.
func NewHTTPSUpgraderFromConfig(cfg *HTTPSUpgraderConfig) *HTTPSUpgrader {
	return &HTTPSUpgrader{
		targetPort: cfg.GetTargetPort(),
	}
}

// HTTPSUpgrader redirects HTTP to HTTPS
type HTTPSUpgrader struct {
	targetPort int32
	log        *logger.Logger
}

// WithTargetPort sets the target port.
func (hu *HTTPSUpgrader) WithTargetPort(targetPort int32) *HTTPSUpgrader {
	hu.targetPort = targetPort
	return hu
}

// TargetPort returns the target port.
func (hu *HTTPSUpgrader) TargetPort() int32 {
	return hu.targetPort
}

// WithLogger sets the logger.
func (hu *HTTPSUpgrader) WithLogger(log *logger.Logger) *HTTPSUpgrader {
	hu.log = log
	return hu
}

// Logger returns the logger.
func (hu *HTTPSUpgrader) Logger() *logger.Logger {
	return hu.log
}

// ServeHTTP redirects HTTP to HTTPS
func (hu *HTTPSUpgrader) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	start := time.Now()
	response := []byte("Upgrade Required")
	if hu.log != nil {
		defer hu.log.Trigger(logger.NewHTTPResponseEvent(req).
			WithStatusCode(http.StatusMovedPermanently).
			WithContentLength(len(response)).
			WithContentType(ContentTypeText).
			WithElapsed(time.Since(start)))
	}

	newURL := *req.URL
	newURL.Scheme = SchemeHTTPS
	if len(newURL.Host) == 0 {
		newURL.Host = req.Host
	}
	if hu.targetPort > 0 {
		if strings.Contains(newURL.Host, ":") {
			newURL.Host = fmt.Sprintf("%s:%d", strings.SplitN(newURL.Host, ":", 2)[0], hu.targetPort)
		} else {
			newURL.Host = fmt.Sprintf("%s:%d", newURL.Host, hu.targetPort)
		}
	}

	http.Redirect(rw, req, newURL.String(), http.StatusMovedPermanently)
}
