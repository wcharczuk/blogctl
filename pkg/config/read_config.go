package config

import (
	"github.com/blend/go-sdk/configutil"
)

// ReadConfig reads a config at a given path as yaml.
func ReadConfig(flags Flags) (cfg Config, configPath string, err error) {
	configPath, err = configutil.Read(&cfg,
		configutil.OptAddPreferredPaths(*flags.ConfigPath),
	)
	if configutil.IsIgnored(err) {
		err = nil
	}
	return
}
