package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CloudWatchManagedInsightRulesResource = "CloudWatchManagedInsightRules"

func init() {
	registry.Register(&registry.Registration{
		Name:     CloudWatchManagedInsightRulesResource,
		Scope:    nuke.Account,
		Resource: &CloudWatchManagedInsightRules{},
		Lister:   &CloudWatchManagedInsightRulesLister{},
	})
}

type CloudWatchManagedInsightRulesLister struct {
	svc CloudWatchV2Client
}

func (l *CloudWatchManagedInsightRulesLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = cloudwatch.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := cloudwatch.NewDescribeInsightRulesPaginator(svc, &cloudwatch.DescribeInsightRulesInput{
		MaxResults: nil,
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range page.InsightRules {
			rule := &page.InsightRules[i]
			if rule.ManagedRule != nil && *rule.ManagedRule {
				resources = append(resources, &CloudWatchManagedInsightRules{
					svc:   svc,
					Name:  rule.Name,
					State: rule.State,
				})
			}
		}
	}

	return resources, nil
}

type CloudWatchManagedInsightRules struct {
	svc   CloudWatchV2Client
	Name  *string
	State *string
}

func (r *CloudWatchManagedInsightRules) Filter() error {
	if r.State != nil && *r.State == "DISABLED" {
		return fmt.Errorf("already disabled")
	}
	return nil
}

func (r *CloudWatchManagedInsightRules) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteInsightRules(ctx, &cloudwatch.DeleteInsightRulesInput{
		RuleNames: []string{*r.Name},
	})
	return err
}

func (r *CloudWatchManagedInsightRules) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CloudWatchManagedInsightRules) String() string {
	return *r.Name
}
