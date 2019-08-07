package config

import (
	"runtime"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
)

// ReadConfig reads a config at a given path as yaml.
func ReadConfig(flags PersistentFlags) (cfg Config, path string, err error) {
	path, err = configutil.Read(&cfg, configutil.OptPaths(*flags.ConfigPath), configutil.OptResolver(func(v interface{}) error {
		c := v.(*Config)
		return configutil.AnyError(
			env.Env().ReadInto(c),
			configutil.SetInt(&c.Parallelism, flagInt(flags.Parallelism), configutil.Int(c.Parallelism), configutil.Int(runtime.NumCPU())),
			configutil.SetStrings(&c.Logger.Flags, flagStrings(flags.LoggerFlags), configutil.Strings(c.Logger.Flags)),
		)
	}))
	if configutil.IsIgnored(err) {
		err = nil
	}
	return
}

func flagInt(flagValue *int) configutil.IntSource {
	return intSource{flagValue}
}

type intSource struct {
	Value *int
}

func (is intSource) Int() (*int, error) {
	return is.Value, nil
}

func flagStrings(flagValue *[]string) configutil.StringsSource {
	return stringsSource{flagValue}
}

type stringsSource struct {
	Value *[]string
}

func (ss stringsSource) Strings() ([]string, error) {
	if ss.Value == nil {
		return nil, nil
	}
	return *ss.Value, nil
}
