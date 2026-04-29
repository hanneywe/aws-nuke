package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func Test_Mock_SSMMaintenanceWindowTask_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSSMV2Client)

	mockClient.On("DescribeMaintenanceWindows", mock.Anything, mock.Anything).
		Return(&ssm.DescribeMaintenanceWindowsOutput{
			WindowIdentities: []ssmtypes.MaintenanceWindowIdentity{
				{WindowId: ptr.String("mw-0123456789abcdef0")},
			},
		}, nil)

	mockClient.On("DescribeMaintenanceWindowTasks", mock.Anything, mock.Anything).
		Return(&ssm.DescribeMaintenanceWindowTasksOutput{
			Tasks: []ssmtypes.MaintenanceWindowTask{
				{
					WindowId:     ptr.String("mw-0123456789abcdef0"),
					WindowTaskId: ptr.String("task-0123456789abcdef0"),
					TaskArn:      ptr.String("arn:aws:ssm:us-east-1:123456789012:task/task-0123456789abcdef0"),
				},
				{
					WindowId:     ptr.String("mw-0123456789abcdef0"),
					WindowTaskId: ptr.String("task-abcdef0123456789a"),
					TaskArn:      ptr.String("arn:aws:ssm:us-east-1:123456789012:task/task-abcdef0123456789a"),
				},
			},
		}, nil)

	lister := &SSMMaintenanceWindowTaskLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSSMV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 2)

	task := resources[0].(*SSMMaintenanceWindowTask)
	assertions.Equal("mw-0123456789abcdef0", *task.WindowID)
	assertions.Equal("task-0123456789abcdef0", *task.WindowTaskID)
	assertions.Equal("arn:aws:ssm:us-east-1:123456789012:task/task-0123456789abcdef0", *task.TaskArn)

	secondTask := resources[1].(*SSMMaintenanceWindowTask)
	assertions.Equal("task-abcdef0123456789a", *secondTask.WindowTaskID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SSMMaintenanceWindowTask_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSSMV2Client)

	mockClient.On("DescribeMaintenanceWindows", mock.Anything, mock.Anything).
		Return(&ssm.DescribeMaintenanceWindowsOutput{
			WindowIdentities: []ssmtypes.MaintenanceWindowIdentity{
				{WindowId: ptr.String("mw-0123456789abcdef0")},
			},
		}, nil)

	mockClient.On("DescribeMaintenanceWindowTasks", mock.Anything, mock.Anything).
		Return(&ssm.DescribeMaintenanceWindowTasksOutput{
			Tasks: []ssmtypes.MaintenanceWindowTask{},
		}, nil)

	lister := &SSMMaintenanceWindowTaskLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSSMV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SSMMaintenanceWindowTask_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSSMV2Client)

	task := &SSMMaintenanceWindowTask{
		svc:          mockClient,
		WindowID:     ptr.String("mw-0123456789abcdef0"),
		WindowTaskID: ptr.String("task-0123456789abcdef0"),
		TaskArn:      ptr.String("arn:aws:ssm:us-east-1:123456789012:task/task-0123456789abcdef0"),
	}

	mockClient.On("DeregisterTaskFromMaintenanceWindow", mock.Anything, &ssm.DeregisterTaskFromMaintenanceWindowInput{
		WindowId:     task.WindowID,
		WindowTaskId: task.WindowTaskID,
	}).Return(&ssm.DeregisterTaskFromMaintenanceWindowOutput{}, nil)

	err := task.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SSMMaintenanceWindowTask_Properties(t *testing.T) {
	assertions := assert.New(t)

	task := SSMMaintenanceWindowTask{
		WindowID:     ptr.String("mw-0123456789abcdef0"),
		WindowTaskID: ptr.String("task-0123456789abcdef0"),
		TaskArn:      ptr.String("arn:aws:ssm:us-east-1:123456789012:task/task-0123456789abcdef0"),
	}

	properties := task.Properties()
	assertions.Equal("mw-0123456789abcdef0", properties.Get("WindowId"))
	assertions.Equal("task-0123456789abcdef0", properties.Get("WindowTaskId"))
	assertions.Equal("arn:aws:ssm:us-east-1:123456789012:task/task-0123456789abcdef0", properties.Get("TaskArn"))
}

func Test_Mock_SSMMaintenanceWindowTask_String(t *testing.T) {
	assertions := assert.New(t)
	task := SSMMaintenanceWindowTask{WindowTaskID: ptr.String("task-0123456789abcdef0")}
	assertions.Equal("task-0123456789abcdef0", task.String())
}
