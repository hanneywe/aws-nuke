package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/appsync"
	appsynctypes "github.com/aws/aws-sdk-go-v2/service/appsync/types"
	"github.com/aws/smithy-go"
)

func Test_Mock_AppSyncAPICache_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppSyncClient)
	mockClient.On("ListGraphqlApis", mock.Anything, mock.Anything).
		Return(&appsync.ListGraphqlApisOutput{
			GraphqlApis: []appsynctypes.GraphqlApi{
				{ApiId: ptr.String("api-1")},
			},
		}, nil)
	mockClient.On("GetApiCache", mock.Anything, &appsync.GetApiCacheInput{ApiId: ptr.String("api-1")}).
		Return(&appsync.GetApiCacheOutput{ApiCache: &appsynctypes.ApiCache{}}, nil)
	lister := &AppSyncAPICacheLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAppSyncListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("api-1", *resources[0].(*AppSyncAPICache).APIID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppSyncAPICache_List_NoCache(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppSyncClient)
	mockClient.On("ListGraphqlApis", mock.Anything, mock.Anything).
		Return(&appsync.ListGraphqlApisOutput{
			GraphqlApis: []appsynctypes.GraphqlApi{
				{ApiId: ptr.String("api-1")},
			},
		}, nil)
	mockClient.On("GetApiCache", mock.Anything, mock.Anything).
		Return(&appsync.GetApiCacheOutput{}, &smithy.GenericAPIError{Code: "NotFoundException", Message: "not found"})
	lister := &AppSyncAPICacheLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAppSyncListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppSyncAPICache_List_Error(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppSyncClient)
	mockClient.On("ListGraphqlApis", mock.Anything, mock.Anything).
		Return(&appsync.ListGraphqlApisOutput{
			GraphqlApis: []appsynctypes.GraphqlApi{
				{ApiId: ptr.String("api-1")},
			},
		}, nil)
	mockClient.On("GetApiCache", mock.Anything, mock.Anything).
		Return(&appsync.GetApiCacheOutput{}, fmt.Errorf("some error"))
	lister := &AppSyncAPICacheLister{svc: mockClient}
	_, err := lister.List(context.TODO(), testAppSyncListerOpts)
	a.Error(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppSyncAPICache_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppSyncClient)
	r := &AppSyncAPICache{svc: mockClient, APIID: ptr.String("api-1")}
	mockClient.On("DeleteApiCache", mock.Anything, &appsync.DeleteApiCacheInput{ApiId: r.APIID}).
		Return(&appsync.DeleteApiCacheOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppSyncAPICache_Properties(t *testing.T) {
	a := assert.New(t)
	r := AppSyncAPICache{APIID: ptr.String("api-1")}
	a.Equal("api-1", r.Properties().Get("ApiId"))
}

func Test_Mock_AppSyncAPICache_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("api-1", (&AppSyncAPICache{APIID: ptr.String("api-1")}).String())
}
