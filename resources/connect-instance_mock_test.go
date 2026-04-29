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

func Test_Mock_ConnectInstance_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	mockClient.
		On("ListInstances", mock.Anything, mock.Anything).
		Return(&connect.ListInstancesOutput{
			InstanceSummaryList: []connecttypes.InstanceSummary{
				{
					Id:            ptr.String("i-12345"),
					InstanceAlias: ptr.String("my-instance"),
				},
			},
		}, nil)

	lister := &ConnectInstanceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	instance := resources[0].(*ConnectInstance)
	a.Equal("my-instance", *instance.InstanceAlias)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectInstance_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	mockClient.
		On("ListInstances", mock.Anything, mock.Anything).
		Return(&connect.ListInstancesOutput{
			InstanceSummaryList: []connecttypes.InstanceSummary{},
		}, nil)

	lister := &ConnectInstanceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectInstance_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	instance := &ConnectInstance{
		svc: mockClient,
		ID:  ptr.String("i-12345"),
	}

	mockClient.
		On("DeleteInstance", mock.Anything, &connect.DeleteInstanceInput{
			InstanceId: instance.ID,
		}).
		Return(&connect.DeleteInstanceOutput{}, nil)

	err := instance.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectInstance_Properties(t *testing.T) {
	a := assert.New(t)

	instance := ConnectInstance{
		ID:            ptr.String("i-12345"),
		InstanceAlias: ptr.String("my-instance"),
	}

	props := instance.Properties()
	a.Equal("i-12345", props.Get("Id"))
	a.Equal("my-instance", props.Get("InstanceAlias"))
}

func Test_Mock_ConnectInstance_String(t *testing.T) {
	a := assert.New(t)

	instance := ConnectInstance{
		ID: ptr.String("i-12345"),
	}

	a.Equal("i-12345", instance.String())
}
