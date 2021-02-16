package config

import (
	"github.com/blend/go-sdk/configutil"
)

// ReadConfig reads a config at a given path as yaml.
func ReadConfig(flags Flags) (cfg Config, configPaths []string, err error) {
	configPaths, err = configutil.Read(&cfg,
		configutil.OptAddPreferredPaths(*flags.ConfigPath),
	)
	if configutil.IsIgnored(err) {
		err = nil
	}
	return
}
