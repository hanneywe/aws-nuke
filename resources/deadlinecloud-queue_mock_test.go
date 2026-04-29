package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/deadline"
	deadlinetypes "github.com/aws/aws-sdk-go-v2/service/deadline/types"
)

func Test_Mock_DeadlineCloudQueue_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDeadlineCloudClient)

	mockClient.
		On("ListFarms", mock.Anything, mock.Anything).
		Return(&deadline.ListFarmsOutput{
			Farms: []deadlinetypes.FarmSummary{
				{
					FarmId:      ptr.String("farm-12345"),
					DisplayName: ptr.String("my-farm"),
				},
			},
		}, nil)

	mockClient.
		On("ListQueues", mock.Anything, mock.Anything).
		Return(&deadline.ListQueuesOutput{
			Queues: []deadlinetypes.QueueSummary{
				{
					FarmId:      ptr.String("farm-12345"),
					QueueId:     ptr.String("queue-12345"),
					DisplayName: ptr.String("my-queue"),
				},
			},
		}, nil)

	lister := &DeadlineCloudQueueLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDeadlineCloudListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	queue := resources[0].(*DeadlineCloudQueue)
	a.Equal("queue-12345", *queue.QueueID)
	a.Equal("my-queue", *queue.DisplayName)
	a.Equal("farm-12345", *queue.FarmID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DeadlineCloudQueue_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDeadlineCloudClient)

	mockClient.
		On("ListFarms", mock.Anything, mock.Anything).
		Return(&deadline.ListFarmsOutput{
			Farms: []deadlinetypes.FarmSummary{},
		}, nil)

	lister := &DeadlineCloudQueueLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDeadlineCloudListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DeadlineCloudQueue_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDeadlineCloudClient)

	queue := &DeadlineCloudQueue{
		svc:     mockClient,
		FarmID:  ptr.String("farm-12345"),
		QueueID: ptr.String("queue-12345"),
	}

	mockClient.
		On("DeleteQueue", mock.Anything, &deadline.DeleteQueueInput{
			FarmId:  queue.FarmID,
			QueueId: queue.QueueID,
		}).
		Return(&deadline.DeleteQueueOutput{}, nil)

	err := queue.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DeadlineCloudQueue_Properties(t *testing.T) {
	a := assert.New(t)

	queue := DeadlineCloudQueue{
		FarmID:      ptr.String("farm-12345"),
		QueueID:     ptr.String("queue-12345"),
		DisplayName: ptr.String("my-queue"),
	}

	props := queue.Properties()
	a.Equal("farm-12345", props.Get("FarmId"))
	a.Equal("queue-12345", props.Get("QueueId"))
	a.Equal("my-queue", props.Get("DisplayName"))
}

func Test_Mock_DeadlineCloudQueue_String(t *testing.T) {
	a := assert.New(t)

	queue := DeadlineCloudQueue{
		QueueID: ptr.String("queue-12345"),
	}

	a.Equal("queue-12345", queue.String())
}
