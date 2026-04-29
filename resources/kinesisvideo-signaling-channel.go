package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/kinesisvideo"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const KinesisVideoSignalingChannelResource = "KinesisVideoSignalingChannel"

func init() {
	registry.Register(&registry.Registration{
		Name:     KinesisVideoSignalingChannelResource,
		Scope:    nuke.Account,
		Resource: &KinesisVideoSignalingChannel{},
		Lister:   &KinesisVideoSignalingChannelLister{},
	})
}

type KinesisVideoSignalingChannelLister struct {
	svc KinesisVideoV2Client
}

func (l *KinesisVideoSignalingChannelLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = kinesisvideo.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &kinesisvideo.ListSignalingChannelsInput{}

	for {
		output, err := svc.ListSignalingChannels(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, channelInfo := range output.ChannelInfoList {
			resources = append(resources, &KinesisVideoSignalingChannel{
				svc:         svc,
				ChannelARN:  channelInfo.ChannelARN,
				ChannelName: channelInfo.ChannelName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

type KinesisVideoSignalingChannel struct {
	svc         KinesisVideoV2Client
	ChannelARN  *string
	ChannelName *string
}

func (r *KinesisVideoSignalingChannel) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteSignalingChannel(ctx, &kinesisvideo.DeleteSignalingChannelInput{
		ChannelARN: r.ChannelARN,
	})
	return err
}

func (r *KinesisVideoSignalingChannel) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *KinesisVideoSignalingChannel) String() string {
	return *r.ChannelName
}
