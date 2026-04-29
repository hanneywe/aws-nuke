package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/iotwireless"
	iotwirelesstypes "github.com/aws/aws-sdk-go-v2/service/iotwireless/types"
)

func Test_Mock_IoTWirelessWirelessGateway_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIoTWirelessClient)

	mockClient.On("ListWirelessGateways", mock.Anything, mock.Anything).
		Return(&iotwireless.ListWirelessGatewaysOutput{
			WirelessGatewayList: []iotwirelesstypes.WirelessGatewayStatistics{
				{
					Id:   ptr.String("gw-11111"),
					Name: ptr.String("my-gateway"),
					Arn:  ptr.String("arn:aws:iotwireless:us-east-1:123456789012:WirelessGateway/gw-11111"),
				},
				{
					Id:   ptr.String("gw-22222"),
					Name: ptr.String("another-gateway"),
					Arn:  ptr.String("arn:aws:iotwireless:us-east-1:123456789012:WirelessGateway/gw-22222"),
				},
			},
		}, nil)

	lister := &IoTWirelessWirelessGatewayLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTWirelessListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 2)

	gw := resources[0].(*IoTWirelessWirelessGateway)
	assertions.Equal("gw-11111", *gw.ID)
	assertions.Equal("my-gateway", *gw.Name)
	assertions.Equal("arn:aws:iotwireless:us-east-1:123456789012:WirelessGateway/gw-11111", *gw.ARN)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTWirelessWirelessGateway_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIoTWirelessClient)

	mockClient.On("ListWirelessGateways", mock.Anything, mock.Anything).
		Return(&iotwireless.ListWirelessGatewaysOutput{
			WirelessGatewayList: []iotwirelesstypes.WirelessGatewayStatistics{},
		}, nil)

	lister := &IoTWirelessWirelessGatewayLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTWirelessListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTWirelessWirelessGateway_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIoTWirelessClient)

	gw := &IoTWirelessWirelessGateway{
		svc: mockClient,
		ID:  ptr.String("gw-11111"),
	}

	mockClient.On("DeleteWirelessGateway", mock.Anything, &iotwireless.DeleteWirelessGatewayInput{
		Id: gw.ID,
	}).Return(&iotwireless.DeleteWirelessGatewayOutput{}, nil)

	assertions.NoError(gw.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTWirelessWirelessGateway_Properties(t *testing.T) {
	assertions := assert.New(t)

	gw := IoTWirelessWirelessGateway{
		ID:   ptr.String("gw-11111"),
		Name: ptr.String("my-gateway"),
		ARN:  ptr.String("arn:aws:iotwireless:us-east-1:123456789012:WirelessGateway/gw-11111"),
	}

	props := gw.Properties()
	assertions.Equal("gw-11111", props.Get("Id"))
	assertions.Equal("my-gateway", props.Get("Name"))
	assertions.Equal("arn:aws:iotwireless:us-east-1:123456789012:WirelessGateway/gw-11111", props.Get("Arn"))
}

func Test_Mock_IoTWirelessWirelessGateway_String(t *testing.T) {
	assertions := assert.New(t)
	gw := IoTWirelessWirelessGateway{ID: ptr.String("gw-11111")}
	assertions.Equal("gw-11111", gw.String())
}
