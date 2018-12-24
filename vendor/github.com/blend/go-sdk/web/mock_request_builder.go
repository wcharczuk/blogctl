package web

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/webutil"
)

// NewMockRequestBuilder returns a new mock request builder for a given app.
func NewMockRequestBuilder(app *App) *MockRequestBuilder {
	return &MockRequestBuilder{
		app:         app,
		verb:        "GET",
		queryString: url.Values{},
		formValues:  url.Values{},
		headers:     http.Header{},
		state:       &SyncState{},
	}
}

// MockRequestBuilder facilitates creating mock requests.
type MockRequestBuilder struct {
	app *App

	verb        string
	path        string
	queryString url.Values
	formValues  url.Values
	headers     http.Header
	cookies     []*http.Cookie
	postBody    []byte

	basicAuthUsername string
	basicAuthPassword string

	postedFiles map[string]PostedFile

	err error

	state State
}

// WithErr sets the error if it is unset.
func (mrb *MockRequestBuilder) WithErr(err error) *MockRequestBuilder {
	if err != nil {
		mrb.err = err
	}
	return mrb
}

// Err rerturns an underlying error
func (mrb *MockRequestBuilder) Err() error {
	return mrb.err
}

// WithBasicAuth sets basic auth credentials for the request.
func (mrb *MockRequestBuilder) WithBasicAuth(username, password string) *MockRequestBuilder {
	mrb.basicAuthUsername, mrb.basicAuthPassword = username, password
	return mrb
}

// Get is a shortcut for WithVerb("GET") WithPathf(pathFormat, args...)
func (mrb *MockRequestBuilder) Get(pathFormat string, args ...interface{}) *MockRequestBuilder {
	return mrb.WithVerb("GET").WithPathf(pathFormat, args...)
}

// Post is a shortcut for WithVerb("POST") WithPathf(pathFormat, args...)
func (mrb *MockRequestBuilder) Post(pathFormat string, args ...interface{}) *MockRequestBuilder {
	return mrb.WithVerb("POST").WithPathf(pathFormat, args...)
}

// Put is a shortcut for WithVerb("PUT") WithPathf(pathFormat, args...)
func (mrb *MockRequestBuilder) Put(pathFormat string, args ...interface{}) *MockRequestBuilder {
	return mrb.WithVerb("PUT").WithPathf(pathFormat, args...)
}

// Patch is a shortcut for WithVerb("PATCH") WithPathf(pathFormat, args...)
func (mrb *MockRequestBuilder) Patch(pathFormat string, args ...interface{}) *MockRequestBuilder {
	return mrb.WithVerb("PATCH").WithPathf(pathFormat, args...)
}

// Delete is a shortcut for WithVerb("DELETE") WithPathf(pathFormat, args...)
func (mrb *MockRequestBuilder) Delete(pathFormat string, args ...interface{}) *MockRequestBuilder {
	return mrb.WithVerb("DELETE").WithPathf(pathFormat, args...)
}

// WithVerb sets the verb for the request.
func (mrb *MockRequestBuilder) WithVerb(verb string) *MockRequestBuilder {
	mrb.verb = strings.ToUpper(verb)
	return mrb
}

// WithPathf sets the path for the request.
func (mrb *MockRequestBuilder) WithPathf(pathFormat string, args ...interface{}) *MockRequestBuilder {
	mrb.path = fmt.Sprintf(pathFormat, args...)

	// url.Parse always includes the '/' path prefix.
	if !strings.HasPrefix(mrb.path, "/") {
		mrb.path = fmt.Sprintf("/%s", mrb.path)
	}

	return mrb
}

// WithQueryString adds a querystring param for the request.
func (mrb *MockRequestBuilder) WithQueryString(key, value string) *MockRequestBuilder {
	mrb.queryString.Add(key, value)
	return mrb
}

// WithFormValue adds a form value for the request.
func (mrb *MockRequestBuilder) WithFormValue(key, value string) *MockRequestBuilder {
	mrb.formValues.Add(key, value)
	return mrb
}

// WithHeader adds a header for the request.
func (mrb *MockRequestBuilder) WithHeader(key, value string) *MockRequestBuilder {
	mrb.headers.Add(key, value)
	return mrb
}

// WithCookie adds a cookie for the request.
func (mrb *MockRequestBuilder) WithCookie(cookie *http.Cookie) *MockRequestBuilder {
	mrb.cookies = append(mrb.cookies, cookie)
	return mrb
}

