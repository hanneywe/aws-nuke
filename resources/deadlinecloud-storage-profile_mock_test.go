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

func Test_Mock_DeadlineCloudStorageProfile_List_One(t *testing.T) {
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
		On("ListStorageProfiles", mock.Anything, mock.Anything).
		Return(&deadline.ListStorageProfilesOutput{
			StorageProfiles: []deadlinetypes.StorageProfileSummary{
				{
					StorageProfileId: ptr.String("sp-12345"),
					DisplayName:      ptr.String("my-storage-profile"),
				},
			},
		}, nil)

	lister := &DeadlineCloudStorageProfileLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDeadlineCloudListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	sp := resources[0].(*DeadlineCloudStorageProfile)
	a.Equal("sp-12345", *sp.StorageProfileID)
	a.Equal("my-storage-profile", *sp.DisplayName)
	a.Equal("farm-12345", *sp.FarmID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DeadlineCloudStorageProfile_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDeadlineCloudClient)

	mockClient.
		On("ListFarms", mock.Anything, mock.Anything).
		Return(&deadline.ListFarmsOutput{
			Farms: []deadlinetypes.FarmSummary{},
		}, nil)

	lister := &DeadlineCloudStorageProfileLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDeadlineCloudListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DeadlineCloudStorageProfile_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDeadlineCloudClient)

	sp := &DeadlineCloudStorageProfile{
		svc:              mockClient,
		FarmID:           ptr.String("farm-12345"),
		StorageProfileID: ptr.String("sp-12345"),
	}

	mockClient.
		On("DeleteStorageProfile", mock.Anything, &deadline.DeleteStorageProfileInput{
			FarmId:           sp.FarmID,
			StorageProfileId: sp.StorageProfileID,
		}).
		Return(&deadline.DeleteStorageProfileOutput{}, nil)

	err := sp.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DeadlineCloudStorageProfile_Properties(t *testing.T) {
	a := assert.New(t)

	sp := DeadlineCloudStorageProfile{
		FarmID:           ptr.String("farm-12345"),
		StorageProfileID: ptr.String("sp-12345"),
		DisplayName:      ptr.String("my-storage-profile"),
	}

	props := sp.Properties()
	a.Equal("farm-12345", props.Get("FarmId"))
	a.Equal("sp-12345", props.Get("StorageProfileId"))
	a.Equal("my-storage-profile", props.Get("DisplayName"))
}

func Test_Mock_DeadlineCloudStorageProfile_String(t *testing.T) {
	a := assert.New(t)

	sp := DeadlineCloudStorageProfile{
		StorageProfileID: ptr.String("sp-12345"),
	}

	a.Equal("sp-12345", sp.String())
}
