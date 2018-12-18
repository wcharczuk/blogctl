package config

// S3 is an optional s3 config.
type S3 struct {
	Region string `json:"region,omitempty" yaml:"region,omitempty"`
	Bucket string `json:"bucket,omitempty" yaml:"bucket,omitempty"`
}

// IsZero returns if the s3 bucket config is set.
func (c S3) IsZero() bool {
	return c.Bucket == ""
}
