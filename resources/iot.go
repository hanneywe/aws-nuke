package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iot"
)

// IoTClient is the interface for the IoT SDK v2 client methods used by new IoT resources.
// Existing IoT resources use SDK v1; this interface is for new SDK v2 resources only.
type IoTClient interface {
	ListBillingGroups(ctx context.Context, params *iot.ListBillingGroupsInput,
		optFns ...func(*iot.Options)) (*iot.ListBillingGroupsOutput, error)
	DeleteBillingGroup(ctx context.Context, params *iot.DeleteBillingGroupInput,
		optFns ...func(*iot.Options)) (*iot.DeleteBillingGroupOutput, error)
	ListCommands(ctx context.Context, params *iot.ListCommandsInput,
		optFns ...func(*iot.Options)) (*iot.ListCommandsOutput, error)
	DeleteCommand(ctx context.Context, params *iot.DeleteCommandInput,
		optFns ...func(*iot.Options)) (*iot.DeleteCommandOutput, error)
	ListJobTemplates(ctx context.Context, params *iot.ListJobTemplatesInput,
		optFns ...func(*iot.Options)) (*iot.ListJobTemplatesOutput, error)
	DeleteJobTemplate(ctx context.Context, params *iot.DeleteJobTemplateInput,
		optFns ...func(*iot.Options)) (*iot.DeleteJobTemplateOutput, error)
	ListDimensions(ctx context.Context, params *iot.ListDimensionsInput,
		optFns ...func(*iot.Options)) (*iot.ListDimensionsOutput, error)
	DeleteDimension(ctx context.Context, params *iot.DeleteDimensionInput,
		optFns ...func(*iot.Options)) (*iot.DeleteDimensionOutput, error)
	ListMitigationActions(ctx context.Context, params *iot.ListMitigationActionsInput,
		optFns ...func(*iot.Options)) (*iot.ListMitigationActionsOutput, error)
	DeleteMitigationAction(ctx context.Context, params *iot.DeleteMitigationActionInput,
		optFns ...func(*iot.Options)) (*iot.DeleteMitigationActionOutput, error)
	ListAuditSuppressions(ctx context.Context, params *iot.ListAuditSuppressionsInput,
		optFns ...func(*iot.Options)) (*iot.ListAuditSuppressionsOutput, error)
	DeleteAuditSuppression(ctx context.Context, params *iot.DeleteAuditSuppressionInput,
		optFns ...func(*iot.Options)) (*iot.DeleteAuditSuppressionOutput, error)
	ListCustomMetrics(ctx context.Context, params *iot.ListCustomMetricsInput,
		optFns ...func(*iot.Options)) (*iot.ListCustomMetricsOutput, error)
	DeleteCustomMetric(ctx context.Context, params *iot.DeleteCustomMetricInput,
		optFns ...func(*iot.Options)) (*iot.DeleteCustomMetricOutput, error)

	ListPolicies(ctx context.Context, params *iot.ListPoliciesInput,
		optFns ...func(*iot.Options)) (*iot.ListPoliciesOutput, error)
	ListPolicyVersions(ctx context.Context, params *iot.ListPolicyVersionsInput,
		optFns ...func(*iot.Options)) (*iot.ListPolicyVersionsOutput, error)
	DeletePolicyVersion(ctx context.Context, params *iot.DeletePolicyVersionInput,
		optFns ...func(*iot.Options)) (*iot.DeletePolicyVersionOutput, error)
}
