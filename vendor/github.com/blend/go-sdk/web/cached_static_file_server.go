package web

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"sync"
)

// NewCachedStaticFileServer returns a new static file cache.
func NewCachedStaticFileServer(searchPaths ...http.FileSystem) *CachedStaticFileServer {
	return &CachedStaticFileServer{
		StaticFileServer: NewStaticFileServer(searchPaths...),
		files:            map[string]*CachedStaticFile{},
	}
}

// CachedStaticFileServer  is a cache of static files.
type CachedStaticFileServer struct {
	sync.Mutex
	*StaticFileServer

	files map[string]*CachedStaticFile
}

// Files returns the underlying file cache.
// Pragma; this should only be used in debugging, as during runtime locks are required to interact with this cache.
func (csfs *CachedStaticFileServer) Files() map[string]*CachedStaticFile {
	return csfs.files
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
	csfs.Lock()
	defer csfs.Unlock()

	for key, values := range csfs.StaticFileServer.headers {
		for _, value := range values {
			r.Response().Header().Set(key, value)
		}
	}

	filepath, err := r.RouteParam("filepath")
	if err != nil {
		return r.DefaultResultProvider().BadRequest(err)
	}
	if file, hasFile := csfs.files[filepath]; hasFile {
		http.ServeContent(r.Response(), r.Request(), filepath, file.ModTime, file.Contents)
		return nil
	}

	file, err := csfs.StaticFileServer.ResolveFile(filepath)
	if err != nil {
		return r.DefaultResultProvider().InternalError(err)
	}

	finfo, err := file.Stat()
	if err != nil {
		return r.DefaultResultProvider().InternalError(err)
	}

	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return r.DefaultResultProvider().InternalError(err)
	}

	csfs.files[filepath] = &CachedStaticFile{
		Path:     filepath,
		Contents: bytes.NewReader(contents),
		ModTime:  finfo.ModTime(),
		Size:     len(contents),
	}
	http.ServeContent(r.Response(), r.Request(), filepath, finfo.ModTime(), bytes.NewReader(contents))
	return nil
}
