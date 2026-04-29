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

const MediaTailorVodSourceResource = "MediaTailorVodSource"

func init() {
	registry.Register(&registry.Registration{
		Name:     MediaTailorVodSourceResource,
		Scope:    nuke.Account,
		Resource: &MediaTailorVodSource{},
		Lister:   &MediaTailorVodSourceLister{},
	})
}

type MediaTailorVodSourceLister struct {
	svc MediaTailorV2Client
}

func (l *MediaTailorVodSourceLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
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
			vodParams := &mediatailor.ListVodSourcesInput{
				SourceLocationName: loc.SourceLocationName,
				MaxResults:         aws.Int32(100),
			}

			for {
				vodOutput, err := svc.ListVodSources(ctx, vodParams)
				if err != nil {
					return nil, err
				}

				for _, item := range vodOutput.Items {
					resources = append(resources, &MediaTailorVodSource{
						svc:                svc,
						SourceLocationName: item.SourceLocationName,
						VodSourceName:      item.VodSourceName,
					})
				}

				if vodOutput.NextToken == nil {
					break
				}

				vodParams.NextToken = vodOutput.NextToken
			}
		}

		if locOutput.NextToken == nil {
			break
		}

		locParams.NextToken = locOutput.NextToken
	}

	return resources, nil
}

type MediaTailorVodSource struct {
	svc                MediaTailorV2Client
	SourceLocationName *string
	VodSourceName      *string
}

func (r *MediaTailorVodSource) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteVodSource(ctx, &mediatailor.DeleteVodSourceInput{
		SourceLocationName: r.SourceLocationName,
		VodSourceName:      r.VodSourceName,
	})
	return err
}

func (r *MediaTailorVodSource) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MediaTailorVodSource) String() string {
	return fmt.Sprintf("%s/%s", *r.SourceLocationName, *r.VodSourceName)
}
