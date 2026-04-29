package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/omics"
	omicstypes "github.com/aws/aws-sdk-go-v2/service/omics/types"
)

func Test_Mock_OmicsReadSet_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockOmicsClient)

	mockClient.
		On("ListSequenceStores", mock.Anything, mock.Anything).
		Return(&omics.ListSequenceStoresOutput{
			SequenceStores: []omicstypes.SequenceStoreDetail{
				{
					Id:   ptr.String("seq-store-1"),
					Name: ptr.String("my-sequence-store"),
				},
			},
		}, nil)

	mockClient.
		On("ListReadSets", mock.Anything, mock.Anything).
		Return(&omics.ListReadSetsOutput{
			ReadSets: []omicstypes.ReadSetListItem{
				{
					Id:              ptr.String("readset-1"),
					Name:            ptr.String("my-read-set"),
					SequenceStoreId: ptr.String("seq-store-1"),
				},
			},
		}, nil)

	lister := &OmicsReadSetLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testOmicsListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	readSet := resources[0].(*OmicsReadSet)
	assertions.Equal("readset-1", *readSet.ID)
	assertions.Equal("my-read-set", *readSet.Name)
	assertions.Equal("seq-store-1", *readSet.SequenceStoreID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OmicsReadSet_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockOmicsClient)

	mockClient.
		On("ListSequenceStores", mock.Anything, mock.Anything).
		Return(&omics.ListSequenceStoresOutput{
			SequenceStores: []omicstypes.SequenceStoreDetail{},
		}, nil)

	lister := &OmicsReadSetLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testOmicsListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OmicsReadSet_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockOmicsClient)

	readSet := &OmicsReadSet{
		svc:             mockClient,
		ID:              ptr.String("readset-1"),
		Name:            ptr.String("my-read-set"),
		SequenceStoreID: ptr.String("seq-store-1"),
	}

	mockClient.
		On("BatchDeleteReadSet", mock.Anything, &omics.BatchDeleteReadSetInput{
			SequenceStoreId: readSet.SequenceStoreID,
			Ids:             []string{"readset-1"},
		}).
		Return(&omics.BatchDeleteReadSetOutput{}, nil)

	err := readSet.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OmicsReadSet_Properties(t *testing.T) {
	assertions := assert.New(t)

	readSet := OmicsReadSet{
		ID:              ptr.String("readset-1"),
		Name:            ptr.String("my-read-set"),
		SequenceStoreID: ptr.String("seq-store-1"),
	}

	properties := readSet.Properties()

	assertions.Equal("readset-1", properties.Get("ID"))
	assertions.Equal("my-read-set", properties.Get("Name"))
	assertions.Equal("seq-store-1", properties.Get("SequenceStoreID"))
}

func Test_Mock_OmicsReadSet_String(t *testing.T) {
	assertions := assert.New(t)

	readSet := OmicsReadSet{
		ID: ptr.String("readset-1"),
	}

	assertions.Equal("readset-1", readSet.String())
}
