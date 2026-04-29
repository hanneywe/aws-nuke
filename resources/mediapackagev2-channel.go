package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/mediapackagev2"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MediaPackageV2ChannelResource = "MediaPackageV2Channel"

func init() {
	registry.Register(&registry.Registration{
		Name:     MediaPackageV2ChannelResource,
		Scope:    nuke.Account,
		Resource: &MediaPackageV2Channel{},
		Lister:   &MediaPackageV2ChannelLister{},
		DependsOn: []string{
			MediaPackageV2OriginEndpointResource,
		},
	})
}

type MediaPackageV2ChannelLister struct {
	svc MediaPackageV2Client
}

func (l *MediaPackageV2ChannelLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = mediapackagev2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	// First list all channel groups
	groupPaginator := mediapackagev2.NewListChannelGroupsPaginator(svc, &mediapackagev2.ListChannelGroupsInput{})
	for groupPaginator.HasMorePages() {
		groupResp, err := groupPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, cg := range groupResp.Items {
			// Then list channels per group
			channelPaginator := mediapackagev2.NewListChannelsPaginator(svc, &mediapackagev2.ListChannelsInput{
				ChannelGroupName: cg.ChannelGroupName,
			})
			for channelPaginator.HasMorePages() {
				channelResp, err := channelPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, ch := range channelResp.Items {
					resources = append(resources, &MediaPackageV2Channel{
						svc:              svc,
						ChannelName:      ch.ChannelName,
						ChannelGroupName: ch.ChannelGroupName,
					})
				}
			}
		}
	}
	return resources, nil
}

type MediaPackageV2Channel struct {
	svc              MediaPackageV2Client
	ChannelName      *string
	ChannelGroupName *string
}

func (r *MediaPackageV2Channel) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteChannel(ctx, &mediapackagev2.DeleteChannelInput{
		ChannelGroupName: r.ChannelGroupName,
		ChannelName:      r.ChannelName,
	})
	return err
}

func (r *MediaPackageV2Channel) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MediaPackageV2Channel) String() string {
	return *r.ChannelName
}
