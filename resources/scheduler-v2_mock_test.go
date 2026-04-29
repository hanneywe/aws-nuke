package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	schedulerv2 "github.com/aws/aws-sdk-go-v2/service/scheduler"
)

type mockSchedulerV2Client struct {
	mock.Mock
}

func (m *mockSchedulerV2Client) ListScheduleGroups(ctx context.Context, params *schedulerv2.ListScheduleGroupsInput,
	_ ...func(*schedulerv2.Options)) (*schedulerv2.ListScheduleGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*schedulerv2.ListScheduleGroupsOutput), args.Error(1)
}

func (m *mockSchedulerV2Client) DeleteScheduleGroup(ctx context.Context, params *schedulerv2.DeleteScheduleGroupInput,
	_ ...func(*schedulerv2.Options)) (*schedulerv2.DeleteScheduleGroupOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*schedulerv2.DeleteScheduleGroupOutput), args.Error(1)
}
