package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/connect"
	connecttypes "github.com/aws/aws-sdk-go-v2/service/connect/types"
)

func Test_Mock_ConnectWorkspace_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	mockClient.
		On("ListInstances", mock.Anything, mock.Anything).
		Return(&connect.ListInstancesOutput{
			InstanceSummaryList: []connecttypes.InstanceSummary{
				{Id: ptr.String("i-12345")},
			},
		}, nil)

	mockClient.
		On("ListWorkspaces", mock.Anything, mock.Anything).
		Return(&connect.ListWorkspacesOutput{
			WorkspaceSummaryList: []connecttypes.WorkspaceSummary{
				{
					Id:   ptr.String("ws-12345"),
					Name: ptr.String("my-workspace"),
				},
			},
		}, nil)

	lister := &ConnectWorkspaceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	ws := resources[0].(*ConnectWorkspace)
	a.Equal("ws-12345", *ws.ID)
	a.Equal("my-workspace", *ws.Name)
	a.Equal("i-12345", *ws.InstanceID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectWorkspace_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	mockClient.
		On("ListInstances", mock.Anything, mock.Anything).
		Return(&connect.ListInstancesOutput{
			InstanceSummaryList: []connecttypes.InstanceSummary{},
		}, nil)

	lister := &ConnectWorkspaceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectWorkspace_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	ws := &ConnectWorkspace{
		svc:        mockClient,
		InstanceID: ptr.String("i-12345"),
		ID:         ptr.String("ws-12345"),
		Name:       ptr.String("my-workspace"),
	}

	mockClient.
		On("DeleteWorkspace", mock.Anything, &connect.DeleteWorkspaceInput{
			InstanceId:  ws.InstanceID,
			WorkspaceId: ws.ID,
		}).
		Return(&connect.DeleteWorkspaceOutput{}, nil)

	err := ws.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectWorkspace_Properties(t *testing.T) {
	a := assert.New(t)

	ws := ConnectWorkspace{
		InstanceID: ptr.String("i-12345"),
		ID:         ptr.String("ws-12345"),
		Name:       ptr.String("my-workspace"),
	}

	props := ws.Properties()
	a.Equal("i-12345", props.Get("InstanceId"))
	a.Equal("ws-12345", props.Get("Id"))
	a.Equal("my-workspace", props.Get("Name"))
}

func Test_Mock_ConnectWorkspace_String(t *testing.T) {
	a := assert.New(t)

	ws := ConnectWorkspace{
		ID: ptr.String("ws-12345"),
	}

	a.Equal("ws-12345", ws.String())
}
