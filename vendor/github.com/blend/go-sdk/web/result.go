package web

// Result is the result of a controller.
type Result interface {
	Render(ctx *Ctx) error
}

// ResultPreRender is a result that has a PreRender step.
type ResultPreRender interface {
	PreRender(ctx *Ctx) error
}

// ResultRenderComplete is a result that has a RenderComplete step.
type ResultRenderComplete interface {
	RenderComplete(ctx *Ctx) error
}

// resultWithLoggedError logs an error before it renders the result.
func resultWithLoggedError(result Result, err error) *loggedErrorResult {
	return &loggedErrorResult{
		Error:  err,
		Result: result,
	}
}

type loggedErrorResult struct {
	Result Result
	Error  error
}

func (ler loggedErrorResult) PreRender(ctx *Ctx) error {
	return ler.Error
}

func (ler loggedErrorResult) Render(ctx *Ctx) error {
	return ler.Result.Render(ctx)
}
