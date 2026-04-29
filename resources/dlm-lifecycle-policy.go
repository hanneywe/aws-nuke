package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dlm"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const DLMLifecyclePolicyResource = "DLMLifecyclePolicy"

func init() {
	registry.Register(&registry.Registration{
		Name:     DLMLifecyclePolicyResource,
		Scope:    nuke.Account,
		Resource: &DLMLifecyclePolicy{},
		Lister:   &DLMLifecyclePolicyLister{},
	})
}

type DLMLifecyclePolicyLister struct {
	svc DlmClient
}

func (l *DLMLifecyclePolicyLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = dlm.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	resp, err := svc.GetLifecyclePolicies(ctx, &dlm.GetLifecyclePoliciesInput{})
	if err != nil {
		return nil, err
	}

	for i := range resp.Policies {
		item := &resp.Policies[i]
		resources = append(resources, &DLMLifecyclePolicy{
			svc:         svc,
			PolicyID:    item.PolicyId,
			Description: item.Description,
			Tags:        item.Tags,
		})
	}

	return resources, nil
}

type DLMLifecyclePolicy struct {
	svc         DlmClient
	PolicyID    *string
	Description *string
	Tags        map[string]string
}

func (r *DLMLifecyclePolicy) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteLifecyclePolicy(ctx, &dlm.DeleteLifecyclePolicyInput{
		PolicyId: r.PolicyID,
	})
	return err
}

func (r *DLMLifecyclePolicy) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *DLMLifecyclePolicy) String() string {
	return *r.PolicyID
}
