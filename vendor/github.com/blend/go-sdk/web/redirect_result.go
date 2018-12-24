package web

import "net/http"

// RedirectResult is a result that should cause the browser to redirect.
type RedirectResult struct {
	Method      string `json:"redirect_method"`
	RedirectURI string `json:"redirect_uri"`
}

// Render writes the result to the response.
func (rr *RedirectResult) Render(ctx *Ctx) error {
	if len(rr.Method) > 0 {
		ctx.Request().Method = rr.Method
		http.Redirect(ctx.Response(), ctx.Request(), rr.RedirectURI, http.StatusFound)
	} else {
		http.Redirect(ctx.Response(), ctx.Request(), rr.RedirectURI, http.StatusTemporaryRedirect)
	}

	return nil
}
