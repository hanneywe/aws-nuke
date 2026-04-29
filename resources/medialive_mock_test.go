package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/medialive"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockMediaLiveClient struct {
	mock.Mock
}

func (m *mockMediaLiveClient) ListNetworks(ctx context.Context,
	params *medialive.ListNetworksInput,
	_ ...func(*medialive.Options)) (*medialive.ListNetworksOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*medialive.ListNetworksOutput), args.Error(1)
}

func (m *mockMediaLiveClient) DeleteNetwork(ctx context.Context,
	params *medialive.DeleteNetworkInput,
	_ ...func(*medialive.Options)) (*medialive.DeleteNetworkOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*medialive.DeleteNetworkOutput), args.Error(1)
}

func (m *mockMediaLiveClient) ListSdiSources(ctx context.Context,
	params *medialive.ListSdiSourcesInput,
	_ ...func(*medialive.Options)) (*medialive.ListSdiSourcesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*medialive.ListSdiSourcesOutput), args.Error(1)
}

func (m *mockMediaLiveClient) DeleteSdiSource(ctx context.Context,
	params *medialive.DeleteSdiSourceInput,
	_ ...func(*medialive.Options)) (*medialive.DeleteSdiSourceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*medialive.DeleteSdiSourceOutput), args.Error(1)
}

func (m *mockMediaLiveClient) ListCloudWatchAlarmTemplateGroups(ctx context.Context,
	params *medialive.ListCloudWatchAlarmTemplateGroupsInput,
	_ ...func(*medialive.Options)) (*medialive.ListCloudWatchAlarmTemplateGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*medialive.ListCloudWatchAlarmTemplateGroupsOutput), args.Error(1)
}

func (m *mockMediaLiveClient) DeleteCloudWatchAlarmTemplateGroup(ctx context.Context,
	params *medialive.DeleteCloudWatchAlarmTemplateGroupInput,
	_ ...func(*medialive.Options)) (*medialive.DeleteCloudWatchAlarmTemplateGroupOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*medialive.DeleteCloudWatchAlarmTemplateGroupOutput), args.Error(1)
}

func (m *mockMediaLiveClient) ListEventBridgeRuleTemplates(ctx context.Context,
	params *medialive.ListEventBridgeRuleTemplatesInput,
	_ ...func(*medialive.Options)) (*medialive.ListEventBridgeRuleTemplatesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*medialive.ListEventBridgeRuleTemplatesOutput), args.Error(1)
}

func (m *mockMediaLiveClient) DeleteEventBridgeRuleTemplate(ctx context.Context,
	params *medialive.DeleteEventBridgeRuleTemplateInput,
	_ ...func(*medialive.Options)) (*medialive.DeleteEventBridgeRuleTemplateOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*medialive.DeleteEventBridgeRuleTemplateOutput), args.Error(1)
}

func (m *mockMediaLiveClient) ListCloudWatchAlarmTemplates(ctx context.Context,
	params *medialive.ListCloudWatchAlarmTemplatesInput,
	_ ...func(*medialive.Options)) (*medialive.ListCloudWatchAlarmTemplatesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*medialive.ListCloudWatchAlarmTemplatesOutput), args.Error(1)
}

func (m *mockMediaLiveClient) DeleteCloudWatchAlarmTemplate(ctx context.Context,
	params *medialive.DeleteCloudWatchAlarmTemplateInput,
	_ ...func(*medialive.Options)) (*medialive.DeleteCloudWatchAlarmTemplateOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*medialive.DeleteCloudWatchAlarmTemplateOutput), args.Error(1)
}

func (m *mockMediaLiveClient) ListEventBridgeRuleTemplateGroups(ctx context.Context,
	params *medialive.ListEventBridgeRuleTemplateGroupsInput,
	_ ...func(*medialive.Options)) (*medialive.ListEventBridgeRuleTemplateGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*medialive.ListEventBridgeRuleTemplateGroupsOutput), args.Error(1)
}

func (m *mockMediaLiveClient) DeleteEventBridgeRuleTemplateGroup(ctx context.Context,
	params *medialive.DeleteEventBridgeRuleTemplateGroupInput,
	_ ...func(*medialive.Options)) (*medialive.DeleteEventBridgeRuleTemplateGroupOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*medialive.DeleteEventBridgeRuleTemplateGroupOutput), args.Error(1)
}

var testMediaLiveListerOpts = &nuke.ListerOpts{}
