package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	mediatailortypes "github.com/aws/aws-sdk-go-v2/service/mediatailor/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testMediaTailorV2ListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_MediaTailorChannel_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockMediaTailorV2Client)

	mockClient.
		On("ListChannels", mock.Anything, mock.Anything).
		Return(
			&mediatailor.ListChannelsOutput{
				Items: []mediatailortypes.Channel{
					{
						ChannelName: ptr.String("test-channel"),
					},
				},
			}, nil,
		)

	lister := &MediaTailorChannelLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testMediaTailorV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	channel := resources[0].(*MediaTailorChannel)
	assertions.Equal("test-channel", *channel.ChannelName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaTailorChannel_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockMediaTailorV2Client)

	mockClient.
		On("ListChannels", mock.Anything, mock.Anything).
		Return(
			&mediatailor.ListChannelsOutput{
				Items: []mediatailortypes.Channel{},
			}, nil,
		)

	lister := &MediaTailorChannelLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testMediaTailorV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaTailorChannel_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockMediaTailorV2Client)

	channel := &MediaTailorChannel{
		svc:         mockClient,
		ChannelName: ptr.String("test-channel"),
	}

	mockClient.
		On(
			"DeleteChannel",
			mock.Anything,
			&mediatailor.DeleteChannelInput{
				ChannelName: channel.ChannelName,
			},
		).
		Return(&mediatailor.DeleteChannelOutput{}, nil)

	err := channel.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaTailorChannel_Properties(t *testing.T) {
	assertions := assert.New(t)

	channel := MediaTailorChannel{
		ChannelName: ptr.String("test-channel"),
	}

	properties := channel.Properties()

	assertions.Equal("test-channel", properties.Get("ChannelName"))
}

func Test_Mock_MediaTailorChannel_String(t *testing.T) {
	assertions := assert.New(t)

	channel := MediaTailorChannel{
		ChannelName: ptr.String("test-channel"),
	}

	assertions.Equal("test-channel", channel.String())
}
