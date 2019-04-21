package web

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"strings"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
)

// New returns a new web app.
func New(options ...Option) *App {
	views := NewViewCache()
	a := App{
		Latch:           async.NewLatch(),
		State:           &SyncState{},
		Statics:         map[string]*StaticFileServer{},
		DefaultHeaders:  DefaultHeaders,
		Views:           views,
		DefaultProvider: views,
	}

	for _, option := range options {
		option(&a)
	}
	return &a
}

// App is the server for the app.
type App struct {
	*async.Latch
	Auth                    AuthManager
	Config                  Config
	Log                     logger.Log
	Views                   *ViewCache
	TLSConfig               *tls.Config
	Server                  *http.Server
	Listener                *net.TCPListener
	DefaultHeaders          map[string]string
	Statics                 map[string]*StaticFileServer
	Routes                  map[string]*RouteNode
	NotFoundHandler         Handler
	MethodNotAllowedHandler Handler
	PanicAction             PanicAction
	DefaultMiddleware       []Middleware
	Tracer                  Tracer
	DefaultProvider         ResultProvider
	State                   *SyncState
}

// CreateServer returns the basic http.Server for the app.
func (a *App) CreateServer() *http.Server {
	return &http.Server{
		Handler:           a,
		TLSConfig:         a.TLSConfig,
		Addr:              a.Config.BindAddrOrDefault(),
		MaxHeaderBytes:    a.Config.MaxHeaderBytesOrDefault(),
		ReadTimeout:       a.Config.ReadTimeoutOrDefault(),
		ReadHeaderTimeout: a.Config.ReadHeaderTimeoutOrDefault(),
		WriteTimeout:      a.Config.WriteTimeoutOrDefault(),
		IdleTimeout:       a.Config.IdleTimeoutOrDefault(),
	}
}

// StartupTasks runs common startup tasks.
func (a *App) StartupTasks() error {
	return a.Views.Initialize()
}

// Start starts the server and binds to the given address.
func (a *App) Start() (err error) {
	// set up the underlying server.
	a.Server = a.CreateServer()

	// initialize the view cache.
	err = a.StartupTasks()
	if err != nil {
		return
	}

	serverProtocol := "http"
	if a.Server.TLSConfig != nil {
		serverProtocol = "https (tls)"
	}

	logger.MaybeInfof(a.Log, "%s server started, listening on %s", serverProtocol, a.Config.BindAddrOrDefault())

	if a.Server.TLSConfig != nil && a.Server.TLSConfig.ClientCAs != nil {
		logger.MaybeInfof(a.Log, "%s using client cert pool with (%d) client certs", serverProtocol, len(a.Server.TLSConfig.ClientCAs.Subjects()))
	}

	var listener net.Listener
	listener, err = net.Listen("tcp", a.Config.BindAddrOrDefault())
	if err != nil {
		err = ex.New(err)
		return
	}
	var ok bool
	a.Listener, ok = listener.(*net.TCPListener)
	if !ok {
		err = ex.New("listener returned was not a net.TCPListener")
		return
	}

	keepAliveListener := TCPKeepAliveListener{a.Listener}
	var shutdownErr error

	a.Started()
	if a.Server.TLSConfig != nil {
		shutdownErr = a.Server.Serve(tls.NewListener(keepAliveListener, a.Server.TLSConfig))
	} else {
		shutdownErr = a.Server.Serve(keepAliveListener)
	}
	if shutdownErr != nil && shutdownErr != http.ErrServerClosed {
		err = ex.New(shutdownErr)
	}
	logger.MaybeInfof(a.Log, "server exited")
	a.Stopped()

	return
}

// Stop stops the server.
func (a *App) Stop() error {
	if !a.CanStop() {
		return ex.New(async.ErrCannotStop)
	}
	a.Stopping()

	ctx := context.Background()
	var cancel context.CancelFunc
	if a.Config.ShutdownGracePeriodOrDefault() > 0 {
		ctx, cancel = context.WithTimeout(ctx, a.Config.ShutdownGracePeriodOrDefault())
		defer cancel()
	}
	logger.MaybeInfof(a.Log, "server shutting down")
	a.Server.SetKeepAlivesEnabled(false)
	if err := a.Server.Shutdown(ctx); err != nil {
		return ex.New(err)
	}

	a.Server = nil
	a.Listener = nil
	logger.MaybeInfof(a.Log, "server shutdown complete")

	return nil
}

// Register registers controllers with the app's router.
func (a *App) Register(controllers ...Controller) {
	for _, c := range controllers {
		c.Register(a)
	}
}

