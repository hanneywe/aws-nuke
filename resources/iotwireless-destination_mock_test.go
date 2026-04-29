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

func Test_Mock_IoTWirelessDestination_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIoTWirelessClient)

	mockClient.On("ListDestinations", mock.Anything, mock.Anything).
		Return(&iotwireless.ListDestinationsOutput{
			DestinationList: []iotwirelesstypes.Destinations{
				{
					Name:           ptr.String("my-destination"),
					ExpressionType: iotwirelesstypes.ExpressionTypeRuleName,
					RoleArn:        ptr.String("arn:aws:iam::123456789012:role/my-role"),
				},
				{
					Name:           ptr.String("another-destination"),
					ExpressionType: iotwirelesstypes.ExpressionTypeMqttTopic,
					RoleArn:        ptr.String("arn:aws:iam::123456789012:role/other-role"),
				},
			},
		}, nil)

	lister := &IoTWirelessDestinationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTWirelessListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 2)

	dest := resources[0].(*IoTWirelessDestination)
	assertions.Equal("my-destination", *dest.Name)
	assertions.Equal(iotwirelesstypes.ExpressionTypeRuleName, dest.ExpressionType)
	assertions.Equal("arn:aws:iam::123456789012:role/my-role", *dest.RoleArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTWirelessDestination_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIoTWirelessClient)

	mockClient.On("ListDestinations", mock.Anything, mock.Anything).
		Return(&iotwireless.ListDestinationsOutput{
			DestinationList: []iotwirelesstypes.Destinations{},
		}, nil)

	lister := &IoTWirelessDestinationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTWirelessListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTWirelessDestination_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIoTWirelessClient)

	dest := &IoTWirelessDestination{
		svc:  mockClient,
		Name: ptr.String("my-destination"),
	}

	mockClient.On("DeleteDestination", mock.Anything, &iotwireless.DeleteDestinationInput{
		Name: dest.Name,
	}).Return(&iotwireless.DeleteDestinationOutput{}, nil)

	assertions.NoError(dest.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTWirelessDestination_Properties(t *testing.T) {
	assertions := assert.New(t)

	dest := IoTWirelessDestination{
		Name:           ptr.String("my-destination"),
		ExpressionType: iotwirelesstypes.ExpressionTypeRuleName,
		RoleArn:        ptr.String("arn:aws:iam::123456789012:role/my-role"),
	}

	props := dest.Properties()
	assertions.Equal("my-destination", props.Get("Name"))
	assertions.Equal("RuleName", props.Get("ExpressionType"))
	assertions.Equal("arn:aws:iam::123456789012:role/my-role", props.Get("RoleArn"))
}

func Test_Mock_IoTWirelessDestination_String(t *testing.T) {
	assertions := assert.New(t)
	dest := IoTWirelessDestination{Name: ptr.String("my-destination")}
	assertions.Equal("my-destination", dest.String())
}
