package config

// S3 is an optional s3 config.
type S3 struct {
	Region string `json:"region" yaml:"region"`
	Bucket string `json:"bucket" yaml:"bucket"`
}

// IsZero returns if the s3 bucket config is set.
func (c S3) IsZero() bool {
	return c.Bucket == ""
}
