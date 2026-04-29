package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ivs"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockIVSClient struct {
	mock.Mock
}

func (m *mockIVSClient) ListChannels(ctx context.Context,
	params *ivs.ListChannelsInput,
	_ ...func(*ivs.Options)) (*ivs.ListChannelsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ivs.ListChannelsOutput), args.Error(1)
}

func (m *mockIVSClient) DeleteChannel(ctx context.Context,
	params *ivs.DeleteChannelInput,
	_ ...func(*ivs.Options)) (*ivs.DeleteChannelOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ivs.DeleteChannelOutput), args.Error(1)
}

var testIVSListerOpts = &nuke.ListerOpts{}
