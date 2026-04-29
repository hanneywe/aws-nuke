package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/deadline"
	deadlinetypes "github.com/aws/aws-sdk-go-v2/service/deadline/types"
	"github.com/aws/smithy-go"

	liberrors "github.com/ekristen/libnuke/pkg/errors"
)

func Test_Mock_DeadlineCloudQueueLimitAssociation_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDeadlineCloudClient)

	mockClient.
		On("ListFarms", mock.Anything, mock.Anything).
		Return(&deadline.ListFarmsOutput{
			Farms: []deadlinetypes.FarmSummary{
				{FarmId: ptr.String("farm-12345"), DisplayName: ptr.String("my-farm")},
			},
		}, nil)

	mockClient.
		On("ListQueueLimitAssociations", mock.Anything, mock.Anything).
		Return(&deadline.ListQueueLimitAssociationsOutput{
			QueueLimitAssociations: []deadlinetypes.QueueLimitAssociationSummary{
				{
					QueueId: ptr.String("queue-12345"),
					LimitId: ptr.String("limit-12345"),
					Status:  deadlinetypes.QueueLimitAssociationStatusActive,
				},
			},
		}, nil)

	lister := &DeadlineCloudQueueLimitAssociationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDeadlineCloudListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	assoc := resources[0].(*DeadlineCloudQueueLimitAssociation)
	a.Equal("queue-12345", *assoc.QueueID)
	a.Equal("limit-12345", *assoc.LimitID)
	a.Equal("farm-12345", *assoc.FarmID)
	a.Equal("ACTIVE", assoc.Status)
	mockClient.AssertExpectations(t)
}

func Test_Mock_DeadlineCloudQueueLimitAssociation_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDeadlineCloudClient)

	mockClient.
		On("ListFarms", mock.Anything, mock.Anything).
		Return(&deadline.ListFarmsOutput{
			Farms: []deadlinetypes.FarmSummary{},
		}, nil)

	lister := &DeadlineCloudQueueLimitAssociationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDeadlineCloudListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func mockGetAssociation(mockClient *mockDeadlineCloudClient, status deadlinetypes.QueueLimitAssociationStatus) {
	mockClient.
		On("GetQueueLimitAssociation", mock.Anything, mock.Anything).
		Return(&deadline.GetQueueLimitAssociationOutput{
			Status: status,
		}, nil)
}

func Test_Mock_DeadlineCloudQueueLimitAssociation_Remove_Active(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDeadlineCloudClient)

	assoc := &DeadlineCloudQueueLimitAssociation{
		svc:     mockClient,
		FarmID:  ptr.String("farm-12345"),
		QueueID: ptr.String("queue-12345"),
		LimitID: ptr.String("limit-12345"),
		Status:  "ACTIVE",
	}

	mockGetAssociation(mockClient, deadlinetypes.QueueLimitAssociationStatusActive)

	mockClient.
		On("UpdateQueueLimitAssociation", mock.Anything, &deadline.UpdateQueueLimitAssociationInput{
			FarmId:  assoc.FarmID,
			QueueId: assoc.QueueID,
			LimitId: assoc.LimitID,
			Status:  deadlinetypes.UpdateQueueLimitAssociationStatusStopLimitUsageAndCancelTasks,
		}).
		Return(&deadline.UpdateQueueLimitAssociationOutput{}, nil)

	err := assoc.Remove(context.TODO())
	a.Error(err)
	var errHold liberrors.ErrHoldResource
	a.ErrorAs(err, &errHold)
	a.Equal(string(deadlinetypes.QueueLimitAssociationStatusStopLimitUsageAndCancelTasks), assoc.Status)
	mockClient.AssertExpectations(t)
}

func Test_Mock_DeadlineCloudQueueLimitAssociation_Remove_Stopping(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDeadlineCloudClient)

	assoc := &DeadlineCloudQueueLimitAssociation{
		svc:     mockClient,
		FarmID:  ptr.String("farm-12345"),
		QueueID: ptr.String("queue-12345"),
		LimitID: ptr.String("limit-12345"),
		Status:  "ACTIVE",
	}

	// API says it's still stopping
	mockGetAssociation(mockClient, deadlinetypes.QueueLimitAssociationStatusStopLimitUsageAndCancelTasks)

	err := assoc.Remove(context.TODO())
	a.Error(err)
	var errHold liberrors.ErrHoldResource
	a.ErrorAs(err, &errHold)
	mockClient.AssertExpectations(t)
}

func Test_Mock_DeadlineCloudQueueLimitAssociation_Remove_Stopped(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDeadlineCloudClient)

	assoc := &DeadlineCloudQueueLimitAssociation{
		svc:     mockClient,
		FarmID:  ptr.String("farm-12345"),
		QueueID: ptr.String("queue-12345"),
		LimitID: ptr.String("limit-12345"),
		Status:  "ACTIVE",
	}

	// API says it's stopped now
	mockGetAssociation(mockClient, deadlinetypes.QueueLimitAssociationStatusStopped)

	mockClient.
		On("DeleteQueueLimitAssociation", mock.Anything, &deadline.DeleteQueueLimitAssociationInput{
			FarmId:  assoc.FarmID,
			QueueId: assoc.QueueID,
			LimitId: assoc.LimitID,
		}).
		Return(&deadline.DeleteQueueLimitAssociationOutput{}, nil)

	err := assoc.Remove(context.TODO())
	a.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_DeadlineCloudQueueLimitAssociation_Remove_Active_ConflictException(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDeadlineCloudClient)

	assoc := &DeadlineCloudQueueLimitAssociation{
		svc:     mockClient,
		FarmID:  ptr.String("farm-12345"),
		QueueID: ptr.String("queue-12345"),
		LimitID: ptr.String("limit-12345"),
		Status:  "ACTIVE",
	}

	mockGetAssociation(mockClient, deadlinetypes.QueueLimitAssociationStatusActive)

	mockClient.
		On("UpdateQueueLimitAssociation", mock.Anything, mock.Anything).
		Return((*deadline.UpdateQueueLimitAssociationOutput)(nil),
			&smithy.GenericAPIError{Code: "ConflictException", Message: "already transitioning"})

	err := assoc.Remove(context.TODO())
	a.Error(err)
	var errHold liberrors.ErrHoldResource
	a.ErrorAs(err, &errHold)
	mockClient.AssertExpectations(t)
}

func Test_Mock_DeadlineCloudQueueLimitAssociation_Properties(t *testing.T) {
	a := assert.New(t)

	assoc := DeadlineCloudQueueLimitAssociation{
		FarmID:  ptr.String("farm-12345"),
		QueueID: ptr.String("queue-12345"),
		LimitID: ptr.String("limit-12345"),
		Status:  "ACTIVE",
	}

	props := assoc.Properties()
	a.Equal("farm-12345", props.Get("FarmId"))
	a.Equal("queue-12345", props.Get("QueueId"))
	a.Equal("limit-12345", props.Get("LimitId"))
	a.Equal("ACTIVE", props.Get("Status"))
}

func Test_Mock_DeadlineCloudQueueLimitAssociation_String(t *testing.T) {
	a := assert.New(t)

	assoc := DeadlineCloudQueueLimitAssociation{
		LimitID: ptr.String("limit-12345"),
	}

	a.Equal("limit-12345", assoc.String())
}
