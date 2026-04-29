package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iot"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IoTDimensionResource = "IoTDimension"

func init() {
	registry.Register(&registry.Registration{
		Name:     IoTDimensionResource,
		Scope:    nuke.Account,
		Resource: &IoTDimension{},
		Lister:   &IoTDimensionLister{},
	})
}

type IoTDimensionLister struct {
	svc IoTClient
}

func (l *IoTDimensionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = iot.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := iot.NewListDimensionsPaginator(svc, &iot.ListDimensionsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, dimensionName := range resp.DimensionNames {
			resources = append(resources, &IoTDimension{
				svc:  svc,
				Name: &dimensionName,
			})
		}
	}

	return resources, nil
}

type IoTDimension struct {
	svc  IoTClient
	Name *string
}

func (r *IoTDimension) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDimension(ctx, &iot.DeleteDimensionInput{
		Name: r.Name,
	})
	return err
}

func (r *IoTDimension) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IoTDimension) String() string {
	return *r.Name
}
