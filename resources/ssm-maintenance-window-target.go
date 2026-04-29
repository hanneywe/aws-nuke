package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SSMMaintenanceWindowTargetResource = "SSMMaintenanceWindowTarget"

func init() {
	registry.Register(&registry.Registration{
		Name:     SSMMaintenanceWindowTargetResource,
		Scope:    nuke.Account,
		Resource: &SSMMaintenanceWindowTarget{},
		Lister:   &SSMMaintenanceWindowTargetLister{},
	})
}

type SSMMaintenanceWindowTargetLister struct {
	svc SSMV2Client
}

func (l *SSMMaintenanceWindowTargetLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
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
			targetParams := &ssm.DescribeMaintenanceWindowTargetsInput{
				WindowId: window.WindowId,
			}
			for {
				targetOutput, err := svc.DescribeMaintenanceWindowTargets(ctx, targetParams)
				if err != nil {
					return nil, err
				}

				for i := range targetOutput.Targets {
					resources = append(resources, &SSMMaintenanceWindowTarget{
						svc:            svc,
						WindowID:       targetOutput.Targets[i].WindowId,
						WindowTargetID: targetOutput.Targets[i].WindowTargetId,
						Name:           targetOutput.Targets[i].Name,
					})
				}

				if targetOutput.NextToken == nil {
					break
				}
				targetParams.NextToken = targetOutput.NextToken
			}
		}

		if windowOutput.NextToken == nil {
			break
		}
		windowParams.NextToken = windowOutput.NextToken
	}

	return resources, nil
}

type SSMMaintenanceWindowTarget struct {
	svc            SSMV2Client
	WindowID       *string `property:"name=WindowId"`
	WindowTargetID *string `property:"name=WindowTargetId"`
	Name           *string
}

func (r *SSMMaintenanceWindowTarget) Remove(ctx context.Context) error {
	_, err := r.svc.DeregisterTargetFromMaintenanceWindow(ctx, &ssm.DeregisterTargetFromMaintenanceWindowInput{
		WindowId:       r.WindowID,
		WindowTargetId: r.WindowTargetID,
	})
	return err
}

func (r *SSMMaintenanceWindowTarget) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SSMMaintenanceWindowTarget) String() string {
	return *r.WindowTargetID
}
