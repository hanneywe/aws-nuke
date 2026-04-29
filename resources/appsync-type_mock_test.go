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

func Test_Mock_AppSyncType_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppSyncClient)
	mockClient.On("ListGraphqlApis", mock.Anything, mock.Anything).
		Return(&appsync.ListGraphqlApisOutput{
			GraphqlApis: []appsynctypes.GraphqlApi{
				{ApiId: ptr.String("api-1")},
			},
		}, nil)
	mockClient.On("ListTypes", mock.Anything, mock.Anything).
		Return(&appsync.ListTypesOutput{
			Types: []appsynctypes.Type{
				{Name: ptr.String("Query")},
			},
		}, nil)
	lister := &AppSyncTypeLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAppSyncListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("Query", *resources[0].(*AppSyncType).TypeName)
	a.Equal("api-1", *resources[0].(*AppSyncType).APIID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppSyncType_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppSyncClient)
	mockClient.On("ListGraphqlApis", mock.Anything, mock.Anything).
		Return(&appsync.ListGraphqlApisOutput{
			GraphqlApis: []appsynctypes.GraphqlApi{
				{ApiId: ptr.String("api-1")},
			},
		}, nil)
	mockClient.On("ListTypes", mock.Anything, mock.Anything).
		Return(&appsync.ListTypesOutput{Types: []appsynctypes.Type{}}, nil)
	lister := &AppSyncTypeLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAppSyncListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppSyncType_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppSyncClient)
	r := &AppSyncType{svc: mockClient, TypeName: ptr.String("Query"), APIID: ptr.String("api-1")}
	mockClient.On("DeleteType", mock.Anything, &appsync.DeleteTypeInput{ApiId: r.APIID, TypeName: r.TypeName}).
		Return(&appsync.DeleteTypeOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppSyncType_Properties(t *testing.T) {
	a := assert.New(t)
	r := AppSyncType{TypeName: ptr.String("Query"), APIID: ptr.String("api-1")}
	a.Equal("Query", r.Properties().Get("TypeName"))
	a.Equal("api-1", r.Properties().Get("ApiId"))
}

func Test_Mock_AppSyncType_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("Query", (&AppSyncType{TypeName: ptr.String("Query")}).String())
}
