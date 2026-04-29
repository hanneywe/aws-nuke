package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/networkmanager"
	nmtypes "github.com/aws/aws-sdk-go-v2/service/networkmanager/types"
)

// --- Listing ---

func Test_Mock_NetworkManagerSite_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNetworkManagerClient)

	// One global network with one site
	mockClient.
		On("DescribeGlobalNetworks", mock.Anything, mock.Anything).
		Return(
			&networkmanager.DescribeGlobalNetworksOutput{
				GlobalNetworks: []nmtypes.GlobalNetwork{
					{
						GlobalNetworkId: ptr.String("gn-1"),
						State:           nmtypes.GlobalNetworkStateAvailable,
					},
				},
			}, nil,
		)

	mockClient.
		On("GetSites", mock.Anything, mock.Anything).
		Return(
			&networkmanager.GetSitesOutput{
				Sites: []nmtypes.Site{
					{
						SiteId:          ptr.String("site-1"),
						GlobalNetworkId: ptr.String("gn-1"),
						Description:     ptr.String("test site"),
						State:           nmtypes.SiteStateAvailable,
						Tags: []nmtypes.Tag{
							{Key: ptr.String("env"), Value: ptr.String("test")},
						},
					},
				},
			}, nil,
		)

	lister := &NetworkManagerSiteLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testNetworkManagerListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	site := resources[0].(*NetworkManagerSite)
	assertions.Equal("site-1", *site.ID)
	assertions.Equal("gn-1", *site.GlobalNetworkID)
	assertions.Equal("test site", *site.Description)
	assertions.Equal("AVAILABLE", *site.State)
	assertions.Equal("test", site.Tags["env"])

	mockClient.AssertExpectations(t)
}

