package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2"
	pinpointtypes "github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testPinpointSMSVoiceV2ListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_PinpointSMSVoiceV2ConfigurationSet_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockPinpointSMSVoiceV2Client)

	mockClient.On("DescribeConfigurationSets", mock.Anything, mock.Anything).
		Return(&pinpointsmsvoicev2.DescribeConfigurationSetsOutput{
			ConfigurationSets: []pinpointtypes.ConfigurationSetInformation{
				{
					ConfigurationSetName: ptr.String("my-config-set"),
					ConfigurationSetArn:  ptr.String("arn:aws:sms-voice:us-east-1:123456789012:configuration-set/my-config-set"),
				},
			},
		}, nil)

	lister := &PinpointSMSVoiceV2ConfigurationSetLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testPinpointSMSVoiceV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	cs := resources[0].(*PinpointSMSVoiceV2ConfigurationSet)
	a.Equal("my-config-set", *cs.ConfigurationSetName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_PinpointSMSVoiceV2ConfigurationSet_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockPinpointSMSVoiceV2Client)

	mockClient.On("DescribeConfigurationSets", mock.Anything, mock.Anything).
		Return(&pinpointsmsvoicev2.DescribeConfigurationSetsOutput{
			ConfigurationSets: []pinpointtypes.ConfigurationSetInformation{},
		}, nil)

	lister := &PinpointSMSVoiceV2ConfigurationSetLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testPinpointSMSVoiceV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_PinpointSMSVoiceV2ConfigurationSet_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockPinpointSMSVoiceV2Client)

	cs := &PinpointSMSVoiceV2ConfigurationSet{
		svc:                  mockClient,
		ConfigurationSetName: ptr.String("my-config-set"),
	}

	mockClient.On("DeleteConfigurationSet", mock.Anything, &pinpointsmsvoicev2.DeleteConfigurationSetInput{
		ConfigurationSetName: cs.ConfigurationSetName,
	}).Return(&pinpointsmsvoicev2.DeleteConfigurationSetOutput{}, nil)

	a.NoError(cs.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_PinpointSMSVoiceV2ConfigurationSet_Properties(t *testing.T) {
	a := assert.New(t)

	cs := PinpointSMSVoiceV2ConfigurationSet{
		ConfigurationSetName: ptr.String("my-config-set"),
		ConfigurationSetArn:  ptr.String("arn:aws:sms-voice:us-east-1:123456789012:configuration-set/my-config-set"),
	}

	props := cs.Properties()
	a.Equal("my-config-set", props.Get("ConfigurationSetName"))
	a.Equal("arn:aws:sms-voice:us-east-1:123456789012:configuration-set/my-config-set", props.Get("ConfigurationSetArn"))
}

func Test_Mock_PinpointSMSVoiceV2ConfigurationSet_String(t *testing.T) {
	a := assert.New(t)
	cs := PinpointSMSVoiceV2ConfigurationSet{ConfigurationSetName: ptr.String("my-config-set")}
	a.Equal("my-config-set", cs.String())
}
