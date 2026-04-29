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

const MediaTailorChannelResource = "MediaTailorChannel"

func init() {
	registry.Register(&registry.Registration{
		Name:     MediaTailorChannelResource,
		Scope:    nuke.Account,
		Resource: &MediaTailorChannel{},
		Lister:   &MediaTailorChannelLister{},
	})
}

type MediaTailorChannelLister struct {
	svc MediaTailorV2Client
}

func (l *MediaTailorChannelLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = mediatailor.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &mediatailor.ListChannelsInput{
		MaxResults: aws.Int32(100),
	}

	for {
		output, err := svc.ListChannels(ctx, params)
		if err != nil {
			return nil, err
		}

		for i := range output.Items {
			channel := &output.Items[i]
			resources = append(resources, &MediaTailorChannel{
				svc:         svc,
				ChannelName: channel.ChannelName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

type MediaTailorChannel struct {
	svc         MediaTailorV2Client
	ChannelName *string
}

func (r *MediaTailorChannel) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteChannel(ctx, &mediatailor.DeleteChannelInput{
		ChannelName: r.ChannelName,
	})
	return err
}

func (r *MediaTailorChannel) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MediaTailorChannel) String() string {
	return *r.ChannelName
}
