package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lightsail"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const LightsailDiskSnapshotResource = "LightsailDiskSnapshot"

func init() {
	registry.Register(&registry.Registration{
		Name:     LightsailDiskSnapshotResource,
		Scope:    nuke.Account,
		Resource: &LightsailDiskSnapshot{},
		Lister:   &LightsailDiskSnapshotLister{},
	})
}

type LightsailDiskSnapshotLister struct {
	svc LightsailClient
}

func (l *LightsailDiskSnapshotLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = lightsail.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	resp, err := svc.GetDiskSnapshots(ctx, &lightsail.GetDiskSnapshotsInput{})
	if err != nil {
		return nil, err
	}

	for i := range resp.DiskSnapshots {
		item := &resp.DiskSnapshots[i]
		resources = append(resources, &LightsailDiskSnapshot{
			svc:              svc,
			DiskSnapshotName: item.Name,
			FromDiskName:     item.FromDiskName,
			SizeInGb:         item.SizeInGb,
		})
	}

	return resources, nil
}

type LightsailDiskSnapshot struct {
	svc              LightsailClient
	DiskSnapshotName *string
	FromDiskName     *string
	SizeInGb         *int32
}

func (r *LightsailDiskSnapshot) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDiskSnapshot(ctx, &lightsail.DeleteDiskSnapshotInput{
		DiskSnapshotName: r.DiskSnapshotName,
	})
	return err
}

func (r *LightsailDiskSnapshot) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *LightsailDiskSnapshot) String() string {
	return *r.DiskSnapshotName
}
