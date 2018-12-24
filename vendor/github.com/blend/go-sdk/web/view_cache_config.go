package web

import "github.com/blend/go-sdk/configutil"

// ViewCacheConfig is a config for the view cache.
type ViewCacheConfig struct {
	Cached                    *bool    `json:"cached,omitempty" yaml:"cached,omitempty" env:"VIEW_CACHE_ENABLED"`
	Paths                     []string `json:"paths,omitempty" yaml:"paths,omitempty" env:"VIEW_CACHE_PATHS,csv"`
	BufferPoolSize            int      `json:"bufferPoolSize,omitempty" yaml:"bufferPoolSize,omitempty"`
	InternalErrorTemplateName string   `json:"internalErrorTemplateName,omitempty" yaml:"internalErrorTemplateName,omitempty"`
	BadRequestTemplateName    string   `json:"badRequestTemplateName,omitempty" yaml:"badRequestTemplateName,omitempty"`
	NotFoundTemplateName      string   `json:"notFoundTemplateName,omitempty" yaml:"notFoundTemplateName,omitempty"`
	NotAuthorizedTemplateName string   `json:"notAuthorizedTemplateName,omitempty" yaml:"notAuthorizedTemplateName,omitempty"`
	StatusTemplateName        string   `json:"statusTemplateName,omitempty" yaml:"statusTemplateName,omitempty"`
}

// GetCached returns if the viewcache should store templates in memory or read from disk.
func (vcc ViewCacheConfig) GetCached(defaults ...bool) bool {
	return configutil.CoalesceBool(vcc.Cached, true, defaults...)
}

// GetPaths returns default view paths.
func (vcc ViewCacheConfig) GetPaths(defaults ...[]string) []string {
	return configutil.CoalesceStrings(vcc.Paths, nil, defaults...)
}

// GetBufferPoolSize gets the buffer pool size or a default.
func (vcc ViewCacheConfig) GetBufferPoolSize(defaults ...int) int {
	return configutil.CoalesceInt(vcc.BufferPoolSize, DefaultViewBufferPoolSize, defaults...)
}

// GetInternalErrorTemplateName returns the internal error template name for the app.
func (vcc ViewCacheConfig) GetInternalErrorTemplateName(defaults ...string) string {
	return configutil.CoalesceString(vcc.InternalErrorTemplateName, DefaultTemplateNameInternalError, defaults...)
}

// GetBadRequestTemplateName returns the bad request template name for the app.
func (vcc ViewCacheConfig) GetBadRequestTemplateName(defaults ...string) string {
	return configutil.CoalesceString(vcc.BadRequestTemplateName, DefaultTemplateNameBadRequest, defaults...)
}

// GetNotFoundTemplateName returns the not found template name for the app.
func (vcc ViewCacheConfig) GetNotFoundTemplateName(defaults ...string) string {
	return configutil.CoalesceString(vcc.NotFoundTemplateName, DefaultTemplateNameNotFound, defaults...)
}

// GetNotAuthorizedTemplateName returns the not authorized template name for the app.
func (vcc ViewCacheConfig) GetNotAuthorizedTemplateName(defaults ...string) string {
	return configutil.CoalesceString(vcc.NotAuthorizedTemplateName, DefaultTemplateNameNotAuthorized, defaults...)
}

// GetStatusTemplateName returns the not authorized template name for the app.
func (vcc ViewCacheConfig) GetStatusTemplateName(defaults ...string) string {
	return configutil.CoalesceString(vcc.StatusTemplateName, DefaultTemplateNameStatus, defaults...)
}
