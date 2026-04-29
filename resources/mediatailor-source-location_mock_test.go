package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	mediatailortypes "github.com/aws/aws-sdk-go-v2/service/mediatailor/types"
)

func Test_Mock_MediaTailorSourceLocation_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaTailorV2Client)

	mockClient.On("ListSourceLocations", mock.Anything, mock.Anything).
		Return(&mediatailor.ListSourceLocationsOutput{
			Items: []mediatailortypes.SourceLocation{
				{
					SourceLocationName: ptr.String("test-location"),
				},
			},
		}, nil)

	lister := &MediaTailorSourceLocationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaTailorV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	loc := resources[0].(*MediaTailorSourceLocation)
	a.Equal("test-location", *loc.SourceLocationName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaTailorSourceLocation_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaTailorV2Client)

	mockClient.On("ListSourceLocations", mock.Anything, mock.Anything).
		Return(&mediatailor.ListSourceLocationsOutput{
			Items: []mediatailortypes.SourceLocation{},
		}, nil)

	lister := &MediaTailorSourceLocationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaTailorV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaTailorSourceLocation_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaTailorV2Client)

	loc := &MediaTailorSourceLocation{
		svc:                mockClient,
		SourceLocationName: ptr.String("test-location"),
	}

	mockClient.On("DeleteSourceLocation", mock.Anything, &mediatailor.DeleteSourceLocationInput{
		SourceLocationName: loc.SourceLocationName,
	}).Return(&mediatailor.DeleteSourceLocationOutput{}, nil)

	err := loc.Remove(context.TODO())
	a.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaTailorSourceLocation_Properties(t *testing.T) {
	a := assert.New(t)

	loc := MediaTailorSourceLocation{
		SourceLocationName: ptr.String("test-location"),
	}

	props := loc.Properties()
	a.Equal("test-location", props.Get("SourceLocationName"))
}

func Test_Mock_MediaTailorSourceLocation_String(t *testing.T) {
	a := assert.New(t)

	loc := MediaTailorSourceLocation{
		SourceLocationName: ptr.String("test-location"),
	}

	a.Equal("test-location", loc.String())
}
