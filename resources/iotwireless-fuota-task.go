package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iotwireless"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IoTWirelessFuotaTaskResource = "IoTWirelessFuotaTask"

func init() {
	registry.Register(&registry.Registration{
		Name:     IoTWirelessFuotaTaskResource,
		Scope:    nuke.Account,
		Resource: &IoTWirelessFuotaTask{},
		Lister:   &IoTWirelessFuotaTaskLister{},
	})
}

type IoTWirelessFuotaTaskLister struct {
	svc IoTWirelessClient
}

func (l *IoTWirelessFuotaTaskLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = iotwireless.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := iotwireless.NewListFuotaTasksPaginator(svc, &iotwireless.ListFuotaTasksInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, fuotaTask := range resp.FuotaTaskList {
			resources = append(resources, &IoTWirelessFuotaTask{
				svc:  svc,
				ID:   fuotaTask.Id,
				Name: fuotaTask.Name,
			})
		}
	}

	return resources, nil
}

type IoTWirelessFuotaTask struct {
	svc  IoTWirelessClient
	ID   *string `property:"name=Id"`
	Name *string
}

func (r *IoTWirelessFuotaTask) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteFuotaTask(ctx, &iotwireless.DeleteFuotaTaskInput{
		Id: r.ID,
	})
	return err
}

func (r *IoTWirelessFuotaTask) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IoTWirelessFuotaTask) String() string {
	return *r.ID
}
