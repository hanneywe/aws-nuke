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

func Test_Mock_OmicsReferenceStore_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockOmicsClient)

	mockClient.
		On("ListReferenceStores", mock.Anything, mock.Anything).
		Return(&omics.ListReferenceStoresOutput{
			ReferenceStores: []omicstypes.ReferenceStoreDetail{
				{
					Id:   ptr.String("ref-store-1"),
					Name: ptr.String("my-reference-store"),
				},
			},
		}, nil)

	lister := &OmicsReferenceStoreLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testOmicsListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	referenceStore := resources[0].(*OmicsReferenceStore)
	assertions.Equal("ref-store-1", *referenceStore.ID)
	assertions.Equal("my-reference-store", *referenceStore.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OmicsReferenceStore_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockOmicsClient)

	mockClient.
		On("ListReferenceStores", mock.Anything, mock.Anything).
		Return(&omics.ListReferenceStoresOutput{
			ReferenceStores: []omicstypes.ReferenceStoreDetail{},
		}, nil)

	lister := &OmicsReferenceStoreLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testOmicsListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OmicsReferenceStore_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockOmicsClient)

	referenceStore := &OmicsReferenceStore{
		svc:  mockClient,
		ID:   ptr.String("ref-store-1"),
		Name: ptr.String("my-reference-store"),
	}

	mockClient.
		On("DeleteReferenceStore", mock.Anything, &omics.DeleteReferenceStoreInput{
			Id: referenceStore.ID,
		}).
		Return(&omics.DeleteReferenceStoreOutput{}, nil)

	err := referenceStore.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OmicsReferenceStore_Properties(t *testing.T) {
	assertions := assert.New(t)

	referenceStore := OmicsReferenceStore{
		ID:   ptr.String("ref-store-1"),
		Name: ptr.String("my-reference-store"),
	}

	properties := referenceStore.Properties()

	assertions.Equal("ref-store-1", properties.Get("ID"))
	assertions.Equal("my-reference-store", properties.Get("Name"))
}

func Test_Mock_OmicsReferenceStore_String(t *testing.T) {
	assertions := assert.New(t)

	referenceStore := OmicsReferenceStore{
		ID: ptr.String("ref-store-1"),
	}

	assertions.Equal("ref-store-1", referenceStore.String())
}
