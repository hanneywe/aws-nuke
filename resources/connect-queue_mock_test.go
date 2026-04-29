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

func Test_Mock_ConnectQueue_List_One(t *testing.T) {
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
		On("ListQueues", mock.Anything, mock.Anything).
		Return(&connect.ListQueuesOutput{
			QueueSummaryList: []connecttypes.QueueSummary{
				{
					Id:        ptr.String("q-12345"),
					Name:      ptr.String("BasicQueue"),
					QueueType: connecttypes.QueueTypeStandard,
				},
			},
		}, nil)

	lister := &ConnectQueueLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	queue := resources[0].(*ConnectQueue)
	a.Equal("q-12345", *queue.ID)
	a.Equal("BasicQueue", *queue.Name)
	a.Equal("i-12345", *queue.InstanceID)
	a.Equal(connecttypes.QueueTypeStandard, queue.QueueType)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectQueue_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	mockClient.
		On("ListInstances", mock.Anything, mock.Anything).
		Return(&connect.ListInstancesOutput{
			InstanceSummaryList: []connecttypes.InstanceSummary{},
		}, nil)

	lister := &ConnectQueueLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectQueue_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	queue := &ConnectQueue{
		svc:        mockClient,
		InstanceID: ptr.String("i-12345"),
		ID:         ptr.String("q-12345"),
		Name:       ptr.String("BasicQueue"),
		QueueType:  connecttypes.QueueTypeStandard,
	}

	mockClient.
		On("DeleteQueue", mock.Anything, &connect.DeleteQueueInput{
			InstanceId: queue.InstanceID,
			QueueId:    queue.ID,
		}).
		Return(&connect.DeleteQueueOutput{}, nil)

	err := queue.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectQueue_Filter_Standard(t *testing.T) {
	a := assert.New(t)

	queue := &ConnectQueue{
		Name:      ptr.String("BasicQueue"),
		QueueType: connecttypes.QueueTypeStandard,
	}

	err := queue.Filter()
	a.NoError(err)
}

func Test_Mock_ConnectQueue_Filter_Agent(t *testing.T) {
	a := assert.New(t)

	queue := &ConnectQueue{
		Name:      ptr.String("AgentQueue"),
		QueueType: connecttypes.QueueTypeAgent,
	}

	err := queue.Filter()
	a.Error(err)
	a.Contains(err.Error(), "cannot delete non-standard queue")
}

func Test_Mock_ConnectQueue_Properties(t *testing.T) {
	a := assert.New(t)

	queue := ConnectQueue{
		InstanceID: ptr.String("i-12345"),
		ID:         ptr.String("q-12345"),
		Name:       ptr.String("BasicQueue"),
		QueueType:  connecttypes.QueueTypeStandard,
	}

	props := queue.Properties()
	a.Equal("i-12345", props.Get("InstanceId"))
	a.Equal("q-12345", props.Get("Id"))
	a.Equal("BasicQueue", props.Get("Name"))
}

func Test_Mock_ConnectQueue_String(t *testing.T) {
	a := assert.New(t)

	queue := ConnectQueue{
		Name: ptr.String("BasicQueue"),
	}

	a.Equal("BasicQueue", queue.String())
}
