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

func Test_Mock_MediaTailorLiveSource_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaTailorV2Client)

	mockClient.On("ListSourceLocations", mock.Anything, mock.Anything).
		Return(&mediatailor.ListSourceLocationsOutput{
			Items: []mediatailortypes.SourceLocation{
				{SourceLocationName: ptr.String("loc-1")},
			},
		}, nil)

	mockClient.On("ListLiveSources", mock.Anything, mock.Anything).
		Return(&mediatailor.ListLiveSourcesOutput{
			Items: []mediatailortypes.LiveSource{
				{
					SourceLocationName: ptr.String("loc-1"),
					LiveSourceName:     ptr.String("live-1"),
				},
			},
		}, nil)

	lister := &MediaTailorLiveSourceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaTailorV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	src := resources[0].(*MediaTailorLiveSource)
	a.Equal("loc-1", *src.SourceLocationName)
	a.Equal("live-1", *src.LiveSourceName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaTailorLiveSource_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaTailorV2Client)

	mockClient.On("ListSourceLocations", mock.Anything, mock.Anything).
		Return(&mediatailor.ListSourceLocationsOutput{
			Items: []mediatailortypes.SourceLocation{},
		}, nil)

	lister := &MediaTailorLiveSourceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaTailorV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaTailorLiveSource_List_MultipleSourceLocations(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaTailorV2Client)

	mockClient.On("ListSourceLocations", mock.Anything, mock.Anything).
		Return(&mediatailor.ListSourceLocationsOutput{
			Items: []mediatailortypes.SourceLocation{
				{SourceLocationName: ptr.String("loc-1")},
				{SourceLocationName: ptr.String("loc-2")},
			},
		}, nil)

	mockClient.On("ListLiveSources", mock.Anything, mock.MatchedBy(func(input *mediatailor.ListLiveSourcesInput) bool {
		return *input.SourceLocationName == "loc-1"
	})).Return(&mediatailor.ListLiveSourcesOutput{
		Items: []mediatailortypes.LiveSource{
			{SourceLocationName: ptr.String("loc-1"), LiveSourceName: ptr.String("live-a")},
		},
	}, nil)

	mockClient.On("ListLiveSources", mock.Anything, mock.MatchedBy(func(input *mediatailor.ListLiveSourcesInput) bool {
		return *input.SourceLocationName == "loc-2"
	})).Return(&mediatailor.ListLiveSourcesOutput{
		Items: []mediatailortypes.LiveSource{
			{SourceLocationName: ptr.String("loc-2"), LiveSourceName: ptr.String("live-b")},
			{SourceLocationName: ptr.String("loc-2"), LiveSourceName: ptr.String("live-c")},
		},
	}, nil)

	lister := &MediaTailorLiveSourceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaTailorV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 3)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaTailorLiveSource_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaTailorV2Client)

	src := &MediaTailorLiveSource{
		svc:                mockClient,
		SourceLocationName: ptr.String("loc-1"),
		LiveSourceName:     ptr.String("live-1"),
	}

	mockClient.On("DeleteLiveSource", mock.Anything, &mediatailor.DeleteLiveSourceInput{
		SourceLocationName: src.SourceLocationName,
		LiveSourceName:     src.LiveSourceName,
	}).Return(&mediatailor.DeleteLiveSourceOutput{}, nil)

	err := src.Remove(context.TODO())
	a.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaTailorLiveSource_Properties(t *testing.T) {
	a := assert.New(t)

	src := MediaTailorLiveSource{
		SourceLocationName: ptr.String("loc-1"),
		LiveSourceName:     ptr.String("live-1"),
	}

	props := src.Properties()
	a.Equal("loc-1", props.Get("SourceLocationName"))
	a.Equal("live-1", props.Get("LiveSourceName"))
}

func Test_Mock_MediaTailorLiveSource_String(t *testing.T) {
	a := assert.New(t)

	src := MediaTailorLiveSource{
		SourceLocationName: ptr.String("loc-1"),
		LiveSourceName:     ptr.String("live-1"),
	}

	a.Equal("loc-1/live-1", src.String())
}
