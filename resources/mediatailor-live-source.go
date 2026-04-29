package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MediaTailorLiveSourceResource = "MediaTailorLiveSource"

func init() {
	registry.Register(&registry.Registration{
		Name:     MediaTailorLiveSourceResource,
		Scope:    nuke.Account,
		Resource: &MediaTailorLiveSource{},
		Lister:   &MediaTailorLiveSourceLister{},
	})
}

type MediaTailorLiveSourceLister struct {
	svc MediaTailorV2Client
}

func (l *MediaTailorLiveSourceLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = mediatailor.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	locParams := &mediatailor.ListSourceLocationsInput{
		MaxResults: aws.Int32(100),
	}

	for {
		locOutput, err := svc.ListSourceLocations(ctx, locParams)
		if err != nil {
			return nil, err
		}

		for _, loc := range locOutput.Items {
			srcParams := &mediatailor.ListLiveSourcesInput{
				SourceLocationName: loc.SourceLocationName,
				MaxResults:         aws.Int32(100),
			}

			for {
				srcOutput, err := svc.ListLiveSources(ctx, srcParams)
				if err != nil {
					return nil, err
				}

				for _, item := range srcOutput.Items {
					resources = append(resources, &MediaTailorLiveSource{
						svc:                svc,
						SourceLocationName: item.SourceLocationName,
						LiveSourceName:     item.LiveSourceName,
					})
				}

				if srcOutput.NextToken == nil {
					break
				}

				srcParams.NextToken = srcOutput.NextToken
			}
		}

		if locOutput.NextToken == nil {
			break
		}

		locParams.NextToken = locOutput.NextToken
	}

	return resources, nil
}

type MediaTailorLiveSource struct {
	svc                MediaTailorV2Client
	SourceLocationName *string
	LiveSourceName     *string
}

func (r *MediaTailorLiveSource) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteLiveSource(ctx, &mediatailor.DeleteLiveSourceInput{
		SourceLocationName: r.SourceLocationName,
		LiveSourceName:     r.LiveSourceName,
	})
	return err
}

func (r *MediaTailorLiveSource) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MediaTailorLiveSource) String() string {
	return fmt.Sprintf("%s/%s", *r.SourceLocationName, *r.LiveSourceName)
}
