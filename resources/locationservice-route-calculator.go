package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/location"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const LocationServiceRouteCalculatorResource = "LocationServiceRouteCalculator"

func init() {
	registry.Register(&registry.Registration{
		Name:     LocationServiceRouteCalculatorResource,
		Scope:    nuke.Account,
		Resource: &LocationServiceRouteCalculator{},
		Lister:   &LocationServiceRouteCalculatorLister{},
	})
}

type LocationServiceRouteCalculatorLister struct {
	svc LocationServiceClient
}

func (l *LocationServiceRouteCalculatorLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = location.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := location.NewListRouteCalculatorsPaginator(svc, &location.ListRouteCalculatorsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, rc := range resp.Entries {
			resources = append(resources, &LocationServiceRouteCalculator{
				svc:            svc,
				CalculatorName: rc.CalculatorName,
			})
		}
	}

	return resources, nil
}

type LocationServiceRouteCalculator struct {
	svc            LocationServiceClient
	CalculatorName *string
}

func (r *LocationServiceRouteCalculator) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteRouteCalculator(ctx, &location.DeleteRouteCalculatorInput{
		CalculatorName: r.CalculatorName,
	})
	return err
}

func (r *LocationServiceRouteCalculator) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *LocationServiceRouteCalculator) String() string {
	return *r.CalculatorName
}
