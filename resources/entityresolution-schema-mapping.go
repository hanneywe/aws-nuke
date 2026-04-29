package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/entityresolution"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EntityResolutionSchemaMappingResource = "EntityResolutionSchemaMapping"

func init() {
	registry.Register(&registry.Registration{
		Name:     EntityResolutionSchemaMappingResource,
		Scope:    nuke.Account,
		Resource: &EntityResolutionSchemaMapping{},
		Lister:   &EntityResolutionSchemaMappingLister{},
	})
}

type EntityResolutionSchemaMappingLister struct {
	svc EntityresolutionClient
}

func (l *EntityResolutionSchemaMappingLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = entityresolution.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	paginator := entityresolution.NewListSchemaMappingsPaginator(svc, &entityresolution.ListSchemaMappingsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.SchemaList {
			item := &resp.SchemaList[i]
			resources = append(resources, &EntityResolutionSchemaMapping{
				svc:        svc,
				SchemaName: item.SchemaName,
				SchemaArn:  item.SchemaArn,
			})
		}
	}

	return resources, nil
}

type EntityResolutionSchemaMapping struct {
	svc        EntityresolutionClient
	SchemaName *string
	SchemaArn  *string
}

func (r *EntityResolutionSchemaMapping) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteSchemaMapping(ctx, &entityresolution.DeleteSchemaMappingInput{
		SchemaName: r.SchemaName,
	})
	return err
}

func (r *EntityResolutionSchemaMapping) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EntityResolutionSchemaMapping) String() string {
	return *r.SchemaName
}
