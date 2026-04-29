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

const (
	TestGlobalNetworkID1   = "gn-1"
	TestGlobalNetworkID2   = "gn-2"
	TestGlobalNetworkPage2 = "gn-page2"
)

// --- Listing ---

func Test_Mock_NetworkManagerConnection_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNetworkManagerClient)

	// One global network with one connection
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
		On("GetConnections", mock.Anything, mock.Anything).
		Return(
			&networkmanager.GetConnectionsOutput{
				Connections: []nmtypes.Connection{
					{
						ConnectionId:      ptr.String("conn-1"),
						GlobalNetworkId:   ptr.String(TestGlobalNetworkID1),
						LinkId:            ptr.String("link-1"),
						ConnectedLinkId:   ptr.String("link-2"),
						DeviceId:          ptr.String("dev-1"),
						ConnectedDeviceId: ptr.String("dev-2"),
						State:             nmtypes.ConnectionStateAvailable,
						Tags: []nmtypes.Tag{
							{Key: ptr.String("env"), Value: ptr.String("test")},
						},
					},
				},
			}, nil,
		)

	lister := &NetworkManagerConnectionLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testNetworkManagerListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	conn := resources[0].(*NetworkManagerConnection)
	assertions.Equal("conn-1", *conn.ID)
	assertions.Equal(TestGlobalNetworkID1, *conn.GlobalNetworkID)
	assertions.Equal("link-1", *conn.LinkID)
	assertions.Equal("link-2", *conn.SecondLinkID)
	assertions.Equal("dev-1", *conn.DeviceID)
	assertions.Equal("dev-2", *conn.SecondDeviceID)
	assertions.Equal("AVAILABLE", *conn.State)
	assertions.Equal("test", conn.Tags["env"])

	mockClient.AssertExpectations(t)
}

