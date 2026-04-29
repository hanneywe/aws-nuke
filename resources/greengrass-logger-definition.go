package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/greengrass"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GreengrassLoggerDefinitionResource = "GreengrassLoggerDefinition"

func init() {
	registry.Register(&registry.Registration{
		Name:     GreengrassLoggerDefinitionResource,
		Scope:    nuke.Account,
		Resource: &GreengrassLoggerDefinition{},
		Lister:   &GreengrassLoggerDefinitionLister{},
	})
}

type GreengrassLoggerDefinitionLister struct {
	svc GreengrassClient
}

func (l *GreengrassLoggerDefinitionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = greengrass.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	var nextToken *string

	for {
		resp, err := svc.ListLoggerDefinitions(ctx, &greengrass.ListLoggerDefinitionsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, def := range resp.Definitions {
			resources = append(resources, &GreengrassLoggerDefinition{
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

type GreengrassLoggerDefinition struct {
	svc  GreengrassClient
	ID   *string `property:"name=Id"`
	Name *string
}

func (r *GreengrassLoggerDefinition) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteLoggerDefinition(ctx, &greengrass.DeleteLoggerDefinitionInput{
		LoggerDefinitionId: r.ID,
	})
	return err
}

func (r *GreengrassLoggerDefinition) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GreengrassLoggerDefinition) String() string {
	return *r.ID
}
