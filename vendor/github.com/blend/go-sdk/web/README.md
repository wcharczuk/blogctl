web
======

`web` is a lightweight framework for building web applications in go. It rolls together very tightly scoped middleware with API endpoint and view endpoint patterns. 

## Requirements

* go 1.8+

## Example

Let's say we have a controller we need to implement:

```go
type FooController struct {}

func (fc FooContoller) Register(app *web.App) {
	app.GET("/bar", fc.bar)
}

func (fc FooController) barHandler(ctx *web.Ctx) web.Result {
	return ctx.Text().Result("bar!")
}
```

Then we would have the following in our `main.go`:

```go
func main() {
	app := web.New()
	app.Register(new(FooController))
	app.Start() // note, this call blocks until the server exists / crashes.
}
```

And that's it! There are options to configure things like the port and tls certificates, but the core use case is to bind
on 8080 or whatever is specified in the `PORT` environment variable. 

## Middleware

If you want to run some steps before controller actions fire (such as for auth etc.) you can add those steps as "middleware". 

```go
	app.GET("/admin/dashboard", c.dashboardAction, middle2, middle1)
```

This will then run `middle1` and then run `middle2`, finally the controller action.
An important detail is that the "cascading" nature of the calls depends on how you structure your middleware functions. If any of the middleware functions
return without calling the `action` parameter, execution stops there and subsequent middleware steps do not get called (ditto the controller action).
This lets us have authentication steps happen in common middlewares before our controller action gets run. It also lets us specify different middlewares per route.

What do `middle1` and `middle2` look like? They are `Middleware`; functions that take an `Action` and return an `Action`.

```go
func middle1(action web.Action) web.Action {
	return func(ctx *web.Ctx) web.Result {
		// check if our "foo" param is correct.
		if ctx.Param("foo") != "bar" { // this is only for illustration, you'd want to do something more secure in practice
			return ctx.DefaultResultProvider().NotAuthorized() //.DefaultResultProvider() can be set by helper middlewares, otherwise defaults to `Text`
		}
		return action(ctx) //we call the input action here
	}
}
```

## Authentication

`go-web` comes built in with some basic handling of authentication and a concept of session. With very basic configuration, middlewares can be added that either require a valid session, or simply read the session and provide it to the downstream controller action.

```go
func main() {
	app := web.New()
	
	// Provide a session validation handler.
	// It is called after the session has been read off the request.
	app.Auth().SetValidateHandler(func(session *web.Session, state web.State) error {
		// here we might want to reach into our permanent session store and make sure the session is still valid
		return nil
	})

	// Provide a redirect handler if a session is required for a route.
	// This is simply a function that takes an incoming request url, and returns what it should be redirected to.
	// You can add as much logic here as you'd like, but below is a simple redirect that points people to a login page
	// with a param denoting where they were trying to go (useful for post-login).
	app.Auth().SetLoginRedirectHandler(func(u *url.URL) *url.URL {
		u.RawQuery = fmt.Sprintf("redirect=%s", url.QueryEscape(u.Path))
		u.Path = fmt.Sprintf("/login")
		return u
	})

	app.POST("/login", func(ctx *web.Ctx) web.Result {
		// if we've already got a session, exit early.
		if ctx.Session() != nil {
			return ctx.RedirectWithMethodf("GET", "/dashboard")
		}
		// audits, other events you might want to trigger.
		// Login(...) will issue cookies, and store the session in the local session cache
		// so we can validate it on subsequent requests.
		ctx.Auth().Login("my user id", ctx)
		return ctx.RedirectWithMethodf("GET", "/dashboard")
	}, web.SessionAware)

	app.GET("/logout", func(ctx *web.Ctx) web.Result {
		// if we don't already have a session, exit early.
		if ctx.Session() == nil {
			return ctx.RedirectWithMethodf("GET", "/login")
		}

		// audits, other events you might want to trigger.
		ctx.Auth().Logout(ctx)
		return ctx.RedirectWithMethodf("GET", "/login")
	}, web.SessionAware)
}
```

## Serving Static Files

You can set a path root to serve static files.

```go
func main() {
	app := web.New()
	app.ServeStatic("/static", "_client/dist")
	app.Start()
}
```

Here we tell the app that we should serve the `_client/dist` directory as a route prefixed by "/static". If we have a file `foo.css` in `_client/dist`, it would
be accessible at `/static/foo.css` in our app. 

You can also have a controller action return a static file:

```go
	app.GET("/thing", func(r *web.Ctx) web.ControllerResult { return r.Static("path/to/my/file") })
```

You can optionally set a static re-write rule (such as if you are cache-breaking assets with timestamps in the filename):

```go
func main() {
	app := web.New()
	app.ServeStatic("/static", "_client/dist")
	app.WithStaticRewriteRule("/static", `^(.*)\.([0-9]+)\.(css|js|html|htm)$`, func(path string, parts ...string) string {
		return fmt.Sprintf("%s.%s", parts[1], parts[3])
	})
	app.Start()
}
```

Here we feed the `WithStaticRewriteRule` function the path (the same path as our static file server, this is important), a regular expression to match, and a special handler function that returns an updated path. 

Note: `parts ...string` is the regular expression sub matches from the expression, with `parts[0]` equal to the full input string. `parts[1]` and `parts[3]` in this case are the nominal root stem, and the extension respecitvely.

You can also set custom headers for static files:

```go
func main() {
	app := web.New()
	app.ServeStatic("/static", "_client/dist")
	app.WithStaticRewriteRule("/static", `^(.*)\.([0-9]+)\.(css|js)$`, func(path string, parts ...string) string {
		return fmt.Sprintf("%s.%s", parts[1], parts[3])
	})
	app.WithStaticHeader("/static", "cache-control", "public,max-age=99999999")	
}
```

This will then set the specified cache headers on response for the static files, this is useful for things like `cache-control`.

You can finally set a static route to run middleware for that route.

```go
func main() {
	app := web.New()
	app.ServeStatic("/static", "_client/dist")
	app.WithStaticRewriteRule("/static", `^(.*)\.([0-9]+)\.(css|js)$`, func(path string, parts ...string) string {
		return fmt.Sprintf("%s.%s", parts[1], parts[3])
	})
	app.WithStaticHeader("/static", "cache-control", "public,max-age=99999999")	
	app.WithStaticMiddleware("/static", web.SessionRequired)
}
```

You would now need to have a valid session to access any of the files under `/static`.

## Benchmarks

Benchmarks are key, obviously, because the ~200us you save choosing a framework won't be wiped out by the 50ms ping time to your servers. 

For a relatively clean implementation (found in `benchmark/main.go`) that uses `go-web`:
```
Running 10s test @ http://localhost:9090/json
  2 threads and 64 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     0.92ms  223.27us  11.82ms   86.00%
    Req/Sec    34.19k     2.18k   40.64k    65.35%
  687017 requests in 10.10s, 203.11MB read
Requests/sec:  68011.73
Transfer/sec:     20.11MB
```

On the same machine, with a very, very bare bones implementation using only built-in stuff in `net/http`:

```
Running 10s test @ http://localhost:9090/json
  2 threads and 64 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     0.93ms  216.64us  10.25ms   89.10%
    Req/Sec    34.22k     2.73k   40.63k    59.90%
  687769 requests in 10.10s, 109.54MB read
Requests/sec:  68091.37
Transfer/sec:     10.84MB
```

The key here is to make sure not to enable logging, because if logging is enabled that throughput gets cut in half. 
