package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/greengrass"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GreengrassFunctionDefinitionResource = "GreengrassFunctionDefinition"

func init() {
	registry.Register(&registry.Registration{
		Name:     GreengrassFunctionDefinitionResource,
		Scope:    nuke.Account,
		Resource: &GreengrassFunctionDefinition{},
		Lister:   &GreengrassFunctionDefinitionLister{},
	})
}

type GreengrassFunctionDefinitionLister struct {
	svc GreengrassClient
}

func (l *GreengrassFunctionDefinitionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = greengrass.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	var nextToken *string

	for {
		resp, err := svc.ListFunctionDefinitions(ctx, &greengrass.ListFunctionDefinitionsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, def := range resp.Definitions {
			resources = append(resources, &GreengrassFunctionDefinition{
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

type GreengrassFunctionDefinition struct {
	svc  GreengrassClient
	ID   *string `property:"name=Id"`
	Name *string
}

func (r *GreengrassFunctionDefinition) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteFunctionDefinition(ctx, &greengrass.DeleteFunctionDefinitionInput{
		FunctionDefinitionId: r.ID,
	})
	return err
}

func (r *GreengrassFunctionDefinition) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GreengrassFunctionDefinition) String() string {
	return *r.ID
}
