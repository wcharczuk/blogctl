package r2

import (
	"io"

	"github.com/blend/go-sdk/webutil"
)

// OptBody sets the post body on the request.
func OptBody(contents io.ReadCloser) Option {
	return RequestOption(webutil.OptBody(contents))
}
