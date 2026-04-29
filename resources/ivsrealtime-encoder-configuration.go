package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ivsrealtime"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IVSRealtimeEncoderConfigurationResource = "IVSRealtimeEncoderConfiguration"

func init() {
	registry.Register(&registry.Registration{
		Name:     IVSRealtimeEncoderConfigurationResource,
		Scope:    nuke.Account,
		Resource: &IVSRealtimeEncoderConfiguration{},
		Lister:   &IVSRealtimeEncoderConfigurationLister{},
	})
}

type IVSRealtimeEncoderConfigurationLister struct {
	svc IVSRealtimeClient
}

func (l *IVSRealtimeEncoderConfigurationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ivsrealtime.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ivsrealtime.NewListEncoderConfigurationsPaginator(svc, &ivsrealtime.ListEncoderConfigurationsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, ec := range resp.EncoderConfigurations {
			resources = append(resources, &IVSRealtimeEncoderConfiguration{
				svc:  svc,
				ARN:  ec.Arn,
				Name: ec.Name,
				Tags: ec.Tags,
			})
		}
	}

	return resources, nil
}

type IVSRealtimeEncoderConfiguration struct {
	svc  IVSRealtimeClient
	ARN  *string
	Name *string
	Tags map[string]string
}

func (r *IVSRealtimeEncoderConfiguration) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteEncoderConfiguration(ctx, &ivsrealtime.DeleteEncoderConfigurationInput{
		Arn: r.ARN,
	})
	return err
}

func (r *IVSRealtimeEncoderConfiguration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IVSRealtimeEncoderConfiguration) String() string {
	return *r.ARN
}
