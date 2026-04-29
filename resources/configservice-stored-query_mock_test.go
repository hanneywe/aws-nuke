package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/configservice"
	configtypes "github.com/aws/aws-sdk-go-v2/service/configservice/types"
)

func Test_Mock_ConfigServiceStoredQuery_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConfigServiceClient)
	mockClient.On("ListStoredQueries", mock.Anything, mock.Anything).
		Return(&configservice.ListStoredQueriesOutput{
			StoredQueryMetadata: []configtypes.StoredQueryMetadata{
				{
					QueryId:   ptr.String("query-abc123"),
					QueryName: ptr.String("my-query"),
					QueryArn:  ptr.String("arn:aws:config:us-east-1:123456789012:stored-query/my-query"),
				},
			},
		}, nil)
	lister := &ConfigServiceStoredQueryLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConfigServiceListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("my-query", resources[0].(*ConfigServiceStoredQuery).String())
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConfigServiceStoredQuery_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConfigServiceClient)
	mockClient.On("ListStoredQueries", mock.Anything, mock.Anything).
		Return(&configservice.ListStoredQueriesOutput{
			StoredQueryMetadata: []configtypes.StoredQueryMetadata{},
		}, nil)
	lister := &ConfigServiceStoredQueryLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConfigServiceListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConfigServiceStoredQuery_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConfigServiceClient)
	r := &ConfigServiceStoredQuery{
		svc:       mockClient,
		QueryID:   ptr.String("query-abc123"),
		QueryName: ptr.String("my-query"),
	}
	mockClient.On("DeleteStoredQuery", mock.Anything, &configservice.DeleteStoredQueryInput{
		QueryName: r.QueryName,
	}).Return(&configservice.DeleteStoredQueryOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConfigServiceStoredQuery_Properties(t *testing.T) {
	a := assert.New(t)
	r := ConfigServiceStoredQuery{
		QueryID:   ptr.String("query-abc123"),
		QueryName: ptr.String("my-query"),
	}
	a.Equal("query-abc123", r.Properties().Get("QueryId"))
	a.Equal("my-query", r.Properties().Get("QueryName"))
}

func Test_Mock_ConfigServiceStoredQuery_String(t *testing.T) {
	a := assert.New(t)
	r := &ConfigServiceStoredQuery{
		QueryID:   ptr.String("query-abc123"),
		QueryName: ptr.String("my-query"),
	}
	a.Equal("my-query", r.String())
}
