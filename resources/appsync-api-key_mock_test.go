package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/appsync"
	appsynctypes "github.com/aws/aws-sdk-go-v2/service/appsync/types"
)

func Test_Mock_AppSyncAPIKey_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppSyncClient)
	mockClient.On("ListGraphqlApis", mock.Anything, mock.Anything).
		Return(&appsync.ListGraphqlApisOutput{
			GraphqlApis: []appsynctypes.GraphqlApi{
				{ApiId: ptr.String("api-1")},
			},
		}, nil)
	mockClient.On("ListApiKeys", mock.Anything, mock.Anything).
		Return(&appsync.ListApiKeysOutput{
			ApiKeys: []appsynctypes.ApiKey{
				{Id: ptr.String("key-1")},
			},
		}, nil)
	lister := &AppSyncAPIKeyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAppSyncListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("key-1", *resources[0].(*AppSyncAPIKey).ID)
	a.Equal("api-1", *resources[0].(*AppSyncAPIKey).APIID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppSyncAPIKey_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppSyncClient)
	mockClient.On("ListGraphqlApis", mock.Anything, mock.Anything).
		Return(&appsync.ListGraphqlApisOutput{
			GraphqlApis: []appsynctypes.GraphqlApi{
				{ApiId: ptr.String("api-1")},
			},
		}, nil)
	mockClient.On("ListApiKeys", mock.Anything, mock.Anything).
		Return(&appsync.ListApiKeysOutput{ApiKeys: []appsynctypes.ApiKey{}}, nil)
	lister := &AppSyncAPIKeyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAppSyncListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppSyncAPIKey_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppSyncClient)
	r := &AppSyncAPIKey{svc: mockClient, ID: ptr.String("key-1"), APIID: ptr.String("api-1")}
	mockClient.On("DeleteApiKey", mock.Anything, &appsync.DeleteApiKeyInput{ApiId: r.APIID, Id: r.ID}).
		Return(&appsync.DeleteApiKeyOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppSyncAPIKey_Properties(t *testing.T) {
	a := assert.New(t)
	r := AppSyncAPIKey{ID: ptr.String("key-1"), APIID: ptr.String("api-1")}
	a.Equal("key-1", r.Properties().Get("Id"))
	a.Equal("api-1", r.Properties().Get("ApiId"))
}

func Test_Mock_AppSyncAPIKey_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("key-1", (&AppSyncAPIKey{ID: ptr.String("key-1")}).String())
}
