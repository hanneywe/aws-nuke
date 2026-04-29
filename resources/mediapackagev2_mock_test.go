package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/mediapackagev2"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockMediaPackageV2Client struct {
	mock.Mock
}

func (m *mockMediaPackageV2Client) ListChannelGroups(
	ctx context.Context, params *mediapackagev2.ListChannelGroupsInput,
	_ ...func(*mediapackagev2.Options),
) (*mediapackagev2.ListChannelGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mediapackagev2.ListChannelGroupsOutput), args.Error(1)
}

func (m *mockMediaPackageV2Client) DeleteChannelGroup(
	ctx context.Context, params *mediapackagev2.DeleteChannelGroupInput,
	_ ...func(*mediapackagev2.Options),
) (*mediapackagev2.DeleteChannelGroupOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mediapackagev2.DeleteChannelGroupOutput), args.Error(1)
}

func (m *mockMediaPackageV2Client) ListChannels(
	ctx context.Context, params *mediapackagev2.ListChannelsInput,
	_ ...func(*mediapackagev2.Options),
) (*mediapackagev2.ListChannelsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mediapackagev2.ListChannelsOutput), args.Error(1)
}

func (m *mockMediaPackageV2Client) DeleteChannel(
	ctx context.Context, params *mediapackagev2.DeleteChannelInput,
	_ ...func(*mediapackagev2.Options),
) (*mediapackagev2.DeleteChannelOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mediapackagev2.DeleteChannelOutput), args.Error(1)
}

func (m *mockMediaPackageV2Client) ListOriginEndpoints(
	ctx context.Context, params *mediapackagev2.ListOriginEndpointsInput,
	_ ...func(*mediapackagev2.Options),
) (*mediapackagev2.ListOriginEndpointsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mediapackagev2.ListOriginEndpointsOutput), args.Error(1)
}

func (m *mockMediaPackageV2Client) DeleteOriginEndpoint(
	ctx context.Context, params *mediapackagev2.DeleteOriginEndpointInput,
	_ ...func(*mediapackagev2.Options),
) (*mediapackagev2.DeleteOriginEndpointOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mediapackagev2.DeleteOriginEndpointOutput), args.Error(1)
}

var testMediaPackageV2ListerOpts = &nuke.ListerOpts{}
