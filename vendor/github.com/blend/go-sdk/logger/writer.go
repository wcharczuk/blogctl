package logger

import (
	"strings"

	"github.com/blend/go-sdk/env"
)

// OutputFormat is a writer output format.
type OutputFormat string

const (
	// OutputFormatJSON is an output format.
	OutputFormatJSON OutputFormat = "json"
	// OutputFormatText is an output format.
	OutputFormatText OutputFormat = "text"
	// Sometime in the future ...
	// OutputFormatProtobuf = "protobuf"
)

// NewWriter creates a new writer based on a given format.
// It reads the writer settings from the environment.
func NewWriter(format OutputFormat) Writer {
	switch OutputFormat(strings.ToLower(string(format))) {
	case OutputFormatJSON:
		return NewJSONWriterFromEnv()
	case OutputFormatText:
		return NewTextWriterFromEnv()
	}

	panic("invalid writer output format")
}

// NewWriterFromEnv returns a new writer based on the environment variable `LOG_FORMAT`.
func NewWriterFromEnv() Writer {
	if format := env.Env().String(EnvVarFormat); len(format) > 0 {
		return NewWriter(OutputFormat(format))
	}
	return NewTextWriterFromEnv()
}
