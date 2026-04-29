package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/devicefarm"
	devicefarmtypes "github.com/aws/aws-sdk-go-v2/service/devicefarm/types"
)

func Test_Mock_DeviceFarmInstanceProfile_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockDeviceFarmClient)

	mockClient.
		On("ListInstanceProfiles", mock.Anything, mock.Anything).
		Return(&devicefarm.ListInstanceProfilesOutput{
			InstanceProfiles: []devicefarmtypes.InstanceProfile{
				{
					Arn:  ptr.String("arn:aws:devicefarm:us-west-2:123456789012:instanceprofile:example-1"),
					Name: ptr.String("my-instance-profile"),
				},
			},
		}, nil)

	lister := &DeviceFarmInstanceProfileLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testDeviceFarmListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	instanceProfile := resources[0].(*DeviceFarmInstanceProfile)
	assertions.Equal("arn:aws:devicefarm:us-west-2:123456789012:instanceprofile:example-1", *instanceProfile.Arn)
	assertions.Equal("my-instance-profile", *instanceProfile.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DeviceFarmInstanceProfile_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockDeviceFarmClient)

	mockClient.
		On("ListInstanceProfiles", mock.Anything, mock.Anything).
		Return(&devicefarm.ListInstanceProfilesOutput{
			InstanceProfiles: []devicefarmtypes.InstanceProfile{},
		}, nil)

	lister := &DeviceFarmInstanceProfileLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testDeviceFarmListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DeviceFarmInstanceProfile_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockDeviceFarmClient)

	instanceProfile := &DeviceFarmInstanceProfile{
		svc:  mockClient,
		Arn:  ptr.String("arn:aws:devicefarm:us-west-2:123456789012:instanceprofile:example-1"),
		Name: ptr.String("my-instance-profile"),
	}

	mockClient.
		On("DeleteInstanceProfile", mock.Anything, &devicefarm.DeleteInstanceProfileInput{
			Arn: instanceProfile.Arn,
		}).
		Return(&devicefarm.DeleteInstanceProfileOutput{}, nil)

	err := instanceProfile.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DeviceFarmInstanceProfile_Properties(t *testing.T) {
	assertions := assert.New(t)

	instanceProfile := DeviceFarmInstanceProfile{
		Arn:  ptr.String("arn:aws:devicefarm:us-west-2:123456789012:instanceprofile:example-1"),
		Name: ptr.String("my-instance-profile"),
	}

	properties := instanceProfile.Properties()

	assertions.Equal("arn:aws:devicefarm:us-west-2:123456789012:instanceprofile:example-1", properties.Get("Arn"))
	assertions.Equal("my-instance-profile", properties.Get("Name"))
}

func Test_Mock_DeviceFarmInstanceProfile_String(t *testing.T) {
	assertions := assert.New(t)

	instanceProfile := DeviceFarmInstanceProfile{
		Name: ptr.String("my-instance-profile"),
	}

	assertions.Equal("my-instance-profile", instanceProfile.String())
}
