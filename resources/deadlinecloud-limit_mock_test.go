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

func Test_Mock_DeadlineCloudLimit_List_One(t *testing.T) {
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
		On("ListLimits", mock.Anything, mock.Anything).
		Return(&deadline.ListLimitsOutput{
			Limits: []deadlinetypes.LimitSummary{
				{
					FarmId:      ptr.String("farm-12345"),
					LimitId:     ptr.String("limit-12345"),
					DisplayName: ptr.String("my-limit"),
				},
			},
		}, nil)

	lister := &DeadlineCloudLimitLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDeadlineCloudListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	limit := resources[0].(*DeadlineCloudLimit)
	a.Equal("limit-12345", *limit.LimitID)
	a.Equal("my-limit", *limit.DisplayName)
	a.Equal("farm-12345", *limit.FarmID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DeadlineCloudLimit_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDeadlineCloudClient)

	mockClient.
		On("ListFarms", mock.Anything, mock.Anything).
		Return(&deadline.ListFarmsOutput{
			Farms: []deadlinetypes.FarmSummary{},
		}, nil)

	lister := &DeadlineCloudLimitLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDeadlineCloudListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DeadlineCloudLimit_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDeadlineCloudClient)

	limit := &DeadlineCloudLimit{
		svc:     mockClient,
		FarmID:  ptr.String("farm-12345"),
		LimitID: ptr.String("limit-12345"),
	}

	mockClient.
		On("DeleteLimit", mock.Anything, &deadline.DeleteLimitInput{
			FarmId:  limit.FarmID,
			LimitId: limit.LimitID,
		}).
		Return(&deadline.DeleteLimitOutput{}, nil)

	err := limit.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DeadlineCloudLimit_Properties(t *testing.T) {
	a := assert.New(t)

	limit := DeadlineCloudLimit{
		FarmID:      ptr.String("farm-12345"),
		LimitID:     ptr.String("limit-12345"),
		DisplayName: ptr.String("my-limit"),
	}

	props := limit.Properties()
	a.Equal("farm-12345", props.Get("FarmId"))
	a.Equal("limit-12345", props.Get("LimitId"))
	a.Equal("my-limit", props.Get("DisplayName"))
}

func Test_Mock_DeadlineCloudLimit_String(t *testing.T) {
	a := assert.New(t)

	limit := DeadlineCloudLimit{
		LimitID: ptr.String("limit-12345"),
	}

	a.Equal("limit-12345", limit.String())
}
