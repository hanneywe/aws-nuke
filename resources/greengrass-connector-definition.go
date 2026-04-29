package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/greengrass"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GreengrassConnectorDefinitionResource = "GreengrassConnectorDefinition"

func init() {
	registry.Register(&registry.Registration{
		Name:     GreengrassConnectorDefinitionResource,
		Scope:    nuke.Account,
		Resource: &GreengrassConnectorDefinition{},
		Lister:   &GreengrassConnectorDefinitionLister{},
	})
}

type GreengrassConnectorDefinitionLister struct {
	svc GreengrassClient
}

func (l *GreengrassConnectorDefinitionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = greengrass.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	var nextToken *string

	for {
		resp, err := svc.ListConnectorDefinitions(ctx, &greengrass.ListConnectorDefinitionsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, def := range resp.Definitions {
			resources = append(resources, &GreengrassConnectorDefinition{
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

type GreengrassConnectorDefinition struct {
	svc  GreengrassClient
	ID   *string `property:"name=Id"`
	Name *string
}

func (r *GreengrassConnectorDefinition) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteConnectorDefinition(ctx, &greengrass.DeleteConnectorDefinitionInput{
		ConnectorDefinitionId: r.ID,
	})
	return err
}

func (r *GreengrassConnectorDefinition) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GreengrassConnectorDefinition) String() string {
	return *r.ID
}
