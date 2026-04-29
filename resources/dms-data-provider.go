package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/databasemigrationservice"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const DMSDataProviderResource = "DMSDataProvider"

func init() {
	registry.Register(&registry.Registration{
		Name:     DMSDataProviderResource,
		Scope:    nuke.Account,
		Resource: &DMSDataProvider{},
		Lister:   &DMSDataProviderLister{},
	})
}

type DMSDataProviderLister struct {
	svc DMSClient
}

func (l *DMSDataProviderLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = databasemigrationservice.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &databasemigrationservice.DescribeDataProvidersInput{}
	for {
		resp, err := svc.DescribeDataProviders(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, dp := range resp.DataProviders {
			resources = append(resources, &DMSDataProvider{
				svc:              svc,
				DataProviderArn:  dp.DataProviderArn,
				DataProviderName: dp.DataProviderName,
			})
		}
		if resp.Marker == nil {
			break
		}
		params.Marker = resp.Marker
	}
	return resources, nil
}

type DMSDataProvider struct {
	svc              DMSClient
	DataProviderArn  *string
	DataProviderName *string
}

func (r *DMSDataProvider) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDataProvider(ctx, &databasemigrationservice.DeleteDataProviderInput{
		DataProviderIdentifier: r.DataProviderArn,
	})
	return err
}

func (r *DMSDataProvider) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *DMSDataProvider) String() string {
	return *r.DataProviderName
}
