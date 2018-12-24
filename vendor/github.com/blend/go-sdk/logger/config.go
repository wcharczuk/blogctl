package logger

import (
	"strings"

	"github.com/blend/go-sdk/env"
)

// NewConfigFromEnv returns a new config from the environment.
func NewConfigFromEnv() (*Config, error) {
	var config Config
	if err := env.Env().ReadInto(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

// MustNewConfigFromEnv returns a new config from the environment,
// and panics if there is an error.
func MustNewConfigFromEnv() *Config {
	cfg, err := NewConfigFromEnv()
	if err != nil {
		panic(err)
	}
	return cfg
}

// Config is the logger config.
type Config struct {
	Heading            string   `json:"heading,omitempty" yaml:"heading,omitempty" env:"LOG_HEADING"`
	OutputFormat       string   `json:"outputFormat,omitempty" yaml:"outputFormat,omitempty" env:"LOG_FORMAT"`
	Flags              []string `json:"flags,omitempty" yaml:"flags,omitempty" env:"LOG_EVENTS,csv"`
	HiddenFlags        []string `json:"hiddenFlags,omitempty" yaml:"hiddenFlags,omitempty" env:"LOG_HIDDEN,csv"`
	RecoverPanics      *bool    `json:"recoverPanics,omitempty" yaml:"recoverPanics,omitempty" env:"LOG_RECOVER"`
	WriteQueueDepth    int      `json:"writeQueueDepth,omitempty" yaml:"writeQueueDepth,omitempty" env:"LOG_WRITE_QUEUE_DEPTH"`
	ListenerQueueDepth int      `json:"listenerQueueDepth,omitempty" yaml:"listenerQueueDepth,omitempty" env:"LOG_LISTENER_QUEUE_DEPTH"`

	Text TextWriterConfig `json:"text,omitempty" yaml:"text,omitempty"`
	JSON JSONWriterConfig `json:"json,omitempty" yaml:"json,omitempty"`
}

// GetHeading returns the writer heading.
func (c Config) GetHeading() string {
	if len(c.Heading) > 0 {
		return c.Heading
	}
	return ""
}

// GetOutputFormat returns the output format.
func (c Config) GetOutputFormat() OutputFormat {
	if len(c.OutputFormat) > 0 {
		return OutputFormat(strings.ToLower(c.OutputFormat))
	}
	return OutputFormatText
}

// GetFlags returns the enabled logger events.
func (c Config) GetFlags() []string {
	if len(c.Flags) > 0 {
		return c.Flags
	}
	return AsStrings(DefaultFlags...)
}

// GetHiddenFlags returns the enabled logger events.
func (c Config) GetHiddenFlags() []string {
	if len(c.HiddenFlags) > 0 {
		return c.HiddenFlags
	}
	return AsStrings(DefaultHiddenFlags...)
}

// GetRecoverPanics returns a field value or a default.
func (c Config) GetRecoverPanics(defaults ...bool) bool {
	if c.RecoverPanics != nil {
		return *c.RecoverPanics
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return DefaultRecoverPanics
}

// GetWriteQueueDepth returns the config queue depth.
func (c Config) GetWriteQueueDepth(defaults ...int) int {
	if c.WriteQueueDepth > 0 {
		return c.WriteQueueDepth
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return DefaultWriteQueueDepth
}

// GetListenerQueueDepth returns the config queue depth.
func (c Config) GetListenerQueueDepth(defaults ...int) int {
	if c.ListenerQueueDepth > 0 {
		return c.ListenerQueueDepth
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return DefaultListenerQueueDepth
}

// GetWriters returns the configured writers
func (c Config) GetWriters() []Writer {
	switch OutputFormat(strings.ToLower(string(c.GetOutputFormat()))) {
	case OutputFormatJSON:
		return []Writer{NewJSONWriterFromConfig(&c.JSON)}
	case OutputFormatText:
		return []Writer{NewTextWriterFromConfig(&c.Text)}
	default:
		return []Writer{NewTextWriterFromConfig(&c.Text)}
	}
}

// NewTextWriterConfigFromEnv returns a new text writer config from the environment.
func NewTextWriterConfigFromEnv() *TextWriterConfig {
	var config TextWriterConfig
	if err := env.Env().ReadInto(&config); err != nil {
		panic(err)
	}
	return &config
}

// TextWriterConfig is the config for a text writer.
type TextWriterConfig struct {
	ShowHeadings  *bool  `json:"showHeadings,omitempty" yaml:"showHeadings,omitempty" env:"LOG_SHOW_HEADINGS"`
	ShowTimestamp *bool  `json:"showTimestamp,omitempty" yaml:"showTimestamp,omitempty" env:"LOG_SHOW_TIMESTAMP"`
	UseColor      *bool  `json:"useColor,omitempty" yaml:"useColor,omitempty" env:"LOG_USE_COLOR"`
	TimeFormat    string `json:"timeFormat,omitempty" yaml:"timeFormat,omitempty" env:"LOG_TIME_FORMAT"`
}

// GetShowHeadings returns a field value or a default.
func (twc TextWriterConfig) GetShowHeadings(defaults ...bool) bool {
	if twc.ShowHeadings != nil {
		return *twc.ShowHeadings
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return DefaultTextWriterShowHeadings
}

// GetShowTimestamp returns a field value or a default.
func (twc TextWriterConfig) GetShowTimestamp(defaults ...bool) bool {
	if twc.ShowTimestamp != nil {
		return *twc.ShowTimestamp
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return DefaultTextWriterShowTimestamp
}

// GetUseColor returns a field value or a default.
func (twc TextWriterConfig) GetUseColor(defaults ...bool) bool {
	if twc.UseColor != nil {
		return *twc.UseColor
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return DefaultTextWriterUseColor
}

// GetTimeFormat returns a field value or a default.
func (twc TextWriterConfig) GetTimeFormat(defaults ...string) string {
	if len(twc.TimeFormat) > 0 {
		return twc.TimeFormat
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return DefaultTextTimeFormat
}

// NewJSONWriterConfigFromEnv returns a new json writer config from the environment.
func NewJSONWriterConfigFromEnv() *JSONWriterConfig {
	var config JSONWriterConfig
	if err := env.Env().ReadInto(&config); err != nil {
		panic(err)
	}
	return &config
}

// JSONWriterConfig is the config for a json writer.
type JSONWriterConfig struct {
	Pretty *bool `json:"pretty,omitempty" yaml:"pretty,omitempty" env:"LOG_JSON_PRETTY"`
}

// GetPretty returns a field value or a default.
func (jwc JSONWriterConfig) GetPretty(defaults ...bool) bool {
	if jwc.Pretty != nil {
		return *jwc.Pretty
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return DefaultJSONWriterPretty
}
