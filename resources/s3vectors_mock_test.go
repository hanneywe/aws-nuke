package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/s3vectors"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockS3VectorsClient struct {
	mock.Mock
}

func (m *mockS3VectorsClient) ListVectorBuckets(
	ctx context.Context, params *s3vectors.ListVectorBucketsInput,
	_ ...func(*s3vectors.Options),
) (*s3vectors.ListVectorBucketsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*s3vectors.ListVectorBucketsOutput), args.Error(1)
}

func (m *mockS3VectorsClient) DeleteVectorBucket(
	ctx context.Context, params *s3vectors.DeleteVectorBucketInput,
	_ ...func(*s3vectors.Options),
) (*s3vectors.DeleteVectorBucketOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*s3vectors.DeleteVectorBucketOutput), args.Error(1)
}

func (m *mockS3VectorsClient) ListIndexes(
	ctx context.Context, params *s3vectors.ListIndexesInput,
	_ ...func(*s3vectors.Options),
) (*s3vectors.ListIndexesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*s3vectors.ListIndexesOutput), args.Error(1)
}

func (m *mockS3VectorsClient) DeleteIndex(
	ctx context.Context, params *s3vectors.DeleteIndexInput,
	_ ...func(*s3vectors.Options),
) (*s3vectors.DeleteIndexOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*s3vectors.DeleteIndexOutput), args.Error(1)
}

var testS3VectorsListerOpts = &nuke.ListerOpts{}
