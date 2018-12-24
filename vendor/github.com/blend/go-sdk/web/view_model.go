package web

import (
	"github.com/blend/go-sdk/env"
)

// ViewModel is a wrapping viewmodel.
type ViewModel struct {
	Env       env.Vars
	Ctx       *Ctx
	ViewModel interface{}
}

// StatusViewModel returns the status view model.
type StatusViewModel struct {
	StatusCode int
	Response   interface{}
}
