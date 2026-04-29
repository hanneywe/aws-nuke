package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/networkmanager"
)

// NetworkManagerClient is a shared interface for Network Manager SDK v2 client methods
// used by multiple resources (Device, Site, Link, Connection).
// It enables mock testing of List and Remove operations.
type NetworkManagerClient interface {
	// Global Networks
	DescribeGlobalNetworks(ctx context.Context, params *networkmanager.DescribeGlobalNetworksInput,
		optFns ...func(*networkmanager.Options)) (*networkmanager.DescribeGlobalNetworksOutput, error)

	// Devices
	GetDevices(ctx context.Context, params *networkmanager.GetDevicesInput,
		optFns ...func(*networkmanager.Options)) (*networkmanager.GetDevicesOutput, error)
	DeleteDevice(ctx context.Context, params *networkmanager.DeleteDeviceInput,
		optFns ...func(*networkmanager.Options)) (*networkmanager.DeleteDeviceOutput, error)

	// Sites
	GetSites(ctx context.Context, params *networkmanager.GetSitesInput,
		optFns ...func(*networkmanager.Options)) (*networkmanager.GetSitesOutput, error)
	DeleteSite(ctx context.Context, params *networkmanager.DeleteSiteInput,
		optFns ...func(*networkmanager.Options)) (*networkmanager.DeleteSiteOutput, error)

	// Links
	GetLinks(ctx context.Context, params *networkmanager.GetLinksInput,
		optFns ...func(*networkmanager.Options)) (*networkmanager.GetLinksOutput, error)
	DeleteLink(ctx context.Context, params *networkmanager.DeleteLinkInput,
		optFns ...func(*networkmanager.Options)) (*networkmanager.DeleteLinkOutput, error)

	// Connections
	GetConnections(ctx context.Context, params *networkmanager.GetConnectionsInput,
		optFns ...func(*networkmanager.Options)) (*networkmanager.GetConnectionsOutput, error)
	DeleteConnection(ctx context.Context, params *networkmanager.DeleteConnectionInput,
		optFns ...func(*networkmanager.Options)) (*networkmanager.DeleteConnectionOutput, error)

	// Transit Gateway Registrations
	GetTransitGatewayRegistrations(ctx context.Context, params *networkmanager.GetTransitGatewayRegistrationsInput,
		optFns ...func(*networkmanager.Options)) (*networkmanager.GetTransitGatewayRegistrationsOutput, error)
	DeregisterTransitGateway(ctx context.Context, params *networkmanager.DeregisterTransitGatewayInput,
		optFns ...func(*networkmanager.Options)) (*networkmanager.DeregisterTransitGatewayOutput, error)
}

// listNetworkManagerGlobalNetworkIDs paginates through all global networks and returns their IDs.
// This helper reduces cyclomatic complexity in individual resource listers.
func listNetworkManagerGlobalNetworkIDs(
	ctx context.Context, svc NetworkManagerClient,
) ([]string, error) {
	var globalNetworkIDs []string
	gnParams := &networkmanager.DescribeGlobalNetworksInput{}
	for {
		resp, err := svc.DescribeGlobalNetworks(ctx, gnParams)
		if err != nil {
			return nil, err
		}

		for _, gn := range resp.GlobalNetworks {
			if gn.GlobalNetworkId != nil {
				globalNetworkIDs = append(globalNetworkIDs, *gn.GlobalNetworkId)
			}
		}

		if resp.NextToken == nil {
			break
		}
		gnParams.NextToken = resp.NextToken
	}
	return globalNetworkIDs, nil
}
