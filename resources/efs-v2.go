package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/efs"
)

// EFSV2Client is an interface for the AWS EFS SDK v2 client methods used by EFS resources.
// This is separate from the existing EFS resources which use the SDK v1 client.
type EFSV2Client interface {
	DescribeAccessPoints(ctx context.Context, params *efs.DescribeAccessPointsInput,
		optFns ...func(*efs.Options)) (*efs.DescribeAccessPointsOutput, error)
	DeleteAccessPoint(ctx context.Context, params *efs.DeleteAccessPointInput,
		optFns ...func(*efs.Options)) (*efs.DeleteAccessPointOutput, error)
	DescribeFileSystems(ctx context.Context, params *efs.DescribeFileSystemsInput,
		optFns ...func(*efs.Options)) (*efs.DescribeFileSystemsOutput, error)
	DescribeBackupPolicy(ctx context.Context, params *efs.DescribeBackupPolicyInput,
		optFns ...func(*efs.Options)) (*efs.DescribeBackupPolicyOutput, error)
	PutBackupPolicy(ctx context.Context, params *efs.PutBackupPolicyInput,
		optFns ...func(*efs.Options)) (*efs.PutBackupPolicyOutput, error)
	ListTagsForResource(ctx context.Context, params *efs.ListTagsForResourceInput,
		optFns ...func(*efs.Options)) (*efs.ListTagsForResourceOutput, error)
	UntagResource(ctx context.Context, params *efs.UntagResourceInput,
		optFns ...func(*efs.Options)) (*efs.UntagResourceOutput, error)
	DescribeLifecycleConfiguration(ctx context.Context, params *efs.DescribeLifecycleConfigurationInput,
		optFns ...func(*efs.Options)) (*efs.DescribeLifecycleConfigurationOutput, error)
	PutLifecycleConfiguration(ctx context.Context, params *efs.PutLifecycleConfigurationInput,
		optFns ...func(*efs.Options)) (*efs.PutLifecycleConfigurationOutput, error)
	DescribeReplicationConfigurations(ctx context.Context, params *efs.DescribeReplicationConfigurationsInput,
		optFns ...func(*efs.Options)) (*efs.DescribeReplicationConfigurationsOutput, error)
	DeleteReplicationConfiguration(ctx context.Context, params *efs.DeleteReplicationConfigurationInput,
		optFns ...func(*efs.Options)) (*efs.DeleteReplicationConfigurationOutput, error)
}
