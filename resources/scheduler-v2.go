package resources

import (
	"context"

	schedulerv2 "github.com/aws/aws-sdk-go-v2/service/scheduler"
)

// SchedulerV2Client is an interface for the AWS EventBridge Scheduler SDK v2 client methods.
// This is separate from the existing scheduler resources which use SDK v1.
type SchedulerV2Client interface {
	ListScheduleGroups(ctx context.Context, params *schedulerv2.ListScheduleGroupsInput,
		optFns ...func(*schedulerv2.Options)) (*schedulerv2.ListScheduleGroupsOutput, error)
	DeleteScheduleGroup(ctx context.Context, params *schedulerv2.DeleteScheduleGroupInput,
		optFns ...func(*schedulerv2.Options)) (*schedulerv2.DeleteScheduleGroupOutput, error)
}
