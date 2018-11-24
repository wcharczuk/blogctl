package s3

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

// IsNotFound returns if an error is a not found error.
func IsNotFound(err error) bool {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return true
			case s3.ErrCodeNoSuchKey:
				return true
			}
		}
	}
	return false
}
