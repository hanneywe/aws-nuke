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

func Test_Mock_CognitoResourceServer_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCognitoClient)
	mockClient.On("ListUserPools", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.ListUserPoolsOutput{
			UserPools: []cognitotypes.UserPoolDescriptionType{
				{Id: ptr.String("us-east-1_abc123")},
			},
		}, nil)
	mockClient.On("ListResourceServers", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.ListResourceServersOutput{
			ResourceServers: []cognitotypes.ResourceServerType{
				{Identifier: ptr.String("my-api"), Name: ptr.String("My API")},
			},
		}, nil)
	lister := &CognitoResourceServerLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCognitoListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	rs := resources[0].(*CognitoResourceServer)
	a.Equal("my-api", *rs.Identifier)
	a.Equal("us-east-1_abc123", *rs.UserPoolID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CognitoResourceServer_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCognitoClient)
	mockClient.On("ListUserPools", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.ListUserPoolsOutput{
			UserPools: []cognitotypes.UserPoolDescriptionType{
				{Id: ptr.String("us-east-1_abc123")},
			},
		}, nil)
	mockClient.On("ListResourceServers", mock.Anything, mock.Anything).
		Return(&cognitoidentityprovider.ListResourceServersOutput{ResourceServers: []cognitotypes.ResourceServerType{}}, nil)
	lister := &CognitoResourceServerLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCognitoListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CognitoResourceServer_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCognitoClient)
	r := &CognitoResourceServer{svc: mockClient, Identifier: ptr.String("my-api"), UserPoolID: ptr.String("us-east-1_abc123")}
	mockClient.On("DeleteResourceServer", mock.Anything, &cognitoidentityprovider.DeleteResourceServerInput{
		Identifier: r.Identifier, UserPoolId: r.UserPoolID,
	}).Return(&cognitoidentityprovider.DeleteResourceServerOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_CognitoResourceServer_Properties(t *testing.T) {
	a := assert.New(t)
	r := CognitoResourceServer{Identifier: ptr.String("my-api"), Name: ptr.String("My API"), UserPoolID: ptr.String("pool-1")}
	a.Equal("my-api", r.Properties().Get("Identifier"))
	a.Equal("My API", r.Properties().Get("Name"))
}

func Test_Mock_CognitoResourceServer_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("My API", (&CognitoResourceServer{Name: ptr.String("My API")}).String())
}
