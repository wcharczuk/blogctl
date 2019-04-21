package web

// Result is the result of a controller.
type Result interface {
	Render(ctx *Ctx) error
}

// ResultPreRender is a result that has a PreRender step.
type ResultPreRender interface {
	PreRender(ctx *Ctx) error
}

// ResultPostRender is a result that has a PostRender step.
type ResultPostRender interface {
	PostRender(ctx *Ctx) error
}

// ResultWithLoggedError logs an error before it renders the result.
func ResultWithLoggedError(result Result, err error) *LoggedErrorResult {
	return &LoggedErrorResult{
		Error:  err,
		Result: result,
	}
}

// LoggedErrorResult is a result that returns an error during the prerender phase.
type LoggedErrorResult struct {
	Result Result
	Error  error
}

// PreRender returns the underlying error.
func (ler LoggedErrorResult) PreRender(ctx *Ctx) error {
	return ler.Error
}

// Render renders the result.
func (ler LoggedErrorResult) Render(ctx *Ctx) error {
	return ler.Result.Render(ctx)
}
