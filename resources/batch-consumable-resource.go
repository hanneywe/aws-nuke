package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/batch"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const BatchConsumableResourceResource = "BatchConsumableResource"

func init() {
	registry.Register(&registry.Registration{
		Name:     BatchConsumableResourceResource,
		Scope:    nuke.Account,
		Resource: &BatchConsumableResource{},
		Lister:   &BatchConsumableResourceLister{},
	})
}

type BatchConsumableResourceLister struct {
	svc BatchClient
}

func (l *BatchConsumableResourceLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = batch.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &batch.ListConsumableResourcesInput{}
	for {
		resp, err := svc.ListConsumableResources(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, cr := range resp.ConsumableResources {
			resources = append(resources, &BatchConsumableResource{
				svc:                    svc,
				ConsumableResourceArn:  cr.ConsumableResourceArn,
				ConsumableResourceName: cr.ConsumableResourceName,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type BatchConsumableResource struct {
	svc                    BatchClient
	ConsumableResourceArn  *string
	ConsumableResourceName *string
}

func (r *BatchConsumableResource) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteConsumableResource(ctx, &batch.DeleteConsumableResourceInput{
		ConsumableResource: r.ConsumableResourceArn,
	})
	return err
}

func (r *BatchConsumableResource) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *BatchConsumableResource) String() string {
	return *r.ConsumableResourceName
}
