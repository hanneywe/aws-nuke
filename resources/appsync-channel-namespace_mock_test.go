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

func Test_Mock_AppSyncChannelNamespace_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppSyncClient)

	mockClient.On("ListApis", mock.Anything, mock.Anything).
		Return(&appsync.ListApisOutput{
			Apis: []appsynctypes.Api{
				{ApiId: ptr.String("api-1")},
			},
		}, nil)

	mockClient.On("ListChannelNamespaces", mock.Anything, mock.Anything).
		Return(&appsync.ListChannelNamespacesOutput{
			ChannelNamespaces: []appsynctypes.ChannelNamespace{
				{Name: ptr.String("test-ns"), Tags: map[string]string{"env": "test"}},
			},
		}, nil)

	lister := &AppSyncChannelNamespaceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAppSyncListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*AppSyncChannelNamespace)
	a.Equal("api-1", *r.APIID)
	a.Equal("test-ns", *r.Name)
	a.Equal("test", r.Tags["env"])
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppSyncChannelNamespace_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppSyncClient)

	mockClient.On("ListApis", mock.Anything, mock.Anything).
		Return(&appsync.ListApisOutput{
			Apis: []appsynctypes.Api{},
		}, nil)

	lister := &AppSyncChannelNamespaceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAppSyncListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppSyncChannelNamespace_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppSyncClient)

	r := &AppSyncChannelNamespace{
		svc:   mockClient,
		APIID: ptr.String("api-1"),
		Name:  ptr.String("test-ns"),
	}

	mockClient.On("DeleteChannelNamespace", mock.Anything,
		&appsync.DeleteChannelNamespaceInput{
			ApiId: r.APIID,
			Name:  r.Name,
		}).Return(&appsync.DeleteChannelNamespaceOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppSyncChannelNamespace_Properties(t *testing.T) {
	a := assert.New(t)
	r := &AppSyncChannelNamespace{
		APIID: ptr.String("api-1"),
		Name:  ptr.String("test-ns"),
	}
	props := r.Properties()
	a.Equal("api-1", props.Get("ApiId"))
	a.Equal("test-ns", props.Get("Name"))
}

func Test_Mock_AppSyncChannelNamespace_String(t *testing.T) {
	a := assert.New(t)
	r := &AppSyncChannelNamespace{
		Name: ptr.String("test-ns"),
	}
	a.Equal("test-ns", r.String())
}
