package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/medialive"
	mltypes "github.com/aws/aws-sdk-go-v2/service/medialive/types"
)

func Test_Mock_MediaLiveSdiSource_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaLiveClient)

	mockClient.On("ListSdiSources", mock.Anything, mock.Anything).
		Return(&medialive.ListSdiSourcesOutput{
			SdiSources: []mltypes.SdiSourceSummary{
				{
					Id:   ptr.String("sdi-123"),
					Name: ptr.String("my-sdi-source"),
				},
			},
		}, nil)

	lister := &MediaLiveSdiSourceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaLiveListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*MediaLiveSdiSource)
	a.Equal("sdi-123", *r.ID)
	a.Equal("my-sdi-source", *r.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaLiveSdiSource_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaLiveClient)

	mockClient.On("ListSdiSources", mock.Anything, mock.Anything).
		Return(&medialive.ListSdiSourcesOutput{
			SdiSources: []mltypes.SdiSourceSummary{},
		}, nil)

	lister := &MediaLiveSdiSourceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaLiveListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaLiveSdiSource_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaLiveClient)

	r := &MediaLiveSdiSource{
		svc: mockClient,
		ID:  ptr.String("sdi-123"),
	}

	mockClient.On("DeleteSdiSource", mock.Anything, &medialive.DeleteSdiSourceInput{
		SdiSourceId: r.ID,
	}).Return(&medialive.DeleteSdiSourceOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaLiveSdiSource_Properties(t *testing.T) {
	a := assert.New(t)

	r := MediaLiveSdiSource{
		ID:   ptr.String("sdi-123"),
		Name: ptr.String("my-sdi-source"),
	}

	props := r.Properties()
	a.Equal("sdi-123", props.Get("ID"))
	a.Equal("my-sdi-source", props.Get("Name"))
}

func Test_Mock_MediaLiveSdiSource_String(t *testing.T) {
	a := assert.New(t)
	r := MediaLiveSdiSource{ID: ptr.String("sdi-123")}
	a.Equal("sdi-123", r.String())
}
