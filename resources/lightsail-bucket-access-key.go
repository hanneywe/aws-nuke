package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lightsail"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const LightsailBucketAccessKeyResource = "LightsailBucketAccessKey"

func init() {
	registry.Register(&registry.Registration{
		Name:     LightsailBucketAccessKeyResource,
		Scope:    nuke.Account,
		Resource: &LightsailBucketAccessKey{},
		Lister:   &LightsailBucketAccessKeyLister{},
	})
}

type LightsailBucketAccessKeyLister struct {
	svc LightsailClient
}

func (l *LightsailBucketAccessKeyLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = lightsail.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	bucketParams := &lightsail.GetBucketsInput{}
	for {
		bucketResp, err := svc.GetBuckets(ctx, bucketParams)
		if err != nil {
			return nil, err
		}
		for i := range bucketResp.Buckets {
			bucket := &bucketResp.Buckets[i]
			keyResp, err := svc.GetBucketAccessKeys(ctx, &lightsail.GetBucketAccessKeysInput{
				BucketName: bucket.Name,
			})
			if err != nil {
				return nil, err
			}
			for _, key := range keyResp.AccessKeys {
				resources = append(resources, &LightsailBucketAccessKey{
					svc:         svc,
					BucketName:  bucket.Name,
					AccessKeyID: key.AccessKeyId,
				})
			}
		}
		if bucketResp.NextPageToken == nil {
			break
		}
		bucketParams.PageToken = bucketResp.NextPageToken
	}

	return resources, nil
}

type LightsailBucketAccessKey struct {
	svc         LightsailClient
	BucketName  *string
	AccessKeyID *string
}

func (r *LightsailBucketAccessKey) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteBucketAccessKey(ctx, &lightsail.DeleteBucketAccessKeyInput{
		BucketName:  r.BucketName,
		AccessKeyId: r.AccessKeyID,
	})
	return err
}

func (r *LightsailBucketAccessKey) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *LightsailBucketAccessKey) String() string {
	return *r.AccessKeyID
}
