package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/kinesisvideo"
)

// KinesisVideoV2Client is the interface for the Kinesis Video SDK v2 client methods.
type KinesisVideoV2Client interface {
	ListSignalingChannels(ctx context.Context, params *kinesisvideo.ListSignalingChannelsInput,
		optFns ...func(*kinesisvideo.Options)) (*kinesisvideo.ListSignalingChannelsOutput, error)
	DeleteSignalingChannel(ctx context.Context, params *kinesisvideo.DeleteSignalingChannelInput,
		optFns ...func(*kinesisvideo.Options)) (*kinesisvideo.DeleteSignalingChannelOutput, error)
}
