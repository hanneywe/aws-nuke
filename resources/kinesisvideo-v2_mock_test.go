package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/kinesisvideo"
)

type mockKinesisVideoV2Client struct {
	mock.Mock
}

func (m *mockKinesisVideoV2Client) ListSignalingChannels(
	ctx context.Context, params *kinesisvideo.ListSignalingChannelsInput,
	_ ...func(*kinesisvideo.Options),
) (*kinesisvideo.ListSignalingChannelsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*kinesisvideo.ListSignalingChannelsOutput), args.Error(1)
}

func (m *mockKinesisVideoV2Client) DeleteSignalingChannel(
	ctx context.Context, params *kinesisvideo.DeleteSignalingChannelInput,
	_ ...func(*kinesisvideo.Options),
) (*kinesisvideo.DeleteSignalingChannelOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*kinesisvideo.DeleteSignalingChannelOutput), args.Error(1)
}
