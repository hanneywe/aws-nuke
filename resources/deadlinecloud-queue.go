package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/deadline"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const DeadlineCloudQueueResource = "DeadlineCloudQueue"

func init() {
	registry.Register(&registry.Registration{
		Name:     DeadlineCloudQueueResource,
		Scope:    nuke.Account,
		Resource: &DeadlineCloudQueue{},
		Lister:   &DeadlineCloudQueueLister{},
		DependsOn: []string{
			DeadlineCloudQueueLimitAssociationResource,
		},
	})
}

type DeadlineCloudQueueLister struct {
	svc DeadlineCloudClient
}

func (l *DeadlineCloudQueueLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = deadline.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	// First list all farms
	farmPaginator := deadline.NewListFarmsPaginator(svc, &deadline.ListFarmsInput{})

	for farmPaginator.HasMorePages() {
		farmResp, err := farmPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, farm := range farmResp.Farms {
			// Then list queues per farm
			queuePaginator := deadline.NewListQueuesPaginator(svc, &deadline.ListQueuesInput{
				FarmId: farm.FarmId,
			})

			for queuePaginator.HasMorePages() {
				queueResp, err := queuePaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, queue := range queueResp.Queues {
					resources = append(resources, &DeadlineCloudQueue{
						svc:         svc,
						FarmID:      farm.FarmId,
						QueueID:     queue.QueueId,
						DisplayName: queue.DisplayName,
					})
				}
			}
		}
	}

	return resources, nil
}

type DeadlineCloudQueue struct {
	svc         DeadlineCloudClient
	FarmID      *string `property:"name=FarmId"`
	QueueID     *string `property:"name=QueueId"`
	DisplayName *string
}

func (r *DeadlineCloudQueue) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteQueue(ctx, &deadline.DeleteQueueInput{
		FarmId:  r.FarmID,
		QueueId: r.QueueID,
	})
	return err
}

func (r *DeadlineCloudQueue) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *DeadlineCloudQueue) String() string {
	return *r.QueueID
}
