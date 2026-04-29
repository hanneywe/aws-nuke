package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/greengrass"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GreengrassResourceDefinitionResource = "GreengrassResourceDefinition"

func init() {
	registry.Register(&registry.Registration{
		Name:     GreengrassResourceDefinitionResource,
		Scope:    nuke.Account,
		Resource: &GreengrassResourceDefinition{},
		Lister:   &GreengrassResourceDefinitionLister{},
	})
}

type GreengrassResourceDefinitionLister struct {
	svc GreengrassClient
}

func (l *GreengrassResourceDefinitionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = greengrass.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	var nextToken *string

	for {
		resp, err := svc.ListResourceDefinitions(ctx, &greengrass.ListResourceDefinitionsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, definition := range resp.Definitions {
			resources = append(resources, &GreengrassResourceDefinition{
				svc:  svc,
				ID:   definition.Id,
				Name: definition.Name,
			})
		}

		if resp.NextToken == nil || *resp.NextToken == "" {
			break
		}
		nextToken = resp.NextToken
	}

	return resources, nil
}

type GreengrassResourceDefinition struct {
	svc  GreengrassClient
	ID   *string `property:"name=Id"`
	Name *string
}

func (r *GreengrassResourceDefinition) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteResourceDefinition(ctx, &greengrass.DeleteResourceDefinitionInput{
		ResourceDefinitionId: r.ID,
	})
	return err
}

func (r *GreengrassResourceDefinition) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GreengrassResourceDefinition) String() string {
	return *r.ID
}
