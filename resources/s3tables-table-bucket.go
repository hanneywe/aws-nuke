package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3tables"
	s3tablestypes "github.com/aws/aws-sdk-go-v2/service/s3tables/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const S3TableBucketResource = "S3TableBucket"

func init() {
	registry.Register(&registry.Registration{
		Name:     S3TableBucketResource,
		Scope:    nuke.Account,
		Resource: &S3TableBucket{},
		Lister:   &S3TableBucketLister{},
	})
}

type S3TableBucketLister struct {
	svc S3TablesClient
}

func (l *S3TableBucketLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = s3tables.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &s3tables.ListTableBucketsInput{}
	for {
		resp, err := svc.ListTableBuckets(ctx, params)
		if err != nil {
			return nil, err
		}

		for i := range resp.TableBuckets {
			resources = append(resources, &S3TableBucket{
				svc:        svc,
				Name:       resp.TableBuckets[i].Name,
				ARN:        resp.TableBuckets[i].Arn,
				CreatedAt:  resp.TableBuckets[i].CreatedAt,
				BucketType: string(resp.TableBuckets[i].Type),
			})
		}

		if resp.ContinuationToken == nil {
			break
		}
		params.ContinuationToken = resp.ContinuationToken
	}

	return resources, nil
}

type S3TableBucket struct {
	svc        S3TablesClient
	Name       *string
	ARN        *string
	CreatedAt  *time.Time
	BucketType string
}

func (r *S3TableBucket) Filter() error {
	if r.BucketType == string(s3tablestypes.TableBucketTypeAws) {
		return fmt.Errorf("cannot delete AWS-managed table bucket")
	}
	return nil
}

func (r *S3TableBucket) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteTableBucket(ctx, &s3tables.DeleteTableBucketInput{
		TableBucketARN: r.ARN,
	})
	return err
}

func (r *S3TableBucket) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *S3TableBucket) String() string {
	return *r.Name
}
