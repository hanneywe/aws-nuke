package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/medialive"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MediaLiveNetworkResource = "MediaLiveNetwork"

func init() {
	registry.Register(&registry.Registration{
		Name:     MediaLiveNetworkResource,
		Scope:    nuke.Account,
		Resource: &MediaLiveNetwork{},
		Lister:   &MediaLiveNetworkLister{},
	})
}

type MediaLiveNetworkLister struct {
	svc MediaLiveClient
}

func (l *MediaLiveNetworkLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = medialive.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := medialive.NewListNetworksPaginator(svc, &medialive.ListNetworksInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Networks {
			resources = append(resources, &MediaLiveNetwork{
				svc:  svc,
				ID:   item.Id,
				Name: item.Name,
			})
		}
	}

	return resources, nil
}

type MediaLiveNetwork struct {
	svc  MediaLiveClient
	ID   *string
	Name *string
}

func (r *MediaLiveNetwork) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteNetwork(ctx, &medialive.DeleteNetworkInput{
		NetworkId: r.ID,
	})
	return err
}

func (r *MediaLiveNetwork) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MediaLiveNetwork) String() string {
	return *r.ID
}
