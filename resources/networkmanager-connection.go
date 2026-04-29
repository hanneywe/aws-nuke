package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/networkmanager"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const NetworkManagerConnectionResource = "NetworkManagerConnection"

func init() {
	registry.Register(&registry.Registration{
		Name:     NetworkManagerConnectionResource,
		Scope:    nuke.Account,
		Resource: &NetworkManagerConnection{},
		Lister:   &NetworkManagerConnectionLister{},
	})
}

type NetworkManagerConnectionLister struct {
	svc NetworkManagerClient
}

func (l *NetworkManagerConnectionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = networkmanager.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	// First, list all global networks
	globalNetworkIDs, err := listNetworkManagerGlobalNetworkIDs(ctx, svc)
	if err != nil {
		return nil, err
	}

	// For each global network, list all connections
	for _, gnID := range globalNetworkIDs {
		connParams := &networkmanager.GetConnectionsInput{
			GlobalNetworkId: &gnID,
		}
		for {
			resp, err := svc.GetConnections(ctx, connParams)
			if err != nil {
				return nil, err
			}

			for _, conn := range resp.Connections {
				tags := make(map[string]string)
				for _, tag := range conn.Tags {
					if tag.Key != nil && tag.Value != nil {
						tags[*tag.Key] = *tag.Value
					}
				}

				var state *string
				if conn.State != "" {
					s := string(conn.State)
					state = &s
				}

				resources = append(resources, &NetworkManagerConnection{
					svc:             svc,
					ID:              conn.ConnectionId,
					GlobalNetworkID: conn.GlobalNetworkId,
					LinkID:          conn.LinkId,
					SecondLinkID:    conn.ConnectedLinkId,
					DeviceID:        conn.DeviceId,
					SecondDeviceID:  conn.ConnectedDeviceId,
					State:           state,
					Tags:            tags,
				})
			}

			if resp.NextToken == nil {
				break
			}
			connParams.NextToken = resp.NextToken
		}
	}

	return resources, nil
}

type NetworkManagerConnection struct {
	svc             NetworkManagerClient
	ID              *string
	GlobalNetworkID *string
	LinkID          *string
	SecondLinkID    *string
	DeviceID        *string
	SecondDeviceID  *string
	State           *string
	Tags            map[string]string
}

func (r *NetworkManagerConnection) Filter() error {
	if r.State != nil && strings.EqualFold(*r.State, "DELETING") {
		return fmt.Errorf("already deleting")
	}
	return nil
}

func (r *NetworkManagerConnection) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteConnection(ctx, &networkmanager.DeleteConnectionInput{
		GlobalNetworkId: r.GlobalNetworkID,
		ConnectionId:    r.ID,
	})
	return err
}

func (r *NetworkManagerConnection) Properties() types.Properties {
	properties := types.NewPropertiesFromStruct(r)
	for key, val := range r.Tags {
		properties.SetTag(&key, val)
	}
	return properties
}

func (r *NetworkManagerConnection) String() string {
	return *r.ID
}
