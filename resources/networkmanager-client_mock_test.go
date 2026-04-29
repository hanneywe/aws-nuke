package resources

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/networkmanager"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testNetworkManagerListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

type mockNetworkManagerClient struct {
	mock.Mock
}

func (m *mockNetworkManagerClient) DescribeGlobalNetworks(ctx context.Context, params *networkmanager.DescribeGlobalNetworksInput,
	_ ...func(*networkmanager.Options)) (*networkmanager.DescribeGlobalNetworksOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*networkmanager.DescribeGlobalNetworksOutput), args.Error(1)
}

func (m *mockNetworkManagerClient) GetDevices(ctx context.Context, params *networkmanager.GetDevicesInput,
	_ ...func(*networkmanager.Options)) (*networkmanager.GetDevicesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*networkmanager.GetDevicesOutput), args.Error(1)
}

func (m *mockNetworkManagerClient) DeleteDevice(ctx context.Context, params *networkmanager.DeleteDeviceInput,
	_ ...func(*networkmanager.Options)) (*networkmanager.DeleteDeviceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*networkmanager.DeleteDeviceOutput), args.Error(1)
}

func (m *mockNetworkManagerClient) GetSites(ctx context.Context, params *networkmanager.GetSitesInput,
	_ ...func(*networkmanager.Options)) (*networkmanager.GetSitesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*networkmanager.GetSitesOutput), args.Error(1)
}

func (m *mockNetworkManagerClient) DeleteSite(ctx context.Context, params *networkmanager.DeleteSiteInput,
	_ ...func(*networkmanager.Options)) (*networkmanager.DeleteSiteOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*networkmanager.DeleteSiteOutput), args.Error(1)
}

func (m *mockNetworkManagerClient) GetLinks(ctx context.Context, params *networkmanager.GetLinksInput,
	_ ...func(*networkmanager.Options)) (*networkmanager.GetLinksOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*networkmanager.GetLinksOutput), args.Error(1)
}

func (m *mockNetworkManagerClient) DeleteLink(ctx context.Context, params *networkmanager.DeleteLinkInput,
	_ ...func(*networkmanager.Options)) (*networkmanager.DeleteLinkOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*networkmanager.DeleteLinkOutput), args.Error(1)
}

func (m *mockNetworkManagerClient) GetConnections(ctx context.Context, params *networkmanager.GetConnectionsInput,
	_ ...func(*networkmanager.Options)) (*networkmanager.GetConnectionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*networkmanager.GetConnectionsOutput), args.Error(1)
}

func (m *mockNetworkManagerClient) DeleteConnection(ctx context.Context, params *networkmanager.DeleteConnectionInput,
	_ ...func(*networkmanager.Options)) (*networkmanager.DeleteConnectionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*networkmanager.DeleteConnectionOutput), args.Error(1)
}

func (m *mockNetworkManagerClient) GetTransitGatewayRegistrations(
	ctx context.Context, params *networkmanager.GetTransitGatewayRegistrationsInput,
	_ ...func(*networkmanager.Options)) (*networkmanager.GetTransitGatewayRegistrationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*networkmanager.GetTransitGatewayRegistrationsOutput), args.Error(1)
}

func (m *mockNetworkManagerClient) DeregisterTransitGateway(ctx context.Context, params *networkmanager.DeregisterTransitGatewayInput,
	_ ...func(*networkmanager.Options)) (*networkmanager.DeregisterTransitGatewayOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*networkmanager.DeregisterTransitGatewayOutput), args.Error(1)
}
