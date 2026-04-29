package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/docdb"
)

// DocDBV2Client is an interface for the AWS DocumentDB SDK v2 client methods used by DocDB global cluster resources.
// It enables mock testing of List, Remove, and Modify operations.
// This is separate from the existing DocDB resources which use the SDK v2 client directly without an interface.
type DocDBV2Client interface {
	DescribeGlobalClusters(ctx context.Context, params *docdb.DescribeGlobalClustersInput,
		optFns ...func(*docdb.Options)) (*docdb.DescribeGlobalClustersOutput, error)
	DeleteGlobalCluster(ctx context.Context, params *docdb.DeleteGlobalClusterInput,
		optFns ...func(*docdb.Options)) (*docdb.DeleteGlobalClusterOutput, error)
	ModifyGlobalCluster(ctx context.Context, params *docdb.ModifyGlobalClusterInput,
		optFns ...func(*docdb.Options)) (*docdb.ModifyGlobalClusterOutput, error)
}
