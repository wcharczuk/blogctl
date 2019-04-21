package web

// NestMiddleware reads the middleware variadic args and organizes the calls recursively in the order they appear.
func NestMiddleware(action Action, middleware ...Middleware) Action {
	if len(middleware) == 0 {
		return action
	}

	var nest = func(a, b Middleware) Middleware {
		if b == nil {
			return a
		}
		return func(inner Action) Action {
			return a(b(inner))
		}
	}

	var outer Middleware
	for _, step := range middleware {
		outer = nest(step, outer)
	}
	return outer(action)
}
