package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/datasync"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const DataSyncLocationResource = "DataSyncLocation"

func init() {
	registry.Register(&registry.Registration{
		Name:     DataSyncLocationResource,
		Scope:    nuke.Account,
		Resource: &DataSyncLocation{},
		Lister:   &DataSyncLocationLister{},
	})
}

type DataSyncLocationLister struct {
	svc DataSyncClient
}

func (l *DataSyncLocationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = datasync.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := datasync.NewListLocationsPaginator(svc, &datasync.ListLocationsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, loc := range resp.Locations {
			resources = append(resources, &DataSyncLocation{
				svc:         svc,
				LocationArn: loc.LocationArn,
				LocationURI: loc.LocationUri,
			})
		}
	}

	return resources, nil
}

type DataSyncLocation struct {
	svc         DataSyncClient
	LocationArn *string
	LocationURI *string `property:"name=LocationUri"`
}

func (r *DataSyncLocation) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteLocation(ctx, &datasync.DeleteLocationInput{
		LocationArn: r.LocationArn,
	})
	return err
}

func (r *DataSyncLocation) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *DataSyncLocation) String() string {
	return *r.LocationArn
}
