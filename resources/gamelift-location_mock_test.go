package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/gamelift"
	gamelifttypes "github.com/aws/aws-sdk-go-v2/service/gamelift/types"
)

func Test_Mock_GameLiftLocation_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGameLiftV2Client)
	mockClient.On("ListLocations", mock.Anything, mock.Anything).
		Return(&gamelift.ListLocationsOutput{
			Locations: []gamelifttypes.LocationModel{
				{LocationName: ptr.String("custom-location-1")},
			},
		}, nil)
	lister := &GameLiftLocationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGameLiftV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("custom-location-1", resources[0].(*GameLiftLocation).String())
	mockClient.AssertExpectations(t)
}

func Test_Mock_GameLiftLocation_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGameLiftV2Client)
	mockClient.On("ListLocations", mock.Anything, mock.Anything).
		Return(&gamelift.ListLocationsOutput{Locations: []gamelifttypes.LocationModel{}}, nil)
	lister := &GameLiftLocationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGameLiftV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GameLiftLocation_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGameLiftV2Client)
	r := &GameLiftLocation{svc: mockClient, LocationName: ptr.String("custom-location-1")}
	mockClient.On("DeleteLocation", mock.Anything, &gamelift.DeleteLocationInput{
		LocationName: r.LocationName,
	}).Return(&gamelift.DeleteLocationOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_GameLiftLocation_Properties(t *testing.T) {
	a := assert.New(t)
	r := GameLiftLocation{LocationName: ptr.String("custom-location-1")}
	a.Equal("custom-location-1", r.Properties().Get("LocationName"))
}

func Test_Mock_GameLiftLocation_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("custom-location-1", (&GameLiftLocation{LocationName: ptr.String("custom-location-1")}).String())
}

func Test_Mock_GameLiftLocation_Filter_AWSManaged(t *testing.T) {
	a := assert.New(t)
	r := GameLiftLocation{LocationName: ptr.String("us-east-1")}
	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "cannot delete AWS-managed location")
}

func Test_Mock_GameLiftLocation_Filter_AWSManagedLocalZone(t *testing.T) {
	a := assert.New(t)
	r := GameLiftLocation{LocationName: ptr.String("us-east-1-chi-1")}
	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "cannot delete AWS-managed location")
}

func Test_Mock_GameLiftLocation_Filter_Custom(t *testing.T) {
	a := assert.New(t)
	r := GameLiftLocation{LocationName: ptr.String("custom-location-1")}
	a.NoError(r.Filter())
}
