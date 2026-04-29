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

func Test_Mock_OmicsReference_List_One(t *testing.T) {
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

	mockClient.
		On("ListReferences", mock.Anything, mock.Anything).
		Return(&omics.ListReferencesOutput{
			References: []omicstypes.ReferenceListItem{
				{
					Id:               ptr.String("ref-1"),
					Name:             ptr.String("my-reference"),
					ReferenceStoreId: ptr.String("ref-store-1"),
				},
			},
		}, nil)

	lister := &OmicsReferenceLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testOmicsListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	reference := resources[0].(*OmicsReference)
	assertions.Equal("ref-1", *reference.ID)
	assertions.Equal("my-reference", *reference.Name)
	assertions.Equal("ref-store-1", *reference.ReferenceStoreID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OmicsReference_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockOmicsClient)

	mockClient.
		On("ListReferenceStores", mock.Anything, mock.Anything).
		Return(&omics.ListReferenceStoresOutput{
			ReferenceStores: []omicstypes.ReferenceStoreDetail{},
		}, nil)

	lister := &OmicsReferenceLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testOmicsListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OmicsReference_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockOmicsClient)

	reference := &OmicsReference{
		svc:              mockClient,
		ID:               ptr.String("ref-1"),
		Name:             ptr.String("my-reference"),
		ReferenceStoreID: ptr.String("ref-store-1"),
	}

	mockClient.
		On("DeleteReference", mock.Anything, &omics.DeleteReferenceInput{
			Id:               reference.ID,
			ReferenceStoreId: reference.ReferenceStoreID,
		}).
		Return(&omics.DeleteReferenceOutput{}, nil)

	err := reference.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OmicsReference_Properties(t *testing.T) {
	assertions := assert.New(t)

	reference := OmicsReference{
		ID:               ptr.String("ref-1"),
		Name:             ptr.String("my-reference"),
		ReferenceStoreID: ptr.String("ref-store-1"),
	}

	properties := reference.Properties()

	assertions.Equal("ref-1", properties.Get("ID"))
	assertions.Equal("my-reference", properties.Get("Name"))
	assertions.Equal("ref-store-1", properties.Get("ReferenceStoreID"))
}

func Test_Mock_OmicsReference_String(t *testing.T) {
	assertions := assert.New(t)

	reference := OmicsReference{
		ID: ptr.String("ref-1"),
	}

	assertions.Equal("ref-1", reference.String())
}
