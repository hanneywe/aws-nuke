package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/servicecatalogappregistry"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const AppRegistryAttributeGroupResource = "AppRegistryAttributeGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     AppRegistryAttributeGroupResource,
		Scope:    nuke.Account,
		Resource: &AppRegistryAttributeGroup{},
		Lister:   &AppRegistryAttributeGroupLister{},
	})
}

type AppRegistryAttributeGroupLister struct {
	svc AppRegistryClient
}

func (l *AppRegistryAttributeGroupLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = servicecatalogappregistry.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &servicecatalogappregistry.ListAttributeGroupsInput{}
	for {
		resp, err := svc.ListAttributeGroups(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, ag := range resp.AttributeGroups {
			resources = append(resources, &AppRegistryAttributeGroup{
				svc:  svc,
				ID:   ag.Id,
				Name: ag.Name,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type AppRegistryAttributeGroup struct {
	svc  AppRegistryClient
	ID   *string `property:"name=Id"`
	Name *string
}

func (r *AppRegistryAttributeGroup) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteAttributeGroup(ctx, &servicecatalogappregistry.DeleteAttributeGroupInput{
		AttributeGroup: r.ID,
	})
	return err
}

func (r *AppRegistryAttributeGroup) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *AppRegistryAttributeGroup) String() string {
	return *r.Name
}