// --------------------------------------------------------------------------------
// Static Result Methods
// --------------------------------------------------------------------------------

// SetStaticRewriteRule adds a rewrite rule for a specific statically served path.
// It mutates the path for the incoming static file request to the fileserver according to the action.
func (a *App) SetStaticRewriteRule(route, match string, action RewriteAction) error {
	mountedRoute := a.formatStaticMountRoute(route)
	if static, hasRoute := a.Statics[mountedRoute]; hasRoute {
		return static.AddRewriteRule(match, action)
	}
	return ex.New("no static fileserver mounted at route", ex.OptMessagef("route: %s", route))
}

// SetStaticHeader adds a header for the given static path.
// These headers are automatically added to any result that the static path fileserver sends.
func (a *App) SetStaticHeader(route, key, value string) error {
	mountedRoute := a.formatStaticMountRoute(route)
	if static, hasRoute := a.Statics[mountedRoute]; hasRoute {
		static.AddHeader(key, value)
		return nil
	}
	return ex.New("no static fileserver mounted at route", ex.OptMessagef("route: %s", mountedRoute))
}

// ServeStatic serves files from the given file system root(s)..
// If the path does not end with "/*filepath" that suffix will be added for you internally.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
func (a *App) ServeStatic(route string, searchPaths []string, middleware ...Middleware) {
	var searchPathFS []http.FileSystem
	for _, searchPath := range searchPaths {
		searchPathFS = append(searchPathFS, http.Dir(searchPath))
	}
	sfs := NewStaticFileServer(searchPathFS...)
	sfs.Middleware = middleware
	sfs.CacheDisabled = true
	mountedRoute := a.formatStaticMountRoute(route)
	a.Statics[mountedRoute] = sfs
	a.Handle("GET", mountedRoute, a.RenderAction(a.Middleware(sfs.Action, middleware...)))
}

// ServeStaticCached serves files from the given file system root(s).
// If the path does not end with "/*filepath" that suffix will be added for you internally.
func (a *App) ServeStaticCached(route string, searchPaths []string, middleware ...Middleware) {
	var searchPathFileSystems []http.FileSystem
	for _, searchPath := range searchPaths {
		searchPathFileSystems = append(searchPathFileSystems, http.Dir(searchPath))
	}
	sfs := NewStaticFileServer(searchPathFileSystems...)
	sfs.Middleware = middleware
	mountedRoute := a.formatStaticMountRoute(route)
	a.Statics[mountedRoute] = sfs
	a.Handle("GET", mountedRoute, a.RenderAction(a.Middleware(sfs.Action, middleware...)))
}

func (a *App) formatStaticMountRoute(route string) string {
	mountedRoute := route
	if !strings.HasSuffix(mountedRoute, "*"+RouteTokenFilepath) {
		if strings.HasSuffix(mountedRoute, "/") {
			mountedRoute = mountedRoute + "*" + RouteTokenFilepath
		} else {
			mountedRoute = mountedRoute + "/*" + RouteTokenFilepath
		}
	}
	return mountedRoute
}

// --------------------------------------------------------------------------------
// Route Registration / HTTP Methods
// --------------------------------------------------------------------------------

// GET registers a GET request handler.
/*
Routes should be registered in the form:

	app.GET("/myroute", myAction, myMiddleware...)

It is important to note that routes are registered in order and
cannot have any wildcards inside the routes.
*/
func (a *App) GET(path string, action Action, middleware ...Middleware) {
	a.Handle("GET", path, a.RenderAction(a.Middleware(action, middleware...)))
}

// OPTIONS registers a OPTIONS request handler.
func (a *App) OPTIONS(path string, action Action, middleware ...Middleware) {
	a.Handle("OPTIONS", path, a.RenderAction(a.Middleware(action, middleware...)))
}

// HEAD registers a HEAD request handler.
func (a *App) HEAD(path string, action Action, middleware ...Middleware) {
	a.Handle("HEAD", path, a.RenderAction(a.Middleware(action, middleware...)))
}

// PUT registers a PUT request handler.
func (a *App) PUT(path string, action Action, middleware ...Middleware) {
	a.Handle("PUT", path, a.RenderAction(a.Middleware(action, middleware...)))
}

// PATCH registers a PATCH request handler.
func (a *App) PATCH(path string, action Action, middleware ...Middleware) {
	a.Handle("PATCH", path, a.RenderAction(a.Middleware(action, middleware...)))
}

