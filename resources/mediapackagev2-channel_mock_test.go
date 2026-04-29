package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/mediapackagev2"
	mediapackagev2types "github.com/aws/aws-sdk-go-v2/service/mediapackagev2/types"
)

func Test_Mock_MediaPackageV2Channel_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaPackageV2Client)

	// First call: list channel groups
	mockClient.On("ListChannelGroups", mock.Anything, mock.Anything).
		Return(&mediapackagev2.ListChannelGroupsOutput{
			Items: []mediapackagev2types.ChannelGroupListConfiguration{
				{ChannelGroupName: ptr.String("my-group")},
			},
		}, nil)

	// Second call: list channels for the group
	mockClient.On("ListChannels", mock.Anything, mock.Anything).
		Return(&mediapackagev2.ListChannelsOutput{
			Items: []mediapackagev2types.ChannelListConfiguration{
				{ChannelName: ptr.String("my-channel"), ChannelGroupName: ptr.String("my-group")},
			},
		}, nil)

	lister := &MediaPackageV2ChannelLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaPackageV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	ch := resources[0].(*MediaPackageV2Channel)
	a.Equal("my-channel", *ch.ChannelName)
	a.Equal("my-group", *ch.ChannelGroupName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaPackageV2Channel_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaPackageV2Client)
	mockClient.On("ListChannelGroups", mock.Anything, mock.Anything).
		Return(&mediapackagev2.ListChannelGroupsOutput{Items: []mediapackagev2types.ChannelGroupListConfiguration{}}, nil)
	lister := &MediaPackageV2ChannelLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaPackageV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaPackageV2Channel_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaPackageV2Client)
	ch := &MediaPackageV2Channel{svc: mockClient, ChannelName: ptr.String("my-channel"), ChannelGroupName: ptr.String("my-group")}
	mockClient.On("DeleteChannel", mock.Anything, &mediapackagev2.DeleteChannelInput{
		ChannelGroupName: ch.ChannelGroupName,
		ChannelName:      ch.ChannelName,
	}).Return(&mediapackagev2.DeleteChannelOutput{}, nil)
	a.NoError(ch.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaPackageV2Channel_Properties(t *testing.T) {
	a := assert.New(t)
	ch := MediaPackageV2Channel{ChannelName: ptr.String("my-channel"), ChannelGroupName: ptr.String("my-group")}
	a.Equal("my-channel", ch.Properties().Get("ChannelName"))
	a.Equal("my-group", ch.Properties().Get("ChannelGroupName"))
}

func Test_Mock_MediaPackageV2Channel_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-channel", (&MediaPackageV2Channel{ChannelName: ptr.String("my-channel")}).String())
}
