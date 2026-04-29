package resources

import (
	"context"
	"testing"
	"time"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2"
	pinpointtypes "github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2/types"
)

func Test_Mock_PinpointSMSVoiceV2ProtectConfigRuleSetOverride_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockPinpointSMSVoiceV2Client)

	mockClient.On("DescribeProtectConfigurations", mock.Anything, mock.Anything).
		Return(&pinpointsmsvoicev2.DescribeProtectConfigurationsOutput{
			ProtectConfigurations: []pinpointtypes.ProtectConfigurationInformation{
				{
					ProtectConfigurationId:  ptr.String("pc-12345"),
					ProtectConfigurationArn: ptr.String("arn:aws:sms-voice:us-east-1:123456789012:protect-configuration/pc-12345"),
				},
			},
		}, nil)

	now := time.Now()
	mockClient.On("ListProtectConfigurationRuleSetNumberOverrides", mock.Anything, mock.Anything).
		Return(&pinpointsmsvoicev2.ListProtectConfigurationRuleSetNumberOverridesOutput{
			ProtectConfigurationArn: ptr.String("arn:aws:sms-voice:us-east-1:123456789012:protect-configuration/pc-12345"),
			ProtectConfigurationId:  ptr.String("pc-12345"),
			RuleSetNumberOverrides: []pinpointtypes.ProtectConfigurationRuleSetNumberOverride{
				{
					DestinationPhoneNumber: ptr.String("+12065551234"),
					Action:                 pinpointtypes.ProtectConfigurationRuleOverrideActionBlock,
					CreatedTimestamp:       &now,
				},
			},
		}, nil)

	lister := &PinpointSMSVoiceV2ProtectConfigRuleSetOverrideLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testPinpointSMSVoiceV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	override := resources[0].(*PinpointSMSVoiceV2ProtectConfigRuleSetOverride)
	a.Equal("pc-12345", *override.ProtectConfigurationID)
	a.Equal("+12065551234", *override.DestinationPhoneNumber)

	mockClient.AssertExpectations(t)
}

func Test_Mock_PinpointSMSVoiceV2ProtectConfigRuleSetOverride_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockPinpointSMSVoiceV2Client)

	mockClient.On("DescribeProtectConfigurations", mock.Anything, mock.Anything).
		Return(&pinpointsmsvoicev2.DescribeProtectConfigurationsOutput{
			ProtectConfigurations: []pinpointtypes.ProtectConfigurationInformation{
				{
					ProtectConfigurationId:  ptr.String("pc-12345"),
					ProtectConfigurationArn: ptr.String("arn:aws:sms-voice:us-east-1:123456789012:protect-configuration/pc-12345"),
				},
			},
		}, nil)

	mockClient.On("ListProtectConfigurationRuleSetNumberOverrides", mock.Anything, mock.Anything).
		Return(&pinpointsmsvoicev2.ListProtectConfigurationRuleSetNumberOverridesOutput{
			ProtectConfigurationArn: ptr.String("arn:aws:sms-voice:us-east-1:123456789012:protect-configuration/pc-12345"),
			ProtectConfigurationId:  ptr.String("pc-12345"),
			RuleSetNumberOverrides:  []pinpointtypes.ProtectConfigurationRuleSetNumberOverride{},
		}, nil)

	lister := &PinpointSMSVoiceV2ProtectConfigRuleSetOverrideLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testPinpointSMSVoiceV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_PinpointSMSVoiceV2ProtectConfigRuleSetOverride_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockPinpointSMSVoiceV2Client)

	override := &PinpointSMSVoiceV2ProtectConfigRuleSetOverride{
		svc:                    mockClient,
		ProtectConfigurationID: ptr.String("pc-12345"),
		DestinationPhoneNumber: ptr.String("+12065551234"),
	}

	mockClient.On("DeleteProtectConfigurationRuleSetNumberOverride", mock.Anything,
		&pinpointsmsvoicev2.DeleteProtectConfigurationRuleSetNumberOverrideInput{
			ProtectConfigurationId: override.ProtectConfigurationID,
			DestinationPhoneNumber: override.DestinationPhoneNumber,
		}).Return(&pinpointsmsvoicev2.DeleteProtectConfigurationRuleSetNumberOverrideOutput{}, nil)

	a.NoError(override.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_PinpointSMSVoiceV2ProtectConfigRuleSetOverride_Properties(t *testing.T) {
	a := assert.New(t)

	override := PinpointSMSVoiceV2ProtectConfigRuleSetOverride{
		ProtectConfigurationID: ptr.String("pc-12345"),
		DestinationPhoneNumber: ptr.String("+12065551234"),
		Action:                 pinpointtypes.ProtectConfigurationRuleOverrideActionBlock,
		IsoCountryCode:         ptr.String("US"),
	}

	props := override.Properties()
	a.Equal("pc-12345", props.Get("ProtectConfigurationId"))
	a.Equal("+12065551234", props.Get("DestinationPhoneNumber"))
	a.Equal("US", props.Get("IsoCountryCode"))
}

func Test_Mock_PinpointSMSVoiceV2ProtectConfigRuleSetOverride_String(t *testing.T) {
	a := assert.New(t)
	override := PinpointSMSVoiceV2ProtectConfigRuleSetOverride{
		ProtectConfigurationID: ptr.String("pc-12345"),
		DestinationPhoneNumber: ptr.String("+12065551234"),
	}
	a.Equal("pc-12345 -> +12065551234", override.String())
}
