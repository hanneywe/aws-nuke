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

func Test_Mock_ConnectAgentStatus_List_One(t *testing.T) {
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
		On("ListAgentStatuses", mock.Anything, mock.Anything).
		Return(&connect.ListAgentStatusesOutput{
			AgentStatusSummaryList: []connecttypes.AgentStatusSummary{
				{
					Id:   ptr.String("as-12345"),
					Name: ptr.String("Available"),
					Type: connecttypes.AgentStatusTypeRoutable,
				},
			},
		}, nil)

	lister := &ConnectAgentStatusLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	status := resources[0].(*ConnectAgentStatus)
	a.Equal("as-12345", *status.ID)
	a.Equal("Available", *status.Name)
	a.Equal("i-12345", *status.InstanceID)
	a.Equal(connecttypes.AgentStatusTypeRoutable, status.Type)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectAgentStatus_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	mockClient.
		On("ListInstances", mock.Anything, mock.Anything).
		Return(&connect.ListInstancesOutput{
			InstanceSummaryList: []connecttypes.InstanceSummary{},
		}, nil)

	lister := &ConnectAgentStatusLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectAgentStatus_Remove(t *testing.T) {
	a := assert.New(t)

	status := &ConnectAgentStatus{
		InstanceID: ptr.String("i-12345"),
		ID:         ptr.String("as-12345"),
		Name:       ptr.String("Available"),
		Type:       connecttypes.AgentStatusTypeRoutable,
	}

	err := status.Remove(context.TODO())
	a.Error(err)
	a.Contains(err.Error(), "cannot be deleted")
}

func Test_Mock_ConnectAgentStatus_Filter_Default(t *testing.T) {
	a := assert.New(t)

	status := &ConnectAgentStatus{
		Name: ptr.String("Available"),
		Type: connecttypes.AgentStatusTypeRoutable,
	}

	err := status.Filter()
	a.Error(err)
	a.Contains(err.Error(), "cannot delete default agent status")
}

func Test_Mock_ConnectAgentStatus_Filter_Custom(t *testing.T) {
	a := assert.New(t)

	status := &ConnectAgentStatus{
		Name: ptr.String("Break"),
		Type: connecttypes.AgentStatusTypeCustom,
	}

	err := status.Filter()
	a.NoError(err)
}

func Test_Mock_ConnectAgentStatus_Properties(t *testing.T) {
	a := assert.New(t)

	status := ConnectAgentStatus{
		InstanceID: ptr.String("i-12345"),
		ID:         ptr.String("as-12345"),
		Name:       ptr.String("Available"),
		Type:       connecttypes.AgentStatusTypeRoutable,
	}

	props := status.Properties()
	a.Equal("i-12345", props.Get("InstanceId"))
	a.Equal("as-12345", props.Get("Id"))
	a.Equal("Available", props.Get("Name"))
}

func Test_Mock_ConnectAgentStatus_String(t *testing.T) {
	a := assert.New(t)

	status := ConnectAgentStatus{
		Name: ptr.String("Available"),
	}

	a.Equal("Available", status.String())
}
