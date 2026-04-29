package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/amplify"
	amplifytypes "github.com/aws/aws-sdk-go-v2/service/amplify/types"
)

func Test_Mock_AmplifyDomainAssociation_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAmplifyClient)

	mockClient.On("ListApps", mock.Anything, mock.Anything).
		Return(&amplify.ListAppsOutput{
			Apps: []amplifytypes.App{
				{AppId: ptr.String("app-123")},
			},
		}, nil)

	mockClient.On("ListDomainAssociations", mock.Anything, mock.Anything).
		Return(&amplify.ListDomainAssociationsOutput{
			DomainAssociations: []amplifytypes.DomainAssociation{
				{
					DomainName:           ptr.String("example.com"),
					DomainAssociationArn: ptr.String("arn:aws:amplify:us-east-1:123456789012:apps/app-123/domains/example.com"),
				},
			},
		}, nil)

	lister := &AmplifyDomainAssociationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAmplifyListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	d := resources[0].(*AmplifyDomainAssociation)
	a.Equal("example.com", *d.DomainName)
	a.Equal("app-123", *d.AppID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AmplifyDomainAssociation_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAmplifyClient)
	mockClient.On("ListApps", mock.Anything, mock.Anything).
		Return(&amplify.ListAppsOutput{Apps: []amplifytypes.App{}}, nil)
	lister := &AmplifyDomainAssociationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAmplifyListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AmplifyDomainAssociation_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAmplifyClient)
	d := &AmplifyDomainAssociation{svc: mockClient, AppID: ptr.String("app-123"), DomainName: ptr.String("example.com")}
	mockClient.On("DeleteDomainAssociation", mock.Anything, &amplify.DeleteDomainAssociationInput{
		AppId: d.AppID, DomainName: d.DomainName,
	}).Return(&amplify.DeleteDomainAssociationOutput{}, nil)
	a.NoError(d.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_AmplifyDomainAssociation_Properties(t *testing.T) {
	a := assert.New(t)
	d := AmplifyDomainAssociation{
		DomainName:           ptr.String("example.com"),
		DomainAssociationArn: ptr.String("arn:aws:amplify:us-east-1:123456789012:apps/app-123/domains/example.com"),
		AppID:                ptr.String("app-123"),
	}
	a.Equal("example.com", d.Properties().Get("DomainName"))
	a.Equal("app-123", d.Properties().Get("AppId"))
	a.Equal("arn:aws:amplify:us-east-1:123456789012:apps/app-123/domains/example.com", d.Properties().Get("DomainAssociationArn"))
}

func Test_Mock_AmplifyDomainAssociation_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("example.com", (&AmplifyDomainAssociation{DomainName: ptr.String("example.com")}).String())
}
