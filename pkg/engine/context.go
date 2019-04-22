package engine

import (
	"context"

	"github.com/wcharczuk/blogctl/pkg/model"
)

type renderContextKey struct{}

// WithRenderContext returns a context with a render context set.
func WithRenderContext(ctx context.Context, rc *model.RenderContext) context.Context {
	return context.WithValue(ctx, renderContextKey{}, rc)
}

// GetRenderContext returns the render context off a context.
func GetRenderContext(ctx context.Context) *model.RenderContext {
	if raw := ctx.Value(renderContextKey{}); raw != nil {
		if typed, ok := raw.(*model.RenderContext); ok {
			return typed
		}
	}
	return nil
}
