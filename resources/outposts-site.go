package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/outposts"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const OutpostsSiteResource = "OutpostsSite"

func init() {
	registry.Register(&registry.Registration{
		Name:     OutpostsSiteResource,
		Scope:    nuke.Account,
		Resource: &OutpostsSite{},
		Lister:   &OutpostsSiteLister{},
	})
}

type OutpostsSiteLister struct {
	svc OutpostsClient
}

func (l *OutpostsSiteLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = outposts.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := outposts.NewListSitesPaginator(svc, &outposts.ListSitesInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, site := range resp.Sites {
			resources = append(resources, &OutpostsSite{
				svc:    svc,
				SiteID: site.SiteId,
				Name:   site.Name,
			})
		}
	}
	return resources, nil
}

type OutpostsSite struct {
	svc    OutpostsClient
	SiteID *string `property:"name=SiteId"`
	Name   *string
}

func (r *OutpostsSite) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteSite(ctx, &outposts.DeleteSiteInput{
		SiteId: r.SiteID,
	})
	return err
}

func (r *OutpostsSite) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *OutpostsSite) String() string {
	return *r.Name
}
