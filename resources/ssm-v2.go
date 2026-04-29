package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// SSMV2Client is an interface for the AWS SSM SDK v2 client methods.
// This is separate from the existing SSM resources which use SDK v1.
type SSMV2Client interface {
	DescribeOpsItems(ctx context.Context, params *ssm.DescribeOpsItemsInput,
		optFns ...func(*ssm.Options)) (*ssm.DescribeOpsItemsOutput, error)
	UpdateOpsItem(ctx context.Context, params *ssm.UpdateOpsItemInput,
		optFns ...func(*ssm.Options)) (*ssm.UpdateOpsItemOutput, error)
	DescribeMaintenanceWindows(ctx context.Context, params *ssm.DescribeMaintenanceWindowsInput,
		optFns ...func(*ssm.Options)) (*ssm.DescribeMaintenanceWindowsOutput, error)
	DescribeMaintenanceWindowTasks(ctx context.Context, params *ssm.DescribeMaintenanceWindowTasksInput,
		optFns ...func(*ssm.Options)) (*ssm.DescribeMaintenanceWindowTasksOutput, error)
	DeregisterTaskFromMaintenanceWindow(ctx context.Context, params *ssm.DeregisterTaskFromMaintenanceWindowInput,
		optFns ...func(*ssm.Options)) (*ssm.DeregisterTaskFromMaintenanceWindowOutput, error)
	DescribeMaintenanceWindowTargets(ctx context.Context, params *ssm.DescribeMaintenanceWindowTargetsInput,
		optFns ...func(*ssm.Options)) (*ssm.DescribeMaintenanceWindowTargetsOutput, error)
	DeregisterTargetFromMaintenanceWindow(ctx context.Context, params *ssm.DeregisterTargetFromMaintenanceWindowInput,
		optFns ...func(*ssm.Options)) (*ssm.DeregisterTargetFromMaintenanceWindowOutput, error)
	DescribePatchBaselines(ctx context.Context, params *ssm.DescribePatchBaselinesInput,
		optFns ...func(*ssm.Options)) (*ssm.DescribePatchBaselinesOutput, error)
	GetDefaultPatchBaseline(ctx context.Context, params *ssm.GetDefaultPatchBaselineInput,
		optFns ...func(*ssm.Options)) (*ssm.GetDefaultPatchBaselineOutput, error)
	RegisterDefaultPatchBaseline(ctx context.Context, params *ssm.RegisterDefaultPatchBaselineInput,
		optFns ...func(*ssm.Options)) (*ssm.RegisterDefaultPatchBaselineOutput, error)
}
