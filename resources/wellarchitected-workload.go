package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/wellarchitected"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const WellArchitectedWorkloadResource = "WellArchitectedWorkload"

func init() {
	registry.Register(&registry.Registration{
		Name:     WellArchitectedWorkloadResource,
		Scope:    nuke.Account,
		Resource: &WellArchitectedWorkload{},
		Lister:   &WellArchitectedWorkloadLister{},
	})
}

type WellArchitectedWorkloadLister struct {
	svc WellArchitectedClient
}

func (l *WellArchitectedWorkloadLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = wellarchitected.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := wellarchitected.NewListWorkloadsPaginator(svc, &wellarchitected.ListWorkloadsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range output.WorkloadSummaries {
			workload := &output.WorkloadSummaries[i]
			resources = append(resources, &WellArchitectedWorkload{
				svc:          svc,
				WorkloadID:   workload.WorkloadId,
				WorkloadName: workload.WorkloadName,
			})
		}
	}

	return resources, nil
}

type WellArchitectedWorkload struct {
	svc          WellArchitectedClient
	WorkloadID   *string `property:"name=WorkloadId"`
	WorkloadName *string
}

func (r *WellArchitectedWorkload) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteWorkload(ctx, &wellarchitected.DeleteWorkloadInput{
		WorkloadId:         r.WorkloadID,
		ClientRequestToken: r.WorkloadID,
	})
	return err
}

func (r *WellArchitectedWorkload) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *WellArchitectedWorkload) String() string {
	return *r.WorkloadID
}
