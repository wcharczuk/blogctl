package web

import "net/http"

// Fileserver is a type that implements the basics of a fileserver.
type Fileserver interface {
	AddHeader(key, value string)
	AddRewriteRule(match string, rewriteAction RewriteAction) error
	SetMiddleware(middleware ...Middleware)
	Headers() http.Header
	RewriteRules() []RewriteRule
	Action(*Ctx) Result
}
