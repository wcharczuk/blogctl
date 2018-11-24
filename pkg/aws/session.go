package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// NewSession creates a new aws session from a config.
func NewSession(cfg *Config) *session.Session {
	if cfg.IsZero() {
		session := session.Must(session.NewSession())
		if cfg.Region != "" {
			session.Config.Region = &cfg.Region
		}
		return session
	}

	awsConfig := &aws.Config{
		Region:      RefStr(cfg.GetRegion()),
		Credentials: credentials.NewStaticCredentials(cfg.GetAccessKeyID(), cfg.GetSecretAccessKey(), cfg.GetToken()),
	}
	return session.Must(session.NewSession(awsConfig))
}
