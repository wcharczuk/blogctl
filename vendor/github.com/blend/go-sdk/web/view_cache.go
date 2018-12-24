package web

import (
	"html/template"
	"net/http"
	"sync"

	"github.com/blend/go-sdk/exception"
	templatehelpers "github.com/blend/go-sdk/template"
)

const (
	// DefaultTemplateNameBadRequest is the default template name for bad request view results.
	DefaultTemplateNameBadRequest = "bad_request"
	// DefaultTemplateNameInternalError is the default template name for internal server error view results.
	DefaultTemplateNameInternalError = "error"
	// DefaultTemplateNameNotFound is the default template name for not found error view results.
	DefaultTemplateNameNotFound = "not_found"
	// DefaultTemplateNameNotAuthorized is the default template name for not authorized error view results.
	DefaultTemplateNameNotAuthorized = "not_authorized"
	// DefaultTemplateNameStatus is the default template name for status view results.
	DefaultTemplateNameStatus = "status"

	// DefaultTemplateBadRequest is a basic view.
	DefaultTemplateBadRequest = `<html><head><style>body { font-family: sans-serif; text-align: center; }</style></head><body><h4>Bad Request</h4></body><pre>{{ .ViewModel }}</pre></html>`
	// DefaultTemplateInternalError is a basic view.
	DefaultTemplateInternalError = `<html><head><style>body { font-family: sans-serif; text-align: center; }</style></head><body><h4>Internal Error</h4><pre>{{ .ViewModel }}</body></html>`
	// DefaultTemplateNotAuthorized is a basic view.
	DefaultTemplateNotAuthorized = `<html><head><style>body { font-family: sans-serif; text-align: center; }</style></head><body><h4>Not Authorized</h4></body></html>`
	// DefaultTemplateNotFound is a basic view.
	DefaultTemplateNotFound = `<html><head><style>body { font-family: sans-serif; text-align: center; }</style></head><body><h4>Not Found</h4></body></html>`
	// DefaultTemplateStatus is a basic view.
	DefaultTemplateStatus = `<html><head><style>body { font-family: sans-serif; text-align: center; }</style></head><body><h4>{{ .ViewModel.StatusCode }}</h4></body><pre>{{ .ViewModel.Response }}</pre></html>`
)

// Assert the view cache is a result provider.
var (
	_ ResultProvider = (*ViewCache)(nil)
)

// NewViewCache returns a new view cache.
func NewViewCache() *ViewCache {
	return &ViewCache{
		viewFuncMap:               template.FuncMap(templatehelpers.ViewFuncs{}.FuncMap()),
		viewCache:                 template.New(""), // an empty template tree.
		bufferPool:                NewBufferPool(32),
		cached:                    true,
		internalErrorTemplateName: DefaultTemplateNameInternalError,
		badRequestTemplateName:    DefaultTemplateNameBadRequest,
		notFoundTemplateName:      DefaultTemplateNameNotFound,
		notAuthorizedTemplateName: DefaultTemplateNameNotAuthorized,
		statusTemplateName:        DefaultTemplateNameStatus,
	}
}

// NewViewCacheFromConfig returns a new view cache from a config.
func NewViewCacheFromConfig(cfg *ViewCacheConfig) *ViewCache {
	return &ViewCache{
		viewFuncMap:               template.FuncMap(templatehelpers.ViewFuncs{}.FuncMap()),
		viewCache:                 template.New(""), // an empty template tree.
		bufferPool:                NewBufferPool(cfg.GetBufferPoolSize()),
		viewPaths:                 cfg.GetPaths(),
		cached:                    cfg.GetCached(),
		internalErrorTemplateName: cfg.GetInternalErrorTemplateName(),
		badRequestTemplateName:    cfg.GetBadRequestTemplateName(),
		notFoundTemplateName:      cfg.GetNotFoundTemplateName(),
		notAuthorizedTemplateName: cfg.GetNotAuthorizedTemplateName(),
		statusTemplateName:        cfg.GetStatusTemplateName(),
	}
}

// ViewCache is the cached views used in view results.
type ViewCache struct {
	viewFuncMap  template.FuncMap
	viewPaths    []string
	viewLiterals []string
	viewCache    *template.Template
	cached       bool

	bufferPool *BufferPool

	initializedLock sync.Mutex
	initialized     bool

	badRequestTemplateName    string
	internalErrorTemplateName string
	notFoundTemplateName      string
	notAuthorizedTemplateName string
	statusTemplateName        string
}

// Initialize caches templates by path.
func (vc *ViewCache) Initialize() error {
	if !vc.initialized {
		vc.initializedLock.Lock()
		defer vc.initializedLock.Unlock()

		if !vc.initialized {
			err := vc.initialize()
			if err != nil {
				return err
			}
			vc.initialized = true
		}
	}

	return nil
}