// WithCookieValue adds a basic name+value cookie for the request.
func (mrb *MockRequestBuilder) WithCookieValue(name, value string) *MockRequestBuilder {
	mrb.cookies = append(mrb.cookies, NewCookie(name, value))
	return mrb
}

// WithPostBody sets the post body for the request.
func (mrb *MockRequestBuilder) WithPostBody(postBody []byte) *MockRequestBuilder {
	mrb.postBody = postBody
	return mrb
}

// WithPostBodyAsJSON sets the post body for the request by serializing an object to JSON.
func (mrb *MockRequestBuilder) WithPostBodyAsJSON(object interface{}) *MockRequestBuilder {
	bytes, _ := json.Marshal(object)
	mrb.postBody = bytes
	return mrb
}

// WithPostedFile includes a file as a post parameter.
func (mrb *MockRequestBuilder) WithPostedFile(postedFile PostedFile) *MockRequestBuilder {
	if mrb.postedFiles == nil {
		mrb.postedFiles = map[string]PostedFile{}
	}
	mrb.postedFiles[postedFile.Key] = postedFile
	return mrb
}

// State returns the underlying state.
func (mrb *MockRequestBuilder) State() State {
	return mrb.state
}

// WithStateValue sets the state for a key to an object.
func (mrb *MockRequestBuilder) WithStateValue(key string, value interface{}) *MockRequestBuilder {
	mrb.state.Set(key, value)
	return mrb
}

// GetStateValue returns an object in the state cache.
func (mrb *MockRequestBuilder) GetStateValue(key string) interface{} {
	if mrb.state == nil {
		return nil
	}
	return mrb.state.Get(key)
}

// Request returns the mock request builder settings as an http.Request.
func (mrb *MockRequestBuilder) Request() (*http.Request, error) {
	req := &http.Request{}

	reqURL, err := url.Parse(fmt.Sprintf("http://localhost%s", mrb.path))
	if err != nil {
		return nil, err
	}

	reqURL.RawQuery = mrb.queryString.Encode()
	req.Method = mrb.verb
	req.URL = reqURL
	req.RequestURI = reqURL.String()
	req.Form = mrb.formValues
	req.Header = http.Header{}
	if len(mrb.basicAuthUsername) > 0 {
		req.SetBasicAuth(mrb.basicAuthUsername, mrb.basicAuthPassword)
	}

	for key, values := range mrb.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	for _, cookie := range mrb.cookies {
		req.AddCookie(cookie)
	}

	if len(mrb.postedFiles) > 0 {
		b := bytes.NewBuffer(nil)
		w := multipart.NewWriter(b)
		for _, file := range mrb.postedFiles {
			fw, err := w.CreateFormFile(file.Key, file.FileName)
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(fw, bytes.NewBuffer(file.Contents))
			if err != nil {
				return nil, err
			}
		}
		// Don't forget to set the content type, this will contain the boundary.
		req.Header.Set("Content-Type", w.FormDataContentType())

		err = w.Close()
		if err != nil {
			return nil, err
		}
		req.Body = ioutil.NopCloser(b)
	} else if len(mrb.postBody) > 0 {
		req.Body = ioutil.NopCloser(bytes.NewBuffer(mrb.postBody))
	}

	return req, nil
}

// LookupRoute returns the corresponding route for the mocked request.
func (mrb *MockRequestBuilder) LookupRoute() (route *Route, params RouteParameters) {
	var tsr bool
	path := mrb.path
	route, params, tsr = mrb.app.Lookup(mrb.verb, path)
	if tsr {
		path = path + "/"
		route, params, tsr = mrb.app.Lookup(mrb.verb, path)
	}
	return
}

// CreateCtx returns the mock request as a request context.
func (mrb *MockRequestBuilder) CreateCtx(p RouteParameters) (*Ctx, error) {
	r, err := mrb.Request()

	if err != nil {
		return nil, err
	}

	route, _ := mrb.LookupRoute()
	return mrb.app.createCtx(webutil.NewMockResponse(new(bytes.Buffer)), r, route, p).WithState(mrb.state), nil
}

// Response runs the mock request.
func (mrb *MockRequestBuilder) Response() (res *http.Response, err error) {
	if mrb.err != nil {
		err = mrb.err
		return
	}

	var req *http.Request
	req, err = mrb.Request()
	if err != nil {
		return
	}

	var route *Route
	var params RouteParameters
	route, params = mrb.LookupRoute()
	if route == nil && mrb.app.notFoundHandler != nil {
		res = mrb.runHandler(mrb.app.notFoundHandler, req, route, params)
		return
	} else if route == nil {
		err = exception.New("No route registered").WithMessagef("verb: %s path: %s", mrb.verb, mrb.path)
		return
	}

	res = mrb.runHandler(route.Handler, req, route, params)
	return
}

