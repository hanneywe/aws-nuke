package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/neptunegraph"
	neptunegraphtypes "github.com/aws/aws-sdk-go-v2/service/neptunegraph/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const NeptuneGraphSnapshotResource = "NeptuneGraphSnapshot"

func init() {
	registry.Register(&registry.Registration{
		Name:     NeptuneGraphSnapshotResource,
		Scope:    nuke.Account,
		Resource: &NeptuneGraphSnapshot{},
		Lister:   &NeptuneGraphSnapshotLister{},
	})
}

type NeptuneGraphSnapshotLister struct {
	svc NeptuneGraphClient
}

func (l *NeptuneGraphSnapshotLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = neptunegraph.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := neptunegraph.NewListGraphSnapshotsPaginator(svc, &neptunegraph.ListGraphSnapshotsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.GraphSnapshots {
			resources = append(resources, &NeptuneGraphSnapshot{
				svc:        svc,
				SnapshotID: item.Id,
				Name:       item.Name,
				Status:     string(item.Status),
			})
		}
	}

	return resources, nil
}

type NeptuneGraphSnapshot struct {
	svc        NeptuneGraphClient
	SnapshotID *string
	Name       *string
	Status     string
}

func (r *NeptuneGraphSnapshot) Filter() error {
	if r.Status == string(neptunegraphtypes.SnapshotStatusDeleting) {
		return fmt.Errorf("already deleting")
	}
	return nil
}

func (r *NeptuneGraphSnapshot) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteGraphSnapshot(ctx, &neptunegraph.DeleteGraphSnapshotInput{
		SnapshotIdentifier: r.SnapshotID,
	})
	return err
}

func (r *NeptuneGraphSnapshot) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *NeptuneGraphSnapshot) String() string {
	return *r.SnapshotID
}
