package logger

import (
	"context"
	"encoding/json"
	"io"

	"github.com/blend/go-sdk/bufferutil"
)

var (
	_ WriteFormatter = (*JSONOutputFormatter)(nil)
)

// NewJSONOutputFormatter returns a new json event formatter.
func NewJSONOutputFormatter(options ...JSONOutputFormatterOption) *JSONOutputFormatter {
	jf := &JSONOutputFormatter{
		BufferPool: bufferutil.NewPool(DefaultBufferPoolSize),
	}

	for _, option := range options {
		option(jf)
	}
	return jf
}

// JSONOutputFormatterOption is an option for json formatters.
type JSONOutputFormatterOption func(*JSONOutputFormatter)

// OptJSONConfig sets a json formatter from a config.
func OptJSONConfig(cfg *JSONConfig) JSONOutputFormatterOption {
	return func(jf *JSONOutputFormatter) {
		jf.Pretty = cfg.Pretty
		jf.PrettyIndent = cfg.PrettyIndentOrDefault()
		jf.PrettyPrefix = cfg.PrettyPrefixOrDefault()
	}
}

// JSONOutputFormatter is a json output formatter.
type JSONOutputFormatter struct {
	BufferPool   *bufferutil.Pool
	Pretty       bool
	PrettyPrefix string
	PrettyIndent string
}

// WriteFormat writes the event to the given output.
func (jw JSONOutputFormatter) WriteFormat(ctx context.Context, output io.Writer, e Event) error {
	buffer := jw.BufferPool.Get()
	defer jw.BufferPool.Put(buffer)

	encoder := json.NewEncoder(buffer)
	if jw.Pretty {
		encoder.SetIndent(jw.PrettyPrefix, jw.PrettyIndent)
	}
	if err := encoder.Encode(e); err != nil {
		return err
	}
	_, err := io.Copy(output, buffer)
	return err
}
