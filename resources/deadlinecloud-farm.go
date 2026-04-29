package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/deadline"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const DeadlineCloudFarmResource = "DeadlineCloudFarm"

func init() {
	registry.Register(&registry.Registration{
		Name:     DeadlineCloudFarmResource,
		Scope:    nuke.Account,
		Resource: &DeadlineCloudFarm{},
		Lister:   &DeadlineCloudFarmLister{},
		DependsOn: []string{
			DeadlineCloudQueueResource,
			DeadlineCloudStorageProfileResource,
			DeadlineCloudLimitResource,
		},
	})
}

type DeadlineCloudFarmLister struct {
	svc DeadlineCloudClient
}

func (l *DeadlineCloudFarmLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = deadline.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := deadline.NewListFarmsPaginator(svc, &deadline.ListFarmsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, farm := range resp.Farms {
			resources = append(resources, &DeadlineCloudFarm{
				svc:         svc,
				FarmID:      farm.FarmId,
				DisplayName: farm.DisplayName,
			})
		}
	}

	return resources, nil
}

type DeadlineCloudFarm struct {
	svc         DeadlineCloudClient
	FarmID      *string `property:"name=FarmId"`
	DisplayName *string
}

func (r *DeadlineCloudFarm) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteFarm(ctx, &deadline.DeleteFarmInput{
		FarmId: r.FarmID,
	})
	return err
}

func (r *DeadlineCloudFarm) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *DeadlineCloudFarm) String() string {
	return *r.FarmID
}
