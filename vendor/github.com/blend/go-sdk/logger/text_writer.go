package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// Asserts text writer is a writer.
var (
	_ Writer = &TextWriter{}
)

// TextWritable is a type with a custom formater for text writing.
type TextWritable interface {
	WriteText(formatter TextFormatter, buf *bytes.Buffer)
}

// FlagTextColorProvider is a function types can implement to provide a color.
type FlagTextColorProvider interface {
	FlagTextColor() AnsiColor
}

// TextFormatter formats text.
type TextFormatter interface {
	Colorize(value string, color AnsiColor) string
	ColorizeStatusCode(code int) string
	ColorizeByStatusCode(code int, value string) string
}

// NewTextWriter returns a new text writer for a given output.
func NewTextWriter(output io.Writer) *TextWriter {
	return &TextWriter{
		output:        NewInterlockedWriter(output),
		bufferPool:    NewBufferPool(DefaultBufferPoolSize),
		showHeadings:  DefaultTextWriterShowHeadings,
		showTimestamp: DefaultTextWriterShowTimestamp,
		useColor:      DefaultTextWriterUseColor,
		timeFormat:    DefaultTextTimeFormat,
	}
}

// NewTextWriterStdout returns a new text writer to stdout/stderr.
func NewTextWriterStdout() *TextWriter {
	return NewTextWriter(os.Stdout).WithErrorOutput(os.Stderr)
}

// NewTextWriterFromEnv returns a new text writer from the environment.
func NewTextWriterFromEnv() *TextWriter {
	return NewTextWriterFromConfig(NewTextWriterConfigFromEnv())
}

// NewTextWriterFromConfig creates a new text writer from a given config.
func NewTextWriterFromConfig(cfg *TextWriterConfig) *TextWriter {
	return &TextWriter{
		output:        NewInterlockedWriter(os.Stdout),
		errorOutput:   NewInterlockedWriter(os.Stderr),
		bufferPool:    NewBufferPool(DefaultBufferPoolSize),
		showTimestamp: cfg.GetShowTimestamp(),
		showHeadings:  cfg.GetShowHeadings(),
		useColor:      cfg.GetUseColor(),
		timeFormat:    cfg.GetTimeFormat(),
	}
}

// TextWriter handles outputting logging events to given writer streams as textual columns.
type TextWriter struct {
	output      io.Writer
	errorOutput io.Writer

	showTimestamp bool
	showHeadings  bool
	useColor      bool

	timeFormat string

	bufferPool *BufferPool
}

// OutputFormat returns the output format.
func (wr *TextWriter) OutputFormat() OutputFormat {
	return OutputFormatText
}

// WithUseColor sets a formatting option.
func (wr *TextWriter) WithUseColor(useColor bool) *TextWriter {
	wr.useColor = useColor
	return wr
}

// UseColor is a formatting option.
func (wr *TextWriter) UseColor() bool {
	return wr.useColor
}

// WithShowTimestamp sets a formatting option.
func (wr *TextWriter) WithShowTimestamp(showTime bool) *TextWriter {
	wr.showTimestamp = showTime
	return wr
}

// ShowTimestamp is a formatting option.
func (wr *TextWriter) ShowTimestamp() bool {
	return wr.showTimestamp
}

// WithShowHeadings sets a formatting option.
func (wr *TextWriter) WithShowHeadings(showHeadings bool) *TextWriter {
	wr.showHeadings = showHeadings
	return wr
}

// ShowHeadings is a formatting option.
func (wr *TextWriter) ShowHeadings() bool {
	return wr.showHeadings
}

// WithTimeFormat sets a formatting option.
func (wr *TextWriter) WithTimeFormat(timeFormat string) *TextWriter {
	wr.timeFormat = timeFormat
	return wr
}

// TimeFormat is a formatting option.
func (wr *TextWriter) TimeFormat() string {
	return wr.timeFormat
}

// Output returns the output.
func (wr *TextWriter) Output() io.Writer {
	return wr.output
}

// WithOutput sets the primary output.
func (wr *TextWriter) WithOutput(output io.Writer) *TextWriter {
	wr.output = NewInterlockedWriter(output)
	return wr
}

// ErrorOutput returns an io.Writer for the error stream.
func (wr *TextWriter) ErrorOutput() io.Writer {
	if wr.errorOutput != nil {
		return wr.errorOutput
	}
	return wr.output
}

// WithErrorOutput sets the error output.
func (wr *TextWriter) WithErrorOutput(errorOutput io.Writer) *TextWriter {
	wr.errorOutput = NewInterlockedWriter(errorOutput)
	return wr
}

