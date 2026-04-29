package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iotwireless"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IoTWirelessWirelessGatewayResource = "IoTWirelessWirelessGateway"

func init() {
	registry.Register(&registry.Registration{
		Name:     IoTWirelessWirelessGatewayResource,
		Scope:    nuke.Account,
		Resource: &IoTWirelessWirelessGateway{},
		Lister:   &IoTWirelessWirelessGatewayLister{},
	})
}

type IoTWirelessWirelessGatewayLister struct {
	svc IoTWirelessClient
}

func (l *IoTWirelessWirelessGatewayLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = iotwireless.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := iotwireless.NewListWirelessGatewaysPaginator(svc, &iotwireless.ListWirelessGatewaysInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, gw := range resp.WirelessGatewayList {
			resources = append(resources, &IoTWirelessWirelessGateway{
				svc:  svc,
				ID:   gw.Id,
				Name: gw.Name,
				ARN:  gw.Arn,
			})
		}
	}

	return resources, nil
}

type IoTWirelessWirelessGateway struct {
	svc  IoTWirelessClient
	ID   *string `property:"name=Id"`
	Name *string
	ARN  *string `property:"name=Arn"`
}

func (r *IoTWirelessWirelessGateway) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteWirelessGateway(ctx, &iotwireless.DeleteWirelessGatewayInput{
		Id: r.ID,
	})
	return err
}

func (r *IoTWirelessWirelessGateway) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IoTWirelessWirelessGateway) String() string {
	return *r.ID
}
