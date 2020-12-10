package aws

import (
	"context"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
)

const (
	// DefaultAWSRegion is a default.
	DefaultAWSRegion = "us-east-1"
)

// NewConfigFromEnv returns a new aws config from the environment.
func NewConfigFromEnv() (*Config, error) {
	var config Config
	if err := env.Env().ReadInto(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Config is a config object.
type Config struct {
	Region          string `json:"region,omitempty" yaml:"region,omitempty"`
	AccessKeyID     string `json:"accessKeyID,omitempty" yaml:"accessKeyID,omitempty"`
	SecretAccessKey string `json:"secretAccessKey,omitempty" yaml:"secretAccessKey,omitempty"`
	SecurityToken   string `json:"securityToken,omitempty" yaml:"securityToken,omitempty"`
}

// Resolve resolves the config.
func (c *Config) Resolve(ctx context.Context) error {
	return configutil.Resolve(ctx,
		configutil.SetString(&c.Region, configutil.String(c.Region), configutil.Env("AWS_REGION"), configutil.String(DefaultAWSRegion)),
		configutil.SetString(&c.AccessKeyID, configutil.String(c.AccessKeyID), configutil.Env("AWS_ACCESS_KEY_ID")),
		configutil.SetString(&c.SecretAccessKey, configutil.String(c.SecretAccessKey), configutil.Env("AWS_SECRET_ACCESS_KEY")),
		configutil.SetString(&c.SecurityToken, configutil.String(c.SecurityToken), configutil.Env("AWS_SECURITY_TOKEN")),
	)
}

// IsZero returns if the config is unset or not.
func (c Config) IsZero() bool {
	return len(c.AccessKeyID) == 0 || len(c.SecretAccessKey) == 0
}
