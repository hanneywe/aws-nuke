package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3vectors"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const S3VectorIndexResource = "S3VectorIndex"

func init() {
	registry.Register(&registry.Registration{
		Name:     S3VectorIndexResource,
		Scope:    nuke.Account,
		Resource: &S3VectorIndex{},
		Lister:   &S3VectorIndexLister{},
	})
}

type S3VectorIndexLister struct {
	svc S3VectorsClient
}

func (l *S3VectorIndexLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = s3vectors.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	// First list all vector buckets
	bucketParams := &s3vectors.ListVectorBucketsInput{}
	for {
		bucketResp, err := svc.ListVectorBuckets(ctx, bucketParams)
		if err != nil {
			return nil, err
		}

		for _, bucket := range bucketResp.VectorBuckets {
			// Then list indexes per bucket
			indexParams := &s3vectors.ListIndexesInput{
				VectorBucketName: bucket.VectorBucketName,
			}
			for {
				indexResp, err := svc.ListIndexes(ctx, indexParams)
				if err != nil {
					return nil, err
				}

				for i := range indexResp.Indexes {
					resources = append(resources, &S3VectorIndex{
						svc:              svc,
						VectorBucketName: indexResp.Indexes[i].VectorBucketName,
						IndexName:        indexResp.Indexes[i].IndexName,
						IndexARN:         indexResp.Indexes[i].IndexArn,
						CreationTime:     indexResp.Indexes[i].CreationTime,
					})
				}

				if indexResp.NextToken == nil {
					break
				}
				indexParams.NextToken = indexResp.NextToken
			}
		}

		if bucketResp.NextToken == nil {
			break
		}
		bucketParams.NextToken = bucketResp.NextToken
	}

	return resources, nil
}

type S3VectorIndex struct {
	svc              S3VectorsClient
	VectorBucketName *string
	IndexName        *string
	IndexARN         *string
	CreationTime     *time.Time
}

func (r *S3VectorIndex) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteIndex(ctx, &s3vectors.DeleteIndexInput{
		VectorBucketName: r.VectorBucketName,
		IndexName:        r.IndexName,
	})
	return err
}

func (r *S3VectorIndex) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *S3VectorIndex) String() string {
	return fmt.Sprintf("%s/%s", *r.VectorBucketName, *r.IndexName)
}
