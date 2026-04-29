package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iotwireless"
)

// IoTWirelessClient is the interface for the IoT Wireless SDK client methods.
type IoTWirelessClient interface {
	ListMulticastGroups(ctx context.Context, params *iotwireless.ListMulticastGroupsInput,
		optFns ...func(*iotwireless.Options)) (*iotwireless.ListMulticastGroupsOutput, error)
	DeleteMulticastGroup(ctx context.Context, params *iotwireless.DeleteMulticastGroupInput,
		optFns ...func(*iotwireless.Options)) (*iotwireless.DeleteMulticastGroupOutput, error)
	ListNetworkAnalyzerConfigurations(ctx context.Context, params *iotwireless.ListNetworkAnalyzerConfigurationsInput,
		optFns ...func(*iotwireless.Options)) (*iotwireless.ListNetworkAnalyzerConfigurationsOutput, error)
	DeleteNetworkAnalyzerConfiguration(ctx context.Context, params *iotwireless.DeleteNetworkAnalyzerConfigurationInput,
		optFns ...func(*iotwireless.Options)) (*iotwireless.DeleteNetworkAnalyzerConfigurationOutput, error)
	ListServiceProfiles(ctx context.Context, params *iotwireless.ListServiceProfilesInput,
		optFns ...func(*iotwireless.Options)) (*iotwireless.ListServiceProfilesOutput, error)
	DeleteServiceProfile(ctx context.Context, params *iotwireless.DeleteServiceProfileInput,
		optFns ...func(*iotwireless.Options)) (*iotwireless.DeleteServiceProfileOutput, error)
	ListDeviceProfiles(ctx context.Context, params *iotwireless.ListDeviceProfilesInput,
		optFns ...func(*iotwireless.Options)) (*iotwireless.ListDeviceProfilesOutput, error)
	DeleteDeviceProfile(ctx context.Context, params *iotwireless.DeleteDeviceProfileInput,
		optFns ...func(*iotwireless.Options)) (*iotwireless.DeleteDeviceProfileOutput, error)
	ListFuotaTasks(ctx context.Context, params *iotwireless.ListFuotaTasksInput,
		optFns ...func(*iotwireless.Options)) (*iotwireless.ListFuotaTasksOutput, error)
	DeleteFuotaTask(ctx context.Context, params *iotwireless.DeleteFuotaTaskInput,
		optFns ...func(*iotwireless.Options)) (*iotwireless.DeleteFuotaTaskOutput, error)
	ListDestinations(ctx context.Context, params *iotwireless.ListDestinationsInput,
		optFns ...func(*iotwireless.Options)) (*iotwireless.ListDestinationsOutput, error)
	DeleteDestination(ctx context.Context, params *iotwireless.DeleteDestinationInput,
		optFns ...func(*iotwireless.Options)) (*iotwireless.DeleteDestinationOutput, error)
	ListWirelessGateways(ctx context.Context, params *iotwireless.ListWirelessGatewaysInput,
		optFns ...func(*iotwireless.Options)) (*iotwireless.ListWirelessGatewaysOutput, error)
	DeleteWirelessGateway(ctx context.Context, params *iotwireless.DeleteWirelessGatewayInput,
		optFns ...func(*iotwireless.Options)) (*iotwireless.DeleteWirelessGatewayOutput, error)
}
