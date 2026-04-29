package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/glue"
	gluetypes "github.com/aws/aws-sdk-go-v2/service/glue/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GlueRegistryResource = "GlueRegistry"

func init() {
	registry.Register(&registry.Registration{
		Name:     GlueRegistryResource,
		Scope:    nuke.Account,
		Resource: &GlueRegistry{},
		Lister:   &GlueRegistryLister{},
	})
}

type GlueRegistryLister struct {
	svc GlueV2Client
}

func (l *GlueRegistryLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = glue.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &glue.ListRegistriesInput{}
	for {
		resp, err := svc.ListRegistries(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, reg := range resp.Registries {
			resources = append(resources, &GlueRegistry{
				svc:          svc,
				RegistryName: reg.RegistryName,
				RegistryArn:  reg.RegistryArn,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type GlueRegistry struct {
	svc          GlueV2Client
	RegistryName *string
	RegistryArn  *string
}

func (r *GlueRegistry) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteRegistry(ctx, &glue.DeleteRegistryInput{
		RegistryId: &gluetypes.RegistryId{
			RegistryArn: r.RegistryArn,
		},
	})
	return err
}

func (r *GlueRegistry) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GlueRegistry) String() string {
	return *r.RegistryName
}
