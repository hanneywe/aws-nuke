package resources

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/elasticache"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testElasticacheListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

type mockElasticacheClient struct {
	mock.Mock
}

func (m *mockElasticacheClient) DescribeServerlessCaches(ctx context.Context, params *elasticache.DescribeServerlessCachesInput,
	_ ...func(*elasticache.Options)) (*elasticache.DescribeServerlessCachesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*elasticache.DescribeServerlessCachesOutput), args.Error(1)
}

func (m *mockElasticacheClient) ModifyServerlessCache(ctx context.Context, params *elasticache.ModifyServerlessCacheInput,
	_ ...func(*elasticache.Options)) (*elasticache.ModifyServerlessCacheOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*elasticache.ModifyServerlessCacheOutput), args.Error(1)
}

func (m *mockElasticacheClient) DeleteServerlessCache(ctx context.Context, params *elasticache.DeleteServerlessCacheInput,
	_ ...func(*elasticache.Options)) (*elasticache.DeleteServerlessCacheOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*elasticache.DeleteServerlessCacheOutput), args.Error(1)
}

func (m *mockElasticacheClient) DescribeSnapshots(ctx context.Context, params *elasticache.DescribeSnapshotsInput,
	_ ...func(*elasticache.Options)) (*elasticache.DescribeSnapshotsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*elasticache.DescribeSnapshotsOutput), args.Error(1)
}

func (m *mockElasticacheClient) ListTagsForResource(ctx context.Context, params *elasticache.ListTagsForResourceInput,
	_ ...func(*elasticache.Options)) (*elasticache.ListTagsForResourceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*elasticache.ListTagsForResourceOutput), args.Error(1)
}

func (m *mockElasticacheClient) DeleteSnapshot(ctx context.Context, params *elasticache.DeleteSnapshotInput,
	_ ...func(*elasticache.Options)) (*elasticache.DeleteSnapshotOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*elasticache.DeleteSnapshotOutput), args.Error(1)
}

func (m *mockElasticacheClient) DescribeServerlessCacheSnapshots(
	ctx context.Context, params *elasticache.DescribeServerlessCacheSnapshotsInput,
	_ ...func(*elasticache.Options),
) (*elasticache.DescribeServerlessCacheSnapshotsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*elasticache.DescribeServerlessCacheSnapshotsOutput), args.Error(1)
}

func (m *mockElasticacheClient) DeleteServerlessCacheSnapshot(
	ctx context.Context, params *elasticache.DeleteServerlessCacheSnapshotInput,
	_ ...func(*elasticache.Options),
) (*elasticache.DeleteServerlessCacheSnapshotOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*elasticache.DeleteServerlessCacheSnapshotOutput), args.Error(1)
}
