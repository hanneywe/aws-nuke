package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ivschat"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockIVSChatClient struct {
	mock.Mock
}

func (m *mockIVSChatClient) ListRooms(ctx context.Context,
	params *ivschat.ListRoomsInput,
	_ ...func(*ivschat.Options)) (*ivschat.ListRoomsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ivschat.ListRoomsOutput), args.Error(1)
}

func (m *mockIVSChatClient) DeleteRoom(ctx context.Context,
	params *ivschat.DeleteRoomInput,
	_ ...func(*ivschat.Options)) (*ivschat.DeleteRoomOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ivschat.DeleteRoomOutput), args.Error(1)
}

func (m *mockIVSChatClient) ListLoggingConfigurations(ctx context.Context,
	params *ivschat.ListLoggingConfigurationsInput,
	_ ...func(*ivschat.Options)) (*ivschat.ListLoggingConfigurationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ivschat.ListLoggingConfigurationsOutput), args.Error(1)
}

func (m *mockIVSChatClient) DeleteLoggingConfiguration(ctx context.Context,
	params *ivschat.DeleteLoggingConfigurationInput,
	_ ...func(*ivschat.Options)) (*ivschat.DeleteLoggingConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ivschat.DeleteLoggingConfigurationOutput), args.Error(1)
}

var testIVSChatListerOpts = &nuke.ListerOpts{}
