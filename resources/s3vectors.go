package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3vectors"
)

type S3VectorsClient interface {
	ListVectorBuckets(ctx context.Context, params *s3vectors.ListVectorBucketsInput,
		optFns ...func(*s3vectors.Options)) (*s3vectors.ListVectorBucketsOutput, error)
	DeleteVectorBucket(ctx context.Context, params *s3vectors.DeleteVectorBucketInput,
		optFns ...func(*s3vectors.Options)) (*s3vectors.DeleteVectorBucketOutput, error)
	ListIndexes(ctx context.Context, params *s3vectors.ListIndexesInput,
		optFns ...func(*s3vectors.Options)) (*s3vectors.ListIndexesOutput, error)
	DeleteIndex(ctx context.Context, params *s3vectors.DeleteIndexInput,
		optFns ...func(*s3vectors.Options)) (*s3vectors.DeleteIndexOutput, error)
}
