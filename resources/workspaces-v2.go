package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/workspaces"
)

type WorkSpacesV2Client interface {
	DescribeIpGroups(ctx context.Context, params *workspaces.DescribeIpGroupsInput,
		optFns ...func(*workspaces.Options)) (*workspaces.DescribeIpGroupsOutput, error)
	DeleteIpGroup(ctx context.Context, params *workspaces.DeleteIpGroupInput,
		optFns ...func(*workspaces.Options)) (*workspaces.DeleteIpGroupOutput, error)
}
