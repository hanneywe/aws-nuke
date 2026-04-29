package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecr"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ECRPullThroughCacheRuleResource = "ECRPullThroughCacheRule"

func init() {
	registry.Register(&registry.Registration{
		Name:     ECRPullThroughCacheRuleResource,
		Scope:    nuke.Account,
		Resource: &ECRPullThroughCacheRule{},
		Lister:   &ECRPullThroughCacheRuleLister{},
	})
}

type ECRPullThroughCacheRuleLister struct {
	svc ECRv2Client
}

func (l *ECRPullThroughCacheRuleLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = ecr.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &ecr.DescribePullThroughCacheRulesInput{}
	for {
		resp, err := svc.DescribePullThroughCacheRules(ctx, params)
		if err != nil {
			return nil, err
		}
		for i := range resp.PullThroughCacheRules {
			resources = append(resources, &ECRPullThroughCacheRule{
				svc:                 svc,
				EcrRepositoryPrefix: resp.PullThroughCacheRules[i].EcrRepositoryPrefix,
				UpstreamRegistryURL: resp.PullThroughCacheRules[i].UpstreamRegistryUrl,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type ECRPullThroughCacheRule struct {
	svc                 ECRv2Client
	EcrRepositoryPrefix *string
	UpstreamRegistryURL *string `property:"name=UpstreamRegistryUrl"`
}

func (r *ECRPullThroughCacheRule) Remove(ctx context.Context) error {
	_, err := r.svc.DeletePullThroughCacheRule(ctx, &ecr.DeletePullThroughCacheRuleInput{
		EcrRepositoryPrefix: r.EcrRepositoryPrefix,
	})
	return err
}

func (r *ECRPullThroughCacheRule) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ECRPullThroughCacheRule) String() string {
	return *r.EcrRepositoryPrefix
}
