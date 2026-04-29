package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/greengrass"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GreengrassDeviceDefinitionResource = "GreengrassDeviceDefinition"

func init() {
	registry.Register(&registry.Registration{
		Name:     GreengrassDeviceDefinitionResource,
		Scope:    nuke.Account,
		Resource: &GreengrassDeviceDefinition{},
		Lister:   &GreengrassDeviceDefinitionLister{},
	})
}

type GreengrassDeviceDefinitionLister struct {
	svc GreengrassClient
}

func (l *GreengrassDeviceDefinitionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = greengrass.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	var nextToken *string

	for {
		resp, err := svc.ListDeviceDefinitions(ctx, &greengrass.ListDeviceDefinitionsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, def := range resp.Definitions {
			resources = append(resources, &GreengrassDeviceDefinition{
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

type GreengrassDeviceDefinition struct {
	svc  GreengrassClient
	ID   *string `property:"name=Id"`
	Name *string
}

func (r *GreengrassDeviceDefinition) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDeviceDefinition(ctx, &greengrass.DeleteDeviceDefinitionInput{
		DeviceDefinitionId: r.ID,
	})
	return err
}

func (r *GreengrassDeviceDefinition) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GreengrassDeviceDefinition) String() string {
	return *r.ID
}
