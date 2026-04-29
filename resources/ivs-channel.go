package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ivs"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IVSChannelResource = "IVSChannel"

func init() {
	registry.Register(&registry.Registration{
		Name:     IVSChannelResource,
		Scope:    nuke.Account,
		Resource: &IVSChannel{},
		Lister:   &IVSChannelLister{},
	})
}

type IVSChannelLister struct {
	svc IVSClient
}

func (l *IVSChannelLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ivs.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ivs.NewListChannelsPaginator(svc, &ivs.ListChannelsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, ch := range resp.Channels {
			resources = append(resources, &IVSChannel{
				svc:  svc,
				ARN:  ch.Arn,
				Name: ch.Name,
				Tags: ch.Tags,
			})
		}
	}

	return resources, nil
}

type IVSChannel struct {
	svc  IVSClient
	ARN  *string
	Name *string
	Tags map[string]string
}

func (r *IVSChannel) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteChannel(ctx, &ivs.DeleteChannelInput{
		Arn: r.ARN,
	})
	return err
}

func (r *IVSChannel) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IVSChannel) String() string {
	return *r.Name
}
