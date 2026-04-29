package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/directconnect"
	directconnecttypes "github.com/aws/aws-sdk-go-v2/service/directconnect/types"
)

func Test_Mock_DirectConnectGateway_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDirectConnectClient)

	mockClient.
		On("DescribeDirectConnectGateways", mock.Anything, mock.Anything).
		Return(&directconnect.DescribeDirectConnectGatewaysOutput{
			DirectConnectGateways: []directconnecttypes.DirectConnectGateway{
				{
					DirectConnectGatewayId:   ptr.String("dgw-12345"),
					DirectConnectGatewayName: ptr.String("my-gateway"),
				},
			},
		}, nil)

	lister := &DirectConnectGatewayLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDirectConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	gw := resources[0].(*DirectConnectGateway)
	a.Equal("dgw-12345", *gw.DirectConnectGatewayID)
	a.Equal("my-gateway", *gw.DirectConnectGatewayName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DirectConnectGateway_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDirectConnectClient)

	mockClient.
		On("DescribeDirectConnectGateways", mock.Anything, mock.Anything).
		Return(&directconnect.DescribeDirectConnectGatewaysOutput{
			DirectConnectGateways: []directconnecttypes.DirectConnectGateway{},
		}, nil)

	lister := &DirectConnectGatewayLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDirectConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DirectConnectGateway_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDirectConnectClient)

	gw := &DirectConnectGateway{
		svc:                    mockClient,
		DirectConnectGatewayID: ptr.String("dgw-12345"),
	}

	mockClient.
		On("DeleteDirectConnectGateway", mock.Anything, &directconnect.DeleteDirectConnectGatewayInput{
			DirectConnectGatewayId: gw.DirectConnectGatewayID,
		}).
		Return(&directconnect.DeleteDirectConnectGatewayOutput{}, nil)

	err := gw.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DirectConnectGateway_Properties(t *testing.T) {
	a := assert.New(t)

	gw := DirectConnectGateway{
		DirectConnectGatewayID:   ptr.String("dgw-12345"),
		DirectConnectGatewayName: ptr.String("my-gateway"),
	}

	props := gw.Properties()
	a.Equal("dgw-12345", props.Get("DirectConnectGatewayId"))
	a.Equal("my-gateway", props.Get("DirectConnectGatewayName"))
}

func Test_Mock_DirectConnectGateway_String(t *testing.T) {
	a := assert.New(t)

	gw := DirectConnectGateway{
		DirectConnectGatewayID: ptr.String("dgw-12345"),
	}

	a.Equal("dgw-12345", gw.String())
}
