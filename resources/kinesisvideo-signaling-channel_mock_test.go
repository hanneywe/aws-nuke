package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/kinesisvideo"
	kinesisvideotypes "github.com/aws/aws-sdk-go-v2/service/kinesisvideo/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testKinesisVideoV2ListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_KinesisVideoSignalingChannel_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockKinesisVideoV2Client)

	mockClient.
		On("ListSignalingChannels", mock.Anything, mock.Anything).
		Return(
			&kinesisvideo.ListSignalingChannelsOutput{
				ChannelInfoList: []kinesisvideotypes.ChannelInfo{
					{
						ChannelARN:  ptr.String("arn:aws:kinesisvideo:us-east-1:123456789012:channel/test-channel/1234567890"),
						ChannelName: ptr.String("test-channel"),
					},
				},
			}, nil,
		)

	lister := &KinesisVideoSignalingChannelLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testKinesisVideoV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	channel := resources[0].(*KinesisVideoSignalingChannel)
	assertions.Equal("arn:aws:kinesisvideo:us-east-1:123456789012:channel/test-channel/1234567890", *channel.ChannelARN)
	assertions.Equal("test-channel", *channel.ChannelName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_KinesisVideoSignalingChannel_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockKinesisVideoV2Client)

	mockClient.
		On("ListSignalingChannels", mock.Anything, mock.Anything).
		Return(
			&kinesisvideo.ListSignalingChannelsOutput{
				ChannelInfoList: []kinesisvideotypes.ChannelInfo{},
			}, nil,
		)

	lister := &KinesisVideoSignalingChannelLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testKinesisVideoV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_KinesisVideoSignalingChannel_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockKinesisVideoV2Client)

	channel := &KinesisVideoSignalingChannel{
		svc:         mockClient,
		ChannelARN:  ptr.String("arn:aws:kinesisvideo:us-east-1:123456789012:channel/test-channel/1234567890"),
		ChannelName: ptr.String("test-channel"),
	}

	mockClient.
		On(
			"DeleteSignalingChannel",
			mock.Anything,
			&kinesisvideo.DeleteSignalingChannelInput{
				ChannelARN: channel.ChannelARN,
			},
		).
		Return(&kinesisvideo.DeleteSignalingChannelOutput{}, nil)

	err := channel.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_KinesisVideoSignalingChannel_Properties(t *testing.T) {
	assertions := assert.New(t)

	channel := KinesisVideoSignalingChannel{
		ChannelARN:  ptr.String("arn:aws:kinesisvideo:us-east-1:123456789012:channel/test-channel/1234567890"),
		ChannelName: ptr.String("test-channel"),
	}

	properties := channel.Properties()

	assertions.Equal("arn:aws:kinesisvideo:us-east-1:123456789012:channel/test-channel/1234567890", properties.Get("ChannelARN"))
	assertions.Equal("test-channel", properties.Get("ChannelName"))
}

func Test_Mock_KinesisVideoSignalingChannel_String(t *testing.T) {
	assertions := assert.New(t)

	channel := KinesisVideoSignalingChannel{
		ChannelARN:  ptr.String("arn:aws:kinesisvideo:us-east-1:123456789012:channel/test-channel/1234567890"),
		ChannelName: ptr.String("test-channel"),
	}

	assertions.Equal("test-channel", channel.String())
}
