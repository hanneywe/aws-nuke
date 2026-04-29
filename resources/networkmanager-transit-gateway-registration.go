package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/networkmanager"
	nmtypes "github.com/aws/aws-sdk-go-v2/service/networkmanager/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const NetworkManagerTransitGatewayRegistrationResource = "NetworkManagerTransitGatewayRegistration"

func init() {
	registry.Register(&registry.Registration{
		Name:     NetworkManagerTransitGatewayRegistrationResource,
		Scope:    nuke.Account,
		Resource: &NetworkManagerTransitGatewayRegistration{},
		Lister:   &NetworkManagerTransitGatewayRegistrationLister{},
	})
}

type NetworkManagerTransitGatewayRegistrationLister struct {
	svc NetworkManagerClient
}

func (l *NetworkManagerTransitGatewayRegistrationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = networkmanager.NewFromConfig(*opts.Config)
	}

	globalNetworkIDs, err := listNetworkManagerGlobalNetworkIDs(ctx, svc)
	if err != nil {
		return nil, err
	}

	var resources []resource.Resource

	for _, gnID := range globalNetworkIDs {
		params := &networkmanager.GetTransitGatewayRegistrationsInput{
			GlobalNetworkId: &gnID,
		}
		for {
			resp, err := svc.GetTransitGatewayRegistrations(ctx, params)
			if err != nil {
				return nil, err
			}

			for i := range resp.TransitGatewayRegistrations {
				reg := &resp.TransitGatewayRegistrations[i]
				var state string
				if reg.State != nil {
					state = string(reg.State.Code)
				}
				resources = append(resources, &NetworkManagerTransitGatewayRegistration{
					svc:               svc,
					GlobalNetworkID:   reg.GlobalNetworkId,
					TransitGatewayArn: reg.TransitGatewayArn,
					State:             state,
				})
			}

			if resp.NextToken == nil {
				break
			}
			params.NextToken = resp.NextToken
		}
	}

	return resources, nil
}

type NetworkManagerTransitGatewayRegistration struct {
	svc               NetworkManagerClient
	GlobalNetworkID   *string
	TransitGatewayArn *string
	State             string
}

func (r *NetworkManagerTransitGatewayRegistration) Remove(ctx context.Context) error {
	_, err := r.svc.DeregisterTransitGateway(ctx, &networkmanager.DeregisterTransitGatewayInput{
		GlobalNetworkId:   r.GlobalNetworkID,
		TransitGatewayArn: r.TransitGatewayArn,
	})
	return err
}

func (r *NetworkManagerTransitGatewayRegistration) Filter() error {
	if r.State == string(nmtypes.TransitGatewayRegistrationStateDeleting) {
		return fmt.Errorf("already deregistering")
	}
	if r.State == string(nmtypes.TransitGatewayRegistrationStateDeleted) {
		return fmt.Errorf("already deregistered")
	}
	return nil
}

func (r *NetworkManagerTransitGatewayRegistration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *NetworkManagerTransitGatewayRegistration) String() string {
	return *r.TransitGatewayArn
}
