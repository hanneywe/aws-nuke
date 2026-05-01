package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CloudFrontRealtimeLogConfigResource = "CloudFrontRealtimeLogConfig"

func init() {
	registry.Register(&registry.Registration{
		Name:     CloudFrontRealtimeLogConfigResource,
		Scope:    nuke.Account,
		Resource: &CloudFrontRealtimeLogConfig{},
		Lister:   &CloudFrontRealtimeLogConfigLister{},
	})
}

type CloudFrontRealtimeLogConfigLister struct {
	svc CloudFrontClient
}

func (l *CloudFrontRealtimeLogConfigLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = cloudfront.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	params := &cloudfront.ListRealtimeLogConfigsInput{
		MaxItems: aws.Int32(100),
	}

	for {
		resp, err := svc.ListRealtimeLogConfigs(ctx, params)
		if err != nil {
			return nil, err
		}

		if resp.RealtimeLogConfigs == nil {
			break
		}

		for _, item := range resp.RealtimeLogConfigs.Items {
			resources = append(resources, &CloudFrontRealtimeLogConfig{
				svc:          svc,
				Name:         item.Name,
				ARN:          item.ARN,
				SamplingRate: item.SamplingRate,
			})
		}

		if resp.RealtimeLogConfigs.NextMarker == nil {
			break
		}
		params.Marker = resp.RealtimeLogConfigs.NextMarker
	}

	return resources, nil
}

type CloudFrontRealtimeLogConfig struct {
	svc          CloudFrontClient
	Name         *string
	ARN          *string
	SamplingRate *int64
}

func (r *CloudFrontRealtimeLogConfig) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteRealtimeLogConfig(ctx, &cloudfront.DeleteRealtimeLogConfigInput{
		ARN: r.ARN,
	})
	return err
}

func (r *CloudFrontRealtimeLogConfig) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CloudFrontRealtimeLogConfig) String() string {
	return *r.Name
}
