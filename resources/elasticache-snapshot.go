package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/elasticache"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ElasticacheSnapshotResource = "ElasticacheSnapshot"

func init() {
	registry.Register(&registry.Registration{
		Name:     ElasticacheSnapshotResource,
		Scope:    nuke.Account,
		Resource: &ElasticacheSnapshot{},
		Lister:   &ElasticacheSnapshotLister{},
	})
}

type ElasticacheSnapshotLister struct {
	svc ElasticacheClient
}

func (l *ElasticacheSnapshotLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = elasticache.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := elasticache.NewDescribeSnapshotsPaginator(svc, &elasticache.DescribeSnapshotsInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.Snapshots {
			snapshot := &resp.Snapshots[i]

			tags := make(map[string]string)
			if snapshot.ARN != nil {
				tagsResp, err := svc.ListTagsForResource(ctx, &elasticache.ListTagsForResourceInput{
					ResourceName: snapshot.ARN,
				})
				if err == nil {
					for _, tag := range tagsResp.TagList {
						tags[*tag.Key] = *tag.Value
					}
				}
			}

			resources = append(resources, &ElasticacheSnapshot{
				svc:            svc,
				SnapshotName:   snapshot.SnapshotName,
				CacheClusterID: snapshot.CacheClusterId,
				Status:         snapshot.SnapshotStatus,
				Tags:           tags,
			})
		}
	}

	return resources, nil
}

type ElasticacheSnapshot struct {
	svc            ElasticacheClient
	SnapshotName   *string
	CacheClusterID *string
	Status         *string
	Tags           map[string]string
}

func (r *ElasticacheSnapshot) Filter() error {
	if r.Status != nil {
		status := strings.ToLower(*r.Status)
		if status == "creating" || status == "deleting" || status == "restoring" {
			return fmt.Errorf("snapshot is in transient state: %s", status)
		}
	}
	if r.SnapshotName != nil && strings.HasPrefix(*r.SnapshotName, "automatic.") {
		return fmt.Errorf("cannot delete system-managed automatic snapshot")
	}
	return nil
}

func (r *ElasticacheSnapshot) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteSnapshot(ctx, &elasticache.DeleteSnapshotInput{
		SnapshotName: r.SnapshotName,
	})
	return err
}

func (r *ElasticacheSnapshot) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ElasticacheSnapshot) String() string {
	return *r.SnapshotName
}
