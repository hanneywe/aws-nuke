package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/mediapackagev2"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MediaPackageV2ChannelGroupResource = "MediaPackageV2ChannelGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     MediaPackageV2ChannelGroupResource,
		Scope:    nuke.Account,
		Resource: &MediaPackageV2ChannelGroup{},
		Lister:   &MediaPackageV2ChannelGroupLister{},
		DependsOn: []string{
			MediaPackageV2ChannelResource,
		},
	})
}

type MediaPackageV2ChannelGroupLister struct {
	svc MediaPackageV2Client
}

func (l *MediaPackageV2ChannelGroupLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = mediapackagev2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := mediapackagev2.NewListChannelGroupsPaginator(svc, &mediapackagev2.ListChannelGroupsInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, cg := range resp.Items {
			resources = append(resources, &MediaPackageV2ChannelGroup{
				svc:              svc,
				ChannelGroupName: cg.ChannelGroupName,
				ARN:              cg.Arn,
			})
		}
	}
	return resources, nil
}

type MediaPackageV2ChannelGroup struct {
	svc              MediaPackageV2Client
	ChannelGroupName *string
	ARN              *string
}

func (r *MediaPackageV2ChannelGroup) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteChannelGroup(ctx, &mediapackagev2.DeleteChannelGroupInput{
		ChannelGroupName: r.ChannelGroupName,
	})
	return err
}

func (r *MediaPackageV2ChannelGroup) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MediaPackageV2ChannelGroup) String() string {
	return *r.ChannelGroupName
}
