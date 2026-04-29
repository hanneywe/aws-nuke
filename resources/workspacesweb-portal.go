package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/workspacesweb"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const WorkSpacesWebPortalResource = "WorkSpacesWebPortal"

func init() {
	registry.Register(&registry.Registration{
		Name:     WorkSpacesWebPortalResource,
		Scope:    nuke.Account,
		Resource: &WorkSpacesWebPortal{},
		Lister:   &WorkSpacesWebPortalLister{},
	})
}

type WorkSpacesWebPortalLister struct {
	svc WorkSpacesWebClient
}

func (l *WorkSpacesWebPortalLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = workspacesweb.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := workspacesweb.NewListPortalsPaginator(svc, &workspacesweb.ListPortalsInput{})
	for paginator.HasMorePages() {
		listPortalsOutput, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for idx := range listPortalsOutput.Portals {
			portal := &listPortalsOutput.Portals[idx]
			resources = append(resources, &WorkSpacesWebPortal{
				svc:         svc,
				PortalArn:   portal.PortalArn,
				DisplayName: portal.DisplayName,
			})
		}
	}

	return resources, nil
}

type WorkSpacesWebPortal struct {
	svc         WorkSpacesWebClient
	PortalArn   *string
	DisplayName *string
}

func (r *WorkSpacesWebPortal) Remove(ctx context.Context) error {
	_, err := r.svc.DeletePortal(ctx, &workspacesweb.DeletePortalInput{
		PortalArn: r.PortalArn,
	})
	return err
}

func (r *WorkSpacesWebPortal) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *WorkSpacesWebPortal) String() string {
	return *r.PortalArn
}
