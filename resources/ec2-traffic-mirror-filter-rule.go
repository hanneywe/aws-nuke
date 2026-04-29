package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EC2TrafficMirrorFilterRuleResource = "EC2TrafficMirrorFilterRule"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2TrafficMirrorFilterRuleResource,
		Scope:    nuke.Account,
		Resource: &EC2TrafficMirrorFilterRule{},
		Lister:   &EC2TrafficMirrorFilterRuleLister{},
	})
}

type EC2TrafficMirrorFilterRuleLister struct {
	svc EC2Client
}

func (l *EC2TrafficMirrorFilterRuleLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	filterPaginator := ec2.NewDescribeTrafficMirrorFiltersPaginator(svc,
		&ec2.DescribeTrafficMirrorFiltersInput{})

	for filterPaginator.HasMorePages() {
		filterResp, err := filterPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, filter := range filterResp.TrafficMirrorFilters {
			rules, err := l.listRulesForFilter(ctx, svc, filter.TrafficMirrorFilterId)
			if err != nil {
				return nil, err
			}
			resources = append(resources, rules...)
		}
	}

	return resources, nil
}

func (l *EC2TrafficMirrorFilterRuleLister) listRulesForFilter(
	ctx context.Context, svc EC2Client, filterID *string,
) ([]resource.Resource, error) {
	var resources []resource.Resource

	input := &ec2.DescribeTrafficMirrorFilterRulesInput{
		TrafficMirrorFilterId: filterID,
	}

	for {
		resp, err := svc.DescribeTrafficMirrorFilterRules(ctx, input)
		if err != nil {
			return nil, err
		}

		for i := range resp.TrafficMirrorFilterRules {
			rule := &resp.TrafficMirrorFilterRules[i]
			resources = append(resources, &EC2TrafficMirrorFilterRule{
				svc:                       svc,
				TrafficMirrorFilterRuleID: rule.TrafficMirrorFilterRuleId,
				TrafficMirrorFilterID:     rule.TrafficMirrorFilterId,
				TrafficDirection:          string(rule.TrafficDirection),
				RuleNumber:                rule.RuleNumber,
				RuleAction:                string(rule.RuleAction),
				Tags:                      rule.Tags,
			})
		}

		if resp.NextToken == nil {
			break
		}
		input.NextToken = resp.NextToken
	}

	return resources, nil
}

type EC2TrafficMirrorFilterRule struct {
	svc                       EC2Client
	TrafficMirrorFilterRuleID *string `property:"name=TrafficMirrorFilterRuleId"`
	TrafficMirrorFilterID     *string `property:"name=TrafficMirrorFilterId"`
	TrafficDirection          string
	RuleNumber                *int32
	RuleAction                string
	Tags                      []ec2types.Tag
}

func (r *EC2TrafficMirrorFilterRule) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteTrafficMirrorFilterRule(ctx, &ec2.DeleteTrafficMirrorFilterRuleInput{
		TrafficMirrorFilterRuleId: r.TrafficMirrorFilterRuleID,
	})
	return err
}

func (r *EC2TrafficMirrorFilterRule) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2TrafficMirrorFilterRule) String() string {
	return fmt.Sprintf("%s -> %s", *r.TrafficMirrorFilterID, *r.TrafficMirrorFilterRuleID)
}
