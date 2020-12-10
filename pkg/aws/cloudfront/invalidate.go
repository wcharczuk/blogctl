package cloudfront

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/uuid"

	"github.com/wcharczuk/blogctl/pkg/aws"
)

// InvalidateMany invalidates a given set of paths for a distribution.
func InvalidateMany(ctx context.Context, session *session.Session, distribution string, items ...string) error {
	_, err := cloudfront.New(session).CreateInvalidationWithContext(ctx, &cloudfront.CreateInvalidationInput{
		DistributionId: aws.RefStr(distribution),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: aws.RefStr(uuid.V4().String()),
			Paths: &cloudfront.Paths{
				Items:    aws.RefStrs(items...),
				Quantity: quantity(items...),
			},
		},
	})
	return ex.New(err)
}

func quantity(items ...string) *int64 {
	l := int64(len(items))
	return &l
}
