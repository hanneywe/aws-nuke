package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/greengrass"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GreengrassCoreDefinitionResource = "GreengrassCoreDefinition"

func init() {
	registry.Register(&registry.Registration{
		Name:     GreengrassCoreDefinitionResource,
		Scope:    nuke.Account,
		Resource: &GreengrassCoreDefinition{},
		Lister:   &GreengrassCoreDefinitionLister{},
	})
}

type GreengrassCoreDefinitionLister struct {
	svc GreengrassClient
}

func (l *GreengrassCoreDefinitionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = greengrass.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	var nextToken *string

	for {
		resp, err := svc.ListCoreDefinitions(ctx, &greengrass.ListCoreDefinitionsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, def := range resp.Definitions {
			resources = append(resources, &GreengrassCoreDefinition{
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

type GreengrassCoreDefinition struct {
	svc  GreengrassClient
	ID   *string `property:"name=Id"`
	Name *string
}

func (r *GreengrassCoreDefinition) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteCoreDefinition(ctx, &greengrass.DeleteCoreDefinitionInput{
		CoreDefinitionId: r.ID,
	})
	return err
}

func (r *GreengrassCoreDefinition) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GreengrassCoreDefinition) String() string {
	return *r.ID
}
