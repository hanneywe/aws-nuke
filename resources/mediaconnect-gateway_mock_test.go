package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/mediaconnect"
	mediaconnecttypes "github.com/aws/aws-sdk-go-v2/service/mediaconnect/types"
)

func Test_Mock_MediaConnectGateway_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockMediaConnectClient)

	mockClient.
		On("ListGateways", mock.Anything, mock.Anything).
		Return(&mediaconnect.ListGatewaysOutput{
			Gateways: []mediaconnecttypes.ListedGateway{
				{
					GatewayArn:   ptr.String("arn:aws:mediaconnect:us-east-1:123456789012:gateway:1-abcdef"),
					Name:         ptr.String("my-gateway"),
					GatewayState: mediaconnecttypes.GatewayStateActive,
				},
			},
		}, nil)

	lister := &MediaConnectGatewayLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testMediaConnectListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	gateway := resources[0].(*MediaConnectGateway)
	assertions.Equal("arn:aws:mediaconnect:us-east-1:123456789012:gateway:1-abcdef", *gateway.GatewayArn)
	assertions.Equal("my-gateway", *gateway.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaConnectGateway_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockMediaConnectClient)

	mockClient.
		On("ListGateways", mock.Anything, mock.Anything).
		Return(&mediaconnect.ListGatewaysOutput{
			Gateways: []mediaconnecttypes.ListedGateway{},
		}, nil)

	lister := &MediaConnectGatewayLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testMediaConnectListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaConnectGateway_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockMediaConnectClient)

	gateway := &MediaConnectGateway{
		svc:        mockClient,
		GatewayArn: ptr.String("arn:aws:mediaconnect:us-east-1:123456789012:gateway:1-abcdef"),
		Name:       ptr.String("my-gateway"),
	}

	mockClient.
		On("DeleteGateway", mock.Anything, &mediaconnect.DeleteGatewayInput{
			GatewayArn: gateway.GatewayArn,
		}).
		Return(&mediaconnect.DeleteGatewayOutput{}, nil)

	err := gateway.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaConnectGateway_Properties(t *testing.T) {
	assertions := assert.New(t)

	gateway := MediaConnectGateway{
		GatewayArn: ptr.String("arn:aws:mediaconnect:us-east-1:123456789012:gateway:1-abcdef"),
		Name:       ptr.String("my-gateway"),
	}

	properties := gateway.Properties()

	assertions.Equal("arn:aws:mediaconnect:us-east-1:123456789012:gateway:1-abcdef", properties.Get("GatewayArn"))
	assertions.Equal("my-gateway", properties.Get("Name"))
}

func Test_Mock_MediaConnectGateway_String(t *testing.T) {
	assertions := assert.New(t)

	gateway := MediaConnectGateway{
		GatewayArn: ptr.String("arn:aws:mediaconnect:us-east-1:123456789012:gateway:1-abcdef"),
	}

	assertions.Equal("arn:aws:mediaconnect:us-east-1:123456789012:gateway:1-abcdef", gateway.String())
}

func Test_Mock_MediaConnectGateway_Filter_Active(t *testing.T) {
	assertions := assert.New(t)

	gateway := MediaConnectGateway{
		GatewayArn:   ptr.String("arn:aws:mediaconnect:us-east-1:123456789012:gateway:1-abcdef"),
		GatewayState: mediaconnecttypes.GatewayStateActive,
	}

	err := gateway.Filter()
	assertions.NoError(err)
}

func Test_Mock_MediaConnectGateway_Filter_Deleting(t *testing.T) {
	assertions := assert.New(t)

	gateway := MediaConnectGateway{
		GatewayArn:   ptr.String("arn:aws:mediaconnect:us-east-1:123456789012:gateway:1-abcdef"),
		GatewayState: mediaconnecttypes.GatewayStateDeleting,
	}

	err := gateway.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "DELETING")
}

func Test_Mock_MediaConnectGateway_Filter_Deleted(t *testing.T) {
	assertions := assert.New(t)

	gateway := MediaConnectGateway{
		GatewayArn:   ptr.String("arn:aws:mediaconnect:us-east-1:123456789012:gateway:1-abcdef"),
		GatewayState: mediaconnecttypes.GatewayStateDeleted,
	}

	err := gateway.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "DELETED")
}
