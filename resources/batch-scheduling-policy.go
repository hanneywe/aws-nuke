package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/batch"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const BatchSchedulingPolicyResource = "BatchSchedulingPolicy"

func init() {
	registry.Register(&registry.Registration{
		Name:     BatchSchedulingPolicyResource,
		Scope:    nuke.Account,
		Resource: &BatchSchedulingPolicy{},
		Lister:   &BatchSchedulingPolicyLister{},
	})
}

type BatchSchedulingPolicyLister struct {
	svc BatchClient
}

func (l *BatchSchedulingPolicyLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = batch.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &batch.ListSchedulingPoliciesInput{}
	for {
		resp, err := svc.ListSchedulingPolicies(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, sp := range resp.SchedulingPolicies {
			resources = append(resources, &BatchSchedulingPolicy{
				svc: svc,
				Arn: sp.Arn,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type BatchSchedulingPolicy struct {
	svc BatchClient
	Arn *string
}

func (r *BatchSchedulingPolicy) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteSchedulingPolicy(ctx, &batch.DeleteSchedulingPolicyInput{
		Arn: r.Arn,
	})
	return err
}

func (r *BatchSchedulingPolicy) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *BatchSchedulingPolicy) String() string {
	return *r.Arn
}
