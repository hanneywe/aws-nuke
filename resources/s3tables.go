package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3tables"
)

type S3TablesClient interface {
	ListTableBuckets(ctx context.Context, params *s3tables.ListTableBucketsInput,
		optFns ...func(*s3tables.Options)) (*s3tables.ListTableBucketsOutput, error)
	DeleteTableBucket(ctx context.Context, params *s3tables.DeleteTableBucketInput,
		optFns ...func(*s3tables.Options)) (*s3tables.DeleteTableBucketOutput, error)
}
