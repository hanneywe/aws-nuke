package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ses"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SESReceiptRuleResource = "SESReceiptRule"

func init() {
	registry.Register(&registry.Registration{
		Name:     SESReceiptRuleResource,
		Scope:    nuke.Account,
		Resource: &SESReceiptRule{},
		Lister:   &SESReceiptRuleLister{},
	})
}

type SESReceiptRuleLister struct {
	svc SESClient
}

func (l *SESReceiptRuleLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = ses.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	ruleSetParams := &ses.ListReceiptRuleSetsInput{}
	for {
		ruleSetResp, err := svc.ListReceiptRuleSets(ctx, ruleSetParams)
		if err != nil {
			return nil, err
		}
		for _, ruleSet := range ruleSetResp.RuleSets {
			ruleResp, err := svc.DescribeReceiptRuleSet(ctx, &ses.DescribeReceiptRuleSetInput{
				RuleSetName: ruleSet.Name,
			})
			if err != nil {
				return nil, err
			}
			for _, rule := range ruleResp.Rules {
				resources = append(resources, &SESReceiptRule{
					svc:         svc,
					RuleSetName: ruleSet.Name,
					RuleName:    rule.Name,
					Enabled:     rule.Enabled,
				})
			}
		}
		if ruleSetResp.NextToken == nil {
			break
		}
		ruleSetParams.NextToken = ruleSetResp.NextToken
	}

	return resources, nil
}

type SESReceiptRule struct {
	svc         SESClient
	RuleSetName *string
	RuleName    *string
	Enabled     bool
}

func (r *SESReceiptRule) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteReceiptRule(ctx, &ses.DeleteReceiptRuleInput{
		RuleName:    r.RuleName,
		RuleSetName: r.RuleSetName,
	})
	return err
}

func (r *SESReceiptRule) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SESReceiptRule) String() string {
	return *r.RuleName
}
