package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ivsrealtime"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IVSRealtimeStorageConfigurationResource = "IVSRealtimeStorageConfiguration"

func init() {
	registry.Register(&registry.Registration{
		Name:     IVSRealtimeStorageConfigurationResource,
		Scope:    nuke.Account,
		Resource: &IVSRealtimeStorageConfiguration{},
		Lister:   &IVSRealtimeStorageConfigurationLister{},
	})
}

type IVSRealtimeStorageConfigurationLister struct {
	svc IVSRealtimeClient
}

func (l *IVSRealtimeStorageConfigurationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = ivsrealtime.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	paginator := ivsrealtime.NewListStorageConfigurationsPaginator(svc, &ivsrealtime.ListStorageConfigurationsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.StorageConfigurations {
			item := &resp.StorageConfigurations[i]
			resources = append(resources, &IVSRealtimeStorageConfiguration{
				svc:  svc,
				ARN:  item.Arn,
				Name: item.Name,
				Tags: item.Tags,
			})
		}
	}

	return resources, nil
}

type IVSRealtimeStorageConfiguration struct {
	svc  IVSRealtimeClient
	ARN  *string
	Name *string
	Tags map[string]string
}

func (r *IVSRealtimeStorageConfiguration) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteStorageConfiguration(ctx, &ivsrealtime.DeleteStorageConfigurationInput{
		Arn: r.ARN,
	})
	return err
}

func (r *IVSRealtimeStorageConfiguration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IVSRealtimeStorageConfiguration) String() string {
	return *r.ARN
}
