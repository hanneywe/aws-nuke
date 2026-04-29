package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/redshift"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const RedshiftSnapshotCopyGrantResource = "RedshiftSnapshotCopyGrant"

func init() {
	registry.Register(&registry.Registration{
		Name:     RedshiftSnapshotCopyGrantResource,
		Scope:    nuke.Account,
		Resource: &RedshiftSnapshotCopyGrant{},
		Lister:   &RedshiftSnapshotCopyGrantLister{},
	})
}

type RedshiftSnapshotCopyGrantLister struct {
	svc RedshiftClient
}

func (l *RedshiftSnapshotCopyGrantLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = redshift.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := redshift.NewDescribeSnapshotCopyGrantsPaginator(svc, &redshift.DescribeSnapshotCopyGrantsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.SnapshotCopyGrants {
			resources = append(resources, &RedshiftSnapshotCopyGrant{
				svc:                   svc,
				SnapshotCopyGrantName: item.SnapshotCopyGrantName,
			})
		}
	}

	return resources, nil
}

type RedshiftSnapshotCopyGrant struct {
	svc                   RedshiftClient
	SnapshotCopyGrantName *string
}

func (r *RedshiftSnapshotCopyGrant) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteSnapshotCopyGrant(ctx, &redshift.DeleteSnapshotCopyGrantInput{
		SnapshotCopyGrantName: r.SnapshotCopyGrantName,
	})
	return err
}

func (r *RedshiftSnapshotCopyGrant) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *RedshiftSnapshotCopyGrant) String() string {
	return *r.SnapshotCopyGrantName
}
