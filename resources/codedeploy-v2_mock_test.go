package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/codedeploy"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockCodeDeployV2Client struct {
	mock.Mock
}

func (m *mockCodeDeployV2Client) ListOnPremisesInstances(
	ctx context.Context, params *codedeploy.ListOnPremisesInstancesInput,
	_ ...func(*codedeploy.Options),
) (*codedeploy.ListOnPremisesInstancesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*codedeploy.ListOnPremisesInstancesOutput), args.Error(1)
}

func (m *mockCodeDeployV2Client) DeregisterOnPremisesInstance(
	ctx context.Context, params *codedeploy.DeregisterOnPremisesInstanceInput,
	_ ...func(*codedeploy.Options),
) (*codedeploy.DeregisterOnPremisesInstanceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*codedeploy.DeregisterOnPremisesInstanceOutput), args.Error(1)
}

var testCodeDeployV2ListerOpts = &nuke.ListerOpts{}
