package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
)

type mockVPCLatticeClient struct {
	mock.Mock
}

func (m *mockVPCLatticeClient) ListServiceNetworks(ctx context.Context, params *vpclattice.ListServiceNetworksInput,
	_ ...func(*vpclattice.Options)) (*vpclattice.ListServiceNetworksOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.ListServiceNetworksOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) ListServices(ctx context.Context, params *vpclattice.ListServicesInput,
	_ ...func(*vpclattice.Options)) (*vpclattice.ListServicesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.ListServicesOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) ListTargetGroups(ctx context.Context, params *vpclattice.ListTargetGroupsInput,
	_ ...func(*vpclattice.Options)) (*vpclattice.ListTargetGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.ListTargetGroupsOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) ListListeners(ctx context.Context, params *vpclattice.ListListenersInput,
	_ ...func(*vpclattice.Options)) (*vpclattice.ListListenersOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.ListListenersOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) ListRules(ctx context.Context, params *vpclattice.ListRulesInput,
	_ ...func(*vpclattice.Options)) (*vpclattice.ListRulesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.ListRulesOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) ListServiceNetworkServiceAssociations(
	ctx context.Context, params *vpclattice.ListServiceNetworkServiceAssociationsInput,
	_ ...func(*vpclattice.Options),
) (*vpclattice.ListServiceNetworkServiceAssociationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.ListServiceNetworkServiceAssociationsOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) ListServiceNetworkVpcAssociations(
	ctx context.Context, params *vpclattice.ListServiceNetworkVpcAssociationsInput,
	_ ...func(*vpclattice.Options),
) (*vpclattice.ListServiceNetworkVpcAssociationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.ListServiceNetworkVpcAssociationsOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) ListResourceConfigurations(ctx context.Context, params *vpclattice.ListResourceConfigurationsInput,
	_ ...func(*vpclattice.Options)) (*vpclattice.ListResourceConfigurationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.ListResourceConfigurationsOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) ListResourceGateways(ctx context.Context, params *vpclattice.ListResourceGatewaysInput,
	_ ...func(*vpclattice.Options)) (*vpclattice.ListResourceGatewaysOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.ListResourceGatewaysOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) ListTagsForResource(ctx context.Context, params *vpclattice.ListTagsForResourceInput,
	_ ...func(*vpclattice.Options)) (*vpclattice.ListTagsForResourceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.ListTagsForResourceOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) GetAuthPolicy(ctx context.Context, params *vpclattice.GetAuthPolicyInput,
	_ ...func(*vpclattice.Options)) (*vpclattice.GetAuthPolicyOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.GetAuthPolicyOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) DeleteServiceNetwork(ctx context.Context, params *vpclattice.DeleteServiceNetworkInput,
	_ ...func(*vpclattice.Options)) (*vpclattice.DeleteServiceNetworkOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.DeleteServiceNetworkOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) DeleteService(ctx context.Context, params *vpclattice.DeleteServiceInput,
	_ ...func(*vpclattice.Options)) (*vpclattice.DeleteServiceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.DeleteServiceOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) DeleteTargetGroup(ctx context.Context, params *vpclattice.DeleteTargetGroupInput,
	_ ...func(*vpclattice.Options)) (*vpclattice.DeleteTargetGroupOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.DeleteTargetGroupOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) DeleteListener(ctx context.Context, params *vpclattice.DeleteListenerInput,
	_ ...func(*vpclattice.Options)) (*vpclattice.DeleteListenerOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.DeleteListenerOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) DeleteRule(ctx context.Context, params *vpclattice.DeleteRuleInput,
	_ ...func(*vpclattice.Options)) (*vpclattice.DeleteRuleOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.DeleteRuleOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) DeleteAuthPolicy(ctx context.Context, params *vpclattice.DeleteAuthPolicyInput,
	_ ...func(*vpclattice.Options)) (*vpclattice.DeleteAuthPolicyOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.DeleteAuthPolicyOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) DeleteServiceNetworkServiceAssociation(
	ctx context.Context, params *vpclattice.DeleteServiceNetworkServiceAssociationInput,
	_ ...func(*vpclattice.Options),
) (*vpclattice.DeleteServiceNetworkServiceAssociationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.DeleteServiceNetworkServiceAssociationOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) DeleteServiceNetworkVpcAssociation(
	ctx context.Context, params *vpclattice.DeleteServiceNetworkVpcAssociationInput,
	_ ...func(*vpclattice.Options),
) (*vpclattice.DeleteServiceNetworkVpcAssociationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.DeleteServiceNetworkVpcAssociationOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) DeleteResourceConfiguration(ctx context.Context, params *vpclattice.DeleteResourceConfigurationInput,
	_ ...func(*vpclattice.Options)) (*vpclattice.DeleteResourceConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.DeleteResourceConfigurationOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) DeleteResourceGateway(ctx context.Context, params *vpclattice.DeleteResourceGatewayInput,
	_ ...func(*vpclattice.Options)) (*vpclattice.DeleteResourceGatewayOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.DeleteResourceGatewayOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) ListTargets(ctx context.Context, params *vpclattice.ListTargetsInput,
	_ ...func(*vpclattice.Options)) (*vpclattice.ListTargetsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.ListTargetsOutput), args.Error(1)
}

func (m *mockVPCLatticeClient) DeregisterTargets(ctx context.Context, params *vpclattice.DeregisterTargetsInput,
	_ ...func(*vpclattice.Options)) (*vpclattice.DeregisterTargetsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*vpclattice.DeregisterTargetsOutput), args.Error(1)
}
