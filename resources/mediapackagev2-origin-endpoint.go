package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/mediapackagev2"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MediaPackageV2OriginEndpointResource = "MediaPackageV2OriginEndpoint"

func init() {
	registry.Register(&registry.Registration{
		Name:     MediaPackageV2OriginEndpointResource,
		Scope:    nuke.Account,
		Resource: &MediaPackageV2OriginEndpoint{},
		Lister:   &MediaPackageV2OriginEndpointLister{},
	})
}

type MediaPackageV2OriginEndpointLister struct {
	svc MediaPackageV2Client
}

func (l *MediaPackageV2OriginEndpointLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = mediapackagev2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	groupPaginator := mediapackagev2.NewListChannelGroupsPaginator(svc, &mediapackagev2.ListChannelGroupsInput{})
	for groupPaginator.HasMorePages() {
		groupResp, err := groupPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, cg := range groupResp.Items {
			channelPaginator := mediapackagev2.NewListChannelsPaginator(svc, &mediapackagev2.ListChannelsInput{
				ChannelGroupName: cg.ChannelGroupName,
			})
			for channelPaginator.HasMorePages() {
				channelResp, err := channelPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, ch := range channelResp.Items {
					epPaginator := mediapackagev2.NewListOriginEndpointsPaginator(svc, &mediapackagev2.ListOriginEndpointsInput{
						ChannelGroupName: cg.ChannelGroupName,
						ChannelName:      ch.ChannelName,
					})
					for epPaginator.HasMorePages() {
						epResp, err := epPaginator.NextPage(ctx)
						if err != nil {
							return nil, err
						}
						for i := range epResp.Items {
							ep := &epResp.Items[i]
							resources = append(resources, &MediaPackageV2OriginEndpoint{
								svc:                svc,
								ChannelGroupName:   cg.ChannelGroupName,
								ChannelName:        ch.ChannelName,
								OriginEndpointName: ep.OriginEndpointName,
							})
						}
					}
				}
			}
		}
	}
	return resources, nil
}

type MediaPackageV2OriginEndpoint struct {
	svc                MediaPackageV2Client
	ChannelGroupName   *string
	ChannelName        *string
	OriginEndpointName *string
}

func (r *MediaPackageV2OriginEndpoint) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteOriginEndpoint(ctx, &mediapackagev2.DeleteOriginEndpointInput{
		ChannelGroupName:   r.ChannelGroupName,
		ChannelName:        r.ChannelName,
		OriginEndpointName: r.OriginEndpointName,
	})
	return err
}

func (r *MediaPackageV2OriginEndpoint) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MediaPackageV2OriginEndpoint) String() string {
	return *r.OriginEndpointName
}
