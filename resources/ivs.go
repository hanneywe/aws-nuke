package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ivs"
)

// IVSClient is the interface for the IVS SDK client methods.
type IVSClient interface {
	ListChannels(ctx context.Context, params *ivs.ListChannelsInput,
		optFns ...func(*ivs.Options)) (*ivs.ListChannelsOutput, error)
	DeleteChannel(ctx context.Context, params *ivs.DeleteChannelInput,
		optFns ...func(*ivs.Options)) (*ivs.DeleteChannelOutput, error)
}
