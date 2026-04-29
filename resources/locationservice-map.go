package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/location"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const LocationServiceMapResource = "LocationServiceMap"

func init() {
	registry.Register(&registry.Registration{
		Name:     LocationServiceMapResource,
		Scope:    nuke.Account,
		Resource: &LocationServiceMap{},
		Lister:   &LocationServiceMapLister{},
	})
}

type LocationServiceMapLister struct {
	svc LocationServiceClient
}

func (l *LocationServiceMapLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = location.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := location.NewListMapsPaginator(svc, &location.ListMapsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, m := range resp.Entries {
			resources = append(resources, &LocationServiceMap{
				svc:     svc,
				MapName: m.MapName,
			})
		}
	}

	return resources, nil
}

type LocationServiceMap struct {
	svc     LocationServiceClient
	MapName *string
}

func (r *LocationServiceMap) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteMap(ctx, &location.DeleteMapInput{
		MapName: r.MapName,
	})
	return err
}

func (r *LocationServiceMap) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *LocationServiceMap) String() string {
	return *r.MapName
}
