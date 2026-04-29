package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/greengrass"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GreengrassSubscriptionDefinitionResource = "GreengrassSubscriptionDefinition"

func init() {
	registry.Register(&registry.Registration{
		Name:     GreengrassSubscriptionDefinitionResource,
		Scope:    nuke.Account,
		Resource: &GreengrassSubscriptionDefinition{},
		Lister:   &GreengrassSubscriptionDefinitionLister{},
	})
}

type GreengrassSubscriptionDefinitionLister struct {
	svc GreengrassClient
}

func (l *GreengrassSubscriptionDefinitionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = greengrass.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	var nextToken *string

	for {
		resp, err := svc.ListSubscriptionDefinitions(ctx, &greengrass.ListSubscriptionDefinitionsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, def := range resp.Definitions {
			resources = append(resources, &GreengrassSubscriptionDefinition{
				svc:  svc,
				ID:   def.Id,
				Name: def.Name,
			})
		}

		if resp.NextToken == nil || *resp.NextToken == "" {
			break
		}
		nextToken = resp.NextToken
	}

	return resources, nil
}

type GreengrassSubscriptionDefinition struct {
	svc  GreengrassClient
	ID   *string `property:"name=Id"`
	Name *string
}

func (r *GreengrassSubscriptionDefinition) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteSubscriptionDefinition(ctx, &greengrass.DeleteSubscriptionDefinitionInput{
		SubscriptionDefinitionId: r.ID,
	})
	return err
}

func (r *GreengrassSubscriptionDefinition) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GreengrassSubscriptionDefinition) String() string {
	return *r.ID
}
