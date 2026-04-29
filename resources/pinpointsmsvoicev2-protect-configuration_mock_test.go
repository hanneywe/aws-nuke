package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2"
	pinpointtypes "github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2/types"

	libsettings "github.com/ekristen/libnuke/pkg/settings"
)

func Test_Mock_PinpointSMSVoiceV2ProtectConfiguration_List_One(t *testing.T) {
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

	lister := &PinpointSMSVoiceV2ProtectConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testPinpointSMSVoiceV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	pc := resources[0].(*PinpointSMSVoiceV2ProtectConfiguration)
	a.Equal("pc-12345", *pc.ProtectConfigurationID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_PinpointSMSVoiceV2ProtectConfiguration_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockPinpointSMSVoiceV2Client)

	mockClient.On("DescribeProtectConfigurations", mock.Anything, mock.Anything).
		Return(&pinpointsmsvoicev2.DescribeProtectConfigurationsOutput{
			ProtectConfigurations: []pinpointtypes.ProtectConfigurationInformation{},
		}, nil)

	lister := &PinpointSMSVoiceV2ProtectConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testPinpointSMSVoiceV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_PinpointSMSVoiceV2ProtectConfiguration_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockPinpointSMSVoiceV2Client)

	pc := &PinpointSMSVoiceV2ProtectConfiguration{
		svc:                    mockClient,
		ProtectConfigurationID: ptr.String("pc-12345"),
	}

	mockClient.On("DeleteProtectConfiguration", mock.Anything,
		&pinpointsmsvoicev2.DeleteProtectConfigurationInput{
			ProtectConfigurationId: pc.ProtectConfigurationID,
		}).Return(&pinpointsmsvoicev2.DeleteProtectConfigurationOutput{}, nil)

	a.NoError(pc.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_PinpointSMSVoiceV2ProtectConfiguration_Properties(t *testing.T) {
	a := assert.New(t)

	pc := PinpointSMSVoiceV2ProtectConfiguration{
		ProtectConfigurationID:  ptr.String("pc-12345"),
		ProtectConfigurationArn: ptr.String("arn:aws:sms-voice:us-east-1:123456789012:protect-configuration/pc-12345"),
	}

	props := pc.Properties()
	a.Equal("pc-12345", props.Get("ProtectConfigurationId"))
	a.Equal("arn:aws:sms-voice:us-east-1:123456789012:protect-configuration/pc-12345",
		props.Get("ProtectConfigurationArn"))
}

func Test_Mock_PinpointSMSVoiceV2ProtectConfiguration_String(t *testing.T) {
	a := assert.New(t)
	pc := PinpointSMSVoiceV2ProtectConfiguration{
		ProtectConfigurationID: ptr.String("pc-12345"),
	}
	a.Equal("pc-12345", pc.String())
}

func Test_Mock_PinpointSMSVoiceV2ProtectConfiguration_Remove_WithDeletionProtection(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockPinpointSMSVoiceV2Client)

	settings := &libsettings.Setting{}
	settings.Set("DisableDeletionProtection", true)

	pc := &PinpointSMSVoiceV2ProtectConfiguration{
		svc:                       mockClient,
		settings:                  settings,
		ProtectConfigurationID:    ptr.String("pc-12345"),
		DeletionProtectionEnabled: ptr.Bool(true),
	}

	mockClient.On("UpdateProtectConfiguration", mock.Anything,
		&pinpointsmsvoicev2.UpdateProtectConfigurationInput{
			ProtectConfigurationId:    pc.ProtectConfigurationID,
			DeletionProtectionEnabled: ptr.Bool(false),
		}).Return(&pinpointsmsvoicev2.UpdateProtectConfigurationOutput{}, nil)

	mockClient.On("DeleteProtectConfiguration", mock.Anything,
		&pinpointsmsvoicev2.DeleteProtectConfigurationInput{
			ProtectConfigurationId: pc.ProtectConfigurationID,
		}).Return(&pinpointsmsvoicev2.DeleteProtectConfigurationOutput{}, nil)

	a.NoError(pc.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_PinpointSMSVoiceV2ProtectConfiguration_Remove_ProtectionNotDisabled(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockPinpointSMSVoiceV2Client)

	settings := &libsettings.Setting{}
	settings.Set("DisableDeletionProtection", false)

	pc := &PinpointSMSVoiceV2ProtectConfiguration{
		svc:                       mockClient,
		settings:                  settings,
		ProtectConfigurationID:    ptr.String("pc-12345"),
		DeletionProtectionEnabled: ptr.Bool(true),
	}

	// Should NOT call UpdateProtectConfiguration since setting is false
	mockClient.On("DeleteProtectConfiguration", mock.Anything,
		&pinpointsmsvoicev2.DeleteProtectConfigurationInput{
			ProtectConfigurationId: pc.ProtectConfigurationID,
		}).Return(&pinpointsmsvoicev2.DeleteProtectConfigurationOutput{}, nil)

	a.NoError(pc.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}
