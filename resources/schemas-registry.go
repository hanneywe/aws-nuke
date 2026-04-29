package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/schemas"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SchemasRegistryResource = "SchemasRegistry"

func init() {
	registry.Register(&registry.Registration{
		Name:     SchemasRegistryResource,
		Scope:    nuke.Account,
		Resource: &SchemasRegistry{},
		Lister:   &SchemasRegistryLister{},
	})
}

type SchemasRegistryLister struct {
	svc SchemasClient
}

func (l *SchemasRegistryLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = schemas.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := schemas.NewListRegistriesPaginator(svc, &schemas.ListRegistriesInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, r := range resp.Registries {
			resources = append(resources, &SchemasRegistry{
				svc:          svc,
				RegistryName: r.RegistryName,
				RegistryArn:  r.RegistryArn,
				Tags:         r.Tags,
			})
		}
	}
	return resources, nil
}

type SchemasRegistry struct {
	svc          SchemasClient
	RegistryName *string
	RegistryArn  *string
	Tags         map[string]string
}

// awsManagedRegistries are registries managed by AWS that cannot be deleted.
var awsManagedRegistries = []string{
	"aws.events",
	"discovered-schemas",
}

func (r *SchemasRegistry) Filter() error {
	for _, managed := range awsManagedRegistries {
		if strings.EqualFold(*r.RegistryName, managed) {
			return fmt.Errorf("cannot delete AWS-managed registry: %s", *r.RegistryName)
		}
	}
	return nil
}

func (r *SchemasRegistry) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteRegistry(ctx, &schemas.DeleteRegistryInput{
		RegistryName: r.RegistryName,
	})
	return err
}

func (r *SchemasRegistry) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SchemasRegistry) String() string {
	return *r.RegistryName
}