func (mrb *MockRequestBuilder) runHandler(handler Handler, req *http.Request, route *Route, params RouteParameters) (res *http.Response) {
	if mrb.app != nil && mrb.app.recoverPanics && mrb.app.panicAction != nil {
		defer func() {
			if r := recover(); r != nil {
				buffer := bytes.NewBuffer(nil)
				w := webutil.NewMockResponse(buffer)
				mrb.app.handlePanic(w, req, r)
				res = mrb.createResponse(buffer, w)
			}
		}()
	}

	buffer := bytes.NewBuffer(nil)
	w := webutil.NewMockResponse(buffer)
	handler(w, req, route, params)
	res = mrb.createResponse(buffer, w)
	return res
}

func (mrb *MockRequestBuilder) createResponse(buffer *bytes.Buffer, w ResponseWriter) *http.Response {
	res := &http.Response{
		Body:          ioutil.NopCloser(bytes.NewBuffer(buffer.Bytes())),
		ContentLength: int64(w.ContentLength()),
		Header:        http.Header{},
		StatusCode:    w.StatusCode(),
		Proto:         "http",
		ProtoMajor:    1,
		ProtoMinor:    1,
	}
	for key, values := range w.Header() {
		for _, value := range values {
			res.Header.Add(key, value)
		}
	}
	return res
}

// JSON executes the mock request and reads the response to the given object as json.
func (mrb *MockRequestBuilder) JSON(object interface{}) error {
	res, err := mrb.Response()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(contents, object)
}

// JSONWithMeta executes the mock request and reads the response to the given object as json.
func (mrb *MockRequestBuilder) JSONWithMeta(object interface{}) (*ResponseMeta, error) {
	res, err := mrb.Response()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return NewResponseMeta(res), json.Unmarshal(contents, object)
}

// XML executes the mock request and reads the response to the given object as json.
func (mrb *MockRequestBuilder) XML(object interface{}) error {
	res, err := mrb.Response()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return xml.Unmarshal(contents, object)
}

// XMLWithMeta executes the mock request and reads the response to the given object as json.
func (mrb *MockRequestBuilder) XMLWithMeta(object interface{}) (*ResponseMeta, error) {
	res, err := mrb.Response()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return NewResponseMeta(res), xml.Unmarshal(contents, object)
}

// Bytes returns the response as bytes.
func (mrb *MockRequestBuilder) Bytes() ([]byte, error) {
	res, err := mrb.Response()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// BytesWithMeta returns the response as bytes with meta information.
func (mrb *MockRequestBuilder) BytesWithMeta() ([]byte, *ResponseMeta, error) {
	res, err := mrb.Response()
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	return contents, NewResponseMeta(res), err
}

// Execute just runs the request.
// It internally calls `Bytes()` which fully consumes the response.
func (mrb *MockRequestBuilder) Execute() error {
	_, err := mrb.Bytes()
	return err
}

// ExecuteWithMeta returns basic metadata for a response.
func (mrb *MockRequestBuilder) ExecuteWithMeta() (*ResponseMeta, error) {
	res, err := mrb.Response()
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
		_, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
	}
	return NewResponseMeta(res), nil
}

// NewRequestMeta returns a new meta object for a request.
func NewRequestMeta(req *http.Request) *RequestMeta {
	return &RequestMeta{
		Verb:    req.Method,
		URL:     req.URL,
		Headers: req.Header,
	}
}

// NewRequestMetaWithBody returns a new meta object for a request and reads the body.
func NewRequestMetaWithBody(req *http.Request) (*RequestMeta, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	return &RequestMeta{
		Verb:    req.Method,
		URL:     req.URL,
		Headers: req.Header,
		Body:    body,
	}, nil
}

// RequestMeta is the metadata for a request.
type RequestMeta struct {
	StartTime time.Time
	Verb      string
	URL       *url.URL
	Headers   http.Header
	Body      []byte
}

// NewResponseMeta creates a new ResponseMeta.
func NewResponseMeta(res *http.Response) *ResponseMeta {
	return &ResponseMeta{
		StatusCode:    res.StatusCode,
		Headers:       res.Header,
		ContentLength: res.ContentLength,
	}
}

// ResponseMeta is a metadata response struct
type ResponseMeta struct {
	StatusCode    int
	ContentLength int64
	Headers       http.Header
}
