package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/neptune"
)

// NeptuneV2Client is an interface for the AWS Neptune SDK v2 client methods used by Neptune sub-resources.
// It enables mock testing of List, Remove, and Modify operations.
// This is separate from the existing Neptune resources which use SDK v1.
type NeptuneV2Client interface {
	DescribeDBParameterGroups(ctx context.Context, params *neptune.DescribeDBParameterGroupsInput,
		optFns ...func(*neptune.Options)) (*neptune.DescribeDBParameterGroupsOutput, error)
	DeleteDBParameterGroup(ctx context.Context, params *neptune.DeleteDBParameterGroupInput,
		optFns ...func(*neptune.Options)) (*neptune.DeleteDBParameterGroupOutput, error)
	DescribeDBClusterParameterGroups(ctx context.Context, params *neptune.DescribeDBClusterParameterGroupsInput,
		optFns ...func(*neptune.Options)) (*neptune.DescribeDBClusterParameterGroupsOutput, error)
	DeleteDBClusterParameterGroup(ctx context.Context, params *neptune.DeleteDBClusterParameterGroupInput,
		optFns ...func(*neptune.Options)) (*neptune.DeleteDBClusterParameterGroupOutput, error)
	DescribeEventSubscriptions(ctx context.Context, params *neptune.DescribeEventSubscriptionsInput,
		optFns ...func(*neptune.Options)) (*neptune.DescribeEventSubscriptionsOutput, error)
	DeleteEventSubscription(ctx context.Context, params *neptune.DeleteEventSubscriptionInput,
		optFns ...func(*neptune.Options)) (*neptune.DeleteEventSubscriptionOutput, error)
	DescribeGlobalClusters(ctx context.Context, params *neptune.DescribeGlobalClustersInput,
		optFns ...func(*neptune.Options)) (*neptune.DescribeGlobalClustersOutput, error)
	DeleteGlobalCluster(ctx context.Context, params *neptune.DeleteGlobalClusterInput,
		optFns ...func(*neptune.Options)) (*neptune.DeleteGlobalClusterOutput, error)
	ModifyGlobalCluster(ctx context.Context, params *neptune.ModifyGlobalClusterInput,
		optFns ...func(*neptune.Options)) (*neptune.ModifyGlobalClusterOutput, error)
	DescribeDBClusterEndpoints(ctx context.Context, params *neptune.DescribeDBClusterEndpointsInput,
		optFns ...func(*neptune.Options)) (*neptune.DescribeDBClusterEndpointsOutput, error)
	DeleteDBClusterEndpoint(ctx context.Context, params *neptune.DeleteDBClusterEndpointInput,
		optFns ...func(*neptune.Options)) (*neptune.DeleteDBClusterEndpointOutput, error)
}
