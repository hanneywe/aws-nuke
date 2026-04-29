package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/directconnect"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const DirectConnectGatewayResource = "DirectConnectGateway"

func init() {
	registry.Register(&registry.Registration{
		Name:     DirectConnectGatewayResource,
		Scope:    nuke.Account,
		Resource: &DirectConnectGateway{},
		Lister:   &DirectConnectGatewayLister{},
	})
}

type DirectConnectGatewayLister struct {
	svc DirectConnectClient
}

func (l *DirectConnectGatewayLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = directconnect.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	var nextToken *string

	for {
		resp, err := svc.DescribeDirectConnectGateways(ctx, &directconnect.DescribeDirectConnectGatewaysInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for i := range resp.DirectConnectGateways {
			gw := resp.DirectConnectGateways[i]
			resources = append(resources, &DirectConnectGateway{
				svc:                      svc,
				DirectConnectGatewayID:   gw.DirectConnectGatewayId,
				DirectConnectGatewayName: gw.DirectConnectGatewayName,
			})
		}

		if resp.NextToken == nil {
			break
		}
		nextToken = resp.NextToken
	}

	return resources, nil
}

type DirectConnectGateway struct {
	svc                      DirectConnectClient
	DirectConnectGatewayID   *string `property:"name=DirectConnectGatewayId"`
	DirectConnectGatewayName *string
}

func (r *DirectConnectGateway) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDirectConnectGateway(ctx, &directconnect.DeleteDirectConnectGatewayInput{
		DirectConnectGatewayId: r.DirectConnectGatewayID,
	})
	return err
}

func (r *DirectConnectGateway) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *DirectConnectGateway) String() string {
	return *r.DirectConnectGatewayID
}
