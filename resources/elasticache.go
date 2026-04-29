package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/elasticache"
)

// ElasticacheClient is an interface for the ElastiCache SDK v2 client methods used by all
// ElastiCache resources that use SDK v2. It enables mock testing of List and Remove operations.
type ElasticacheClient interface {
	// Listing
	DescribeServerlessCaches(ctx context.Context, params *elasticache.DescribeServerlessCachesInput,
		optFns ...func(*elasticache.Options)) (*elasticache.DescribeServerlessCachesOutput, error)
	DescribeSnapshots(ctx context.Context, params *elasticache.DescribeSnapshotsInput,
		optFns ...func(*elasticache.Options)) (*elasticache.DescribeSnapshotsOutput, error)
	ListTagsForResource(ctx context.Context, params *elasticache.ListTagsForResourceInput,
		optFns ...func(*elasticache.Options)) (*elasticache.ListTagsForResourceOutput, error)

	// Modification
	ModifyServerlessCache(ctx context.Context, params *elasticache.ModifyServerlessCacheInput,
		optFns ...func(*elasticache.Options)) (*elasticache.ModifyServerlessCacheOutput, error)

	// Deletion
	DeleteServerlessCache(ctx context.Context, params *elasticache.DeleteServerlessCacheInput,
		optFns ...func(*elasticache.Options)) (*elasticache.DeleteServerlessCacheOutput, error)
	DeleteSnapshot(ctx context.Context, params *elasticache.DeleteSnapshotInput,
		optFns ...func(*elasticache.Options)) (*elasticache.DeleteSnapshotOutput, error)

	DescribeServerlessCacheSnapshots(ctx context.Context, params *elasticache.DescribeServerlessCacheSnapshotsInput,
		optFns ...func(*elasticache.Options)) (*elasticache.DescribeServerlessCacheSnapshotsOutput, error)
	DeleteServerlessCacheSnapshot(ctx context.Context, params *elasticache.DeleteServerlessCacheSnapshotInput,
		optFns ...func(*elasticache.Options)) (*elasticache.DeleteServerlessCacheSnapshotOutput, error)
}
