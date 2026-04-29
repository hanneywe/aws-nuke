package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/deadline"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const DeadlineCloudLimitResource = "DeadlineCloudLimit"

func init() {
	registry.Register(&registry.Registration{
		Name:     DeadlineCloudLimitResource,
		Scope:    nuke.Account,
		Resource: &DeadlineCloudLimit{},
		Lister:   &DeadlineCloudLimitLister{},
		DependsOn: []string{
			DeadlineCloudQueueLimitAssociationResource,
		},
	})
}

type DeadlineCloudLimitLister struct {
	svc DeadlineCloudClient
}

func (l *DeadlineCloudLimitLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = deadline.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	farmPaginator := deadline.NewListFarmsPaginator(svc, &deadline.ListFarmsInput{})

	for farmPaginator.HasMorePages() {
		farmResp, err := farmPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, farm := range farmResp.Farms {
			limitPaginator := deadline.NewListLimitsPaginator(svc, &deadline.ListLimitsInput{
				FarmId: farm.FarmId,
			})

			for limitPaginator.HasMorePages() {
				limitResp, err := limitPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, limit := range limitResp.Limits {
					resources = append(resources, &DeadlineCloudLimit{
						svc:         svc,
						FarmID:      farm.FarmId,
						LimitID:     limit.LimitId,
						DisplayName: limit.DisplayName,
					})
				}
			}
		}
	}

	return resources, nil
}

type DeadlineCloudLimit struct {
	svc         DeadlineCloudClient
	FarmID      *string `property:"name=FarmId"`
	LimitID     *string `property:"name=LimitId"`
	DisplayName *string
}

func (r *DeadlineCloudLimit) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteLimit(ctx, &deadline.DeleteLimitInput{
		FarmId:  r.FarmID,
		LimitId: r.LimitID,
	})
	return err
}

func (r *DeadlineCloudLimit) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *DeadlineCloudLimit) String() string {
	return *r.LimitID
}
