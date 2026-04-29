package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/medialive"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MediaLiveSdiSourceResource = "MediaLiveSdiSource"

func init() {
	registry.Register(&registry.Registration{
		Name:     MediaLiveSdiSourceResource,
		Scope:    nuke.Account,
		Resource: &MediaLiveSdiSource{},
		Lister:   &MediaLiveSdiSourceLister{},
	})
}

type MediaLiveSdiSourceLister struct {
	svc MediaLiveClient
}

func (l *MediaLiveSdiSourceLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = medialive.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := medialive.NewListSdiSourcesPaginator(svc, &medialive.ListSdiSourcesInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.SdiSources {
			resources = append(resources, &MediaLiveSdiSource{
				svc:  svc,
				ID:   item.Id,
				Name: item.Name,
			})
		}
	}

	return resources, nil
}

type MediaLiveSdiSource struct {
	svc  MediaLiveClient
	ID   *string
	Name *string
}

func (r *MediaLiveSdiSource) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteSdiSource(ctx, &medialive.DeleteSdiSourceInput{
		SdiSourceId: r.ID,
	})
	return err
}

func (r *MediaLiveSdiSource) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MediaLiveSdiSource) String() string {
	return *r.ID
}
