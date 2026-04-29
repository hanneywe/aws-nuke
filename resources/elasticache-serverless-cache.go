package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const (
	ElasticacheServerlessCacheResource       = "ElasticacheServerlessCache"
	ElasticacheServerlessCacheStatusDeleting = "deleting"
)

func init() {
	registry.Register(&registry.Registration{
		Name:     ElasticacheServerlessCacheResource,
		Scope:    nuke.Account,
		Resource: &ElasticacheServerlessCache{},
		Lister:   &ElasticacheServerlessCacheLister{},
	})
}

type ElasticacheServerlessCacheLister struct {
	svc ElasticacheClient
}

func (l *ElasticacheServerlessCacheLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = elasticache.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := elasticache.NewDescribeServerlessCachesPaginator(svc, &elasticache.DescribeServerlessCachesInput{
		MaxResults: aws.Int32(50),
	})
	for paginator.HasMorePages() {
		describeCachesOutput, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for index := range describeCachesOutput.ServerlessCaches {
			serverlessCache := &describeCachesOutput.ServerlessCaches[index]
			resources = append(resources, &ElasticacheServerlessCache{
				svc:                 svc,
				ServerlessCacheName: serverlessCache.ServerlessCacheName,
				ARN:                 serverlessCache.ARN,
				Status:              serverlessCache.Status,
			})
		}
	}

	return resources, nil
}

type ElasticacheServerlessCache struct {
	svc                 ElasticacheClient
	ServerlessCacheName *string
	ARN                 *string
	Status              *string
}

func (r *ElasticacheServerlessCache) Filter() error {
	if r.Status != nil && *r.Status == ElasticacheServerlessCacheStatusDeleting {
		return fmt.Errorf("already deleting")
	}
	return nil
}

func (r *ElasticacheServerlessCache) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteServerlessCache(ctx, &elasticache.DeleteServerlessCacheInput{
		ServerlessCacheName: r.ServerlessCacheName,
	})
	return err
}

func (r *ElasticacheServerlessCache) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ElasticacheServerlessCache) String() string {
	return *r.ServerlessCacheName
}
