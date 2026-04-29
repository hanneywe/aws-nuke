package resources

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/elasticache"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ElastiCacheServerlessCacheSnapshotResource = "ElastiCacheServerlessCacheSnapshot"

func init() {
	registry.Register(&registry.Registration{
		Name:     ElastiCacheServerlessCacheSnapshotResource,
		Scope:    nuke.Account,
		Resource: &ElastiCacheServerlessCacheSnapshot{},
		Lister:   &ElastiCacheServerlessCacheSnapshotLister{},
	})
}

type ElastiCacheServerlessCacheSnapshotLister struct {
	svc ElasticacheClient
}

func (l *ElastiCacheServerlessCacheSnapshotLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = elasticache.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	paginator := elasticache.NewDescribeServerlessCacheSnapshotsPaginator(svc, &elasticache.DescribeServerlessCacheSnapshotsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.ServerlessCacheSnapshots {
			item := &resp.ServerlessCacheSnapshots[i]
			resources = append(resources, &ElastiCacheServerlessCacheSnapshot{
				svc:                         svc,
				ServerlessCacheSnapshotName: item.ServerlessCacheSnapshotName,
				ARN:                         item.ARN,
				Status:                      item.Status,
				SnapshotType:                item.SnapshotType,
				CreateTime:                  item.CreateTime,
			})
		}
	}

	return resources, nil
}

type ElastiCacheServerlessCacheSnapshot struct {
	svc                         ElasticacheClient
	ServerlessCacheSnapshotName *string
	ARN                         *string
	Status                      *string
	SnapshotType                *string
	CreateTime                  *time.Time
}

func (r *ElastiCacheServerlessCacheSnapshot) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteServerlessCacheSnapshot(ctx, &elasticache.DeleteServerlessCacheSnapshotInput{
		ServerlessCacheSnapshotName: r.ServerlessCacheSnapshotName,
	})
	return err
}

func (r *ElastiCacheServerlessCacheSnapshot) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ElastiCacheServerlessCacheSnapshot) String() string {
	return *r.ServerlessCacheSnapshotName
}
