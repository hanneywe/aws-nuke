package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/connect"
	connecttypes "github.com/aws/aws-sdk-go-v2/service/connect/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ConnectQueueResource = "ConnectQueue"

func init() {
	registry.Register(&registry.Registration{
		Name:     ConnectQueueResource,
		Scope:    nuke.Account,
		Resource: &ConnectQueue{},
		Lister:   &ConnectQueueLister{},
	})
}

type ConnectQueueLister struct {
	svc ConnectClient
}

func (l *ConnectQueueLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = connect.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	instancePaginator := connect.NewListInstancesPaginator(svc, &connect.ListInstancesInput{})
	for instancePaginator.HasMorePages() {
		instanceResp, err := instancePaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, instance := range instanceResp.InstanceSummaryList {
			queuePaginator := connect.NewListQueuesPaginator(svc, &connect.ListQueuesInput{
				InstanceId: instance.Id,
			})
			for queuePaginator.HasMorePages() {
				queueResp, err := queuePaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, queue := range queueResp.QueueSummaryList {
					resources = append(resources, &ConnectQueue{
						svc:        svc,
						InstanceID: instance.Id,
						ID:         queue.Id,
						Name:       queue.Name,
						QueueType:  queue.QueueType,
					})
				}
			}
		}
	}

	return resources, nil
}

type ConnectQueue struct {
	svc        ConnectClient
	InstanceID *string `property:"name=InstanceId"`
	ID         *string `property:"name=Id"`
	Name       *string
	QueueType  connecttypes.QueueType
}

func (r *ConnectQueue) Filter() error {
	if r.QueueType != connecttypes.QueueTypeStandard {
		return fmt.Errorf("cannot delete non-standard queue of type %s", r.QueueType)
	}
	return nil
}

func (r *ConnectQueue) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteQueue(ctx, &connect.DeleteQueueInput{
		InstanceId: r.InstanceID,
		QueueId:    r.ID,
	})
	return err
}

func (r *ConnectQueue) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ConnectQueue) String() string {
	return *r.Name
}
