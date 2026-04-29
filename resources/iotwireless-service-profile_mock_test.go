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

func Test_Mock_IoTWirelessServiceProfile_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTWirelessClient)

	mockClient.On("ListServiceProfiles", mock.Anything, mock.Anything).
		Return(&iotwireless.ListServiceProfilesOutput{
			ServiceProfileList: []iotwirelesstypes.ServiceProfile{
				{
					Id:   ptr.String("sp-12345"),
					Name: ptr.String("my-profile"),
				},
			},
		}, nil)

	lister := &IoTWirelessServiceProfileLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTWirelessListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	sp := resources[0].(*IoTWirelessServiceProfile)
	a.Equal("sp-12345", *sp.ID)
	a.Equal("my-profile", *sp.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTWirelessServiceProfile_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTWirelessClient)

	mockClient.On("ListServiceProfiles", mock.Anything, mock.Anything).
		Return(&iotwireless.ListServiceProfilesOutput{
			ServiceProfileList: []iotwirelesstypes.ServiceProfile{},
		}, nil)

	lister := &IoTWirelessServiceProfileLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTWirelessListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTWirelessServiceProfile_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTWirelessClient)

	sp := &IoTWirelessServiceProfile{
		svc: mockClient,
		ID:  ptr.String("sp-12345"),
	}

	mockClient.On("DeleteServiceProfile", mock.Anything, &iotwireless.DeleteServiceProfileInput{
		Id: sp.ID,
	}).Return(&iotwireless.DeleteServiceProfileOutput{}, nil)

	a.NoError(sp.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTWirelessServiceProfile_Properties(t *testing.T) {
	a := assert.New(t)

	sp := IoTWirelessServiceProfile{
		ID:   ptr.String("sp-12345"),
		Name: ptr.String("my-profile"),
	}

	props := sp.Properties()
	a.Equal("sp-12345", props.Get("Id"))
	a.Equal("my-profile", props.Get("Name"))
}

func Test_Mock_IoTWirelessServiceProfile_String(t *testing.T) {
	a := assert.New(t)
	sp := IoTWirelessServiceProfile{ID: ptr.String("sp-12345")}
	a.Equal("sp-12345", sp.String())
}
