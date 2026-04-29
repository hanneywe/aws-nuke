package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	cognitotypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func Test_Mock_CognitoManagedLoginBranding_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockCognitoClient)

	mockClient.On("ListUserPools", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.ListUserPoolsOutput{
			UserPools: []cognitotypes.UserPoolDescriptionType{
				{Id: ptr.String("us-east-1_pool1")},
			},
		}, nil)

	mockClient.On("ListUserPoolClients", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.ListUserPoolClientsOutput{
			UserPoolClients: []cognitotypes.UserPoolClientDescription{
				{ClientId: ptr.String("client-123")},
			},
		}, nil)

	mockClient.On("DescribeManagedLoginBrandingByClient", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.DescribeManagedLoginBrandingByClientOutput{
			ManagedLoginBranding: &cognitotypes.ManagedLoginBrandingType{
				ManagedLoginBrandingId: ptr.String("branding-456"),
				UserPoolId:             ptr.String("us-east-1_pool1"),
			},
		}, nil)

	lister := &CognitoManagedLoginBrandingLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCognitoListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	branding := resources[0].(*CognitoManagedLoginBranding)
	assertions.Equal("branding-456", *branding.ManagedLoginBrandingID)
	assertions.Equal("us-east-1_pool1", *branding.UserPoolID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CognitoManagedLoginBranding_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockCognitoClient)

	mockClient.On("ListUserPools", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.ListUserPoolsOutput{
			UserPools: []cognitotypes.UserPoolDescriptionType{
				{Id: ptr.String("us-east-1_pool1")},
			},
		}, nil)

	mockClient.On("ListUserPoolClients", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.ListUserPoolClientsOutput{
			UserPoolClients: []cognitotypes.UserPoolClientDescription{},
		}, nil)

	lister := &CognitoManagedLoginBrandingLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCognitoListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CognitoManagedLoginBranding_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockCognitoClient)

	branding := &CognitoManagedLoginBranding{
		svc:                    mockClient,
		ManagedLoginBrandingID: ptr.String("branding-456"),
		UserPoolID:             ptr.String("us-east-1_pool1"),
	}

	mockClient.On("DeleteManagedLoginBranding", mock.Anything, &cognitoidentityprovider.DeleteManagedLoginBrandingInput{
		ManagedLoginBrandingId: branding.ManagedLoginBrandingID,
		UserPoolId:             branding.UserPoolID,
	}).Return(&cognitoidentityprovider.DeleteManagedLoginBrandingOutput{}, nil)

	err := branding.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CognitoManagedLoginBranding_Properties(t *testing.T) {
	assertions := assert.New(t)

	branding := CognitoManagedLoginBranding{
		ManagedLoginBrandingID: ptr.String("branding-456"),
		UserPoolID:             ptr.String("us-east-1_pool1"),
	}

	properties := branding.Properties()
	assertions.Equal("branding-456", properties.Get("ManagedLoginBrandingId"))
	assertions.Equal("us-east-1_pool1", properties.Get("UserPoolId"))
}

func Test_Mock_CognitoManagedLoginBranding_String(t *testing.T) {
	assertions := assert.New(t)
	branding := CognitoManagedLoginBranding{ManagedLoginBrandingID: ptr.String("branding-456")}
	assertions.Equal("branding-456", branding.String())
}
