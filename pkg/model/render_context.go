package model

// RenderContext is the full context for a particular render.
type RenderContext struct {
	Data     *Data    `json:"data"`
	Partials []string `json:"partials"`
	Stats    Stats    `json:"stats"`
}
