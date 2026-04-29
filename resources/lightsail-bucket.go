package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const LightsailBucketResource = "LightsailBucket"

func init() {
	registry.Register(&registry.Registration{
		Name:     LightsailBucketResource,
		Scope:    nuke.Account,
		Resource: &LightsailBucket{},
		Lister:   &LightsailBucketLister{},
	})
}

type LightsailBucketLister struct {
	svc LightsailClient
}

func (l *LightsailBucketLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = lightsail.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &lightsail.GetBucketsInput{}

	for {
		resp, err := svc.GetBuckets(ctx, params)
		if err != nil {
			return nil, err
		}

		for i := range resp.Buckets {
			bucket := &resp.Buckets[i]
			resources = append(resources, &LightsailBucket{
				svc:        svc,
				BucketName: bucket.Name,
				Arn:        bucket.Arn,
			})
		}

		if resp.NextPageToken == nil {
			break
		}
		params.PageToken = resp.NextPageToken
	}

	return resources, nil
}

type LightsailBucket struct {
	svc        LightsailClient
	BucketName *string
	Arn        *string
}

func (r *LightsailBucket) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteBucket(ctx, &lightsail.DeleteBucketInput{
		BucketName:  r.BucketName,
		ForceDelete: aws.Bool(true),
	})
	return err
}

func (r *LightsailBucket) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *LightsailBucket) String() string {
	return *r.BucketName
}
