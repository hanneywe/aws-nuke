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

func Test_Mock_MediaTailorVodSource_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaTailorV2Client)

	mockClient.On("ListSourceLocations", mock.Anything, mock.Anything).
		Return(&mediatailor.ListSourceLocationsOutput{
			Items: []mediatailortypes.SourceLocation{
				{SourceLocationName: ptr.String("loc-1")},
			},
		}, nil)

	mockClient.On("ListVodSources", mock.Anything, mock.Anything).
		Return(&mediatailor.ListVodSourcesOutput{
			Items: []mediatailortypes.VodSource{
				{
					SourceLocationName: ptr.String("loc-1"),
					VodSourceName:      ptr.String("vod-1"),
				},
			},
		}, nil)

	lister := &MediaTailorVodSourceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaTailorV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	src := resources[0].(*MediaTailorVodSource)
	a.Equal("loc-1", *src.SourceLocationName)
	a.Equal("vod-1", *src.VodSourceName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaTailorVodSource_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaTailorV2Client)

	mockClient.On("ListSourceLocations", mock.Anything, mock.Anything).
		Return(&mediatailor.ListSourceLocationsOutput{
			Items: []mediatailortypes.SourceLocation{},
		}, nil)

	lister := &MediaTailorVodSourceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaTailorV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaTailorVodSource_List_MultipleSourceLocations(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaTailorV2Client)

	mockClient.On("ListSourceLocations", mock.Anything, mock.Anything).
		Return(&mediatailor.ListSourceLocationsOutput{
			Items: []mediatailortypes.SourceLocation{
				{SourceLocationName: ptr.String("loc-1")},
				{SourceLocationName: ptr.String("loc-2")},
			},
		}, nil)

	mockClient.On("ListVodSources", mock.Anything, mock.MatchedBy(func(input *mediatailor.ListVodSourcesInput) bool {
		return *input.SourceLocationName == "loc-1"
	})).Return(&mediatailor.ListVodSourcesOutput{
		Items: []mediatailortypes.VodSource{
			{SourceLocationName: ptr.String("loc-1"), VodSourceName: ptr.String("vod-a")},
		},
	}, nil)

	mockClient.On("ListVodSources", mock.Anything, mock.MatchedBy(func(input *mediatailor.ListVodSourcesInput) bool {
		return *input.SourceLocationName == "loc-2"
	})).Return(&mediatailor.ListVodSourcesOutput{
		Items: []mediatailortypes.VodSource{
			{SourceLocationName: ptr.String("loc-2"), VodSourceName: ptr.String("vod-b")},
			{SourceLocationName: ptr.String("loc-2"), VodSourceName: ptr.String("vod-c")},
		},
	}, nil)

	lister := &MediaTailorVodSourceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaTailorV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 3)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaTailorVodSource_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaTailorV2Client)

	src := &MediaTailorVodSource{
		svc:                mockClient,
		SourceLocationName: ptr.String("loc-1"),
		VodSourceName:      ptr.String("vod-1"),
	}

	mockClient.On("DeleteVodSource", mock.Anything, &mediatailor.DeleteVodSourceInput{
		SourceLocationName: src.SourceLocationName,
		VodSourceName:      src.VodSourceName,
	}).Return(&mediatailor.DeleteVodSourceOutput{}, nil)

	err := src.Remove(context.TODO())
	a.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaTailorVodSource_Properties(t *testing.T) {
	a := assert.New(t)

	src := MediaTailorVodSource{
		SourceLocationName: ptr.String("loc-1"),
		VodSourceName:      ptr.String("vod-1"),
	}

	props := src.Properties()
	a.Equal("loc-1", props.Get("SourceLocationName"))
	a.Equal("vod-1", props.Get("VodSourceName"))
}

func Test_Mock_MediaTailorVodSource_String(t *testing.T) {
	a := assert.New(t)

	src := MediaTailorVodSource{
		SourceLocationName: ptr.String("loc-1"),
		VodSourceName:      ptr.String("vod-1"),
	}

	a.Equal("loc-1/vod-1", src.String())
}
