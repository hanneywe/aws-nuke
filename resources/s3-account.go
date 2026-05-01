package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3control"
)

type S3AccountClient interface {
	GetPublicAccessBlock(ctx context.Context, params *s3control.GetPublicAccessBlockInput,
		optFns ...func(*s3control.Options)) (*s3control.GetPublicAccessBlockOutput, error)
	DeletePublicAccessBlock(ctx context.Context, params *s3control.DeletePublicAccessBlockInput,
		optFns ...func(*s3control.Options)) (*s3control.DeletePublicAccessBlockOutput, error)
}
