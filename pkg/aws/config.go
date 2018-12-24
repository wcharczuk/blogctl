package aws

import (
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
	Region          string `json:"region,omitempty" yaml:"region,omitempty" env:"AWS_REGION"`
	AccessKeyID     string `json:"accessKeyID,omitempty" yaml:"accessKeyID,omitempty" env:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `json:"secretAccessKey,omitempty" yaml:"secretAccessKey,omitempty" env:"AWS_SECRET_ACCESS_KEY"`
	SecurityToken   string `json:"securityToken,omitempty" yaml:"securityToken,omitempty" env:"AWS_SECURITY_TOKEN"`
}

// IsZero returns if the config is unset or not.
func (a Config) IsZero() bool {
	return len(a.AccessKeyID) == 0 || len(a.SecretAccessKey) == 0
}

// GetRegion gets a property or a default.
func (a Config) GetRegion(defaults ...string) string {
	return configutil.CoalesceString(a.Region, DefaultAWSRegion, defaults...)
}

// GetAccessKeyID gets a property or a default.
func (a Config) GetAccessKeyID(defaults ...string) string {
	return configutil.CoalesceString(a.AccessKeyID, "", defaults...)
}

// GetSecretAccessKey gets a property or a default.
func (a Config) GetSecretAccessKey(defaults ...string) string {
	return configutil.CoalesceString(a.SecretAccessKey, "", defaults...)
}

// GetToken returns a secret access token or a default.
func (a Config) GetToken(defaults ...string) string {
	return configutil.CoalesceString(a.SecurityToken, "", defaults...)
}
