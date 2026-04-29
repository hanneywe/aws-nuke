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

func Test_Mock_NetworkManagerLink_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNetworkManagerClient)

	// One global network with one link
	mockClient.
		On("DescribeGlobalNetworks", mock.Anything, mock.Anything).
		Return(
			&networkmanager.DescribeGlobalNetworksOutput{
				GlobalNetworks: []nmtypes.GlobalNetwork{
					{
						GlobalNetworkId: ptr.String(TestGlobalNetworkID1),
						State:           nmtypes.GlobalNetworkStateAvailable,
					},
				},
			}, nil,
		)

	mockClient.
		On("GetLinks", mock.Anything, mock.Anything).
		Return(
			&networkmanager.GetLinksOutput{
				Links: []nmtypes.Link{
					{
						LinkId:          ptr.String("link-1"),
						GlobalNetworkId: ptr.String(TestGlobalNetworkID1),
						SiteId:          ptr.String("site-1"),
						Description:     ptr.String("test link"),
						State:           nmtypes.LinkStateAvailable,
						Tags: []nmtypes.Tag{
							{Key: ptr.String("env"), Value: ptr.String("test")},
						},
					},
				},
			}, nil,
		)

	lister := &NetworkManagerLinkLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testNetworkManagerListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	link := resources[0].(*NetworkManagerLink)
	assertions.Equal("link-1", *link.ID)
	assertions.Equal(TestGlobalNetworkID1, *link.GlobalNetworkID)
	assertions.Equal("site-1", *link.SiteID)
	assertions.Equal("test link", *link.Description)
	assertions.Equal("AVAILABLE", *link.State)
	assertions.Equal("test", link.Tags["env"])

	mockClient.AssertExpectations(t)
}

