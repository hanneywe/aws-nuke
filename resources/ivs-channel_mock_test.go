package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ivs"
	ivstypes "github.com/aws/aws-sdk-go-v2/service/ivs/types"
)

func Test_Mock_IVSChannel_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIVSClient)

	mockClient.On("ListChannels", mock.Anything, mock.Anything).
		Return(&ivs.ListChannelsOutput{
			Channels: []ivstypes.ChannelSummary{
				{
					Arn:  ptr.String("arn:aws:ivs:us-east-1:123456789012:channel/abc123"),
					Name: ptr.String("my-channel"),
					Tags: map[string]string{"env": "test"},
				},
			},
		}, nil)

	lister := &IVSChannelLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIVSListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	ch := resources[0].(*IVSChannel)
	a.Equal("my-channel", *ch.Name)
	a.Equal("test", ch.Tags["env"])

	mockClient.AssertExpectations(t)
}

func Test_Mock_IVSChannel_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIVSClient)

	mockClient.On("ListChannels", mock.Anything, mock.Anything).
		Return(&ivs.ListChannelsOutput{
			Channels: []ivstypes.ChannelSummary{},
		}, nil)

	lister := &IVSChannelLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIVSListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IVSChannel_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIVSClient)

	ch := &IVSChannel{
		svc: mockClient,
		ARN: ptr.String("arn:aws:ivs:us-east-1:123456789012:channel/abc123"),
	}

	mockClient.On("DeleteChannel", mock.Anything, &ivs.DeleteChannelInput{
		Arn: ch.ARN,
	}).Return(&ivs.DeleteChannelOutput{}, nil)

	a.NoError(ch.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IVSChannel_Properties(t *testing.T) {
	a := assert.New(t)

	ch := IVSChannel{
		ARN:  ptr.String("arn:aws:ivs:us-east-1:123456789012:channel/abc123"),
		Name: ptr.String("my-channel"),
		Tags: map[string]string{"env": "test"},
	}

	props := ch.Properties()
	a.Equal("my-channel", props.Get("Name"))
	a.Equal("test", props.Get("tag:env"))
}

func Test_Mock_IVSChannel_String(t *testing.T) {
	a := assert.New(t)
	ch := IVSChannel{Name: ptr.String("my-channel")}
	a.Equal("my-channel", ch.String())
}