// POST registers a POST request actions.
func (a *App) POST(path string, action Action, middleware ...Middleware) {
	a.Handle("POST", path, a.RenderAction(a.Middleware(action, middleware...)))
}

// DELETE registers a DELETE request handler.
func (a *App) DELETE(path string, action Action, middleware ...Middleware) {
	a.Handle("DELETE", path, a.RenderAction(a.Middleware(action, middleware...)))
}

// Handle adds a raw handler at a given method and path.
func (a *App) Handle(method, path string, handler Handler) {
	if len(path) == 0 {
		panic("path must not be empty")
	}
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}
	if a.Routes == nil {
		a.Routes = make(map[string]*RouteNode)
	}

	root := a.Routes[method]
	if root == nil {
		root = new(RouteNode)
		a.Routes[method] = root
	}
	root.addRoute(method, path, handler)
}

// Lookup finds the route data for a given method and path.
func (a *App) Lookup(method, path string) (route *Route, params RouteParameters, skipSlashRedirect bool) {
	if root := a.Routes[method]; root != nil {
		return root.getValue(path)
	}
	return nil, nil, false
}

// --------------------------------------------------------------------------------
// Request Pipeline
// --------------------------------------------------------------------------------

// ServeHTTP makes the router implement the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !a.Config.DisablePanicRecovery {
		defer a.recover(w, req)
	}

	path := req.URL.Path
	if root := a.Routes[req.Method]; root != nil {
		if route, params, tsr := root.getValue(path); route != nil {
			route.Handler(w, req, route, params)
			return
		} else if req.Method != MethodConnect && path != "/" {
			code := http.StatusMovedPermanently // 301 // Permanent redirect, request with GET method
			if req.Method != MethodGet {
				code = http.StatusTemporaryRedirect // 307
			}

			if tsr && !a.Config.SkipRedirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					req.URL.Path = path[:len(path)-1]
				} else {
					req.URL.Path = path + "/"
				}
				http.Redirect(w, req, req.URL.String(), code)
				return
			}
		}
	}

	if req.Method == MethodOptions {
		// Handle OPTIONS requests
		if a.Config.HandleOptions {
			if allow := a.allowed(path, req.Method); len(allow) > 0 {
				w.Header().Set(HeaderAllow, allow)
				return
			}
		}
	} else {
		// Handle 405
		if a.Config.HandleMethodNotAllowed {
			if allow := a.allowed(path, req.Method); len(allow) > 0 {
				w.Header().Set(HeaderAllow, allow)
				if a.MethodNotAllowedHandler != nil {
					a.MethodNotAllowedHandler(w, req, nil, nil)
				} else {
					http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				}
				return
			}
		}
	}

	// Handle 404
	if a.NotFoundHandler != nil {
		a.NotFoundHandler(w, req, nil, nil)
	} else {
		http.NotFound(w, req)
	}
}

// RenderAction is the translation step from Action to Handler.
// this is where the bulk of the "pipeline" happens.
func (a *App) RenderAction(action Action) Handler {
	return func(w http.ResponseWriter, r *http.Request, route *Route, p RouteParameters) {
		var err error
		var tf TraceFinisher

		var response ResponseWriter
		if strings.Contains(r.Header.Get(HeaderAcceptEncoding), ContentEncodingGZIP) {
			w.Header().Set(HeaderContentEncoding, ContentEncodingGZIP)
			response = NewCompressedResponseWriter(w)
		} else {
			w.Header().Set(HeaderContentEncoding, ContentEncodingIdentity)
			response = NewRawResponseWriter(w)
		}

		ctx := a.createCtx(response, r, route, p)
		ctx.onRequestStart()
		if a.Tracer != nil {
			tf = a.Tracer.Start(ctx)
		}
		if a.Log != nil {
			a.Log.Trigger(r.Context(), a.httpRequestEvent(ctx))
		}

		if len(a.DefaultHeaders) > 0 {
			for key, value := range a.DefaultHeaders {
				response.Header().Set(key, value)
			}
		}
		result := action(ctx)
		if result != nil {
			// check for a prerender step
			if typed, ok := result.(ResultPreRender); ok {
				if preRenderErr := typed.PreRender(ctx); preRenderErr != nil {
					err = ex.Nest(err, preRenderErr)
				}
			}

			// do the render, log any errors emitted
			if resultErr := result.Render(ctx); resultErr != nil {
				err = ex.Nest(err, resultErr)
			}

			// check for a render complete step
			// typically this is used to render error results if there was a problem rendering
			// the result.
			if typed, ok := result.(ResultPostRender); ok {
				if postRenderErr := typed.PostRender(ctx); postRenderErr != nil {
					err = ex.Nest(err, postRenderErr)
				}
			}
		}

		ctx.onRequestFinish()
		response.Close()

		if err != nil {
			a.logFatal(err, r)
		}
		if a.Log != nil {
			a.Log.Trigger(r.Context(), a.httpResponseEvent(ctx))
		}
		if tf != nil {
			tf.Finish(ctx, err)
		}
	}
}

