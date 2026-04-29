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

func Test_Mock_IoTWirelessMulticastGroup_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTWirelessClient)

	mockClient.On("ListMulticastGroups", mock.Anything, mock.Anything).
		Return(&iotwireless.ListMulticastGroupsOutput{
			MulticastGroupList: []iotwirelesstypes.MulticastGroup{
				{
					Id:   ptr.String("mg-12345"),
					Name: ptr.String("my-group"),
				},
			},
		}, nil)

	lister := &IoTWirelessMulticastGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTWirelessListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	mg := resources[0].(*IoTWirelessMulticastGroup)
	a.Equal("mg-12345", *mg.ID)
	a.Equal("my-group", *mg.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTWirelessMulticastGroup_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTWirelessClient)

	mockClient.On("ListMulticastGroups", mock.Anything, mock.Anything).
		Return(&iotwireless.ListMulticastGroupsOutput{
			MulticastGroupList: []iotwirelesstypes.MulticastGroup{},
		}, nil)

	lister := &IoTWirelessMulticastGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTWirelessListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTWirelessMulticastGroup_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTWirelessClient)

	mg := &IoTWirelessMulticastGroup{
		svc: mockClient,
		ID:  ptr.String("mg-12345"),
	}

	mockClient.On("DeleteMulticastGroup", mock.Anything, &iotwireless.DeleteMulticastGroupInput{
		Id: mg.ID,
	}).Return(&iotwireless.DeleteMulticastGroupOutput{}, nil)

	a.NoError(mg.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTWirelessMulticastGroup_Properties(t *testing.T) {
	a := assert.New(t)

	mg := IoTWirelessMulticastGroup{
		ID:   ptr.String("mg-12345"),
		Name: ptr.String("my-group"),
	}

	props := mg.Properties()
	a.Equal("mg-12345", props.Get("Id"))
	a.Equal("my-group", props.Get("Name"))
}

func Test_Mock_IoTWirelessMulticastGroup_String(t *testing.T) {
	a := assert.New(t)
	mg := IoTWirelessMulticastGroup{ID: ptr.String("mg-12345")}
	a.Equal("mg-12345", mg.String())
}
