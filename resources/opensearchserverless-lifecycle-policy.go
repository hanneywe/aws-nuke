package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/opensearchserverless"
	osstypes "github.com/aws/aws-sdk-go-v2/service/opensearchserverless/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const OpenSearchServerlessLifecyclePolicyResource = "OpenSearchServerlessLifecyclePolicy"

func init() {
	registry.Register(&registry.Registration{
		Name:     OpenSearchServerlessLifecyclePolicyResource,
		Scope:    nuke.Account,
		Resource: &OpenSearchServerlessLifecyclePolicy{},
		Lister:   &OpenSearchServerlessLifecyclePolicyLister{},
	})
}

type OpenSearchServerlessLifecyclePolicyLister struct {
	svc OpenSearchServerlessClient
}

func (l *OpenSearchServerlessLifecyclePolicyLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = opensearchserverless.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	// ListLifecyclePolicies requires a Type parameter. Iterate over all known types.
	for _, policyType := range osstypes.LifecyclePolicyTypeRetention.Values() {
		paginator := opensearchserverless.NewListLifecyclePoliciesPaginator(svc, &opensearchserverless.ListLifecyclePoliciesInput{
			Type: policyType,
		})

		for paginator.HasMorePages() {
			resp, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, item := range resp.LifecyclePolicySummaries {
				resources = append(resources, &OpenSearchServerlessLifecyclePolicy{
					svc:  svc,
					Name: item.Name,
					Type: item.Type,
				})
			}
		}
	}

	return resources, nil
}

type OpenSearchServerlessLifecyclePolicy struct {
	svc  OpenSearchServerlessClient
	Name *string
	Type osstypes.LifecyclePolicyType
}

func (r *OpenSearchServerlessLifecyclePolicy) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteLifecyclePolicy(ctx, &opensearchserverless.DeleteLifecyclePolicyInput{
		Name: r.Name,
		Type: r.Type,
	})
	return err
}

func (r *OpenSearchServerlessLifecyclePolicy) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *OpenSearchServerlessLifecyclePolicy) String() string {
	return *r.Name
}
