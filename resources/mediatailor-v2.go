package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
)

// MediaTailorV2Client is the interface for the MediaTailor SDK v2 client methods.
type MediaTailorV2Client interface {
	ListChannels(ctx context.Context, params *mediatailor.ListChannelsInput,
		optFns ...func(*mediatailor.Options)) (*mediatailor.ListChannelsOutput, error)
	DeleteChannel(ctx context.Context, params *mediatailor.DeleteChannelInput,
		optFns ...func(*mediatailor.Options)) (*mediatailor.DeleteChannelOutput, error)
	ListSourceLocations(ctx context.Context, params *mediatailor.ListSourceLocationsInput,
		optFns ...func(*mediatailor.Options)) (*mediatailor.ListSourceLocationsOutput, error)
	DeleteSourceLocation(ctx context.Context, params *mediatailor.DeleteSourceLocationInput,
		optFns ...func(*mediatailor.Options)) (*mediatailor.DeleteSourceLocationOutput, error)
	ListLiveSources(ctx context.Context, params *mediatailor.ListLiveSourcesInput,
		optFns ...func(*mediatailor.Options)) (*mediatailor.ListLiveSourcesOutput, error)
	DeleteLiveSource(ctx context.Context, params *mediatailor.DeleteLiveSourceInput,
		optFns ...func(*mediatailor.Options)) (*mediatailor.DeleteLiveSourceOutput, error)
	ListVodSources(ctx context.Context, params *mediatailor.ListVodSourcesInput,
		optFns ...func(*mediatailor.Options)) (*mediatailor.ListVodSourcesOutput, error)
	DeleteVodSource(ctx context.Context, params *mediatailor.DeleteVodSourceInput,
		optFns ...func(*mediatailor.Options)) (*mediatailor.DeleteVodSourceOutput, error)
}
