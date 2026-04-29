package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/workspacesweb"
	workspaceswebtypes "github.com/aws/aws-sdk-go-v2/service/workspacesweb/types"
)

func Test_Mock_WorkSpacesWebDataProtectionSettings_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockWorkSpacesWebClient)

	mockClient.On("ListDataProtectionSettings", mock.Anything, mock.Anything).
		Return(&workspacesweb.ListDataProtectionSettingsOutput{
			DataProtectionSettings: []workspaceswebtypes.DataProtectionSettingsSummary{
				{
					DataProtectionSettingsArn: ptr.String("arn:aws:workspaces-web:us-east-1:123456789012:data-protection-settings/dps-12345"),
					DisplayName:               ptr.String("my-dps"),
				},
			},
		}, nil)

	lister := &WorkSpacesWebDataProtectionSettingsLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testWorkSpacesWebListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	dps := resources[0].(*WorkSpacesWebDataProtectionSettings)
	assertions.Equal("arn:aws:workspaces-web:us-east-1:123456789012:data-protection-settings/dps-12345", *dps.DataProtectionSettingsArn)
	assertions.Equal("my-dps", *dps.DisplayName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_WorkSpacesWebDataProtectionSettings_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockWorkSpacesWebClient)

	mockClient.On("ListDataProtectionSettings", mock.Anything, mock.Anything).
		Return(&workspacesweb.ListDataProtectionSettingsOutput{
			DataProtectionSettings: []workspaceswebtypes.DataProtectionSettingsSummary{},
		}, nil)

	lister := &WorkSpacesWebDataProtectionSettingsLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testWorkSpacesWebListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_WorkSpacesWebDataProtectionSettings_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockWorkSpacesWebClient)

	dps := &WorkSpacesWebDataProtectionSettings{
		svc:                       mockClient,
		DataProtectionSettingsArn: ptr.String("arn:aws:workspaces-web:us-east-1:123456789012:data-protection-settings/dps-12345"),
	}

	mockClient.On("DeleteDataProtectionSettings", mock.Anything, &workspacesweb.DeleteDataProtectionSettingsInput{
		DataProtectionSettingsArn: dps.DataProtectionSettingsArn,
	}).Return(&workspacesweb.DeleteDataProtectionSettingsOutput{}, nil)

	err := dps.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_WorkSpacesWebDataProtectionSettings_Properties(t *testing.T) {
	assertions := assert.New(t)

	dps := WorkSpacesWebDataProtectionSettings{
		DataProtectionSettingsArn: ptr.String("arn:aws:workspaces-web:us-east-1:123456789012:data-protection-settings/dps-12345"),
		DisplayName:               ptr.String("my-dps"),
	}

	properties := dps.Properties()
	assertions.Equal(
		"arn:aws:workspaces-web:us-east-1:123456789012:data-protection-settings/dps-12345",
		properties.Get("DataProtectionSettingsArn"),
	)
	assertions.Equal("my-dps", properties.Get("DisplayName"))
}

func Test_Mock_WorkSpacesWebDataProtectionSettings_String(t *testing.T) {
	assertions := assert.New(t)
	testArn := "arn:aws:workspaces-web:us-east-1:123456789012:data-protection-settings/dps-12345"
	dps := WorkSpacesWebDataProtectionSettings{DataProtectionSettingsArn: ptr.String(testArn)}
	assertions.Equal(testArn, dps.String())
}