func Test_Mock_NetworkManagerLink_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNetworkManagerClient)

	// No global networks means no links
	mockClient.
		On("DescribeGlobalNetworks", mock.Anything, mock.Anything).
		Return(
			&networkmanager.DescribeGlobalNetworksOutput{
				GlobalNetworks: []nmtypes.GlobalNetwork{},
			}, nil,
		)

	lister := &NetworkManagerLinkLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testNetworkManagerListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NetworkManagerLink_List_MultiPage(t *testing.T) {
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
					{GlobalNetworkId: ptr.String(TestGlobalNetworkID1), State: nmtypes.GlobalNetworkStateAvailable},
				},
				NextToken: ptr.String(TestGlobalNetworkPage2),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"DescribeGlobalNetworks",
			mock.Anything,
			mock.MatchedBy(func(input *networkmanager.DescribeGlobalNetworksInput) bool {
				return input.NextToken != nil && *input.NextToken == TestGlobalNetworkPage2
			}),
		).
		Return(
			&networkmanager.DescribeGlobalNetworksOutput{
				GlobalNetworks: []nmtypes.GlobalNetwork{
					{GlobalNetworkId: ptr.String(TestGlobalNetworkID2), State: nmtypes.GlobalNetworkStateAvailable},
				},
			}, nil,
		).
		Once()

	// gn-1: 100 links on page 1, 1 on page 2
	firstPageLinks := make([]nmtypes.Link, 100)
	for i := range firstPageLinks {
		firstPageLinks[i] = nmtypes.Link{
			LinkId:          ptr.String(fmt.Sprintf("link-%d", i)),
			GlobalNetworkId: ptr.String(TestGlobalNetworkID1),
			State:           nmtypes.LinkStateAvailable,
		}
	}

	mockClient.
		On(
			"GetLinks",
			mock.Anything,
			mock.MatchedBy(func(input *networkmanager.GetLinksInput) bool {
				return *input.GlobalNetworkId == TestGlobalNetworkID1 && input.NextToken == nil
			}),
		).
		Return(
			&networkmanager.GetLinksOutput{
				Links:     firstPageLinks,
				NextToken: ptr.String("link-page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"GetLinks",
			mock.Anything,
			mock.MatchedBy(func(input *networkmanager.GetLinksInput) bool {
				return *input.GlobalNetworkId == TestGlobalNetworkID1 && input.NextToken != nil
			}),
		).
		Return(
			&networkmanager.GetLinksOutput{
				Links: []nmtypes.Link{
					{
						LinkId:          ptr.String("link-100"),
						GlobalNetworkId: ptr.String(TestGlobalNetworkID1),
						State:           nmtypes.LinkStateAvailable,
					},
				},
			}, nil,
		).
		Once()

	// gn-2: no links
	mockClient.
		On(
			"GetLinks",
			mock.Anything,
			mock.MatchedBy(func(input *networkmanager.GetLinksInput) bool {
				return *input.GlobalNetworkId == TestGlobalNetworkID2
			}),
		).
		Return(
			&networkmanager.GetLinksOutput{
				Links: []nmtypes.Link{},
			}, nil,
		).
		Once()

	lister := &NetworkManagerLinkLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testNetworkManagerListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

// --- Removal ---

func Test_Mock_NetworkManagerLink_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNetworkManagerClient)

	link := &NetworkManagerLink{
		svc:             mockClient,
		ID:              ptr.String("link-1"),
		GlobalNetworkID: ptr.String(TestGlobalNetworkID1),
	}

	mockClient.
		On(
			"DeleteLink",
			mock.Anything,
			&networkmanager.DeleteLinkInput{
				GlobalNetworkId: link.GlobalNetworkID,
				LinkId:          link.ID,
			},
		).
		Return(&networkmanager.DeleteLinkOutput{}, nil)

	err := link.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

// --- Properties ---

func Test_Mock_NetworkManagerLink_Properties(t *testing.T) {
	assertions := assert.New(t)

	link := NetworkManagerLink{
		ID:              ptr.String("link-12345"),
		GlobalNetworkID: ptr.String("gn-99"),
		SiteID:          ptr.String("site-1"),
		Description:     ptr.String("primary link"),
		State:           ptr.String("AVAILABLE"),
		Tags: map[string]string{
			"Environment": "production",
			"Team":        "platform",
		},
	}

	properties := link.Properties()

	assertions.Equal("link-12345", properties.Get("ID"))
	assertions.Equal("gn-99", properties.Get("GlobalNetworkID"))
	assertions.Equal("site-1", properties.Get("SiteID"))
	assertions.Equal("primary link", properties.Get("Description"))
	assertions.Equal("AVAILABLE", properties.Get("State"))
	assertions.Equal("production", properties.Get("tag:Environment"))
	assertions.Equal("platform", properties.Get("tag:Team"))
}

func Test_Mock_NetworkManagerLink_Properties_EmptyTags(t *testing.T) {
	assertions := assert.New(t)

	link := NetworkManagerLink{
		ID:              ptr.String("link-99999"),
		GlobalNetworkID: ptr.String(TestGlobalNetworkID1),
		State:           ptr.String("AVAILABLE"),
		Tags:            map[string]string{},
	}

	properties := link.Properties()

	assertions.Equal("link-99999", properties.Get("ID"))
	assertions.Equal(TestGlobalNetworkID1, properties.Get("GlobalNetworkID"))
	assertions.Equal("AVAILABLE", properties.Get("State"))
}

// --- Display ---

func Test_Mock_NetworkManagerLink_String(t *testing.T) {
	assertions := assert.New(t)

	link := NetworkManagerLink{
		ID: ptr.String("link-abc123"),
	}

	assertions.Equal("link-abc123", link.String())
}

// --- Filter ---

func Test_Mock_NetworkManagerLink_Filter_Deleting(t *testing.T) {
	assertions := assert.New(t)

	link := NetworkManagerLink{
		ID:    ptr.String("link-1"),
		State: ptr.String("DELETING"),
	}

	err := link.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "already deleting")
}

func Test_Mock_NetworkManagerLink_Filter_Available(t *testing.T) {
	assertions := assert.New(t)

	link := NetworkManagerLink{
		ID:    ptr.String("link-1"),
		State: ptr.String("AVAILABLE"),
	}

	err := link.Filter()
	assertions.NoError(err)
}
