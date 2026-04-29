package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/mediapackagevod"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MediaPackageVODPackagingGroupResource = "MediaPackageVODPackagingGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     MediaPackageVODPackagingGroupResource,
		Scope:    nuke.Account,
		Resource: &MediaPackageVODPackagingGroup{},
		Lister:   &MediaPackageVODPackagingGroupLister{},
	})
}

type MediaPackageVODPackagingGroupLister struct {
	svc MediaPackageVODClient
}

func (l *MediaPackageVODPackagingGroupLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = mediapackagevod.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := mediapackagevod.NewListPackagingGroupsPaginator(svc, &mediapackagevod.ListPackagingGroupsInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, pg := range resp.PackagingGroups {
			resources = append(resources, &MediaPackageVODPackagingGroup{
				svc: svc,
				ID:  pg.Id,
				Arn: pg.Arn,
			})
		}
	}
	return resources, nil
}

type MediaPackageVODPackagingGroup struct {
	svc MediaPackageVODClient
	ID  *string `property:"name=Id"`
	Arn *string
}

func (r *MediaPackageVODPackagingGroup) Remove(ctx context.Context) error {
	_, err := r.svc.DeletePackagingGroup(ctx, &mediapackagevod.DeletePackagingGroupInput{
		Id: r.ID,
	})
	return err
}

func (r *MediaPackageVODPackagingGroup) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MediaPackageVODPackagingGroup) String() string {
	return *r.ID
}
