package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MediaTailorSourceLocationResource = "MediaTailorSourceLocation"

func init() {
	registry.Register(&registry.Registration{
		Name:     MediaTailorSourceLocationResource,
		Scope:    nuke.Account,
		Resource: &MediaTailorSourceLocation{},
		Lister:   &MediaTailorSourceLocationLister{},
		DependsOn: []string{
			MediaTailorLiveSourceResource,
			MediaTailorVodSourceResource,
		},
	})
}

type MediaTailorSourceLocationLister struct {
	svc MediaTailorV2Client
}

func (l *MediaTailorSourceLocationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = mediatailor.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &mediatailor.ListSourceLocationsInput{
		MaxResults: aws.Int32(100),
	}

	for {
		output, err := svc.ListSourceLocations(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, item := range output.Items {
			resources = append(resources, &MediaTailorSourceLocation{
				svc:                svc,
				SourceLocationName: item.SourceLocationName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

type MediaTailorSourceLocation struct {
	svc                MediaTailorV2Client
	SourceLocationName *string
}

func (r *MediaTailorSourceLocation) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteSourceLocation(ctx, &mediatailor.DeleteSourceLocationInput{
		SourceLocationName: r.SourceLocationName,
	})
	return err
}

func (r *MediaTailorSourceLocation) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MediaTailorSourceLocation) String() string {
	return *r.SourceLocationName
}
