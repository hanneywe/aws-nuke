package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/glue"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GlueDataQualityRulesetResource = "GlueDataQualityRuleset"

func init() {
	registry.Register(&registry.Registration{
		Name:     GlueDataQualityRulesetResource,
		Scope:    nuke.Account,
		Resource: &GlueDataQualityRuleset{},
		Lister:   &GlueDataQualityRulesetLister{},
	})
}

type GlueDataQualityRulesetLister struct {
	svc GlueV2Client
}

func (l *GlueDataQualityRulesetLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = glue.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &glue.ListDataQualityRulesetsInput{}
	for {
		resp, err := svc.ListDataQualityRulesets(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, ruleset := range resp.Rulesets {
			resources = append(resources, &GlueDataQualityRuleset{
				svc:         svc,
				Name:        ruleset.Name,
				Description: ruleset.Description,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type GlueDataQualityRuleset struct {
	svc         GlueV2Client
	Name        *string
	Description *string
}

func (r *GlueDataQualityRuleset) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDataQualityRuleset(ctx, &glue.DeleteDataQualityRulesetInput{
		Name: r.Name,
	})
	return err
}

func (r *GlueDataQualityRuleset) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GlueDataQualityRuleset) String() string {
	return *r.Name
}
