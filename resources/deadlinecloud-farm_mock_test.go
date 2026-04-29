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

func Test_Mock_DeadlineCloudFarm_List_One(t *testing.T) {
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

	lister := &DeadlineCloudFarmLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDeadlineCloudListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	farm := resources[0].(*DeadlineCloudFarm)
	a.Equal("farm-12345", *farm.FarmID)
	a.Equal("my-farm", *farm.DisplayName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DeadlineCloudFarm_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDeadlineCloudClient)

	mockClient.
		On("ListFarms", mock.Anything, mock.Anything).
		Return(&deadline.ListFarmsOutput{
			Farms: []deadlinetypes.FarmSummary{},
		}, nil)

	lister := &DeadlineCloudFarmLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDeadlineCloudListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DeadlineCloudFarm_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDeadlineCloudClient)

	farm := &DeadlineCloudFarm{
		svc:    mockClient,
		FarmID: ptr.String("farm-12345"),
	}

	mockClient.
		On("DeleteFarm", mock.Anything, &deadline.DeleteFarmInput{
			FarmId: farm.FarmID,
		}).
		Return(&deadline.DeleteFarmOutput{}, nil)

	err := farm.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DeadlineCloudFarm_Properties(t *testing.T) {
	a := assert.New(t)

	farm := DeadlineCloudFarm{
		FarmID:      ptr.String("farm-12345"),
		DisplayName: ptr.String("my-farm"),
	}

	props := farm.Properties()
	a.Equal("farm-12345", props.Get("FarmId"))
	a.Equal("my-farm", props.Get("DisplayName"))
}

func Test_Mock_DeadlineCloudFarm_String(t *testing.T) {
	a := assert.New(t)

	farm := DeadlineCloudFarm{
		FarmID: ptr.String("farm-12345"),
	}

	a.Equal("farm-12345", farm.String())
}
