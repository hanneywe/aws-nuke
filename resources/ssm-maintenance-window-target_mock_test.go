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

func Test_Mock_SSMMaintenanceWindowTarget_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSSMV2Client)

	mockClient.On("DescribeMaintenanceWindows", mock.Anything, mock.Anything).
		Return(&ssm.DescribeMaintenanceWindowsOutput{
			WindowIdentities: []ssmtypes.MaintenanceWindowIdentity{
				{WindowId: ptr.String("mw-0123456789abcdef0")},
			},
		}, nil)

	mockClient.On("DescribeMaintenanceWindowTargets", mock.Anything, mock.Anything).
		Return(&ssm.DescribeMaintenanceWindowTargetsOutput{
			Targets: []ssmtypes.MaintenanceWindowTarget{
				{
					WindowId:       ptr.String("mw-0123456789abcdef0"),
					WindowTargetId: ptr.String("target-0123456789abcdef0"),
					Name:           ptr.String("my-target"),
				},
				{
					WindowId:       ptr.String("mw-0123456789abcdef0"),
					WindowTargetId: ptr.String("target-abcdef0123456789a"),
					Name:           ptr.String("my-target-2"),
				},
			},
		}, nil)

	lister := &SSMMaintenanceWindowTargetLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSSMV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 2)

	target := resources[0].(*SSMMaintenanceWindowTarget)
	assertions.Equal("mw-0123456789abcdef0", *target.WindowID)
	assertions.Equal("target-0123456789abcdef0", *target.WindowTargetID)
	assertions.Equal("my-target", *target.Name)

	secondTarget := resources[1].(*SSMMaintenanceWindowTarget)
	assertions.Equal("target-abcdef0123456789a", *secondTarget.WindowTargetID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SSMMaintenanceWindowTarget_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSSMV2Client)

	mockClient.On("DescribeMaintenanceWindows", mock.Anything, mock.Anything).
		Return(&ssm.DescribeMaintenanceWindowsOutput{
			WindowIdentities: []ssmtypes.MaintenanceWindowIdentity{
				{WindowId: ptr.String("mw-0123456789abcdef0")},
			},
		}, nil)

	mockClient.On("DescribeMaintenanceWindowTargets", mock.Anything, mock.Anything).
		Return(&ssm.DescribeMaintenanceWindowTargetsOutput{
			Targets: []ssmtypes.MaintenanceWindowTarget{},
		}, nil)

	lister := &SSMMaintenanceWindowTargetLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSSMV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SSMMaintenanceWindowTarget_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSSMV2Client)

	target := &SSMMaintenanceWindowTarget{
		svc:            mockClient,
		WindowID:       ptr.String("mw-0123456789abcdef0"),
		WindowTargetID: ptr.String("target-0123456789abcdef0"),
		Name:           ptr.String("my-target"),
	}

	mockClient.On("DeregisterTargetFromMaintenanceWindow", mock.Anything, &ssm.DeregisterTargetFromMaintenanceWindowInput{
		WindowId:       target.WindowID,
		WindowTargetId: target.WindowTargetID,
	}).Return(&ssm.DeregisterTargetFromMaintenanceWindowOutput{}, nil)

	err := target.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SSMMaintenanceWindowTarget_Properties(t *testing.T) {
	assertions := assert.New(t)

	target := SSMMaintenanceWindowTarget{
		WindowID:       ptr.String("mw-0123456789abcdef0"),
		WindowTargetID: ptr.String("target-0123456789abcdef0"),
		Name:           ptr.String("my-target"),
	}

	properties := target.Properties()
	assertions.Equal("mw-0123456789abcdef0", properties.Get("WindowId"))
	assertions.Equal("target-0123456789abcdef0", properties.Get("WindowTargetId"))
	assertions.Equal("my-target", properties.Get("Name"))
}

func Test_Mock_SSMMaintenanceWindowTarget_String(t *testing.T) {
	assertions := assert.New(t)
	target := SSMMaintenanceWindowTarget{WindowTargetID: ptr.String("target-0123456789abcdef0")}
	assertions.Equal("target-0123456789abcdef0", target.String())
}