func Test_Mock_NetworkManagerConnection_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNetworkManagerClient)

	// No global networks means no connections
	mockClient.
		On("DescribeGlobalNetworks", mock.Anything, mock.Anything).
		Return(
			&networkmanager.DescribeGlobalNetworksOutput{
				GlobalNetworks: []nmtypes.GlobalNetwork{},
			}, nil,
		)

	lister := &NetworkManagerConnectionLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testNetworkManagerListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NetworkManagerConnection_List_MultiPage(t *testing.T) {
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

	// gn-1: 100 connections on page 1, 1 on page 2
	firstPageConns := make([]nmtypes.Connection, 100)
	for i := range firstPageConns {
		firstPageConns[i] = nmtypes.Connection{
			ConnectionId:    ptr.String(fmt.Sprintf("conn-%d", i)),
			GlobalNetworkId: ptr.String(TestGlobalNetworkID1),
			State:           nmtypes.ConnectionStateAvailable,
		}
	}

	mockClient.
		On(
			"GetConnections",
			mock.Anything,
			mock.MatchedBy(func(input *networkmanager.GetConnectionsInput) bool {
				return *input.GlobalNetworkId == TestGlobalNetworkID1 && input.NextToken == nil
			}),
		).
		Return(
			&networkmanager.GetConnectionsOutput{
				Connections: firstPageConns,
				NextToken:   ptr.String("conn-page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"GetConnections",
			mock.Anything,
			mock.MatchedBy(func(input *networkmanager.GetConnectionsInput) bool {
				return *input.GlobalNetworkId == TestGlobalNetworkID1 && input.NextToken != nil
			}),
		).
		Return(
			&networkmanager.GetConnectionsOutput{
				Connections: []nmtypes.Connection{
					{
						ConnectionId:    ptr.String("conn-100"),
						GlobalNetworkId: ptr.String(TestGlobalNetworkID1),
						State:           nmtypes.ConnectionStateAvailable,
					},
				},
			}, nil,
		).
		Once()

	// gn-2: no connections
	mockClient.
		On(
			"GetConnections",
			mock.Anything,
			mock.MatchedBy(func(input *networkmanager.GetConnectionsInput) bool {
				return *input.GlobalNetworkId == TestGlobalNetworkID2
			}),
		).
		Return(
			&networkmanager.GetConnectionsOutput{
				Connections: []nmtypes.Connection{},
			}, nil,
		).
		Once()

	lister := &NetworkManagerConnectionLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testNetworkManagerListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

// --- Removal ---

func Test_Mock_NetworkManagerConnection_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNetworkManagerClient)

	conn := &NetworkManagerConnection{
		svc:             mockClient,
		ID:              ptr.String("conn-1"),
		GlobalNetworkID: ptr.String(TestGlobalNetworkID1),
	}

	mockClient.
		On(
			"DeleteConnection",
			mock.Anything,
			&networkmanager.DeleteConnectionInput{
				GlobalNetworkId: conn.GlobalNetworkID,
				ConnectionId:    conn.ID,
			},
		).
		Return(&networkmanager.DeleteConnectionOutput{}, nil)

	err := conn.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

// --- Properties ---

func Test_Mock_NetworkManagerConnection_Properties(t *testing.T) {
	assertions := assert.New(t)

	conn := NetworkManagerConnection{
		ID:              ptr.String("conn-12345"),
		GlobalNetworkID: ptr.String("gn-99"),
		LinkID:          ptr.String("link-1"),
		SecondLinkID:    ptr.String("link-2"),
		DeviceID:        ptr.String("dev-1"),
		SecondDeviceID:  ptr.String("dev-2"),
		State:           ptr.String("AVAILABLE"),
		Tags: map[string]string{
			"Environment": "production",
			"Team":        "platform",
		},
	}

	properties := conn.Properties()

	assertions.Equal("conn-12345", properties.Get("ID"))
	assertions.Equal("gn-99", properties.Get("GlobalNetworkID"))
	assertions.Equal("link-1", properties.Get("LinkID"))
	assertions.Equal("link-2", properties.Get("SecondLinkID"))
	assertions.Equal("dev-1", properties.Get("DeviceID"))
	assertions.Equal("dev-2", properties.Get("SecondDeviceID"))
	assertions.Equal("AVAILABLE", properties.Get("State"))
	assertions.Equal("production", properties.Get("tag:Environment"))
	assertions.Equal("platform", properties.Get("tag:Team"))
}

func Test_Mock_NetworkManagerConnection_Properties_EmptyTags(t *testing.T) {
	assertions := assert.New(t)

	conn := NetworkManagerConnection{
		ID:              ptr.String("conn-99999"),
		GlobalNetworkID: ptr.String(TestGlobalNetworkID1),
		State:           ptr.String("AVAILABLE"),
		Tags:            map[string]string{},
	}

	properties := conn.Properties()

	assertions.Equal("conn-99999", properties.Get("ID"))
	assertions.Equal(TestGlobalNetworkID1, properties.Get("GlobalNetworkID"))
	assertions.Equal("AVAILABLE", properties.Get("State"))
}

// --- Display ---

func Test_Mock_NetworkManagerConnection_String(t *testing.T) {
	assertions := assert.New(t)

	conn := NetworkManagerConnection{
		ID: ptr.String("conn-abc123"),
	}

	assertions.Equal("conn-abc123", conn.String())
}

// --- Filter ---

func Test_Mock_NetworkManagerConnection_Filter_Deleting(t *testing.T) {
	assertions := assert.New(t)

	conn := NetworkManagerConnection{
		ID:    ptr.String("conn-1"),
		State: ptr.String("DELETING"),
	}

	err := conn.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "already deleting")
}

func Test_Mock_NetworkManagerConnection_Filter_Available(t *testing.T) {
	assertions := assert.New(t)

	conn := NetworkManagerConnection{
		ID:    ptr.String("conn-1"),
		State: ptr.String("AVAILABLE"),
	}

	err := conn.Filter()
	assertions.NoError(err)
}
