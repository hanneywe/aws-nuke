package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2"
	pinpointtypes "github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2/types"
)

func Test_Mock_PinpointSMSVoiceV2EventDestination_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockPinpointSMSVoiceV2Client)

	mockClient.On("DescribeConfigurationSets", mock.Anything, mock.Anything).
		Return(&pinpointsmsvoicev2.DescribeConfigurationSetsOutput{
			ConfigurationSets: []pinpointtypes.ConfigurationSetInformation{
				{
					ConfigurationSetName: ptr.String("my-config-set"),
					ConfigurationSetArn:  ptr.String("arn:aws:sms-voice:us-east-1:123456789012:configuration-set/my-config-set"),
					EventDestinations: []pinpointtypes.EventDestination{
						{
							EventDestinationName: ptr.String("my-event-dest"),
							Enabled:              ptr.Bool(true),
							MatchingEventTypes:   []pinpointtypes.EventType{pinpointtypes.EventTypeAll},
						},
					},
				},
			},
		}, nil)

	lister := &PinpointSMSVoiceV2EventDestinationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testPinpointSMSVoiceV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	ed := resources[0].(*PinpointSMSVoiceV2EventDestination)
	a.Equal("my-event-dest", *ed.EventDestinationName)
	a.Equal("my-config-set", *ed.ConfigurationSetName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_PinpointSMSVoiceV2EventDestination_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockPinpointSMSVoiceV2Client)

	mockClient.On("DescribeConfigurationSets", mock.Anything, mock.Anything).
		Return(&pinpointsmsvoicev2.DescribeConfigurationSetsOutput{
			ConfigurationSets: []pinpointtypes.ConfigurationSetInformation{
				{
					ConfigurationSetName: ptr.String("my-config-set"),
					ConfigurationSetArn:  ptr.String("arn:aws:sms-voice:us-east-1:123456789012:configuration-set/my-config-set"),
					EventDestinations:    []pinpointtypes.EventDestination{},
				},
			},
		}, nil)

	lister := &PinpointSMSVoiceV2EventDestinationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testPinpointSMSVoiceV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_PinpointSMSVoiceV2EventDestination_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockPinpointSMSVoiceV2Client)

	ed := &PinpointSMSVoiceV2EventDestination{
		svc:                  mockClient,
		EventDestinationName: ptr.String("my-event-dest"),
		ConfigurationSetName: ptr.String("my-config-set"),
	}

	mockClient.On("DeleteEventDestination", mock.Anything, &pinpointsmsvoicev2.DeleteEventDestinationInput{
		ConfigurationSetName: ed.ConfigurationSetName,
		EventDestinationName: ed.EventDestinationName,
	}).Return(&pinpointsmsvoicev2.DeleteEventDestinationOutput{}, nil)

	a.NoError(ed.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_PinpointSMSVoiceV2EventDestination_Properties(t *testing.T) {
	a := assert.New(t)

	ed := PinpointSMSVoiceV2EventDestination{
		EventDestinationName: ptr.String("my-event-dest"),
		ConfigurationSetName: ptr.String("my-config-set"),
	}

	props := ed.Properties()
	a.Equal("my-event-dest", props.Get("EventDestinationName"))
	a.Equal("my-config-set", props.Get("ConfigurationSetName"))
}

func Test_Mock_PinpointSMSVoiceV2EventDestination_String(t *testing.T) {
	a := assert.New(t)
	ed := PinpointSMSVoiceV2EventDestination{
		EventDestinationName: ptr.String("my-event-dest"),
		ConfigurationSetName: ptr.String("my-config-set"),
	}
	a.Equal("my-config-set -> my-event-dest", ed.String())
}
