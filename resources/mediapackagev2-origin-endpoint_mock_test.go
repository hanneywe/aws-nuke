package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/mediapackagev2"
	mptypes "github.com/aws/aws-sdk-go-v2/service/mediapackagev2/types"
)

func Test_Mock_MediaPackageV2OriginEndpoint_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaPackageV2Client)

	mockClient.On("ListChannelGroups", mock.Anything, mock.Anything).
		Return(&mediapackagev2.ListChannelGroupsOutput{
			Items: []mptypes.ChannelGroupListConfiguration{
				{ChannelGroupName: ptr.String("cg-1")},
			},
		}, nil)

	mockClient.On("ListChannels", mock.Anything, mock.Anything).
		Return(&mediapackagev2.ListChannelsOutput{
			Items: []mptypes.ChannelListConfiguration{
				{ChannelGroupName: ptr.String("cg-1"), ChannelName: ptr.String("ch-1")},
			},
		}, nil)

	mockClient.On("ListOriginEndpoints", mock.Anything, mock.Anything).
		Return(&mediapackagev2.ListOriginEndpointsOutput{
			Items: []mptypes.OriginEndpointListConfiguration{
				{
					ChannelGroupName:   ptr.String("cg-1"),
					ChannelName:        ptr.String("ch-1"),
					OriginEndpointName: ptr.String("ep-1"),
				},
			},
		}, nil)

	lister := &MediaPackageV2OriginEndpointLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaPackageV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*MediaPackageV2OriginEndpoint)
	a.Equal("ep-1", *r.OriginEndpointName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaPackageV2OriginEndpoint_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaPackageV2Client)

	mockClient.On("ListChannelGroups", mock.Anything, mock.Anything).
		Return(&mediapackagev2.ListChannelGroupsOutput{
			Items: []mptypes.ChannelGroupListConfiguration{},
		}, nil)

	lister := &MediaPackageV2OriginEndpointLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaPackageV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaPackageV2OriginEndpoint_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaPackageV2Client)

	r := &MediaPackageV2OriginEndpoint{
		svc:                mockClient,
		ChannelGroupName:   ptr.String("cg-1"),
		ChannelName:        ptr.String("ch-1"),
		OriginEndpointName: ptr.String("ep-1"),
	}

	mockClient.On("DeleteOriginEndpoint", mock.Anything,
		&mediapackagev2.DeleteOriginEndpointInput{
			ChannelGroupName:   r.ChannelGroupName,
			ChannelName:        r.ChannelName,
			OriginEndpointName: r.OriginEndpointName,
		}).Return(&mediapackagev2.DeleteOriginEndpointOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaPackageV2OriginEndpoint_Properties(t *testing.T) {
	a := assert.New(t)
	r := &MediaPackageV2OriginEndpoint{
		ChannelGroupName:   ptr.String("cg-1"),
		ChannelName:        ptr.String("ch-1"),
		OriginEndpointName: ptr.String("ep-1"),
	}
	props := r.Properties()
	a.Equal("cg-1", props.Get("ChannelGroupName"))
	a.Equal("ch-1", props.Get("ChannelName"))
	a.Equal("ep-1", props.Get("OriginEndpointName"))
}

func Test_Mock_MediaPackageV2OriginEndpoint_String(t *testing.T) {
	a := assert.New(t)
	r := &MediaPackageV2OriginEndpoint{OriginEndpointName: ptr.String("ep-1")}
	a.Equal("ep-1", r.String())
}
