package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/configservice"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ConfigServiceRetentionConfigurationResource = "ConfigServiceRetentionConfiguration"

func init() {
	registry.Register(&registry.Registration{
		Name:     ConfigServiceRetentionConfigurationResource,
		Scope:    nuke.Account,
		Resource: &ConfigServiceRetentionConfiguration{},
		Lister:   &ConfigServiceRetentionConfigurationLister{},
	})
}

type ConfigServiceRetentionConfigurationLister struct {
	svc ConfigServiceClient
}

func (l *ConfigServiceRetentionConfigurationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = configservice.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &configservice.DescribeRetentionConfigurationsInput{}
	for {
		resp, err := svc.DescribeRetentionConfigurations(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, rc := range resp.RetentionConfigurations {
			resources = append(resources, &ConfigServiceRetentionConfiguration{
				svc:  svc,
				Name: rc.Name,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type ConfigServiceRetentionConfiguration struct {
	svc  ConfigServiceClient
	Name *string
}

func (r *ConfigServiceRetentionConfiguration) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteRetentionConfiguration(ctx, &configservice.DeleteRetentionConfigurationInput{
		RetentionConfigurationName: r.Name,
	})
	return err
}

func (r *ConfigServiceRetentionConfiguration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ConfigServiceRetentionConfiguration) String() string {
	return *r.Name
}
