package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/medialive"
	mltypes "github.com/aws/aws-sdk-go-v2/service/medialive/types"
)

func Test_Mock_MediaLiveNetwork_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaLiveClient)

	mockClient.On("ListNetworks", mock.Anything, mock.Anything).
		Return(&medialive.ListNetworksOutput{
			Networks: []mltypes.DescribeNetworkSummary{
				{
					Id:   ptr.String("net-123"),
					Name: ptr.String("my-network"),
				},
			},
		}, nil)

	lister := &MediaLiveNetworkLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaLiveListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*MediaLiveNetwork)
	a.Equal("net-123", *r.ID)
	a.Equal("my-network", *r.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaLiveNetwork_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaLiveClient)

	mockClient.On("ListNetworks", mock.Anything, mock.Anything).
		Return(&medialive.ListNetworksOutput{
			Networks: []mltypes.DescribeNetworkSummary{},
		}, nil)

	lister := &MediaLiveNetworkLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaLiveListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaLiveNetwork_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaLiveClient)

	r := &MediaLiveNetwork{
		svc: mockClient,
		ID:  ptr.String("net-123"),
	}

	mockClient.On("DeleteNetwork", mock.Anything, &medialive.DeleteNetworkInput{
		NetworkId: r.ID,
	}).Return(&medialive.DeleteNetworkOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaLiveNetwork_Properties(t *testing.T) {
	a := assert.New(t)

	r := MediaLiveNetwork{
		ID:   ptr.String("net-123"),
		Name: ptr.String("my-network"),
	}

	props := r.Properties()
	a.Equal("net-123", props.Get("ID"))
	a.Equal("my-network", props.Get("Name"))
}

func Test_Mock_MediaLiveNetwork_String(t *testing.T) {
	a := assert.New(t)
	r := MediaLiveNetwork{ID: ptr.String("net-123")}
	a.Equal("net-123", r.String())
}
