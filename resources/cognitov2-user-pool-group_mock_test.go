package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	cognitoidentityprovidertypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func Test_Mock_CognitoUserPoolGroup_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCognitoidentityproviderClient)

	mockClient.On("ListUserPools", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.ListUserPoolsOutput{
			UserPools: []cognitoidentityprovidertypes.UserPoolDescriptionType{
				{Id: ptr.String("test-userpoolid")},
			},
		}, nil)

	mockClient.On("ListGroups", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.ListGroupsOutput{
			Groups: []cognitoidentityprovidertypes.GroupType{
				{GroupName: ptr.String("test-groupname"), Description: ptr.String("test-description")},
			},
		}, nil)

	lister := &CognitoUserPoolGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCognitoidentityproviderListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*CognitoUserPoolGroup)
	a.Equal("test-userpoolid", *r.UserPoolID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CognitoUserPoolGroup_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCognitoidentityproviderClient)

	mockClient.On("ListUserPools", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.ListUserPoolsOutput{
			UserPools: []cognitoidentityprovidertypes.UserPoolDescriptionType{},
		}, nil)

	lister := &CognitoUserPoolGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCognitoidentityproviderListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CognitoUserPoolGroup_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCognitoidentityproviderClient)

	r := &CognitoUserPoolGroup{
		svc:        mockClient,
		GroupName:  ptr.String("test-groupname"),
		UserPoolID: ptr.String("test-userpoolid"),
	}

	mockClient.On("DeleteGroup", mock.Anything,
		&cognitoidentityprovider.DeleteGroupInput{
			GroupName:  r.GroupName,
			UserPoolId: r.UserPoolID,
		}).Return(&cognitoidentityprovider.DeleteGroupOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_CognitoUserPoolGroup_Properties(t *testing.T) {
	a := assert.New(t)
	r := &CognitoUserPoolGroup{
		UserPoolID:  ptr.String("test-userpoolid"),
		GroupName:   ptr.String("test-groupname"),
		Description: ptr.String("test-description"),
	}
	props := r.Properties()
	a.Equal("test-userpoolid", props.Get("UserPoolId"))
	a.Equal("test-groupname", props.Get("GroupName"))
	a.Equal("test-description", props.Get("Description"))
}

func Test_Mock_CognitoUserPoolGroup_String(t *testing.T) {
	a := assert.New(t)
	r := &CognitoUserPoolGroup{
		GroupName: ptr.String("test-groupname"),
	}
	a.Equal("test-groupname", r.String())
}
