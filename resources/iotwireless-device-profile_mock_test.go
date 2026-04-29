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

func Test_Mock_IoTWirelessDeviceProfile_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIoTWirelessClient)

	mockClient.On("ListDeviceProfiles", mock.Anything, mock.Anything).
		Return(&iotwireless.ListDeviceProfilesOutput{
			DeviceProfileList: []iotwirelesstypes.DeviceProfile{
				{
					Id:   ptr.String("dp-11111"),
					Name: ptr.String("my-device-profile"),
				},
				{
					Id:   ptr.String("dp-22222"),
					Name: ptr.String("another-device-profile"),
				},
			},
		}, nil)

	lister := &IoTWirelessDeviceProfileLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTWirelessListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 2)

	deviceProfile := resources[0].(*IoTWirelessDeviceProfile)
	assertions.Equal("dp-11111", *deviceProfile.ID)
	assertions.Equal("my-device-profile", *deviceProfile.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTWirelessDeviceProfile_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIoTWirelessClient)

	mockClient.On("ListDeviceProfiles", mock.Anything, mock.Anything).
		Return(&iotwireless.ListDeviceProfilesOutput{
			DeviceProfileList: []iotwirelesstypes.DeviceProfile{},
		}, nil)

	lister := &IoTWirelessDeviceProfileLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTWirelessListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTWirelessDeviceProfile_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIoTWirelessClient)

	deviceProfile := &IoTWirelessDeviceProfile{
		svc: mockClient,
		ID:  ptr.String("dp-11111"),
	}

	mockClient.On("DeleteDeviceProfile", mock.Anything, &iotwireless.DeleteDeviceProfileInput{
		Id: deviceProfile.ID,
	}).Return(&iotwireless.DeleteDeviceProfileOutput{}, nil)

	assertions.NoError(deviceProfile.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTWirelessDeviceProfile_Properties(t *testing.T) {
	assertions := assert.New(t)

	deviceProfile := IoTWirelessDeviceProfile{
		ID:   ptr.String("dp-11111"),
		Name: ptr.String("my-device-profile"),
	}

	props := deviceProfile.Properties()
	assertions.Equal("dp-11111", props.Get("Id"))
	assertions.Equal("my-device-profile", props.Get("Name"))
}

func Test_Mock_IoTWirelessDeviceProfile_String(t *testing.T) {
	assertions := assert.New(t)
	deviceProfile := IoTWirelessDeviceProfile{ID: ptr.String("dp-11111")}
	assertions.Equal("dp-11111", deviceProfile.String())
}
