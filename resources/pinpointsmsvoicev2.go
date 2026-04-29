package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2"
)

// PinpointSMSVoiceV2Client is the interface for the Pinpoint SMS Voice V2 SDK client methods.
type PinpointSMSVoiceV2Client interface {
	DescribeConfigurationSets(ctx context.Context, params *pinpointsmsvoicev2.DescribeConfigurationSetsInput,
		optFns ...func(*pinpointsmsvoicev2.Options)) (*pinpointsmsvoicev2.DescribeConfigurationSetsOutput, error)
	DeleteConfigurationSet(ctx context.Context, params *pinpointsmsvoicev2.DeleteConfigurationSetInput,
		optFns ...func(*pinpointsmsvoicev2.Options)) (*pinpointsmsvoicev2.DeleteConfigurationSetOutput, error)
	DescribeOptOutLists(ctx context.Context, params *pinpointsmsvoicev2.DescribeOptOutListsInput,
		optFns ...func(*pinpointsmsvoicev2.Options)) (*pinpointsmsvoicev2.DescribeOptOutListsOutput, error)
	DeleteOptOutList(ctx context.Context, params *pinpointsmsvoicev2.DeleteOptOutListInput,
		optFns ...func(*pinpointsmsvoicev2.Options)) (*pinpointsmsvoicev2.DeleteOptOutListOutput, error)
	DescribeProtectConfigurations(ctx context.Context, params *pinpointsmsvoicev2.DescribeProtectConfigurationsInput,
		optFns ...func(*pinpointsmsvoicev2.Options)) (*pinpointsmsvoicev2.DescribeProtectConfigurationsOutput, error)
	DeleteProtectConfiguration(ctx context.Context, params *pinpointsmsvoicev2.DeleteProtectConfigurationInput,
		optFns ...func(*pinpointsmsvoicev2.Options)) (*pinpointsmsvoicev2.DeleteProtectConfigurationOutput, error)
	UpdateProtectConfiguration(ctx context.Context, params *pinpointsmsvoicev2.UpdateProtectConfigurationInput,
		optFns ...func(*pinpointsmsvoicev2.Options)) (*pinpointsmsvoicev2.UpdateProtectConfigurationOutput, error)
	DeleteEventDestination(ctx context.Context, params *pinpointsmsvoicev2.DeleteEventDestinationInput,
		optFns ...func(*pinpointsmsvoicev2.Options)) (*pinpointsmsvoicev2.DeleteEventDestinationOutput, error)
	DescribeOptedOutNumbers(ctx context.Context, params *pinpointsmsvoicev2.DescribeOptedOutNumbersInput,
		optFns ...func(*pinpointsmsvoicev2.Options)) (*pinpointsmsvoicev2.DescribeOptedOutNumbersOutput, error)
	DeleteOptedOutNumber(ctx context.Context, params *pinpointsmsvoicev2.DeleteOptedOutNumberInput,
		optFns ...func(*pinpointsmsvoicev2.Options)) (*pinpointsmsvoicev2.DeleteOptedOutNumberOutput, error)
	ListProtectConfigurationRuleSetNumberOverrides(ctx context.Context,
		params *pinpointsmsvoicev2.ListProtectConfigurationRuleSetNumberOverridesInput,
		optFns ...func(*pinpointsmsvoicev2.Options),
	) (*pinpointsmsvoicev2.ListProtectConfigurationRuleSetNumberOverridesOutput, error)
	DeleteProtectConfigurationRuleSetNumberOverride(ctx context.Context,
		params *pinpointsmsvoicev2.DeleteProtectConfigurationRuleSetNumberOverrideInput,
		optFns ...func(*pinpointsmsvoicev2.Options),
	) (*pinpointsmsvoicev2.DeleteProtectConfigurationRuleSetNumberOverrideOutput, error)
}
