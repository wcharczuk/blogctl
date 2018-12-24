package web

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/blend/go-sdk/logger"
)

// NewCachedStaticFileServer returns a new static file cache.
func NewCachedStaticFileServer(fs http.FileSystem) *CachedStaticFileServer {
	return &CachedStaticFileServer{
		fileSystem: fs,
		files:      map[string]*CachedStaticFile{},
	}
}

// CachedStaticFileServer  is a cache of static files.
type CachedStaticFileServer struct {
	log          *logger.Logger
	fileSystem   http.FileSystem
	syncRoot     sync.Mutex
	rewriteRules []RewriteRule
	headers      http.Header
	middleware   Action
	files        map[string]*CachedStaticFile
}

// Log returns a logger reference.
func (csfs *CachedStaticFileServer) Log() *logger.Logger {
	return csfs.log
}

// WithLogger sets the logger reference for the static file cache.
func (csfs *CachedStaticFileServer) WithLogger(log *logger.Logger) *CachedStaticFileServer {
	csfs.log = log
	return csfs
}

// Files returns the underlying file cache.
// Pragma; this should only be used in debugging, as during runtime locks are required to interact with this cache.
func (csfs *CachedStaticFileServer) Files() map[string]*CachedStaticFile {
	return csfs.files
}

// AddHeader adds a header to the static cache results.
func (csfs *CachedStaticFileServer) AddHeader(key, value string) {
	if csfs.headers == nil {
		csfs.headers = http.Header{}
	}
	csfs.headers[key] = append(csfs.headers[key], value)
}

// Headers returns the headers for the static server.
func (csfs *CachedStaticFileServer) Headers() http.Header {
	return csfs.headers
}

// AddRewriteRule adds a static re-write rule.
func (csfs *CachedStaticFileServer) AddRewriteRule(match string, action RewriteAction) error {
	expr, err := regexp.Compile(match)
	if err != nil {
		return err
	}
	csfs.rewriteRules = append(csfs.rewriteRules, RewriteRule{
		MatchExpression: match,
		expr:            expr,
		Action:          action,
	})
	return nil
}

// RewriteRules returns the rewrite rules
func (csfs *CachedStaticFileServer) RewriteRules() []RewriteRule {
	return csfs.rewriteRules
}

// SetMiddleware sets the middlewares.
func (csfs *CachedStaticFileServer) SetMiddleware(middlewares ...Middleware) {
	csfs.middleware = NestMiddleware(csfs.ServeFile, middlewares...)
}

// GetCachedFile returns a file from the filesystem at a given path.
func (csfs *CachedStaticFileServer) GetCachedFile(filepath string) (*CachedStaticFile, error) {
	csfs.syncRoot.Lock()
	defer csfs.syncRoot.Unlock()

	if file, hasFile := csfs.files[filepath]; hasFile {
		return file, nil
	}

	f, err := csfs.fileSystem.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		return nil, err
	}

	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	newFile := &CachedStaticFile{
		Path:     filepath,
		Size:     len(contents),
		ModTime:  d.ModTime(),
		Contents: bytes.NewReader(contents),
	}

	csfs.files[filepath] = newFile
	return newFile, nil
}

// Action is the entrypoint for the static server.
func (csfs *CachedStaticFileServer) Action(r *Ctx) Result {
	if csfs.middleware != nil {
		return csfs.middleware(r)
	}
	return csfs.ServeFile(r)
}

// ServeFile writes the file to the response.
func (csfs *CachedStaticFileServer) ServeFile(r *Ctx) Result {
	filePath, err := r.RouteParam(RouteTokenFilepath)
	if err != nil {
		return r.DefaultResultProvider().InternalError(err)
	}

	for key, values := range csfs.headers {
		for _, value := range values {
			r.Response().Header().Set(key, value)
		}
	}

	for _, rule := range csfs.rewriteRules {
		if matched, newFilePath := rule.Apply(filePath); matched {
			filePath = newFilePath
		}
	}

	f, err := csfs.GetCachedFile(filePath)
	if f == nil || os.IsNotExist(err) {
		return r.DefaultResultProvider().NotFound()
	}
	if err != nil {
		return r.DefaultResultProvider().InternalError(err)
	}

	http.ServeContent(r.Response(), r.Request(), filePath, f.ModTime, f.Contents)
	return nil
}
