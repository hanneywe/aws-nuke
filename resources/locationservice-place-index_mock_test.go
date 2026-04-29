package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/location"
	locationtypes "github.com/aws/aws-sdk-go-v2/service/location/types"
)

func Test_Mock_LocationServicePlaceIndex_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLocationServiceClient)

	mockClient.On("ListPlaceIndexes", mock.Anything, mock.Anything).
		Return(&location.ListPlaceIndexesOutput{
			Entries: []locationtypes.ListPlaceIndexesResponseEntry{
				{IndexName: ptr.String("test-value"), DataSource: ptr.String("test-value"), Description: ptr.String("test-value")},
			},
		}, nil)

	lister := &LocationServicePlaceIndexLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLocationServiceListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*LocationServicePlaceIndex)
	a.Equal("test-value", *r.IndexName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_LocationServicePlaceIndex_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLocationServiceClient)

	mockClient.On("ListPlaceIndexes", mock.Anything, mock.Anything).
		Return(&location.ListPlaceIndexesOutput{
			Entries: []locationtypes.ListPlaceIndexesResponseEntry{},
		}, nil)

	lister := &LocationServicePlaceIndexLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLocationServiceListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_LocationServicePlaceIndex_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLocationServiceClient)

	r := &LocationServicePlaceIndex{
		svc:       mockClient,
		IndexName: ptr.String("test-indexname"),
	}

	mockClient.On("DeletePlaceIndex", mock.Anything,
		&location.DeletePlaceIndexInput{
			IndexName: r.IndexName,
		}).Return(&location.DeletePlaceIndexOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_LocationServicePlaceIndex_Properties(t *testing.T) {
	a := assert.New(t)
	r := &LocationServicePlaceIndex{
		IndexName:   ptr.String("test-indexname"),
		DataSource:  ptr.String("test-datasource"),
		Description: ptr.String("test-description"),
	}
	props := r.Properties()
	a.Equal("test-indexname", props.Get("IndexName"))
	a.Equal("test-datasource", props.Get("DataSource"))
	a.Equal("test-description", props.Get("Description"))
}

func Test_Mock_LocationServicePlaceIndex_String(t *testing.T) {
	a := assert.New(t)
	r := &LocationServicePlaceIndex{
		IndexName: ptr.String("test-indexname"),
	}
	a.Equal("test-indexname", r.String())
}
