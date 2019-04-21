r2
==

This is meant to be an experiment for an options based api for sending request. It is not stable and should only be used on an experimental basis.

## Philosophy

Departing from "Fluent APIs", `go-sdk/r2` investigates what an "Options" based api for making http requests would look like.

Funadmentally, it means taking code that looks like:

```golang
res, err := request.New().
	WithVerb("POST").
	MustWithURL("https://www.google.com/robots.txt").
	WithHeaderValue("X-Authorization", "none").
	WithHeaderValue(request.HeaderContentType, request.ContentTypeApplicationJSON).
	WithBody([]byte(`{"status":"maybe?"}`)).
	Execute()
```

And refactors that to:

```golang
res, err := r2.New("https://www.google.com/robots.txt",
	r2.OptPost(),
	r2.OptHeaderValue("X-Authorization", "none").
	r2.OptHeaderValue(request.HeaderContentType, request.ContentTypeApplicationJSON),
	r2.OptBody([]byte(`{"status":"maybe?"}`)).Do()
```

The key difference here is making use of a variadic list of "Options" which are really just functions that satisfy the signature `func(*r2.Request) error`. This lets developers _extend_ the possible options that can be specified, vs. having a strictly hard coded list hung off the `request.Request` object, which require a PR to make changes to.

## Usage

R2 uses a different paradigm from `go-sdk/request`; instead of chaining calls with a "fluent" api, options can be provided in a variadic list. This lets users extend the possible options as necessary.

## Example

```golang
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/blend/go-sdk/r2"
)

func CustomOption() r2.Option {
	return func(r *r2.Request) {
		r.Client.Timeout = 10 * time.Millisecond
	}
}

func main() {
	_, err := r2.New("https://google.com",
		r2.Get(),
		r2.Timeout(500*time.Millisecond),
		r2.Header("X-Sent-By", "go-sdk/request2"),
		r2.CookieValue("ssid", "baileydog01"),
		CustomOption(),
	).CopyTo(os.Stdout)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
```