func Test_Mock_NetworkManagerSite_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNetworkManagerClient)

	// No global networks means no sites
	mockClient.
		On("DescribeGlobalNetworks", mock.Anything, mock.Anything).
		Return(
			&networkmanager.DescribeGlobalNetworksOutput{
				GlobalNetworks: []nmtypes.GlobalNetwork{},
			}, nil,
		)

	lister := &NetworkManagerSiteLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testNetworkManagerListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NetworkManagerSite_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNetworkManagerClient)

	// Paginated global networks: page 1 returns one GN with a next token, page 2 returns another GN
	mockClient.
		On(
			"DescribeGlobalNetworks",
			mock.Anything,
			mock.MatchedBy(func(input *networkmanager.DescribeGlobalNetworksInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&networkmanager.DescribeGlobalNetworksOutput{
				GlobalNetworks: []nmtypes.GlobalNetwork{
					{GlobalNetworkId: ptr.String("gn-1"), State: nmtypes.GlobalNetworkStateAvailable},
				},
				NextToken: ptr.String("gn-page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"DescribeGlobalNetworks",
			mock.Anything,
			mock.MatchedBy(func(input *networkmanager.DescribeGlobalNetworksInput) bool {
				return input.NextToken != nil && *input.NextToken == "gn-page2"
			}),
		).
		Return(
			&networkmanager.DescribeGlobalNetworksOutput{
				GlobalNetworks: []nmtypes.GlobalNetwork{
					{GlobalNetworkId: ptr.String("gn-2"), State: nmtypes.GlobalNetworkStateAvailable},
				},
			}, nil,
		).
		Once()

	// gn-1: 100 sites on page 1, 1 on page 2
	firstPageSites := make([]nmtypes.Site, 100)
	for i := range firstPageSites {
		firstPageSites[i] = nmtypes.Site{
			SiteId:          ptr.String(fmt.Sprintf("site-%d", i)),
			GlobalNetworkId: ptr.String("gn-1"),
			State:           nmtypes.SiteStateAvailable,
		}
	}

	mockClient.
		On(
			"GetSites",
			mock.Anything,
			mock.MatchedBy(func(input *networkmanager.GetSitesInput) bool {
				return *input.GlobalNetworkId == "gn-1" && input.NextToken == nil
			}),
		).
		Return(
			&networkmanager.GetSitesOutput{
				Sites:     firstPageSites,
				NextToken: ptr.String("site-page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"GetSites",
			mock.Anything,
			mock.MatchedBy(func(input *networkmanager.GetSitesInput) bool {
				return *input.GlobalNetworkId == "gn-1" && input.NextToken != nil
			}),
		).
		Return(
			&networkmanager.GetSitesOutput{
				Sites: []nmtypes.Site{
					{
						SiteId:          ptr.String("site-100"),
						GlobalNetworkId: ptr.String("gn-1"),
						State:           nmtypes.SiteStateAvailable,
					},
				},
			}, nil,
		).
		Once()

	// gn-2: no sites
	mockClient.
		On(
			"GetSites",
			mock.Anything,
			mock.MatchedBy(func(input *networkmanager.GetSitesInput) bool {
				return *input.GlobalNetworkId == "gn-2"
			}),
		).
		Return(
			&networkmanager.GetSitesOutput{
				Sites: []nmtypes.Site{},
			}, nil,
		).
		Once()

	lister := &NetworkManagerSiteLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testNetworkManagerListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

// --- Removal ---

func Test_Mock_NetworkManagerSite_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNetworkManagerClient)

	site := &NetworkManagerSite{
		svc:             mockClient,
		ID:              ptr.String("site-1"),
		GlobalNetworkID: ptr.String("gn-1"),
	}

	mockClient.
		On(
			"DeleteSite",
			mock.Anything,
			&networkmanager.DeleteSiteInput{
				GlobalNetworkId: site.GlobalNetworkID,
				SiteId:          site.ID,
			},
		).
		Return(&networkmanager.DeleteSiteOutput{}, nil)

	err := site.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

// --- Properties ---

func Test_Mock_NetworkManagerSite_Properties(t *testing.T) {
	assertions := assert.New(t)

	site := NetworkManagerSite{
		ID:              ptr.String("site-12345"),
		GlobalNetworkID: ptr.String("gn-99"),
		Description:     ptr.String("primary site"),
		State:           ptr.String("AVAILABLE"),
		Tags: map[string]string{
			"Environment": "production",
			"Team":        "platform",
		},
	}

	properties := site.Properties()

	assertions.Equal("site-12345", properties.Get("ID"))
	assertions.Equal("gn-99", properties.Get("GlobalNetworkID"))
	assertions.Equal("primary site", properties.Get("Description"))
	assertions.Equal("AVAILABLE", properties.Get("State"))
	assertions.Equal("production", properties.Get("tag:Environment"))
	assertions.Equal("platform", properties.Get("tag:Team"))
}

func Test_Mock_NetworkManagerSite_Properties_EmptyTags(t *testing.T) {
	assertions := assert.New(t)

	site := NetworkManagerSite{
		ID:              ptr.String("site-99999"),
		GlobalNetworkID: ptr.String("gn-1"),
		State:           ptr.String("AVAILABLE"),
		Tags:            map[string]string{},
	}

	properties := site.Properties()

	assertions.Equal("site-99999", properties.Get("ID"))
	assertions.Equal("gn-1", properties.Get("GlobalNetworkID"))
	assertions.Equal("AVAILABLE", properties.Get("State"))
}

// --- Display ---

func Test_Mock_NetworkManagerSite_String(t *testing.T) {
	assertions := assert.New(t)

	site := NetworkManagerSite{
		ID: ptr.String("site-abc123"),
	}

	assertions.Equal("site-abc123", site.String())
}

// --- Filter ---

func Test_Mock_NetworkManagerSite_Filter_Deleting(t *testing.T) {
	assertions := assert.New(t)

	site := NetworkManagerSite{
		ID:    ptr.String("site-1"),
		State: ptr.String("DELETING"),
	}

	err := site.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "already deleting")
}

func Test_Mock_NetworkManagerSite_Filter_Available(t *testing.T) {
	assertions := assert.New(t)

	site := NetworkManagerSite{
		ID:    ptr.String("site-1"),
		State: ptr.String("AVAILABLE"),
	}

	err := site.Filter()
	assertions.NoError(err)
}
