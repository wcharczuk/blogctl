package web

import (
	"net/http"
	"os"
	"regexp"

	"github.com/blend/go-sdk/logger"
)

// NewStaticFileServer returns a new static file cache.
func NewStaticFileServer(fs http.FileSystem) *StaticFileServer {
	return &StaticFileServer{
		fileSystem: fs,
	}
}

// StaticFileServer is a cache of static files.
type StaticFileServer struct {
	log          logger.Log
	fileSystem   http.FileSystem
	rewriteRules []RewriteRule
	middleware   Action
	headers      http.Header
}

// Log returns a logger reference.
func (sc *StaticFileServer) Log() logger.Log {
	return sc.log
}

// WithLogger sets the logger reference for the static file cache.
func (sc *StaticFileServer) WithLogger(log logger.Log) *StaticFileServer {
	sc.log = log
	return sc
}

// AddHeader adds a header to the static cache results.
func (sc *StaticFileServer) AddHeader(key, value string) {
	if sc.headers == nil {
		sc.headers = http.Header{}
	}
	sc.headers[key] = append(sc.headers[key], value)
}

// Headers returns the headers for the static server.
func (sc *StaticFileServer) Headers() http.Header {
	return sc.headers
}

// AddRewriteRule adds a static re-write rule.
func (sc *StaticFileServer) AddRewriteRule(match string, action RewriteAction) error {
	expr, err := regexp.Compile(match)
	if err != nil {
		return err
	}
	sc.rewriteRules = append(sc.rewriteRules, RewriteRule{
		MatchExpression: match,
		expr:            expr,
		Action:          action,
	})
	return nil
}

// SetMiddleware sets the middlewares.
func (sc *StaticFileServer) SetMiddleware(middlewares ...Middleware) {
	sc.middleware = NestMiddleware(sc.ServeFile, middlewares...)
}

// RewriteRules returns the rewrite rules
func (sc *StaticFileServer) RewriteRules() []RewriteRule {
	return sc.rewriteRules
}

// Action is the entrypoint for the static server.
// It will run middleware if specified before serving the file.
func (sc *StaticFileServer) Action(r *Ctx) Result {
	if sc.middleware != nil {
		return sc.middleware(r)
	}
	return sc.ServeFile(r)
}

// ServeFile writes the file to the response without running middleware.
func (sc *StaticFileServer) ServeFile(r *Ctx) Result {
	for key, values := range sc.headers {
		for _, value := range values {
			r.Response().Header().Set(key, value)
		}
	}

	filePath, err := r.RouteParam("filepath")
	if err != nil {
		return r.DefaultResultProvider().InternalError(err)
	}

	for _, rule := range sc.rewriteRules {
		if matched, newFilePath := rule.Apply(filePath); matched {
			filePath = newFilePath
		}
	}

	f, err := sc.fileSystem.Open(filePath)
	if f == nil || os.IsNotExist(err) {
		return r.DefaultResultProvider().NotFound()
	}
	if err != nil {
		return r.DefaultResultProvider().InternalError(err)
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		return r.DefaultResultProvider().InternalError(err)
	}

	http.ServeContent(r.Response(), r.Request(), filePath, d.ModTime(), f)
	return nil

}
