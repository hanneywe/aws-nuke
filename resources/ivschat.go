package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ivschat"
)

// IVSChatClient is the interface for the IVS Chat SDK client methods.
type IVSChatClient interface {
	ListRooms(ctx context.Context, params *ivschat.ListRoomsInput,
		optFns ...func(*ivschat.Options)) (*ivschat.ListRoomsOutput, error)
	DeleteRoom(ctx context.Context, params *ivschat.DeleteRoomInput,
		optFns ...func(*ivschat.Options)) (*ivschat.DeleteRoomOutput, error)
	ListLoggingConfigurations(ctx context.Context, params *ivschat.ListLoggingConfigurationsInput,
		optFns ...func(*ivschat.Options)) (*ivschat.ListLoggingConfigurationsOutput, error)
	DeleteLoggingConfiguration(ctx context.Context, params *ivschat.DeleteLoggingConfigurationInput,
		optFns ...func(*ivschat.Options)) (*ivschat.DeleteLoggingConfigurationOutput, error)
}
