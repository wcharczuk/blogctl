package config

// Cloudfront represents cloudfront options.
type Cloudfront struct {
	Distribution string `json:"distribution,omitempty" yaml:"distribution,omitempty"`
}

// IsZero returns if the config is set or not.
func (cf Cloudfront) IsZero() bool {
	return cf.Distribution == ""
}
