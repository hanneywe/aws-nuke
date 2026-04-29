package resources

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3vectors"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const S3VectorBucketResource = "S3VectorBucket"

func init() {
	registry.Register(&registry.Registration{
		Name:     S3VectorBucketResource,
		Scope:    nuke.Account,
		Resource: &S3VectorBucket{},
		Lister:   &S3VectorBucketLister{},
		DependsOn: []string{
			S3VectorIndexResource,
		},
	})
}

type S3VectorBucketLister struct {
	svc S3VectorsClient
}

func (l *S3VectorBucketLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = s3vectors.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &s3vectors.ListVectorBucketsInput{}
	for {
		resp, err := svc.ListVectorBuckets(ctx, params)
		if err != nil {
			return nil, err
		}

		for i := range resp.VectorBuckets {
			resources = append(resources, &S3VectorBucket{
				svc:          svc,
				Name:         resp.VectorBuckets[i].VectorBucketName,
				ARN:          resp.VectorBuckets[i].VectorBucketArn,
				CreationTime: resp.VectorBuckets[i].CreationTime,
			})
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type S3VectorBucket struct {
	svc          S3VectorsClient
	Name         *string
	ARN          *string
	CreationTime *time.Time
}

func (r *S3VectorBucket) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteVectorBucket(ctx, &s3vectors.DeleteVectorBucketInput{
		VectorBucketName: r.Name,
	})
	return err
}

func (r *S3VectorBucket) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *S3VectorBucket) String() string {
	return *r.Name
}
