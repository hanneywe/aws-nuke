package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testCloudFrontListerOpts = &nuke.ListerOpts{
	Config: &aws.Config{
		Region: "us-east-1",
	},
}

type mockCloudFrontClient struct {
	mock.Mock
}

func (m *mockCloudFrontClient) ListDistributions(ctx context.Context, params *cloudfront.ListDistributionsInput,
	_ ...func(*cloudfront.Options)) (*cloudfront.ListDistributionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudfront.ListDistributionsOutput), args.Error(1)
}

func (m *mockCloudFrontClient) ListTagsForResource(ctx context.Context, params *cloudfront.ListTagsForResourceInput,
	_ ...func(*cloudfront.Options)) (*cloudfront.ListTagsForResourceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudfront.ListTagsForResourceOutput), args.Error(1)
}

func (m *mockCloudFrontClient) GetDistributionConfig(ctx context.Context, params *cloudfront.GetDistributionConfigInput,
	_ ...func(*cloudfront.Options)) (*cloudfront.GetDistributionConfigOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudfront.GetDistributionConfigOutput), args.Error(1)
}

func (m *mockCloudFrontClient) UpdateDistribution(ctx context.Context, params *cloudfront.UpdateDistributionInput,
	_ ...func(*cloudfront.Options)) (*cloudfront.UpdateDistributionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudfront.UpdateDistributionOutput), args.Error(1)
}

func (m *mockCloudFrontClient) DeleteDistribution(ctx context.Context, params *cloudfront.DeleteDistributionInput,
	_ ...func(*cloudfront.Options)) (*cloudfront.DeleteDistributionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudfront.DeleteDistributionOutput), args.Error(1)
}

func (m *mockCloudFrontClient) ListKeyValueStores(ctx context.Context, params *cloudfront.ListKeyValueStoresInput,
	_ ...func(*cloudfront.Options)) (*cloudfront.ListKeyValueStoresOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudfront.ListKeyValueStoresOutput), args.Error(1)
}

func (m *mockCloudFrontClient) DescribeKeyValueStore(ctx context.Context, params *cloudfront.DescribeKeyValueStoreInput,
	_ ...func(*cloudfront.Options)) (*cloudfront.DescribeKeyValueStoreOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudfront.DescribeKeyValueStoreOutput), args.Error(1)
}

func (m *mockCloudFrontClient) DeleteKeyValueStore(ctx context.Context, params *cloudfront.DeleteKeyValueStoreInput,
	_ ...func(*cloudfront.Options)) (*cloudfront.DeleteKeyValueStoreOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudfront.DeleteKeyValueStoreOutput), args.Error(1)
}

func (m *mockCloudFrontClient) ListRealtimeLogConfigs(ctx context.Context, params *cloudfront.ListRealtimeLogConfigsInput,
	_ ...func(*cloudfront.Options)) (*cloudfront.ListRealtimeLogConfigsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudfront.ListRealtimeLogConfigsOutput), args.Error(1)
}

func (m *mockCloudFrontClient) DeleteRealtimeLogConfig(ctx context.Context, params *cloudfront.DeleteRealtimeLogConfigInput,
	_ ...func(*cloudfront.Options)) (*cloudfront.DeleteRealtimeLogConfigOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudfront.DeleteRealtimeLogConfigOutput), args.Error(1)
}
