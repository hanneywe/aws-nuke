package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/s3tables"
)

type mockS3TablesClient struct {
	mock.Mock
}

func (m *mockS3TablesClient) ListTableBuckets(
	ctx context.Context, params *s3tables.ListTableBucketsInput,
	_ ...func(*s3tables.Options),
) (*s3tables.ListTableBucketsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*s3tables.ListTableBucketsOutput), args.Error(1)
}

func (m *mockS3TablesClient) DeleteTableBucket(
	ctx context.Context, params *s3tables.DeleteTableBucketInput,
	_ ...func(*s3tables.Options),
) (*s3tables.DeleteTableBucketOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*s3tables.DeleteTableBucketOutput), args.Error(1)
}
