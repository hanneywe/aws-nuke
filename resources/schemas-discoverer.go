package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/schemas"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SchemasDiscovererResource = "SchemasDiscoverer"

func init() {
	registry.Register(&registry.Registration{
		Name:     SchemasDiscovererResource,
		Scope:    nuke.Account,
		Resource: &SchemasDiscoverer{},
		Lister:   &SchemasDiscovererLister{},
	})
}

type SchemasDiscovererLister struct {
	svc SchemasClient
}

func (l *SchemasDiscovererLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = schemas.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &schemas.ListDiscoverersInput{}
	for {
		resp, err := svc.ListDiscoverers(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, d := range resp.Discoverers {
			resources = append(resources, &SchemasDiscoverer{
				svc:          svc,
				DiscovererID: d.DiscovererId,
				SourceArn:    d.SourceArn,
				Tags:         d.Tags,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type SchemasDiscoverer struct {
	svc          SchemasClient
	DiscovererID *string
	SourceArn    *string
	Tags         map[string]string
}

func (r *SchemasDiscoverer) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDiscoverer(ctx, &schemas.DeleteDiscovererInput{
		DiscovererId: r.DiscovererID,
	})
	return err
}

func (r *SchemasDiscoverer) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SchemasDiscoverer) String() string {
	return *r.DiscovererID
}
