package web

// Middleware are steps that run in order before a given action.
type Middleware func(Action) Action

// Action is the function signature for controller actions.
type Action func(*Ctx) Result

// PanicAction is a receiver for app.PanicHandler.
type PanicAction func(*Ctx, interface{}) Result
