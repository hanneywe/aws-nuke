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

func Test_Mock_LocationServiceMap_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLocationServiceClient)

	mockClient.On("ListMaps", mock.Anything, mock.Anything).
		Return(&location.ListMapsOutput{
			Entries: []locationtypes.ListMapsResponseEntry{
				{MapName: ptr.String("my-map")},
			},
		}, nil)

	lister := &LocationServiceMapLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLocationServiceListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	m := resources[0].(*LocationServiceMap)
	a.Equal("my-map", *m.MapName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LocationServiceMap_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLocationServiceClient)

	mockClient.On("ListMaps", mock.Anything, mock.Anything).
		Return(&location.ListMapsOutput{
			Entries: []locationtypes.ListMapsResponseEntry{},
		}, nil)

	lister := &LocationServiceMapLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLocationServiceListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LocationServiceMap_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLocationServiceClient)

	m := &LocationServiceMap{
		svc:     mockClient,
		MapName: ptr.String("my-map"),
	}

	mockClient.On("DeleteMap", mock.Anything, &location.DeleteMapInput{
		MapName: m.MapName,
	}).Return(&location.DeleteMapOutput{}, nil)

	a.NoError(m.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_LocationServiceMap_Properties(t *testing.T) {
	a := assert.New(t)
	m := LocationServiceMap{MapName: ptr.String("my-map")}
	props := m.Properties()
	a.Equal("my-map", props.Get("MapName"))
}

func Test_Mock_LocationServiceMap_String(t *testing.T) {
	a := assert.New(t)
	m := LocationServiceMap{MapName: ptr.String("my-map")}
	a.Equal("my-map", m.String())
}
