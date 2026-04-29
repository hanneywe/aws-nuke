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

func Test_Mock_ConnectQuickConnect_List_One(t *testing.T) {
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
		On("ListQuickConnects", mock.Anything, mock.Anything).
		Return(&connect.ListQuickConnectsOutput{
			QuickConnectSummaryList: []connecttypes.QuickConnectSummary{
				{
					Id:   ptr.String("qc-12345"),
					Name: ptr.String("my-quick-connect"),
				},
			},
		}, nil)

	lister := &ConnectQuickConnectLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	qc := resources[0].(*ConnectQuickConnect)
	a.Equal("qc-12345", *qc.ID)
	a.Equal("my-quick-connect", *qc.Name)
	a.Equal("i-12345", *qc.InstanceID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectQuickConnect_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	mockClient.
		On("ListInstances", mock.Anything, mock.Anything).
		Return(&connect.ListInstancesOutput{
			InstanceSummaryList: []connecttypes.InstanceSummary{},
		}, nil)

	lister := &ConnectQuickConnectLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectQuickConnect_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	qc := &ConnectQuickConnect{
		svc:        mockClient,
		InstanceID: ptr.String("i-12345"),
		ID:         ptr.String("qc-12345"),
		Name:       ptr.String("my-quick-connect"),
	}

	mockClient.
		On("DeleteQuickConnect", mock.Anything, &connect.DeleteQuickConnectInput{
			InstanceId:     qc.InstanceID,
			QuickConnectId: qc.ID,
		}).
		Return(&connect.DeleteQuickConnectOutput{}, nil)

	err := qc.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectQuickConnect_Properties(t *testing.T) {
	a := assert.New(t)

	qc := ConnectQuickConnect{
		InstanceID: ptr.String("i-12345"),
		ID:         ptr.String("qc-12345"),
		Name:       ptr.String("my-quick-connect"),
	}

	props := qc.Properties()
	a.Equal("i-12345", props.Get("InstanceId"))
	a.Equal("qc-12345", props.Get("Id"))
	a.Equal("my-quick-connect", props.Get("Name"))
}

func Test_Mock_ConnectQuickConnect_String(t *testing.T) {
	a := assert.New(t)

	qc := ConnectQuickConnect{
		Name: ptr.String("my-quick-connect"),
	}

	a.Equal("my-quick-connect", qc.String())
}
