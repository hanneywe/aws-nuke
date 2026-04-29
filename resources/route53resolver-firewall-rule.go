package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/route53resolver"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const Route53ResolverFirewallRuleResource = "Route53ResolverFirewallRule"

func init() {
	registry.Register(&registry.Registration{
		Name:     Route53ResolverFirewallRuleResource,
		Scope:    nuke.Account,
		Resource: &Route53ResolverFirewallRuleV2{},
		Lister:   &Route53ResolverFirewallRuleLister{},
	})
}

type Route53ResolverFirewallRuleLister struct {
	svc Route53ResolverV2Client
}

func (l *Route53ResolverFirewallRuleLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = route53resolver.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	ruleGroupParams := &route53resolver.ListFirewallRuleGroupsInput{}
	for {
		ruleGroupOutput, err := svc.ListFirewallRuleGroups(ctx, ruleGroupParams)
		if err != nil {
			return nil, err
		}

		for _, ruleGroup := range ruleGroupOutput.FirewallRuleGroups {
			ruleParams := &route53resolver.ListFirewallRulesInput{
				FirewallRuleGroupId: ruleGroup.Id,
			}
			for {
				ruleOutput, err := svc.ListFirewallRules(ctx, ruleParams)
				if err != nil {
					return nil, err
				}

				for i := range ruleOutput.FirewallRules {
					resources = append(resources, &Route53ResolverFirewallRuleV2{
						svc:                  svc,
						FirewallRuleGroupID:  ruleOutput.FirewallRules[i].FirewallRuleGroupId,
						FirewallDomainListID: ruleOutput.FirewallRules[i].FirewallDomainListId,
						Name:                 ruleOutput.FirewallRules[i].Name,
					})
				}

				if ruleOutput.NextToken == nil {
					break
				}
				ruleParams.NextToken = ruleOutput.NextToken
			}
		}

		if ruleGroupOutput.NextToken == nil {
			break
		}
		ruleGroupParams.NextToken = ruleGroupOutput.NextToken
	}

	return resources, nil
}

// Route53ResolverFirewallRuleV2 is the standalone resource for individual firewall rules.
// Named with V2 suffix to avoid conflict with the existing Route53ResolverFirewallRule helper struct
// in route53-resolver-firewall-rule-group.go.
type Route53ResolverFirewallRuleV2 struct {
	svc                  Route53ResolverV2Client
	FirewallRuleGroupID  *string
	FirewallDomainListID *string
	Name                 *string
}

func (r *Route53ResolverFirewallRuleV2) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteFirewallRule(ctx, &route53resolver.DeleteFirewallRuleInput{
		FirewallRuleGroupId:  r.FirewallRuleGroupID,
		FirewallDomainListId: r.FirewallDomainListID,
	})
	return err
}

func (r *Route53ResolverFirewallRuleV2) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *Route53ResolverFirewallRuleV2) String() string {
	return *r.Name
}
