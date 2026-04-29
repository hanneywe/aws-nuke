package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/glue"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GlueUsageProfileResource = "GlueUsageProfile"

func init() {
	registry.Register(&registry.Registration{
		Name:     GlueUsageProfileResource,
		Scope:    nuke.Account,
		Resource: &GlueUsageProfile{},
		Lister:   &GlueUsageProfileLister{},
	})
}

type GlueUsageProfileLister struct {
	svc GlueV2Client
}

func (l *GlueUsageProfileLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = glue.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &glue.ListUsageProfilesInput{}
	for {
		resp, err := svc.ListUsageProfiles(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, profile := range resp.Profiles {
			resources = append(resources, &GlueUsageProfile{
				svc:  svc,
				Name: profile.Name,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type GlueUsageProfile struct {
	svc  GlueV2Client
	Name *string
}

func (r *GlueUsageProfile) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteUsageProfile(ctx, &glue.DeleteUsageProfileInput{
		Name: r.Name,
	})
	return err
}

func (r *GlueUsageProfile) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GlueUsageProfile) String() string {
	return *r.Name
}
