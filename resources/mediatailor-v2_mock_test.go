package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
)

type mockMediaTailorV2Client struct {
	mock.Mock
}

func (m *mockMediaTailorV2Client) ListChannels(
	ctx context.Context, params *mediatailor.ListChannelsInput,
	_ ...func(*mediatailor.Options),
) (*mediatailor.ListChannelsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mediatailor.ListChannelsOutput), args.Error(1)
}

func (m *mockMediaTailorV2Client) DeleteChannel(
	ctx context.Context, params *mediatailor.DeleteChannelInput,
	_ ...func(*mediatailor.Options),
) (*mediatailor.DeleteChannelOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mediatailor.DeleteChannelOutput), args.Error(1)
}

func (m *mockMediaTailorV2Client) ListSourceLocations(
	ctx context.Context, params *mediatailor.ListSourceLocationsInput,
	_ ...func(*mediatailor.Options),
) (*mediatailor.ListSourceLocationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mediatailor.ListSourceLocationsOutput), args.Error(1)
}

func (m *mockMediaTailorV2Client) DeleteSourceLocation(
	ctx context.Context, params *mediatailor.DeleteSourceLocationInput,
	_ ...func(*mediatailor.Options),
) (*mediatailor.DeleteSourceLocationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mediatailor.DeleteSourceLocationOutput), args.Error(1)
}

func (m *mockMediaTailorV2Client) ListLiveSources(
	ctx context.Context, params *mediatailor.ListLiveSourcesInput,
	_ ...func(*mediatailor.Options),
) (*mediatailor.ListLiveSourcesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mediatailor.ListLiveSourcesOutput), args.Error(1)
}

func (m *mockMediaTailorV2Client) DeleteLiveSource(
	ctx context.Context, params *mediatailor.DeleteLiveSourceInput,
	_ ...func(*mediatailor.Options),
) (*mediatailor.DeleteLiveSourceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mediatailor.DeleteLiveSourceOutput), args.Error(1)
}

func (m *mockMediaTailorV2Client) ListVodSources(
	ctx context.Context, params *mediatailor.ListVodSourcesInput,
	_ ...func(*mediatailor.Options),
) (*mediatailor.ListVodSourcesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mediatailor.ListVodSourcesOutput), args.Error(1)
}

func (m *mockMediaTailorV2Client) DeleteVodSource(
	ctx context.Context, params *mediatailor.DeleteVodSourceInput,
	_ ...func(*mediatailor.Options),
) (*mediatailor.DeleteVodSourceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mediatailor.DeleteVodSourceOutput), args.Error(1)
}
