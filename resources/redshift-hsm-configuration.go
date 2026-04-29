package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/redshift"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const RedshiftHsmConfigurationResource = "RedshiftHsmConfiguration"

func init() {
	registry.Register(&registry.Registration{
		Name:     RedshiftHsmConfigurationResource,
		Scope:    nuke.Account,
		Resource: &RedshiftHsmConfiguration{},
		Lister:   &RedshiftHsmConfigurationLister{},
	})
}

type RedshiftHsmConfigurationLister struct {
	svc RedshiftClient
}

func (l *RedshiftHsmConfigurationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = redshift.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := redshift.NewDescribeHsmConfigurationsPaginator(svc, &redshift.DescribeHsmConfigurationsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, hsmConfiguration := range output.HsmConfigurations {
			resources = append(resources, &RedshiftHsmConfiguration{
				svc:                        svc,
				HsmConfigurationIdentifier: hsmConfiguration.HsmConfigurationIdentifier,
			})
		}
	}

	return resources, nil
}

type RedshiftHsmConfiguration struct {
	svc                        RedshiftClient
	HsmConfigurationIdentifier *string
}

func (r *RedshiftHsmConfiguration) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteHsmConfiguration(ctx, &redshift.DeleteHsmConfigurationInput{
		HsmConfigurationIdentifier: r.HsmConfigurationIdentifier,
	})
	return err
}

func (r *RedshiftHsmConfiguration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *RedshiftHsmConfiguration) String() string {
	return *r.HsmConfigurationIdentifier
}
