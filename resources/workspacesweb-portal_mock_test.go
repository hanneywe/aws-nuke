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

func Test_Mock_WorkSpacesWebPortal_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockWorkSpacesWebClient)

	mockClient.
		On("ListPortals", mock.Anything, mock.Anything).
		Return(&workspacesweb.ListPortalsOutput{
			Portals: []workspaceswebtypes.PortalSummary{
				{
					PortalArn:   ptr.String("arn:aws:workspaces-web:us-east-1:123456789012:portal/portal-1"),
					DisplayName: ptr.String("my-portal"),
				},
			},
		}, nil)

	lister := &WorkSpacesWebPortalLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testWorkSpacesWebListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	portal := resources[0].(*WorkSpacesWebPortal)
	assertions.Equal("arn:aws:workspaces-web:us-east-1:123456789012:portal/portal-1", *portal.PortalArn)
	assertions.Equal("my-portal", *portal.DisplayName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_WorkSpacesWebPortal_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockWorkSpacesWebClient)

	mockClient.
		On("ListPortals", mock.Anything, mock.Anything).
		Return(&workspacesweb.ListPortalsOutput{
			Portals: []workspaceswebtypes.PortalSummary{},
		}, nil)

	lister := &WorkSpacesWebPortalLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testWorkSpacesWebListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_WorkSpacesWebPortal_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockWorkSpacesWebClient)

	portal := &WorkSpacesWebPortal{
		svc:         mockClient,
		PortalArn:   ptr.String("arn:aws:workspaces-web:us-east-1:123456789012:portal/portal-1"),
		DisplayName: ptr.String("my-portal"),
	}

	mockClient.
		On("DeletePortal", mock.Anything, &workspacesweb.DeletePortalInput{
			PortalArn: portal.PortalArn,
		}).
		Return(&workspacesweb.DeletePortalOutput{}, nil)

	err := portal.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_WorkSpacesWebPortal_Properties(t *testing.T) {
	assertions := assert.New(t)

	portal := WorkSpacesWebPortal{
		PortalArn:   ptr.String("arn:aws:workspaces-web:us-east-1:123456789012:portal/portal-1"),
		DisplayName: ptr.String("my-portal"),
	}

	properties := portal.Properties()

	assertions.Equal("arn:aws:workspaces-web:us-east-1:123456789012:portal/portal-1", properties.Get("PortalArn"))
	assertions.Equal("my-portal", properties.Get("DisplayName"))
}

func Test_Mock_WorkSpacesWebPortal_String(t *testing.T) {
	assertions := assert.New(t)

	portal := WorkSpacesWebPortal{
		PortalArn: ptr.String("arn:aws:workspaces-web:us-east-1:123456789012:portal/portal-1"),
	}

	assertions.Equal("arn:aws:workspaces-web:us-east-1:123456789012:portal/portal-1", portal.String())
}
