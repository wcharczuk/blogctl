package config

import (
	"runtime"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
)

// ReadConfig reads a config at a given path as yaml.
func ReadConfig(flags *PersistentFlags) (cfg Config, path string, err error) {
	path, err = configutil.Read(&cfg, configutil.OptPaths(*flags.ConfigPath), configutil.OptResolver(func(untyped interface{}) error {
		if untyped == nil || flags == nil {
			return nil
		}

		c, ok := untyped.(*Config)
		if !ok || c == nil {
			return nil
		}

		return configutil.AnyError(
			env.Env().ReadInto(c),
			configutil.SetInt(&c.Parallelism, configutil.Int(*flags.Parallelism), configutil.Int(c.Parallelism), configutil.Int(runtime.NumCPU())),
			configutil.SetStrings(&c.Logger.Flags, configutil.Strings(*flags.LoggerFlags), configutil.Strings(c.Logger.Flags)),
		)
	}))
	if configutil.IsIgnored(err) {
		err = nil
	}
	return
}