// Middleware wraps an action with a given set of middleware, including app level default middleware.
func (a *App) Middleware(action Action, middleware ...Middleware) Action {
	if len(middleware) == 0 && len(a.DefaultMiddleware) == 0 {
		return action
	}

	finalMiddleware := make([]Middleware, len(middleware)+len(a.DefaultMiddleware))
	cursor := len(finalMiddleware) - 1
	for i := len(a.DefaultMiddleware) - 1; i >= 0; i-- {
		finalMiddleware[cursor] = a.DefaultMiddleware[i]
		cursor--
	}

	for i := len(middleware) - 1; i >= 0; i-- {
		finalMiddleware[cursor] = middleware[i]
		cursor--
	}

	return NestMiddleware(action, finalMiddleware...)
}

//
// internal helpers
//

func (a *App) createCtx(w ResponseWriter, r *http.Request, route *Route, p RouteParameters, extra ...CtxOption) *Ctx {
	options := []CtxOption{
		OptCtxApp(a),
		OptCtxRoute(route),
		OptCtxRouteParams(p),
		OptCtxState(a.State.Copy()),
		OptCtxTracer(a.Tracer),
		OptCtxViews(a.Views),
		OptCtxAuth(a.Auth),
		OptCtxLog(a.Log),
		OptCtxDefaultProvider(a.DefaultProvider),
	}
	return NewCtx(w, r, append(options, extra...)...)
}

func (a *App) allowed(path, reqMethod string) (allow string) {
	if path == "*" { // server-wide
		for method := range a.Routes {
			if method == "OPTIONS" {
				continue
			}

			// add request method to list of allowed methods
			if len(allow) == 0 {
				allow = method
			} else {
				allow += ", " + method
			}
		}
		return
	}
	for method := range a.Routes {
		// Skip the requested method - we already tried this one
		if method == reqMethod || method == "OPTIONS" {
			continue
		}

		handle, _, _ := a.Routes[method].getValue(path)
		if handle != nil {
			// add request method to list of allowed methods
			if len(allow) == 0 {
				allow = method
			} else {
				allow += ", " + method
			}
		}
	}
	if len(allow) > 0 {
		allow += ", OPTIONS"
	}
	return
}

func (a *App) httpRequestEvent(ctx *Ctx) *logger.HTTPRequestEvent {
	event := logger.NewHTTPRequestEvent(ctx.Request)
	if ctx.Route != nil {
		event.Route = ctx.Route.String()
	}
	return event
}

func (a *App) httpResponseEvent(ctx *Ctx) *logger.HTTPResponseEvent {
	event := logger.NewHTTPResponseEvent(ctx.Request,
		logger.OptHTTPResponseStatusCode(ctx.Response.StatusCode()),
		logger.OptHTTPResponseContentLength(ctx.Response.ContentLength()),
		logger.OptHTTPResponseElapsed(ctx.Elapsed()),
	)

	if ctx.Route != nil {
		event.Route = ctx.Route.String()
	}

	if ctx.Response.Header() != nil {
		event.ContentType = ctx.Response.Header().Get(HeaderContentType)
		event.ContentEncoding = ctx.Response.Header().Get(HeaderContentEncoding)
	}
	return event
}

func (a *App) recover(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		err := ex.New(rcv)
		a.logFatal(err, req)
		if a.PanicAction != nil {
			a.handlePanic(w, req, rcv)
		} else {
			http.Error(w, "an internal server error occurred", http.StatusInternalServerError)
		}
	}
}

func (a *App) handlePanic(w http.ResponseWriter, r *http.Request, err interface{}) {
	a.RenderAction(func(ctx *Ctx) Result {
		if a.Log != nil {
			a.Log.Fatalf("%v", err)
		}
		return a.PanicAction(ctx, err)
	})(w, r, nil, nil)
}

func (a *App) logFatal(err error, req *http.Request) {
	if a.Log == nil {
		return
	}
	if err != nil {
		a.Log.Trigger(req.Context(), logger.NewErrorEvent(logger.Fatal, err, logger.OptErrorEventState(req)))
	}
}
