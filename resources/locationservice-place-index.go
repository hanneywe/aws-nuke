package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/location"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const LocationServicePlaceIndexResource = "LocationServicePlaceIndex"

func init() {
	registry.Register(&registry.Registration{
		Name:     LocationServicePlaceIndexResource,
		Scope:    nuke.Account,
		Resource: &LocationServicePlaceIndex{},
		Lister:   &LocationServicePlaceIndexLister{},
	})
}

type LocationServicePlaceIndexLister struct {
	svc LocationServiceClient
}

func (l *LocationServicePlaceIndexLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = location.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	paginator := location.NewListPlaceIndexesPaginator(svc, &location.ListPlaceIndexesInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.Entries {
			item := &resp.Entries[i]
			resources = append(resources, &LocationServicePlaceIndex{
				svc:         svc,
				IndexName:   item.IndexName,
				DataSource:  item.DataSource,
				Description: item.Description,
			})
		}
	}

	return resources, nil
}

type LocationServicePlaceIndex struct {
	svc         LocationServiceClient
	IndexName   *string
	DataSource  *string
	Description *string
}

func (r *LocationServicePlaceIndex) Remove(ctx context.Context) error {
	_, err := r.svc.DeletePlaceIndex(ctx, &location.DeletePlaceIndexInput{
		IndexName: r.IndexName,
	})
	return err
}

func (r *LocationServicePlaceIndex) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *LocationServicePlaceIndex) String() string {
	return *r.IndexName
}