// Parse parses the view tree.
func (vc *ViewCache) Parse() (views *template.Template, err error) {
	views = template.New("").Funcs(vc.viewFuncMap)
	if len(vc.viewPaths) > 0 {
		views, err = views.ParseFiles(vc.viewPaths...)
		if err != nil {
			err = exception.New(err)
			return
		}
	}

	if len(vc.viewLiterals) > 0 {
		for _, viewLiteral := range vc.viewLiterals {
			views, err = views.Parse(viewLiteral)
			if err != nil {
				err = exception.New(err)
				return
			}
		}
	}
	return
}

// Lookup looks up a view.
func (vc *ViewCache) Lookup(name string) (*template.Template, error) {
	views, err := vc.Templates()
	if err != nil {
		return nil, err
	}

	return views.Lookup(name), nil
}

// SetTemplates sets the view cache for the app.
func (vc *ViewCache) SetTemplates(viewCache *template.Template) {
	vc.viewCache = viewCache
}

// ----------------------------------------------------------------------
// results
// ----------------------------------------------------------------------

// BadRequest returns a view result.
func (vc *ViewCache) BadRequest(err error) Result {
	t, viewErr := vc.Lookup(vc.BadRequestTemplateName())
	if viewErr != nil {
		return vc.viewError(viewErr)
	}
	if t == nil {
		t, _ = template.New("default").Parse(DefaultTemplateBadRequest)
	}

	return &ViewResult{
		ViewName:   vc.BadRequestTemplateName(),
		StatusCode: http.StatusBadRequest,
		ViewModel:  err,
		Template:   t,
		Views:      vc,
	}
}

// InternalError returns a view result.
func (vc *ViewCache) InternalError(err error) Result {
	t, viewErr := vc.Lookup(vc.InternalErrorTemplateName())
	if viewErr != nil {
		return vc.viewError(viewErr)
	}
	if t == nil {
		t, _ = template.New("").Parse(DefaultTemplateInternalError)
	}
	return resultWithLoggedError(&ViewResult{
		ViewName:   vc.InternalErrorTemplateName(),
		StatusCode: http.StatusInternalServerError,
		ViewModel:  err,
		Template:   t,
		Views:      vc,
	}, err)
}

// NotFound returns a view result.
func (vc *ViewCache) NotFound() Result {
	t, viewErr := vc.Lookup(vc.NotFoundTemplateName())
	if viewErr != nil {
		return vc.viewError(viewErr)
	}
	if t == nil {
		t, _ = template.New("").Parse(DefaultTemplateNotFound)
	}
	return &ViewResult{
		ViewName:   vc.NotFoundTemplateName(),
		StatusCode: http.StatusNotFound,
		Template:   t,
		Views:      vc,
	}
}

// NotAuthorized returns a view result.
func (vc *ViewCache) NotAuthorized() Result {
	t, err := vc.Lookup(vc.NotAuthorizedTemplateName())
	if err != nil {
		return vc.viewError(err)
	}
	if t == nil {
		t, _ = template.New("").Parse(DefaultTemplateNotAuthorized)
	}

	return &ViewResult{
		ViewName:   vc.NotAuthorizedTemplateName(),
		StatusCode: http.StatusForbidden,
		Template:   t,
		Views:      vc,
	}
}

// Status returns a status view result.
func (vc *ViewCache) Status(statusCode int, response ...interface{}) Result {
	t, viewErr := vc.Lookup(vc.StatusTemplateName())
	if viewErr != nil {
		return vc.viewError(viewErr)
	}
	if t == nil {
		t, _ = template.New("").Parse(DefaultTemplateStatus)
	}

	return &ViewResult{
		Views:      vc,
		ViewName:   vc.StatusTemplateName(),
		StatusCode: statusCode,
		Template:   t,
		ViewModel:  StatusViewModel{StatusCode: statusCode, Response: ResultOrDefault(http.StatusText(statusCode), response...)},
	}
}

// View returns a view result.
func (vc *ViewCache) View(viewName string, viewModel interface{}) Result {
	return vc.ViewStatus(http.StatusOK, viewName, viewModel)
}

// ViewStatus returns a view result with a given status code..
func (vc *ViewCache) ViewStatus(statusCode int, viewName string, viewModel interface{}) Result {
	t, err := vc.Lookup(viewName)
	if err != nil {
		return vc.viewError(err)
	}
	if t == nil {
		return vc.InternalError(exception.New(ErrUnsetViewTemplate).WithMessagef("viewname: %s", viewName))
	}

	return &ViewResult{
		ViewName:   viewName,
		StatusCode: statusCode,
		ViewModel:  viewModel,
		Template:   t,
		Views:      vc,
	}
}

