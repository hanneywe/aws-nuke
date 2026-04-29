package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/rds"
)

// RDSV2Client is an interface for the AWS RDS SDK v2 client methods used by RDS sub-resources.
// This is separate from the existing RDS resources which use SDK v1.
type RDSV2Client interface {
	DescribeDBSecurityGroups(ctx context.Context, params *rds.DescribeDBSecurityGroupsInput,
		optFns ...func(*rds.Options)) (*rds.DescribeDBSecurityGroupsOutput, error)
	DeleteDBSecurityGroup(ctx context.Context, params *rds.DeleteDBSecurityGroupInput,
		optFns ...func(*rds.Options)) (*rds.DeleteDBSecurityGroupOutput, error)
}
