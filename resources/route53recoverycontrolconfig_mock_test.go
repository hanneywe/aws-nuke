package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/route53recoverycontrolconfig"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockRoute53RecoveryControlConfigClient struct {
	mock.Mock
}

func (m *mockRoute53RecoveryControlConfigClient) ListClusters(
	ctx context.Context, params *route53recoverycontrolconfig.ListClustersInput,
	_ ...func(*route53recoverycontrolconfig.Options),
) (*route53recoverycontrolconfig.ListClustersOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*route53recoverycontrolconfig.ListClustersOutput), args.Error(1)
}

func (m *mockRoute53RecoveryControlConfigClient) DeleteCluster(
	ctx context.Context, params *route53recoverycontrolconfig.DeleteClusterInput,
	_ ...func(*route53recoverycontrolconfig.Options),
) (*route53recoverycontrolconfig.DeleteClusterOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*route53recoverycontrolconfig.DeleteClusterOutput), args.Error(1)
}

func (m *mockRoute53RecoveryControlConfigClient) ListControlPanels(
	ctx context.Context, params *route53recoverycontrolconfig.ListControlPanelsInput,
	_ ...func(*route53recoverycontrolconfig.Options),
) (*route53recoverycontrolconfig.ListControlPanelsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*route53recoverycontrolconfig.ListControlPanelsOutput), args.Error(1)
}

func (m *mockRoute53RecoveryControlConfigClient) DeleteControlPanel(
	ctx context.Context, params *route53recoverycontrolconfig.DeleteControlPanelInput,
	_ ...func(*route53recoverycontrolconfig.Options),
) (*route53recoverycontrolconfig.DeleteControlPanelOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*route53recoverycontrolconfig.DeleteControlPanelOutput), args.Error(1)
}

func (m *mockRoute53RecoveryControlConfigClient) ListRoutingControls(
	ctx context.Context, params *route53recoverycontrolconfig.ListRoutingControlsInput,
	_ ...func(*route53recoverycontrolconfig.Options),
) (*route53recoverycontrolconfig.ListRoutingControlsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*route53recoverycontrolconfig.ListRoutingControlsOutput), args.Error(1)
}

func (m *mockRoute53RecoveryControlConfigClient) DeleteRoutingControl(
	ctx context.Context, params *route53recoverycontrolconfig.DeleteRoutingControlInput,
	_ ...func(*route53recoverycontrolconfig.Options),
) (*route53recoverycontrolconfig.DeleteRoutingControlOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*route53recoverycontrolconfig.DeleteRoutingControlOutput), args.Error(1)
}

var testRoute53RecoveryControlConfigListerOpts = &nuke.ListerOpts{}
