package web

import "net/http"

// RawResult is for when you just want to dump bytes.
type RawResult struct {
	StatusCode  int
	ContentType string
	Response    []byte
}

// Render renders the result.
func (rr *RawResult) Render(ctx *Ctx) error {
	if len(rr.ContentType) != 0 {
		ctx.Response().Header().Set("Content-Type", rr.ContentType)
	}
	if rr.StatusCode == 0 {
		ctx.Response().WriteHeader(http.StatusOK)
	} else {
		ctx.Response().WriteHeader(rr.StatusCode)
	}
	_, err := ctx.Response().Write(rr.Response)
	return err
}
