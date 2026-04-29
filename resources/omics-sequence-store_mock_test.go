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

func Test_Mock_OmicsSequenceStore_List_One(t *testing.T) {
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

	lister := &OmicsSequenceStoreLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testOmicsListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	sequenceStore := resources[0].(*OmicsSequenceStore)
	assertions.Equal("seq-store-1", *sequenceStore.ID)
	assertions.Equal("my-sequence-store", *sequenceStore.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OmicsSequenceStore_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockOmicsClient)

	mockClient.
		On("ListSequenceStores", mock.Anything, mock.Anything).
		Return(&omics.ListSequenceStoresOutput{
			SequenceStores: []omicstypes.SequenceStoreDetail{},
		}, nil)

	lister := &OmicsSequenceStoreLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testOmicsListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OmicsSequenceStore_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockOmicsClient)

	sequenceStore := &OmicsSequenceStore{
		svc:  mockClient,
		ID:   ptr.String("seq-store-1"),
		Name: ptr.String("my-sequence-store"),
	}

	mockClient.
		On("DeleteSequenceStore", mock.Anything, &omics.DeleteSequenceStoreInput{
			Id: sequenceStore.ID,
		}).
		Return(&omics.DeleteSequenceStoreOutput{}, nil)

	err := sequenceStore.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OmicsSequenceStore_Properties(t *testing.T) {
	assertions := assert.New(t)

	sequenceStore := OmicsSequenceStore{
		ID:   ptr.String("seq-store-1"),
		Name: ptr.String("my-sequence-store"),
	}

	properties := sequenceStore.Properties()

	assertions.Equal("seq-store-1", properties.Get("ID"))
	assertions.Equal("my-sequence-store", properties.Get("Name"))
}

func Test_Mock_OmicsSequenceStore_String(t *testing.T) {
	assertions := assert.New(t)

	sequenceStore := OmicsSequenceStore{
		ID: ptr.String("seq-store-1"),
	}

	assertions.Equal("seq-store-1", sequenceStore.String())
}
