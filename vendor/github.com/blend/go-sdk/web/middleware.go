package web

// ViewProviderAsDefault sets the context.DefaultResultProvider() equal to context.View().
func ViewProviderAsDefault(action Action) Action {
	return func(ctx *Ctx) Result {
		return action(ctx.WithDefaultResultProvider(ctx.View()))
	}
}

// JSONProviderAsDefault sets the context.DefaultResultProvider() equal to context.JSON().
func JSONProviderAsDefault(action Action) Action {
	return func(ctx *Ctx) Result {
		return action(ctx.WithDefaultResultProvider(ctx.JSON()))
	}
}

// XMLProviderAsDefault sets the context.DefaultResultProvider() equal to context.XML().
func XMLProviderAsDefault(action Action) Action {
	return func(ctx *Ctx) Result {
		return action(ctx.WithDefaultResultProvider(ctx.XML()))
	}
}

// TextProviderAsDefault sets the context.DefaultResultProvider() equal to context.Text().
func TextProviderAsDefault(action Action) Action {
	return func(ctx *Ctx) Result {
		return action(ctx.WithDefaultResultProvider(ctx.Text()))
	}
}
