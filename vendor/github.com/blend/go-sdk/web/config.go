package web

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/webutil"
)

// NewConfigFromEnv returns a new config from the environment.
func NewConfigFromEnv() (*Config, error) {
	var cfg Config
	if err := env.Env().ReadInto(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// MustNewConfigFromEnv returns a new config from environment variables,
// and will panic on error.
func MustNewConfigFromEnv() *Config {
	cfg, err := NewConfigFromEnv()
	if err != nil {
		panic(err)
	}
	return cfg
}

// Config is an object used to set up a web app.
type Config struct {
	Port     int32  `json:"port,omitempty" yaml:"port,omitempty" env:"PORT"`
	BindAddr string `json:"bindAddr,omitempty" yaml:"bindAddr,omitempty" env:"BIND_ADDR"`
	BaseURL  string `json:"baseURL,omitempty" yaml:"baseURL,omitempty" env:"BASE_URL"`

	RedirectTrailingSlash  *bool `json:"redirectTrailingSlash,omitempty" yaml:"redirectTrailingSlash,omitempty"`
	HandleOptions          *bool `json:"handleOptions,omitempty" yaml:"handleOptions,omitempty"`
	HandleMethodNotAllowed *bool `json:"handleMethodNotAllowed,omitempty" yaml:"handleMethodNotAllowed,omitempty"`
	RecoverPanics          *bool `json:"recoverPanics,omitempty" yaml:"recoverPanics,omitempty"`

	// AuthManagerMode is a mode designation for the auth manager.
	AuthManagerMode string `json:"authManagerMode" yaml:"authManagerMode"`
	// AuthSecret is a secret key to use with auth management.
	AuthSecret string `json:"authSecret" yaml:"authSecret" env:"AUTH_SECRET"`
	// SessionTimeout is a fixed duration to use when calculating hard or rolling deadlines.
	SessionTimeout time.Duration `json:"sessionTimeout,omitempty" yaml:"sessionTimeout,omitempty" env:"SESSION_TIMEOUT"`
	// SessionTimeoutIsAbsolute determines if the session timeout is a hard deadline or if it gets pushed forward with usage.
	// The default is to use a hard deadline.
	SessionTimeoutIsAbsolute *bool `json:"sessionTimeoutIsAbsolute,omitempty" yaml:"sessionTimeoutIsAbsolute,omitempty" env:"SESSION_TIMEOUT_ABSOLUTE"`
	// CookieHTTPS determines if we should flip the `https only` flag on issued cookies.
	CookieHTTPSOnly *bool `json:"cookieHTTPSOnly,omitempty" yaml:"cookieHTTPSOnly,omitempty" env:"COOKIE_HTTPS_ONLY"`
	// CookieName is the name of the cookie to issue with sessions.
	CookieName string `json:"cookieName,omitempty" yaml:"cookieName,omitempty" env:"COOKIE_NAME"`
	// CookiePath is the path on the cookie to issue with sessions.
	CookiePath string `json:"cookiePath,omitempty" yaml:"cookiePath,omitempty" env:"COOKIE_PATH"`

	// DefaultHeaders are included on any responses. The app ships with a set of default headers, which you can augment with this property.
	DefaultHeaders map[string]string `json:"defaultHeaders,omitempty" yaml:"defaultHeaders,omitempty"`

	MaxHeaderBytes    int           `json:"maxHeaderBytes,omitempty" yaml:"maxHeaderBytes,omitempty" env:"MAX_HEADER_BYTES"`
	ReadTimeout       time.Duration `json:"readTimeout,omitempty" yaml:"readTimeout,omitempty" env:"READ_HEADER_TIMEOUT"`
	ReadHeaderTimeout time.Duration `json:"readHeaderTimeout,omitempty" yaml:"readHeaderTimeout,omitempty" env:"READ_HEADER_TIMEOUT"`
	WriteTimeout      time.Duration `json:"writeTimeout,omitempty" yaml:"writeTimeout,omitempty" env:"WRITE_TIMEOUT"`
	IdleTimeout       time.Duration `json:"idleTimeout,omitempty" yaml:"idleTimeout,omitempty" env:"IDLE_TIMEOUT"`

	ShutdownGracePeriod time.Duration `json:"shutdownGracePeriod" yaml:"shutdownGracePeriod" env:"SHUTDOWN_GRACE_PERIOD"`

	HSTS  HSTSConfig      `json:"hsts,omitempty" yaml:"hsts,omitempty"`
	TLS   TLSConfig       `json:"tls,omitempty" yaml:"tls,omitempty"`
	Views ViewCacheConfig `json:"views,omitempty" yaml:"views,omitempty"`

	Healthz HealthzConfig `json:"healthz,omitempty" yaml:"healthz,omitempty"`
}

// GetBindAddr configutil.Coalesces the bind addr, the port, or the default.
func (c Config) GetBindAddr(defaults ...string) string {
	if len(c.BindAddr) > 0 {
		return c.BindAddr
	}
	if c.Port > 0 {
		return fmt.Sprintf(":%d", c.Port)
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return DefaultBindAddr
}

// GetPort returns the int32 port for a given config.
// This is useful in things like kubernetes pod templates.
// If the config .Port is unset, it will parse the .BindAddr,
// or the DefaultBindAddr for the port number.
func (c Config) GetPort(defaults ...int32) int32 {
	if c.Port > 0 {
		return c.Port
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	if len(c.BindAddr) > 0 {
		return webutil.PortFromBindAddr(c.BindAddr)
	}
	return webutil.PortFromBindAddr(DefaultBindAddr)
}

// GetBaseURL gets a property.
func (c Config) GetBaseURL(defaults ...string) string {
	return configutil.CoalesceString(c.BaseURL, "", defaults...)
}

// GetRedirectTrailingSlash returns if we automatically redirect for a missing trailing slash.
func (c Config) GetRedirectTrailingSlash(defaults ...bool) bool {
	return configutil.CoalesceBool(c.RedirectTrailingSlash, DefaultRedirectTrailingSlash, defaults...)
}

// GetHandleOptions returns if we should handle OPTIONS verb requests.
func (c Config) GetHandleOptions(defaults ...bool) bool {
	return configutil.CoalesceBool(c.HandleOptions, DefaultHandleOptions, defaults...)
}

// GetHandleMethodNotAllowed returns if we should handle method not allowed results.
func (c Config) GetHandleMethodNotAllowed(defaults ...bool) bool {
	return configutil.CoalesceBool(c.HandleMethodNotAllowed, DefaultHandleMethodNotAllowed, defaults...)
}

// GetRecoverPanics returns if we should recover panics or not.
func (c Config) GetRecoverPanics(defaults ...bool) bool {
	return configutil.CoalesceBool(c.RecoverPanics, DefaultRecoverPanics, defaults...)
}

// GetDefaultHeaders returns the default headers from the config.
func (c Config) GetDefaultHeaders(inherited ...map[string]string) map[string]string {
	output := map[string]string{}
	if len(inherited) > 0 {
		for _, set := range inherited {
			for key, value := range set {
				output[key] = value
			}
		}
	}
	for key, value := range c.DefaultHeaders {
		output[key] = value
	}
	return output
}

// ListenTLS returns if the server will directly serve requests with tls.
func (c Config) ListenTLS() bool {
	return c.TLS.HasKeyPair()
}

// BaseURLIsSecureScheme returns if the base url starts with a secure scheme.
func (c Config) BaseURLIsSecureScheme() bool {
	baseURL := c.GetBaseURL()
	if len(baseURL) == 0 {
		return false
	}
	return strings.HasPrefix(strings.ToLower(baseURL), SchemeHTTPS) || strings.HasPrefix(strings.ToLower(baseURL), SchemeSPDY)
}

// IsSecure returns if the config specifies the app will eventually be handling https requests.
func (c Config) IsSecure() bool {
	return c.ListenTLS() || c.BaseURLIsSecureScheme()
}

// GetAuthManagerMode returns the auth manager mode.
func (c Config) GetAuthManagerMode(inherited ...string) AuthManagerMode {
	return AuthManagerMode(configutil.CoalesceString(c.AuthManagerMode, string(AuthManagerModeServer), inherited...))
}

// GetAuthSecret returns a property or a default.
func (c Config) GetAuthSecret(defaults ...[]byte) []byte {
	decoded, err := base64.StdEncoding.DecodeString(c.AuthSecret)
	if err != nil {
		panic(err)
	}
	return decoded
}

// GetSessionTimeout returns a property or a default.
func (c Config) GetSessionTimeout(defaults ...time.Duration) time.Duration {
	return configutil.CoalesceDuration(c.SessionTimeout, DefaultSessionTimeout, defaults...)
}

// GetSessionTimeoutIsAbsolute returns a property or a default.
func (c Config) GetSessionTimeoutIsAbsolute(defaults ...bool) bool {
	return configutil.CoalesceBool(c.SessionTimeoutIsAbsolute, DefaultSessionTimeoutIsAbsolute, defaults...)
}

// GetCookieHTTPSOnly returns a property or a default.
func (c Config) GetCookieHTTPSOnly(defaults ...bool) bool {
	return configutil.CoalesceBool(c.CookieHTTPSOnly, c.IsSecure(), defaults...)
}

// GetCookieName returns a property or a default.
func (c Config) GetCookieName(defaults ...string) string {
	return configutil.CoalesceString(c.CookieName, DefaultCookieName, defaults...)
}

// GetCookiePath returns a property or a default.
func (c Config) GetCookiePath(defaults ...string) string {
	return configutil.CoalesceString(c.CookiePath, DefaultCookiePath, defaults...)
}

// GetMaxHeaderBytes returns the maximum header size in bytes or a default.
func (c Config) GetMaxHeaderBytes(defaults ...int) int {
	return configutil.CoalesceInt(c.MaxHeaderBytes, DefaultMaxHeaderBytes, defaults...)
}

// GetReadTimeout gets a property.
func (c Config) GetReadTimeout(defaults ...time.Duration) time.Duration {
	return configutil.CoalesceDuration(c.ReadTimeout, DefaultReadTimeout, defaults...)
}

// GetReadHeaderTimeout gets a property.
func (c Config) GetReadHeaderTimeout(defaults ...time.Duration) time.Duration {
	return configutil.CoalesceDuration(c.ReadHeaderTimeout, DefaultReadHeaderTimeout, defaults...)
}

// GetWriteTimeout gets a property.
func (c Config) GetWriteTimeout(defaults ...time.Duration) time.Duration {
	return configutil.CoalesceDuration(c.WriteTimeout, DefaultWriteTimeout, defaults...)
}

// GetIdleTimeout gets a property.
func (c Config) GetIdleTimeout(defaults ...time.Duration) time.Duration {
	return configutil.CoalesceDuration(c.IdleTimeout, DefaultIdleTimeout, defaults...)
}

// GetShutdownGracePeriod gets the shutdown grace period.
func (c Config) GetShutdownGracePeriod(defaults ...time.Duration) time.Duration {
	return configutil.CoalesceDuration(c.ShutdownGracePeriod, DefaultShutdownGracePeriod, defaults...)
}
