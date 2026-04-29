package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/iotwireless"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockIoTWirelessClient struct {
	mock.Mock
}

func (m *mockIoTWirelessClient) ListMulticastGroups(ctx context.Context,
	params *iotwireless.ListMulticastGroupsInput,
	_ ...func(*iotwireless.Options)) (*iotwireless.ListMulticastGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iotwireless.ListMulticastGroupsOutput), args.Error(1)
}

func (m *mockIoTWirelessClient) DeleteMulticastGroup(ctx context.Context,
	params *iotwireless.DeleteMulticastGroupInput,
	_ ...func(*iotwireless.Options)) (*iotwireless.DeleteMulticastGroupOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iotwireless.DeleteMulticastGroupOutput), args.Error(1)
}

func (m *mockIoTWirelessClient) ListNetworkAnalyzerConfigurations(ctx context.Context,
	params *iotwireless.ListNetworkAnalyzerConfigurationsInput,
	_ ...func(*iotwireless.Options)) (*iotwireless.ListNetworkAnalyzerConfigurationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iotwireless.ListNetworkAnalyzerConfigurationsOutput), args.Error(1)
}

func (m *mockIoTWirelessClient) DeleteNetworkAnalyzerConfiguration(ctx context.Context,
	params *iotwireless.DeleteNetworkAnalyzerConfigurationInput,
	_ ...func(*iotwireless.Options)) (*iotwireless.DeleteNetworkAnalyzerConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iotwireless.DeleteNetworkAnalyzerConfigurationOutput), args.Error(1)
}

func (m *mockIoTWirelessClient) ListServiceProfiles(ctx context.Context,
	params *iotwireless.ListServiceProfilesInput,
	_ ...func(*iotwireless.Options)) (*iotwireless.ListServiceProfilesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iotwireless.ListServiceProfilesOutput), args.Error(1)
}

func (m *mockIoTWirelessClient) DeleteServiceProfile(ctx context.Context,
	params *iotwireless.DeleteServiceProfileInput,
	_ ...func(*iotwireless.Options)) (*iotwireless.DeleteServiceProfileOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iotwireless.DeleteServiceProfileOutput), args.Error(1)
}

func (m *mockIoTWirelessClient) ListDeviceProfiles(ctx context.Context,
	params *iotwireless.ListDeviceProfilesInput,
	_ ...func(*iotwireless.Options)) (*iotwireless.ListDeviceProfilesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iotwireless.ListDeviceProfilesOutput), args.Error(1)
}

func (m *mockIoTWirelessClient) DeleteDeviceProfile(ctx context.Context,
	params *iotwireless.DeleteDeviceProfileInput,
	_ ...func(*iotwireless.Options)) (*iotwireless.DeleteDeviceProfileOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iotwireless.DeleteDeviceProfileOutput), args.Error(1)
}

func (m *mockIoTWirelessClient) ListFuotaTasks(ctx context.Context,
	params *iotwireless.ListFuotaTasksInput,
	_ ...func(*iotwireless.Options)) (*iotwireless.ListFuotaTasksOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iotwireless.ListFuotaTasksOutput), args.Error(1)
}

func (m *mockIoTWirelessClient) DeleteFuotaTask(ctx context.Context,
	params *iotwireless.DeleteFuotaTaskInput,
	_ ...func(*iotwireless.Options)) (*iotwireless.DeleteFuotaTaskOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iotwireless.DeleteFuotaTaskOutput), args.Error(1)
}

func (m *mockIoTWirelessClient) ListDestinations(ctx context.Context,
	params *iotwireless.ListDestinationsInput,
	_ ...func(*iotwireless.Options)) (*iotwireless.ListDestinationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iotwireless.ListDestinationsOutput), args.Error(1)
}

func (m *mockIoTWirelessClient) DeleteDestination(ctx context.Context,
	params *iotwireless.DeleteDestinationInput,
	_ ...func(*iotwireless.Options)) (*iotwireless.DeleteDestinationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iotwireless.DeleteDestinationOutput), args.Error(1)
}

func (m *mockIoTWirelessClient) ListWirelessGateways(ctx context.Context,
	params *iotwireless.ListWirelessGatewaysInput,
	_ ...func(*iotwireless.Options)) (*iotwireless.ListWirelessGatewaysOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iotwireless.ListWirelessGatewaysOutput), args.Error(1)
}

func (m *mockIoTWirelessClient) DeleteWirelessGateway(ctx context.Context,
	params *iotwireless.DeleteWirelessGatewayInput,
	_ ...func(*iotwireless.Options)) (*iotwireless.DeleteWirelessGatewayOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iotwireless.DeleteWirelessGatewayOutput), args.Error(1)
}

var testIoTWirelessListerOpts = &nuke.ListerOpts{}
