package web

import (
	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
)

// NewHTTPSUpgraderConfigFromEnv returns an https upgrader config populated from the environment.
func NewHTTPSUpgraderConfigFromEnv() *HTTPSUpgraderConfig {
	var cfg HTTPSUpgraderConfig
	if err := env.Env().ReadInto(&cfg); err != nil {
		panic(err)
	}
	return &cfg
}

// HTTPSUpgraderConfig is the config for the https upgrader server.
type HTTPSUpgraderConfig struct {
	TargetPort int32 `json:"targetPort" yaml:"targetPort" env:"UPGRADE_TARGET_PORT"`
}

// GetTargetPort gets the target port.
// It defaults to unset, i.e. use the https default of 443.
func (c HTTPSUpgraderConfig) GetTargetPort(defaults ...int32) int32 {
	return configutil.CoalesceInt32(c.TargetPort, 0, defaults...)
}
