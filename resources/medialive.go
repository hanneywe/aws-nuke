package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/medialive"
)

// MediaLiveClient is the interface for the MediaLive SDK v2 client methods used by new MediaLive resources.
// Existing MediaLive resources use SDK v1; this interface is for new SDK v2 resources only.
type MediaLiveClient interface {
	ListNetworks(ctx context.Context, params *medialive.ListNetworksInput,
		optFns ...func(*medialive.Options)) (*medialive.ListNetworksOutput, error)
	DeleteNetwork(ctx context.Context, params *medialive.DeleteNetworkInput,
		optFns ...func(*medialive.Options)) (*medialive.DeleteNetworkOutput, error)
	ListSdiSources(ctx context.Context, params *medialive.ListSdiSourcesInput,
		optFns ...func(*medialive.Options)) (*medialive.ListSdiSourcesOutput, error)
	DeleteSdiSource(ctx context.Context, params *medialive.DeleteSdiSourceInput,
		optFns ...func(*medialive.Options)) (*medialive.DeleteSdiSourceOutput, error)
	ListCloudWatchAlarmTemplateGroups(ctx context.Context, params *medialive.ListCloudWatchAlarmTemplateGroupsInput,
		optFns ...func(*medialive.Options)) (*medialive.ListCloudWatchAlarmTemplateGroupsOutput, error)
	DeleteCloudWatchAlarmTemplateGroup(ctx context.Context, params *medialive.DeleteCloudWatchAlarmTemplateGroupInput,
		optFns ...func(*medialive.Options)) (*medialive.DeleteCloudWatchAlarmTemplateGroupOutput, error)
	ListCloudWatchAlarmTemplates(ctx context.Context, params *medialive.ListCloudWatchAlarmTemplatesInput,
		optFns ...func(*medialive.Options)) (*medialive.ListCloudWatchAlarmTemplatesOutput, error)
	DeleteCloudWatchAlarmTemplate(ctx context.Context, params *medialive.DeleteCloudWatchAlarmTemplateInput,
		optFns ...func(*medialive.Options)) (*medialive.DeleteCloudWatchAlarmTemplateOutput, error)
	ListEventBridgeRuleTemplates(ctx context.Context, params *medialive.ListEventBridgeRuleTemplatesInput,
		optFns ...func(*medialive.Options)) (*medialive.ListEventBridgeRuleTemplatesOutput, error)
	DeleteEventBridgeRuleTemplate(ctx context.Context, params *medialive.DeleteEventBridgeRuleTemplateInput,
		optFns ...func(*medialive.Options)) (*medialive.DeleteEventBridgeRuleTemplateOutput, error)
	ListEventBridgeRuleTemplateGroups(ctx context.Context, params *medialive.ListEventBridgeRuleTemplateGroupsInput,
		optFns ...func(*medialive.Options)) (*medialive.ListEventBridgeRuleTemplateGroupsOutput, error)
	DeleteEventBridgeRuleTemplateGroup(ctx context.Context, params *medialive.DeleteEventBridgeRuleTemplateGroupInput,
		optFns ...func(*medialive.Options)) (*medialive.DeleteEventBridgeRuleTemplateGroupOutput, error)
}
