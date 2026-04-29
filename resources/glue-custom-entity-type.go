package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/glue"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GlueCustomEntityTypeResource = "GlueCustomEntityType"

func init() {
	registry.Register(&registry.Registration{
		Name:     GlueCustomEntityTypeResource,
		Scope:    nuke.Account,
		Resource: &GlueCustomEntityType{},
		Lister:   &GlueCustomEntityTypeLister{},
	})
}

type GlueCustomEntityTypeLister struct {
	svc GlueV2Client
}

func (l *GlueCustomEntityTypeLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = glue.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &glue.ListCustomEntityTypesInput{}
	for {
		resp, err := svc.ListCustomEntityTypes(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, cet := range resp.CustomEntityTypes {
			resources = append(resources, &GlueCustomEntityType{
				svc:  svc,
				Name: cet.Name,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type GlueCustomEntityType struct {
	svc  GlueV2Client
	Name *string
}

func (r *GlueCustomEntityType) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteCustomEntityType(ctx, &glue.DeleteCustomEntityTypeInput{
		Name: r.Name,
	})
	return err
}

func (r *GlueCustomEntityType) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GlueCustomEntityType) String() string {
	return *r.Name
}
