package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SSMMaintenanceWindowTaskResource = "SSMMaintenanceWindowTask"

func init() {
	registry.Register(&registry.Registration{
		Name:     SSMMaintenanceWindowTaskResource,
		Scope:    nuke.Account,
		Resource: &SSMMaintenanceWindowTask{},
		Lister:   &SSMMaintenanceWindowTaskLister{},
	})
}

type SSMMaintenanceWindowTaskLister struct {
	svc SSMV2Client
}

func (l *SSMMaintenanceWindowTaskLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ssm.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	windowParams := &ssm.DescribeMaintenanceWindowsInput{}
	for {
		windowOutput, err := svc.DescribeMaintenanceWindows(ctx, windowParams)
		if err != nil {
			return nil, err
		}

		for _, window := range windowOutput.WindowIdentities {
			taskParams := &ssm.DescribeMaintenanceWindowTasksInput{
				WindowId: window.WindowId,
			}
			for {
				taskOutput, err := svc.DescribeMaintenanceWindowTasks(ctx, taskParams)
				if err != nil {
					return nil, err
				}

				for i := range taskOutput.Tasks {
					resources = append(resources, &SSMMaintenanceWindowTask{
						svc:          svc,
						WindowID:     taskOutput.Tasks[i].WindowId,
						WindowTaskID: taskOutput.Tasks[i].WindowTaskId,
						TaskArn:      taskOutput.Tasks[i].TaskArn,
					})
				}

				if taskOutput.NextToken == nil {
					break
				}
				taskParams.NextToken = taskOutput.NextToken
			}
		}

		if windowOutput.NextToken == nil {
			break
		}
		windowParams.NextToken = windowOutput.NextToken
	}

	return resources, nil
}

type SSMMaintenanceWindowTask struct {
	svc          SSMV2Client
	WindowID     *string `property:"name=WindowId"`
	WindowTaskID *string `property:"name=WindowTaskId"`
	TaskArn      *string
}

func (r *SSMMaintenanceWindowTask) Remove(ctx context.Context) error {
	_, err := r.svc.DeregisterTaskFromMaintenanceWindow(ctx, &ssm.DeregisterTaskFromMaintenanceWindowInput{
		WindowId:     r.WindowID,
		WindowTaskId: r.WindowTaskID,
	})
	return err
}

func (r *SSMMaintenanceWindowTask) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SSMMaintenanceWindowTask) String() string {
	return *r.WindowTaskID
}
