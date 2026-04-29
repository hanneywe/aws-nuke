package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ivsrealtime"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IVSRealtimeIngestConfigurationResource = "IVSRealtimeIngestConfiguration"

func init() {
	registry.Register(&registry.Registration{
		Name:     IVSRealtimeIngestConfigurationResource,
		Scope:    nuke.Account,
		Resource: &IVSRealtimeIngestConfiguration{},
		Lister:   &IVSRealtimeIngestConfigurationLister{},
	})
}

type IVSRealtimeIngestConfigurationLister struct {
	svc IVSRealtimeClient
}

func (l *IVSRealtimeIngestConfigurationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ivsrealtime.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ivsrealtime.NewListIngestConfigurationsPaginator(svc, &ivsrealtime.ListIngestConfigurationsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, ingestConfig := range resp.IngestConfigurations {
			resources = append(resources, &IVSRealtimeIngestConfiguration{
				svc:  svc,
				Arn:  ingestConfig.Arn,
				Name: ingestConfig.Name,
			})
		}
	}

	return resources, nil
}

type IVSRealtimeIngestConfiguration struct {
	svc  IVSRealtimeClient
	Arn  *string
	Name *string
}

func (r *IVSRealtimeIngestConfiguration) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteIngestConfiguration(ctx, &ivsrealtime.DeleteIngestConfigurationInput{
		Arn: r.Arn,
	})
	return err
}

func (r *IVSRealtimeIngestConfiguration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IVSRealtimeIngestConfiguration) String() string {
	return *r.Arn
}
