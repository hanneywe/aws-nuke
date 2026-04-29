package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type mockSSMV2Client struct {
	mock.Mock
}

func (m *mockSSMV2Client) DescribeOpsItems(ctx context.Context, params *ssm.DescribeOpsItemsInput,
	_ ...func(*ssm.Options)) (*ssm.DescribeOpsItemsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ssm.DescribeOpsItemsOutput), args.Error(1)
}

func (m *mockSSMV2Client) UpdateOpsItem(ctx context.Context, params *ssm.UpdateOpsItemInput,
	_ ...func(*ssm.Options)) (*ssm.UpdateOpsItemOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ssm.UpdateOpsItemOutput), args.Error(1)
}

func (m *mockSSMV2Client) DescribeMaintenanceWindows(ctx context.Context, params *ssm.DescribeMaintenanceWindowsInput,
	_ ...func(*ssm.Options)) (*ssm.DescribeMaintenanceWindowsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ssm.DescribeMaintenanceWindowsOutput), args.Error(1)
}

func (m *mockSSMV2Client) DescribeMaintenanceWindowTasks(ctx context.Context, params *ssm.DescribeMaintenanceWindowTasksInput,
	_ ...func(*ssm.Options)) (*ssm.DescribeMaintenanceWindowTasksOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ssm.DescribeMaintenanceWindowTasksOutput), args.Error(1)
}

func (m *mockSSMV2Client) DeregisterTaskFromMaintenanceWindow(ctx context.Context, params *ssm.DeregisterTaskFromMaintenanceWindowInput,
	_ ...func(*ssm.Options)) (*ssm.DeregisterTaskFromMaintenanceWindowOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ssm.DeregisterTaskFromMaintenanceWindowOutput), args.Error(1)
}

func (m *mockSSMV2Client) DescribeMaintenanceWindowTargets(ctx context.Context, params *ssm.DescribeMaintenanceWindowTargetsInput,
	_ ...func(*ssm.Options)) (*ssm.DescribeMaintenanceWindowTargetsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ssm.DescribeMaintenanceWindowTargetsOutput), args.Error(1)
}

func (m *mockSSMV2Client) DeregisterTargetFromMaintenanceWindow(ctx context.Context, params *ssm.DeregisterTargetFromMaintenanceWindowInput,
	_ ...func(*ssm.Options)) (*ssm.DeregisterTargetFromMaintenanceWindowOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ssm.DeregisterTargetFromMaintenanceWindowOutput), args.Error(1)
}

func (m *mockSSMV2Client) DescribePatchBaselines(ctx context.Context, params *ssm.DescribePatchBaselinesInput,
	_ ...func(*ssm.Options)) (*ssm.DescribePatchBaselinesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ssm.DescribePatchBaselinesOutput), args.Error(1)
}

func (m *mockSSMV2Client) GetDefaultPatchBaseline(ctx context.Context, params *ssm.GetDefaultPatchBaselineInput,
	_ ...func(*ssm.Options)) (*ssm.GetDefaultPatchBaselineOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ssm.GetDefaultPatchBaselineOutput), args.Error(1)
}

func (m *mockSSMV2Client) RegisterDefaultPatchBaseline(ctx context.Context, params *ssm.RegisterDefaultPatchBaselineInput,
	_ ...func(*ssm.Options)) (*ssm.RegisterDefaultPatchBaselineOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ssm.RegisterDefaultPatchBaselineOutput), args.Error(1)
}
