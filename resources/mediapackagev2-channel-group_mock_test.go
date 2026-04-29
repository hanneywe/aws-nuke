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

func Test_Mock_MediaPackageV2ChannelGroup_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaPackageV2Client)
	mockClient.On("ListChannelGroups", mock.Anything, mock.Anything).
		Return(&mediapackagev2.ListChannelGroupsOutput{
			Items: []mediapackagev2types.ChannelGroupListConfiguration{
				{ChannelGroupName: ptr.String("my-group"), Arn: ptr.String("arn:aws:mediapackagev2:us-east-1:123456789012:channelGroup/my-group")},
			},
		}, nil)
	lister := &MediaPackageV2ChannelGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaPackageV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	cg := resources[0].(*MediaPackageV2ChannelGroup)
	a.Equal("my-group", *cg.ChannelGroupName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaPackageV2ChannelGroup_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaPackageV2Client)
	mockClient.On("ListChannelGroups", mock.Anything, mock.Anything).
		Return(&mediapackagev2.ListChannelGroupsOutput{Items: []mediapackagev2types.ChannelGroupListConfiguration{}}, nil)
	lister := &MediaPackageV2ChannelGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaPackageV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaPackageV2ChannelGroup_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaPackageV2Client)
	cg := &MediaPackageV2ChannelGroup{svc: mockClient, ChannelGroupName: ptr.String("my-group")}
	mockClient.On("DeleteChannelGroup", mock.Anything, &mediapackagev2.DeleteChannelGroupInput{ChannelGroupName: cg.ChannelGroupName}).
		Return(&mediapackagev2.DeleteChannelGroupOutput{}, nil)
	a.NoError(cg.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaPackageV2ChannelGroup_Properties(t *testing.T) {
	a := assert.New(t)
	cg := MediaPackageV2ChannelGroup{
		ChannelGroupName: ptr.String("my-group"),
		ARN:              ptr.String("arn:aws:mediapackagev2:us-east-1:123456789012:channelGroup/my-group"),
	}
	a.Equal("my-group", cg.Properties().Get("ChannelGroupName"))
}

func Test_Mock_MediaPackageV2ChannelGroup_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-group", (&MediaPackageV2ChannelGroup{ChannelGroupName: ptr.String("my-group")}).String())
}
