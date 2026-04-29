package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/mailmanager"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MailManagerRuleSetResource = "MailManagerRuleSet"

func init() {
	registry.Register(&registry.Registration{
		Name:     MailManagerRuleSetResource,
		Scope:    nuke.Account,
		Resource: &MailManagerRuleSet{},
		Lister:   &MailManagerRuleSetLister{},
	})
}

type MailManagerRuleSetLister struct {
	svc MailManagerClient
}

func (l *MailManagerRuleSetLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = mailmanager.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &mailmanager.ListRuleSetsInput{}

	for {
		output, err := svc.ListRuleSets(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, ruleSet := range output.RuleSets {
			resources = append(resources, &MailManagerRuleSet{
				svc:         svc,
				RuleSetID:   ruleSet.RuleSetId,
				RuleSetName: ruleSet.RuleSetName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

type MailManagerRuleSet struct {
	svc         MailManagerClient
	RuleSetID   *string `property:"name=RuleSetId"`
	RuleSetName *string
}

func (r *MailManagerRuleSet) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteRuleSet(ctx, &mailmanager.DeleteRuleSetInput{
		RuleSetId: r.RuleSetID,
	})
	return err
}

func (r *MailManagerRuleSet) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MailManagerRuleSet) String() string {
	return *r.RuleSetName
}
