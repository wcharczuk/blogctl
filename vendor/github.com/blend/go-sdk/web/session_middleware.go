package web

// SessionAware is an action that injects the session into the context, it acquires a read lock on session.
func SessionAware(action Action) Action {
	return func(ctx *Ctx) Result {
		session, err := ctx.Auth().VerifySession(ctx)
		if err != nil && !IsErrSessionInvalid(err) {
			return ctx.DefaultResultProvider().InternalError(err)
		}
		ctx.WithSession(session)
		return action(ctx)
	}
}

// SessionRequired is an action that requires a session to be present
// or identified in some form on the request, and acquires a read lock on session.
func SessionRequired(action Action) Action {
	return func(ctx *Ctx) Result {
		session, err := ctx.Auth().VerifySession(ctx)
		if err != nil && !IsErrSessionInvalid(err) {
			return ctx.DefaultResultProvider().InternalError(err)
		}
		if session == nil {
			return ctx.Auth().LoginRedirect(ctx)
		}
		ctx.WithSession(session)
		return action(ctx)
	}
}

// SessionMiddleware implements a custom notAuthorized action.
func SessionMiddleware(notAuthorized Action) Middleware {
	return func(action Action) Action {
		return func(ctx *Ctx) Result {
			session, err := ctx.Auth().VerifySession(ctx)
			if err != nil && !IsErrSessionInvalid(err) {
				return ctx.DefaultResultProvider().InternalError(err)
			}

			if session == nil {
				if notAuthorized != nil {
					return notAuthorized(ctx)
				}
				return ctx.Auth().LoginRedirect(ctx)
			}
			ctx.WithSession(session)
			return action(ctx)
		}
	}
}
