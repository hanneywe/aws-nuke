package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/configservice"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ConfigServiceConfigurationAggregatorResource = "ConfigServiceConfigurationAggregator"

func init() {
	registry.Register(&registry.Registration{
		Name:     ConfigServiceConfigurationAggregatorResource,
		Scope:    nuke.Account,
		Resource: &ConfigServiceConfigurationAggregator{},
		Lister:   &ConfigServiceConfigurationAggregatorLister{},
	})
}

type ConfigServiceConfigurationAggregatorLister struct {
	svc ConfigServiceClient
}

func (l *ConfigServiceConfigurationAggregatorLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = configservice.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &configservice.DescribeConfigurationAggregatorsInput{}
	for {
		resp, err := svc.DescribeConfigurationAggregators(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, agg := range resp.ConfigurationAggregators {
			tags := make(map[string]string)
			if agg.ConfigurationAggregatorArn != nil {
				tagsResp, err := svc.ListTagsForResource(ctx, &configservice.ListTagsForResourceInput{
					ResourceArn: agg.ConfigurationAggregatorArn,
				})
				if err == nil {
					for _, t := range tagsResp.Tags {
						if t.Key != nil && t.Value != nil {
							tags[*t.Key] = *t.Value
						}
					}
				}
			}
			resources = append(resources, &ConfigServiceConfigurationAggregator{
				svc:  svc,
				Name: agg.ConfigurationAggregatorName,
				Tags: tags,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type ConfigServiceConfigurationAggregator struct {
	svc  ConfigServiceClient
	Name *string
	Tags map[string]string
}

func (r *ConfigServiceConfigurationAggregator) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteConfigurationAggregator(ctx, &configservice.DeleteConfigurationAggregatorInput{
		ConfigurationAggregatorName: r.Name,
	})
	return err
}

func (r *ConfigServiceConfigurationAggregator) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ConfigServiceConfigurationAggregator) String() string {
	return *r.Name
}
