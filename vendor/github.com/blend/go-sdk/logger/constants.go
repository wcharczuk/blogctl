package logger

import "time"

const (
	// Gigabyte is an SI unit.
	Gigabyte int = 1 << 30
	// Megabyte is an SI unit.
	Megabyte int = 1 << 20
	// Kilobyte is an SI unit.
	Kilobyte int = 1 << 10
)

const (
	// LoggerStarted is a logger state.
	LoggerStarted int32 = 0
	// LoggerStopping is a logger state.
	LoggerStopping int32 = 1
	// LoggerStopped is a logger state.
	LoggerStopped int32 = 2
)

const (
	// DefaultBufferPoolSize is the default buffer pool size.
	DefaultBufferPoolSize = 1 << 8 // 256

	// DefaultTextTimeFormat is the default time format.
	DefaultTextTimeFormat = time.RFC3339Nano

	// DefaultTextWriterUseColor is a default setting for writers.
	DefaultTextWriterUseColor = true
	// DefaultTextWriterShowHeadings is a default setting for writers.
	DefaultTextWriterShowHeadings = true
	// DefaultTextWriterShowTimestamp is a default setting for writers.
	DefaultTextWriterShowTimestamp = true
)

var (
	// DefaultFlags are the default flags.
	DefaultFlags = []Flag{Fatal, Error, Warning, Info, HTTPResponse}
	// DefaultFlagSet is the default verbosity for a diagnostics agent inited from the environment.
	DefaultFlagSet = NewFlagSet(DefaultFlags...)

	// DefaultHiddenFlags are the default hidden flags.
	DefaultHiddenFlags []Flag
)

const (
	// DefaultWriteQueueDepth  is the default depth per listener to queue work.
	// It's currently set to 256k entries.
	DefaultWriteQueueDepth = 1 << 18

	// DefaultListenerQueueDepth is the default depth per listener to queue work.
	// It's currently set to 256k entries.
	DefaultListenerQueueDepth = 1 << 10
)
