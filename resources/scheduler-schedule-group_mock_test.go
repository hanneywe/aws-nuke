package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	schedulerv2 "github.com/aws/aws-sdk-go-v2/service/scheduler"
	schedulertypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testSchedulerV2ListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_SchedulerScheduleGroup_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSchedulerV2Client)

	mockClient.On("ListScheduleGroups", mock.Anything, mock.Anything).
		Return(&schedulerv2.ListScheduleGroupsOutput{
			ScheduleGroups: []schedulertypes.ScheduleGroupSummary{
				{
					Name:  ptr.String("my-group"),
					Arn:   ptr.String("arn:aws:scheduler:us-east-1:123456789012:schedule-group/my-group"),
					State: schedulertypes.ScheduleGroupStateActive,
				},
			},
		}, nil)

	lister := &SchedulerScheduleGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSchedulerV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	scheduleGroup := resources[0].(*SchedulerScheduleGroup)
	assertions.Equal("my-group", *scheduleGroup.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SchedulerScheduleGroup_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSchedulerV2Client)

	mockClient.On("ListScheduleGroups", mock.Anything, mock.Anything).
		Return(&schedulerv2.ListScheduleGroupsOutput{
			ScheduleGroups: []schedulertypes.ScheduleGroupSummary{},
		}, nil)

	lister := &SchedulerScheduleGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSchedulerV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SchedulerScheduleGroup_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSchedulerV2Client)

	scheduleGroup := &SchedulerScheduleGroup{
		svc:  mockClient,
		Name: ptr.String("my-group"),
	}

	mockClient.On("DeleteScheduleGroup", mock.Anything, &schedulerv2.DeleteScheduleGroupInput{
		Name: scheduleGroup.Name,
	}).Return(&schedulerv2.DeleteScheduleGroupOutput{}, nil)

	err := scheduleGroup.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SchedulerScheduleGroup_Properties(t *testing.T) {
	assertions := assert.New(t)

	scheduleGroup := SchedulerScheduleGroup{
		Name:  ptr.String("my-group"),
		Arn:   ptr.String("arn:aws:scheduler:us-east-1:123456789012:schedule-group/my-group"),
		State: ptr.String("ACTIVE"),
	}

	properties := scheduleGroup.Properties()
	assertions.Equal("my-group", properties.Get("Name"))
	assertions.Equal("arn:aws:scheduler:us-east-1:123456789012:schedule-group/my-group", properties.Get("Arn"))
}

func Test_Mock_SchedulerScheduleGroup_String(t *testing.T) {
	assertions := assert.New(t)
	scheduleGroup := SchedulerScheduleGroup{Name: ptr.String("my-group")}
	assertions.Equal("my-group", scheduleGroup.String())
}

func Test_Mock_SchedulerScheduleGroup_Filter(t *testing.T) {
	assertions := assert.New(t)

	defaultGroup := SchedulerScheduleGroup{Name: ptr.String("default"), State: ptr.String("ACTIVE")}
	assertions.Error(defaultGroup.Filter())

	deletingGroup := SchedulerScheduleGroup{Name: ptr.String("my-group"), State: ptr.String("DELETING")}
	assertions.Error(deletingGroup.Filter())

	activeGroup := SchedulerScheduleGroup{Name: ptr.String("my-group"), State: ptr.String("ACTIVE")}
	assertions.NoError(activeGroup.Filter())
}
