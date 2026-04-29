package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iotwireless"
	iotwirelesstypes "github.com/aws/aws-sdk-go-v2/service/iotwireless/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IoTWirelessDestinationResource = "IoTWirelessDestination"

func init() {
	registry.Register(&registry.Registration{
		Name:     IoTWirelessDestinationResource,
		Scope:    nuke.Account,
		Resource: &IoTWirelessDestination{},
		Lister:   &IoTWirelessDestinationLister{},
	})
}

type IoTWirelessDestinationLister struct {
	svc IoTWirelessClient
}

func (l *IoTWirelessDestinationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = iotwireless.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := iotwireless.NewListDestinationsPaginator(svc, &iotwireless.ListDestinationsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, dest := range resp.DestinationList {
			resources = append(resources, &IoTWirelessDestination{
				svc:            svc,
				Name:           dest.Name,
				ExpressionType: dest.ExpressionType,
				RoleArn:        dest.RoleArn,
			})
		}
	}

	return resources, nil
}

type IoTWirelessDestination struct {
	svc            IoTWirelessClient
	Name           *string
	ExpressionType iotwirelesstypes.ExpressionType
	RoleArn        *string
}

func (r *IoTWirelessDestination) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDestination(ctx, &iotwireless.DeleteDestinationInput{
		Name: r.Name,
	})
	return err
}

func (r *IoTWirelessDestination) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IoTWirelessDestination) String() string {
	return *r.Name
}
