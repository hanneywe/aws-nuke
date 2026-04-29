package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/outposts"
	outpoststypes "github.com/aws/aws-sdk-go-v2/service/outposts/types"
)

func Test_Mock_OutpostsSite_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockOutpostsClient)
	mockClient.On("ListSites", mock.Anything, mock.Anything).
		Return(&outposts.ListSitesOutput{
			Sites: []outpoststypes.Site{
				{SiteId: ptr.String("os-12345"), Name: ptr.String("my-site")},
			},
		}, nil)
	lister := &OutpostsSiteLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testOutpostsListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	site := resources[0].(*OutpostsSite)
	a.Equal("os-12345", *site.SiteID)
	a.Equal("my-site", *site.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_OutpostsSite_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockOutpostsClient)
	mockClient.On("ListSites", mock.Anything, mock.Anything).
		Return(&outposts.ListSitesOutput{Sites: []outpoststypes.Site{}}, nil)
	lister := &OutpostsSiteLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testOutpostsListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_OutpostsSite_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockOutpostsClient)
	site := &OutpostsSite{svc: mockClient, SiteID: ptr.String("os-12345"), Name: ptr.String("my-site")}
	mockClient.On("DeleteSite", mock.Anything, &outposts.DeleteSiteInput{SiteId: site.SiteID}).
		Return(&outposts.DeleteSiteOutput{}, nil)
	a.NoError(site.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_OutpostsSite_Properties(t *testing.T) {
	a := assert.New(t)
	site := OutpostsSite{SiteID: ptr.String("os-12345"), Name: ptr.String("my-site")}
	a.Equal("os-12345", site.Properties().Get("SiteId"))
	a.Equal("my-site", site.Properties().Get("Name"))
}

func Test_Mock_OutpostsSite_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-site", (&OutpostsSite{Name: ptr.String("my-site")}).String())
}