// ----------------------------------------------------------------------
// properties
// ----------------------------------------------------------------------

// AddPaths adds paths to the view collection.
func (vc *ViewCache) AddPaths(paths ...string) {
	vc.viewPaths = append(vc.viewPaths, paths...)
}

// AddLiterals adds view literal strings to the view collection.
func (vc *ViewCache) AddLiterals(views ...string) {
	vc.viewLiterals = append(vc.viewLiterals, views...)
}

// SetPaths sets the view paths outright.
func (vc *ViewCache) SetPaths(paths ...string) {
	vc.viewPaths = paths
}

// SetLiterals sets the raw views outright.
func (vc *ViewCache) SetLiterals(viewLiterals ...string) {
	vc.viewLiterals = viewLiterals
}

// Literals returns the view literals.
func (vc *ViewCache) Literals() []string {
	return vc.viewLiterals
}

// Paths returns the view paths.
func (vc *ViewCache) Paths() []string {
	return vc.viewPaths
}

// FuncMap returns the global view func map.
func (vc *ViewCache) FuncMap() template.FuncMap {
	return vc.viewFuncMap
}

// Templates gets the view cache for the app.
func (vc *ViewCache) Templates() (*template.Template, error) {
	if vc.cached {
		return vc.viewCache, nil
	}
	return vc.Parse()
}

// WithBadRequestTemplateName sets the bad request template.
func (vc *ViewCache) WithBadRequestTemplateName(templateName string) *ViewCache {
	vc.badRequestTemplateName = templateName
	return vc
}

// BadRequestTemplateName returns the bad request template.
func (vc *ViewCache) BadRequestTemplateName() string {
	return vc.badRequestTemplateName
}

// WithInternalErrorTemplateName sets the bad request template.
func (vc *ViewCache) WithInternalErrorTemplateName(templateName string) *ViewCache {
	vc.internalErrorTemplateName = templateName
	return vc
}

// InternalErrorTemplateName returns the bad request template.
func (vc *ViewCache) InternalErrorTemplateName() string {
	return vc.internalErrorTemplateName
}

// WithNotFoundTemplateName sets the not found request template name.
func (vc *ViewCache) WithNotFoundTemplateName(templateName string) *ViewCache {
	vc.notFoundTemplateName = templateName
	return vc
}

// NotFoundTemplateName returns the not found template name.
func (vc *ViewCache) NotFoundTemplateName() string {
	return vc.notFoundTemplateName
}

// WithNotAuthorizedTemplateName sets the not authorized template name.
func (vc *ViewCache) WithNotAuthorizedTemplateName(templateName string) *ViewCache {
	vc.notAuthorizedTemplateName = templateName
	return vc
}

// NotAuthorizedTemplateName returns the not authorized template name.
func (vc *ViewCache) NotAuthorizedTemplateName() string {
	return vc.notAuthorizedTemplateName
}

// WithStatusTemplateName sets the status templatename .
func (vc *ViewCache) WithStatusTemplateName(templateName string) *ViewCache {
	vc.statusTemplateName = templateName
	return vc
}

// StatusTemplateName returns the status template name.
func (vc *ViewCache) StatusTemplateName() string {
	return vc.statusTemplateName
}

// Initialized returns if the viewcache is initialized.
func (vc *ViewCache) Initialized() bool {
	return vc.initialized
}

// WithCached sets if we should cache views once they're compiled, or always read them from disk.
// Cached == True, use in memory storage for views
// Cached == False, read the file from disk every time we want to render the view.
func (vc *ViewCache) WithCached(cached bool) *ViewCache {
	vc.cached = cached
	return vc
}

// Cached indicates if the cache is enabled, or if we skip parsing views each load.
// Cached == True, use in memory storage for views
// Cached == False, read the file from disk every time we want to render the view.
func (vc *ViewCache) Cached() bool {
	return vc.cached
}

// ----------------------------------------------------------------------
// helpers
// ----------------------------------------------------------------------

func (vc *ViewCache) viewError(err error) Result {
	t, _ := template.New("").Parse(DefaultTemplateInternalError)
	return &ViewResult{
		ViewName:   DefaultTemplateNameInternalError,
		StatusCode: http.StatusInternalServerError,
		ViewModel:  err,
		Template:   t,
		Views:      vc,
	}
}

func (vc *ViewCache) initialize() error {
	if len(vc.viewPaths) == 0 && len(vc.viewLiterals) == 0 {
		return nil
	}

	views, err := vc.Parse()
	if err != nil {
		return err
	}
	vc.viewCache = views
	return nil
}
