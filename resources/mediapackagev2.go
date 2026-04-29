package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/mediapackagev2"
)

// MediaPackageV2Client is the interface for the MediaPackage V2 SDK client methods.
type MediaPackageV2Client interface {
	ListChannelGroups(ctx context.Context, params *mediapackagev2.ListChannelGroupsInput,
		optFns ...func(*mediapackagev2.Options)) (*mediapackagev2.ListChannelGroupsOutput, error)
	DeleteChannelGroup(ctx context.Context, params *mediapackagev2.DeleteChannelGroupInput,
		optFns ...func(*mediapackagev2.Options)) (*mediapackagev2.DeleteChannelGroupOutput, error)
	ListChannels(ctx context.Context, params *mediapackagev2.ListChannelsInput,
		optFns ...func(*mediapackagev2.Options)) (*mediapackagev2.ListChannelsOutput, error)
	DeleteChannel(ctx context.Context, params *mediapackagev2.DeleteChannelInput,
		optFns ...func(*mediapackagev2.Options)) (*mediapackagev2.DeleteChannelOutput, error)
	ListOriginEndpoints(ctx context.Context, params *mediapackagev2.ListOriginEndpointsInput,
		optFns ...func(*mediapackagev2.Options)) (*mediapackagev2.ListOriginEndpointsOutput, error)
	DeleteOriginEndpoint(ctx context.Context, params *mediapackagev2.DeleteOriginEndpointInput,
		optFns ...func(*mediapackagev2.Options)) (*mediapackagev2.DeleteOriginEndpointOutput, error)
}