// Colorize (optionally) applies a color to a string.
func (wr *TextWriter) Colorize(value string, color AnsiColor) string {
	if wr.useColor {
		return color.Apply(value)
	}
	return value
}

// ColorizeStatusCode adds color to a status code.
func (wr *TextWriter) ColorizeStatusCode(statusCode int) string {
	if wr.useColor {
		return ColorizeStatusCode(statusCode)
	}
	return strconv.Itoa(statusCode)
}

// ColorizeByStatusCode colorizes a string by a status code (green, yellow, red).
func (wr *TextWriter) ColorizeByStatusCode(statusCode int, value string) string {
	if wr.useColor {
		return ColorizeByStatusCode(statusCode, value)
	}
	return value
}

// FormatFlag formats the flag portion of the message.
func (wr *TextWriter) FormatFlag(flag Flag, color AnsiColor) string {
	return fmt.Sprintf("[%s]", wr.Colorize(string(flag), color))
}

// FormatEntity formats the flag portion of the message.
func (wr *TextWriter) FormatEntity(entity string, color AnsiColor) string {
	return fmt.Sprintf("[%s]", wr.Colorize(entity, color))
}

// FormatHeadings returns the headings section of the message.
func (wr *TextWriter) FormatHeadings(headings ...string) string {
	if len(headings) == 0 {
		return ""
	}
	if len(headings) == 1 {
		return fmt.Sprintf("[%s]", wr.Colorize(headings[0], ColorBlue))
	}
	if wr.useColor {
		for index := 0; index < len(headings); index++ {
			headings[index] = wr.Colorize(headings[index], ColorBlue)
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(headings, " > "))
}

// FormatTimestamp returns a new timestamp string.
func (wr *TextWriter) FormatTimestamp(optionalTime ...time.Time) string {
	timeFormat := DefaultTextTimeFormat
	if len(wr.timeFormat) > 0 {
		timeFormat = wr.timeFormat
	}
	var value string
	if len(optionalTime) > 0 {
		value = optionalTime[0].Format(timeFormat)
	} else {
		value = time.Now().UTC().Format(timeFormat)
	}
	return wr.Colorize(fmt.Sprintf("%-27s", value), ColorGray)
}

// GetBuffer returns a leased buffer from the buffer pool.
func (wr *TextWriter) GetBuffer() *bytes.Buffer {
	return wr.bufferPool.Get()
}

// PutBuffer adds the leased buffer back to the pool.
// It Should be called in conjunction with `GetBuffer`.
func (wr *TextWriter) PutBuffer(buffer *bytes.Buffer) {
	wr.bufferPool.Put(buffer)
}

// Write writes to stdout.
func (wr *TextWriter) Write(e Event) error {
	return wr.write(wr.Output(), e)
}

// WriteError writes to stderr (or stdout if .errorOutput is unset).
func (wr *TextWriter) WriteError(e Event) error {
	return wr.write(wr.ErrorOutput(), e)
}

func (wr *TextWriter) write(output io.Writer, e Event) error {
	buf := wr.bufferPool.Get()
	defer wr.bufferPool.Put(buf)

	if wr.showTimestamp {
		buf.WriteString(wr.FormatTimestamp(e.Timestamp()))
		buf.WriteRune(RuneSpace)
	}

	if typed, isTyped := e.(EventHeadings); isTyped {
		headings := typed.Headings()
		if len(headings) > 0 {
			buf.WriteString(wr.FormatHeadings(headings...))
			buf.WriteRune(RuneSpace)
		}
	}

	if typed, isTyped := e.(FlagTextColorProvider); isTyped {
		if flagColor := typed.FlagTextColor(); len(flagColor) > 0 {
			buf.WriteString(wr.FormatFlag(e.Flag(), flagColor))
		} else {
			buf.WriteString(wr.FormatFlag(e.Flag(), GetFlagTextColor(e.Flag())))
		}
	} else {
		buf.WriteString(wr.FormatFlag(e.Flag(), GetFlagTextColor(e.Flag())))
	}
	buf.WriteRune(RuneSpace)

	if typed, isTyped := e.(EventEntity); isTyped {
		if len(typed.Entity()) > 0 {
			buf.WriteString(wr.FormatEntity(typed.Entity(), ColorBlue))
			buf.WriteRune(RuneSpace)
		}
	}

	if typed, isTyped := e.(TextWritable); isTyped {
		typed.WriteText(wr, buf)
	} else if typed, isTyped := e.(fmt.Stringer); isTyped {
		buf.WriteString(typed.String())
	}

	buf.WriteRune(RuneNewline)
	_, err := buf.WriteTo(output)
	return err
}
