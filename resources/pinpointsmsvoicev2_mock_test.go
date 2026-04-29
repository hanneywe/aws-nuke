package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2"
)

type mockPinpointSMSVoiceV2Client struct {
	mock.Mock
}

func (m *mockPinpointSMSVoiceV2Client) DescribeConfigurationSets(
	ctx context.Context, params *pinpointsmsvoicev2.DescribeConfigurationSetsInput,
	_ ...func(*pinpointsmsvoicev2.Options),
) (*pinpointsmsvoicev2.DescribeConfigurationSetsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*pinpointsmsvoicev2.DescribeConfigurationSetsOutput), args.Error(1)
}

func (m *mockPinpointSMSVoiceV2Client) DeleteConfigurationSet(
	ctx context.Context, params *pinpointsmsvoicev2.DeleteConfigurationSetInput,
	_ ...func(*pinpointsmsvoicev2.Options),
) (*pinpointsmsvoicev2.DeleteConfigurationSetOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*pinpointsmsvoicev2.DeleteConfigurationSetOutput), args.Error(1)
}

func (m *mockPinpointSMSVoiceV2Client) DescribeOptOutLists(
	ctx context.Context, params *pinpointsmsvoicev2.DescribeOptOutListsInput,
	_ ...func(*pinpointsmsvoicev2.Options),
) (*pinpointsmsvoicev2.DescribeOptOutListsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*pinpointsmsvoicev2.DescribeOptOutListsOutput), args.Error(1)
}

func (m *mockPinpointSMSVoiceV2Client) DeleteOptOutList(
	ctx context.Context, params *pinpointsmsvoicev2.DeleteOptOutListInput,
	_ ...func(*pinpointsmsvoicev2.Options),
) (*pinpointsmsvoicev2.DeleteOptOutListOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*pinpointsmsvoicev2.DeleteOptOutListOutput), args.Error(1)
}

func (m *mockPinpointSMSVoiceV2Client) DescribeProtectConfigurations(
	ctx context.Context, params *pinpointsmsvoicev2.DescribeProtectConfigurationsInput,
	_ ...func(*pinpointsmsvoicev2.Options),
) (*pinpointsmsvoicev2.DescribeProtectConfigurationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*pinpointsmsvoicev2.DescribeProtectConfigurationsOutput), args.Error(1)
}

func (m *mockPinpointSMSVoiceV2Client) DeleteProtectConfiguration(
	ctx context.Context, params *pinpointsmsvoicev2.DeleteProtectConfigurationInput,
	_ ...func(*pinpointsmsvoicev2.Options),
) (*pinpointsmsvoicev2.DeleteProtectConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*pinpointsmsvoicev2.DeleteProtectConfigurationOutput), args.Error(1)
}

func (m *mockPinpointSMSVoiceV2Client) UpdateProtectConfiguration(
	ctx context.Context, params *pinpointsmsvoicev2.UpdateProtectConfigurationInput,
	_ ...func(*pinpointsmsvoicev2.Options),
) (*pinpointsmsvoicev2.UpdateProtectConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*pinpointsmsvoicev2.UpdateProtectConfigurationOutput), args.Error(1)
}

func (m *mockPinpointSMSVoiceV2Client) DeleteEventDestination(
	ctx context.Context, params *pinpointsmsvoicev2.DeleteEventDestinationInput,
	_ ...func(*pinpointsmsvoicev2.Options),
) (*pinpointsmsvoicev2.DeleteEventDestinationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*pinpointsmsvoicev2.DeleteEventDestinationOutput), args.Error(1)
}

func (m *mockPinpointSMSVoiceV2Client) DescribeOptedOutNumbers(
	ctx context.Context, params *pinpointsmsvoicev2.DescribeOptedOutNumbersInput,
	_ ...func(*pinpointsmsvoicev2.Options),
) (*pinpointsmsvoicev2.DescribeOptedOutNumbersOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*pinpointsmsvoicev2.DescribeOptedOutNumbersOutput), args.Error(1)
}

func (m *mockPinpointSMSVoiceV2Client) DeleteOptedOutNumber(
	ctx context.Context, params *pinpointsmsvoicev2.DeleteOptedOutNumberInput,
	_ ...func(*pinpointsmsvoicev2.Options),
) (*pinpointsmsvoicev2.DeleteOptedOutNumberOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*pinpointsmsvoicev2.DeleteOptedOutNumberOutput), args.Error(1)
}

func (m *mockPinpointSMSVoiceV2Client) ListProtectConfigurationRuleSetNumberOverrides(
	ctx context.Context, params *pinpointsmsvoicev2.ListProtectConfigurationRuleSetNumberOverridesInput,
	_ ...func(*pinpointsmsvoicev2.Options),
) (*pinpointsmsvoicev2.ListProtectConfigurationRuleSetNumberOverridesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*pinpointsmsvoicev2.ListProtectConfigurationRuleSetNumberOverridesOutput), args.Error(1)
}

func (m *mockPinpointSMSVoiceV2Client) DeleteProtectConfigurationRuleSetNumberOverride(
	ctx context.Context, params *pinpointsmsvoicev2.DeleteProtectConfigurationRuleSetNumberOverrideInput,
	_ ...func(*pinpointsmsvoicev2.Options),
) (*pinpointsmsvoicev2.DeleteProtectConfigurationRuleSetNumberOverrideOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*pinpointsmsvoicev2.DeleteProtectConfigurationRuleSetNumberOverrideOutput), args.Error(1)
}
