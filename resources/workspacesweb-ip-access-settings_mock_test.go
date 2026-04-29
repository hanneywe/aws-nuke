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

func Test_Mock_WorkSpacesWebIPAccessSettings_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockWorkSpacesWebClient)

	mockClient.
		On("ListIpAccessSettings", mock.Anything, mock.Anything).
		Return(&workspacesweb.ListIpAccessSettingsOutput{
			IpAccessSettings: []workspaceswebtypes.IpAccessSettingsSummary{
				{
					IpAccessSettingsArn: ptr.String("arn:aws:workspaces-web:us-east-1:123456789012:ipAccessSettings/ip-1"),
					DisplayName:         ptr.String("my-ip-access-settings"),
				},
			},
		}, nil)

	lister := &WorkSpacesWebIPAccessSettingsLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testWorkSpacesWebListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	ipAccessSettings := resources[0].(*WorkSpacesWebIPAccessSettings)
	assertions.Equal("arn:aws:workspaces-web:us-east-1:123456789012:ipAccessSettings/ip-1", *ipAccessSettings.IPAccessSettingsArn)
	assertions.Equal("my-ip-access-settings", *ipAccessSettings.DisplayName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_WorkSpacesWebIPAccessSettings_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockWorkSpacesWebClient)

	mockClient.
		On("ListIpAccessSettings", mock.Anything, mock.Anything).
		Return(&workspacesweb.ListIpAccessSettingsOutput{
			IpAccessSettings: []workspaceswebtypes.IpAccessSettingsSummary{},
		}, nil)

	lister := &WorkSpacesWebIPAccessSettingsLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testWorkSpacesWebListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_WorkSpacesWebIPAccessSettings_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockWorkSpacesWebClient)

	ipAccessSettings := &WorkSpacesWebIPAccessSettings{
		svc:                 mockClient,
		IPAccessSettingsArn: ptr.String("arn:aws:workspaces-web:us-east-1:123456789012:ipAccessSettings/ip-1"),
		DisplayName:         ptr.String("my-ip-access-settings"),
	}

	mockClient.
		On("DeleteIpAccessSettings", mock.Anything, &workspacesweb.DeleteIpAccessSettingsInput{
			IpAccessSettingsArn: ipAccessSettings.IPAccessSettingsArn,
		}).
		Return(&workspacesweb.DeleteIpAccessSettingsOutput{}, nil)

	err := ipAccessSettings.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_WorkSpacesWebIPAccessSettings_Properties(t *testing.T) {
	assertions := assert.New(t)

	ipAccessSettings := WorkSpacesWebIPAccessSettings{
		IPAccessSettingsArn: ptr.String("arn:aws:workspaces-web:us-east-1:123456789012:ipAccessSettings/ip-1"),
		DisplayName:         ptr.String("my-ip-access-settings"),
	}

	properties := ipAccessSettings.Properties()

	assertions.Equal("arn:aws:workspaces-web:us-east-1:123456789012:ipAccessSettings/ip-1", properties.Get("IPAccessSettingsArn"))
	assertions.Equal("my-ip-access-settings", properties.Get("DisplayName"))
}

func Test_Mock_WorkSpacesWebIPAccessSettings_String(t *testing.T) {
	assertions := assert.New(t)

	ipAccessSettings := WorkSpacesWebIPAccessSettings{
		IPAccessSettingsArn: ptr.String("arn:aws:workspaces-web:us-east-1:123456789012:ipAccessSettings/ip-1"),
	}

	assertions.Equal("arn:aws:workspaces-web:us-east-1:123456789012:ipAccessSettings/ip-1", ipAccessSettings.String())
}
