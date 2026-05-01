package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
)

type CloudFrontClient interface {
	ListDistributions(ctx context.Context, params *cloudfront.ListDistributionsInput,
		optFns ...func(*cloudfront.Options)) (*cloudfront.ListDistributionsOutput, error)
	ListTagsForResource(ctx context.Context, params *cloudfront.ListTagsForResourceInput,
		optFns ...func(*cloudfront.Options)) (*cloudfront.ListTagsForResourceOutput, error)
	GetDistributionConfig(ctx context.Context, params *cloudfront.GetDistributionConfigInput,
		optFns ...func(*cloudfront.Options)) (*cloudfront.GetDistributionConfigOutput, error)
	UpdateDistribution(ctx context.Context, params *cloudfront.UpdateDistributionInput,
		optFns ...func(*cloudfront.Options)) (*cloudfront.UpdateDistributionOutput, error)
	DeleteDistribution(ctx context.Context, params *cloudfront.DeleteDistributionInput,
		optFns ...func(*cloudfront.Options)) (*cloudfront.DeleteDistributionOutput, error)
	ListKeyValueStores(ctx context.Context, params *cloudfront.ListKeyValueStoresInput,
		optFns ...func(*cloudfront.Options)) (*cloudfront.ListKeyValueStoresOutput, error)
	DescribeKeyValueStore(ctx context.Context, params *cloudfront.DescribeKeyValueStoreInput,
		optFns ...func(*cloudfront.Options)) (*cloudfront.DescribeKeyValueStoreOutput, error)
	DeleteKeyValueStore(ctx context.Context, params *cloudfront.DeleteKeyValueStoreInput,
		optFns ...func(*cloudfront.Options)) (*cloudfront.DeleteKeyValueStoreOutput, error)

	ListRealtimeLogConfigs(ctx context.Context, params *cloudfront.ListRealtimeLogConfigsInput,
		optFns ...func(*cloudfront.Options)) (*cloudfront.ListRealtimeLogConfigsOutput, error)
	DeleteRealtimeLogConfig(ctx context.Context, params *cloudfront.DeleteRealtimeLogConfigInput,
		optFns ...func(*cloudfront.Options)) (*cloudfront.DeleteRealtimeLogConfigOutput, error)
}
