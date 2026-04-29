package resources

import (
	"context"
	"fmt"

	schedulerv2 "github.com/aws/aws-sdk-go-v2/service/scheduler"
	schedulertypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SchedulerScheduleGroupResource = "SchedulerScheduleGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     SchedulerScheduleGroupResource,
		Scope:    nuke.Account,
		Resource: &SchedulerScheduleGroup{},
		Lister:   &SchedulerScheduleGroupLister{},
	})
}

type SchedulerScheduleGroupLister struct {
	svc SchedulerV2Client
}

func (l *SchedulerScheduleGroupLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = schedulerv2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &schedulerv2.ListScheduleGroupsInput{}
	for {
		output, err := svc.ListScheduleGroups(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, scheduleGroup := range output.ScheduleGroups {
			groupState := string(scheduleGroup.State)
			resources = append(resources, &SchedulerScheduleGroup{
				svc:   svc,
				Name:  scheduleGroup.Name,
				Arn:   scheduleGroup.Arn,
				State: &groupState,
			})
		}

		if output.NextToken == nil {
			break
		}
		params.NextToken = output.NextToken
	}

	return resources, nil
}

type SchedulerScheduleGroup struct {
	svc   SchedulerV2Client
	Name  *string
	Arn   *string
	State *string
}

func (r *SchedulerScheduleGroup) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteScheduleGroup(ctx, &schedulerv2.DeleteScheduleGroupInput{
		Name: r.Name,
	})
	return err
}

func (r *SchedulerScheduleGroup) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SchedulerScheduleGroup) String() string {
	return *r.Name
}

func (r *SchedulerScheduleGroup) Filter() error {
	if r.Name != nil && *r.Name == "default" {
		return fmt.Errorf("cannot delete default schedule group")
	}
	if r.State != nil && schedulertypes.ScheduleGroupState(*r.State) == schedulertypes.ScheduleGroupStateDeleting {
		return fmt.Errorf("already %s", *r.State)
	}
	return nil
}
