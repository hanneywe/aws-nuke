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

const NetworkManagerDeviceResource = "NetworkManagerDevice"

func init() {
	registry.Register(&registry.Registration{
		Name:     NetworkManagerDeviceResource,
		Scope:    nuke.Account,
		Resource: &NetworkManagerDevice{},
		Lister:   &NetworkManagerDeviceLister{},
	})
}

type NetworkManagerDeviceLister struct {
	svc NetworkManagerClient
}

func (l *NetworkManagerDeviceLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
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

	// For each global network, list all devices
	for _, gnID := range globalNetworkIDs {
		devParams := &networkmanager.GetDevicesInput{
			GlobalNetworkId: &gnID,
		}
		for {
			resp, err := svc.GetDevices(ctx, devParams)
			if err != nil {
				return nil, err
			}

			for i := range resp.Devices {
				tags := make(map[string]string)
				for _, tag := range resp.Devices[i].Tags {
					if tag.Key != nil && tag.Value != nil {
						tags[*tag.Key] = *tag.Value
					}
				}

				var state *string
				if resp.Devices[i].State != "" {
					s := string(resp.Devices[i].State)
					state = &s
				}

				resources = append(resources, &NetworkManagerDevice{
					svc:             svc,
					ID:              resp.Devices[i].DeviceId,
					GlobalNetworkID: resp.Devices[i].GlobalNetworkId,
					Description:     resp.Devices[i].Description,
					State:           state,
					Tags:            tags,
				})
			}

			if resp.NextToken == nil {
				break
			}
			devParams.NextToken = resp.NextToken
		}
	}

	return resources, nil
}

type NetworkManagerDevice struct {
	svc             NetworkManagerClient
	ID              *string
	GlobalNetworkID *string
	Description     *string
	State           *string
	Tags            map[string]string
}

func (r *NetworkManagerDevice) Filter() error {
	if r.State != nil && strings.EqualFold(*r.State, "DELETING") {
		return fmt.Errorf("already deleting")
	}
	return nil
}

func (r *NetworkManagerDevice) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDevice(ctx, &networkmanager.DeleteDeviceInput{
		GlobalNetworkId: r.GlobalNetworkID,
		DeviceId:        r.ID,
	})
	return err
}

func (r *NetworkManagerDevice) Properties() types.Properties {
	properties := types.NewPropertiesFromStruct(r)
	for key, val := range r.Tags {
		properties.SetTag(&key, val)
	}
	return properties
}

func (r *NetworkManagerDevice) String() string {
	return *r.ID
}
