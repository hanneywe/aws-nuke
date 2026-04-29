package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/configservice"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ConfigServiceOrganizationConfigRuleResource = "ConfigServiceOrganizationConfigRule"

func init() {
	registry.Register(&registry.Registration{
		Name:     ConfigServiceOrganizationConfigRuleResource,
		Scope:    nuke.Account,
		Resource: &ConfigServiceOrganizationConfigRule{},
		Lister:   &ConfigServiceOrganizationConfigRuleLister{},
	})
}

type ConfigServiceOrganizationConfigRuleLister struct {
	svc ConfigServiceClient
}

func (l *ConfigServiceOrganizationConfigRuleLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = configservice.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &configservice.DescribeOrganizationConfigRulesInput{}
	for {
		resp, err := svc.DescribeOrganizationConfigRules(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, rule := range resp.OrganizationConfigRules {
			resources = append(resources, &ConfigServiceOrganizationConfigRule{
				svc:  svc,
				Name: rule.OrganizationConfigRuleName,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type ConfigServiceOrganizationConfigRule struct {
	svc  ConfigServiceClient
	Name *string
}

func (r *ConfigServiceOrganizationConfigRule) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteOrganizationConfigRule(ctx, &configservice.DeleteOrganizationConfigRuleInput{
		OrganizationConfigRuleName: r.Name,
	})
	return err
}

func (r *ConfigServiceOrganizationConfigRule) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ConfigServiceOrganizationConfigRule) String() string {
	return *r.Name
}
