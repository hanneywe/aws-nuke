package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/workspaces"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockWorkSpacesV2Client struct {
	mock.Mock
}

func (m *mockWorkSpacesV2Client) DescribeIpGroups(
	ctx context.Context, params *workspaces.DescribeIpGroupsInput,
	_ ...func(*workspaces.Options),
) (*workspaces.DescribeIpGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*workspaces.DescribeIpGroupsOutput), args.Error(1)
}

func (m *mockWorkSpacesV2Client) DeleteIpGroup(
	ctx context.Context, params *workspaces.DeleteIpGroupInput,
	_ ...func(*workspaces.Options),
) (*workspaces.DeleteIpGroupOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*workspaces.DeleteIpGroupOutput), args.Error(1)
}

var testWorkSpacesV2ListerOpts = &nuke.ListerOpts{}
