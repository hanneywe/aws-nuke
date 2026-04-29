package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/synthetics"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockSyntheticsClient struct {
	mock.Mock
}

func (m *mockSyntheticsClient) ListGroups(
	ctx context.Context, params *synthetics.ListGroupsInput,
	_ ...func(*synthetics.Options),
) (*synthetics.ListGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*synthetics.ListGroupsOutput), args.Error(1)
}

func (m *mockSyntheticsClient) DeleteGroup(
	ctx context.Context, params *synthetics.DeleteGroupInput,
	_ ...func(*synthetics.Options),
) (*synthetics.DeleteGroupOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*synthetics.DeleteGroupOutput), args.Error(1)
}

var testSyntheticsListerOpts = &nuke.ListerOpts{}
