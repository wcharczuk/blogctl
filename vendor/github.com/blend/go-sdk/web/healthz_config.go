package web

import (
	"time"

	"github.com/blend/go-sdk/configutil"
)

// HealthzConfig is the healthz config.
type HealthzConfig struct {
	BindAddr         string        `json:"bindAddr" yaml:"bindAddr" env:"HEALTHZ_BIND_ADDR"`
	GracePeriod      time.Duration `json:"gracePeriod" yaml:"gracePeriod"`
	FailureThreshold int           `json:"failureThreshold" yaml:"failureThreshold" env:"READY_FAILURE_THRESHOLD"`
	RecoverPanics    *bool         `json:"recoverPanics" yaml:"recoverPanics"`

	MaxHeaderBytes    int           `json:"maxHeaderBytes,omitempty" yaml:"maxHeaderBytes,omitempty" env:"MAX_HEADER_BYTES"`
	ReadTimeout       time.Duration `json:"readTimeout,omitempty" yaml:"readTimeout,omitempty" env:"READ_HEADER_TIMEOUT"`
	ReadHeaderTimeout time.Duration `json:"readHeaderTimeout,omitempty" yaml:"readHeaderTimeout,omitempty" env:"READ_HEADER_TIMEOUT"`
	WriteTimeout      time.Duration `json:"writeTimeout,omitempty" yaml:"writeTimeout,omitempty" env:"WRITE_TIMEOUT"`
	IdleTimeout       time.Duration `json:"idleTimeout,omitempty" yaml:"idleTimeout,omitempty" env:"IDLE_TIMEOUT"`
}

// GetBindAddr gets the bind address.
func (hzc HealthzConfig) GetBindAddr(defaults ...string) string {
	return configutil.CoalesceString(hzc.BindAddr, DefaultHealthzBindAddr, defaults...)
}

// GetGracePeriod gets a grace period or a default.
func (hzc HealthzConfig) GetGracePeriod(defaults ...time.Duration) time.Duration {
	return configutil.CoalesceDuration(hzc.GracePeriod, DefaultShutdownGracePeriod, defaults...)
}

// GetRecoverPanics gets recover panics or a default.
func (hzc HealthzConfig) GetRecoverPanics(defaults ...bool) bool {
	return configutil.CoalesceBool(hzc.RecoverPanics, DefaultRecoverPanics, defaults...)
}

// GetFailureThreshold gets the failure threshold or a default.
func (hzc HealthzConfig) GetFailureThreshold(defaults ...int) int {
	return configutil.CoalesceInt(hzc.FailureThreshold, DefaultHealthzFailureThreshold, defaults...)
}

// GetMaxHeaderBytes returns the maximum header size in bytes or a default.
func (hzc HealthzConfig) GetMaxHeaderBytes(defaults ...int) int {
	return configutil.CoalesceInt(hzc.MaxHeaderBytes, DefaultMaxHeaderBytes, defaults...)
}

// GetReadTimeout gets a property.
func (hzc HealthzConfig) GetReadTimeout(defaults ...time.Duration) time.Duration {
	return configutil.CoalesceDuration(hzc.ReadTimeout, DefaultReadTimeout, defaults...)
}

// GetReadHeaderTimeout gets a property.
func (hzc HealthzConfig) GetReadHeaderTimeout(defaults ...time.Duration) time.Duration {
	return configutil.CoalesceDuration(hzc.ReadHeaderTimeout, DefaultReadHeaderTimeout, defaults...)
}

// GetWriteTimeout gets a property.
func (hzc HealthzConfig) GetWriteTimeout(defaults ...time.Duration) time.Duration {
	return configutil.CoalesceDuration(hzc.WriteTimeout, DefaultWriteTimeout, defaults...)
}

// GetIdleTimeout gets a property.
func (hzc HealthzConfig) GetIdleTimeout(defaults ...time.Duration) time.Duration {
	return configutil.CoalesceDuration(hzc.IdleTimeout, DefaultIdleTimeout, defaults...)
}
