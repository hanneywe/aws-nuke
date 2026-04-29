package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/iot"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockIoTClient struct {
	mock.Mock
}

func (m *mockIoTClient) ListBillingGroups(ctx context.Context,
	params *iot.ListBillingGroupsInput,
	_ ...func(*iot.Options)) (*iot.ListBillingGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iot.ListBillingGroupsOutput), args.Error(1)
}

func (m *mockIoTClient) DeleteBillingGroup(ctx context.Context,
	params *iot.DeleteBillingGroupInput,
	_ ...func(*iot.Options)) (*iot.DeleteBillingGroupOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iot.DeleteBillingGroupOutput), args.Error(1)
}

func (m *mockIoTClient) ListCommands(ctx context.Context,
	params *iot.ListCommandsInput,
	_ ...func(*iot.Options)) (*iot.ListCommandsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iot.ListCommandsOutput), args.Error(1)
}

func (m *mockIoTClient) DeleteCommand(ctx context.Context,
	params *iot.DeleteCommandInput,
	_ ...func(*iot.Options)) (*iot.DeleteCommandOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iot.DeleteCommandOutput), args.Error(1)
}

func (m *mockIoTClient) ListJobTemplates(ctx context.Context,
	params *iot.ListJobTemplatesInput,
	_ ...func(*iot.Options)) (*iot.ListJobTemplatesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iot.ListJobTemplatesOutput), args.Error(1)
}

func (m *mockIoTClient) DeleteJobTemplate(ctx context.Context,
	params *iot.DeleteJobTemplateInput,
	_ ...func(*iot.Options)) (*iot.DeleteJobTemplateOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iot.DeleteJobTemplateOutput), args.Error(1)
}

func (m *mockIoTClient) ListDimensions(ctx context.Context,
	params *iot.ListDimensionsInput,
	_ ...func(*iot.Options)) (*iot.ListDimensionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iot.ListDimensionsOutput), args.Error(1)
}

func (m *mockIoTClient) DeleteDimension(ctx context.Context,
	params *iot.DeleteDimensionInput,
	_ ...func(*iot.Options)) (*iot.DeleteDimensionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iot.DeleteDimensionOutput), args.Error(1)
}

func (m *mockIoTClient) ListMitigationActions(ctx context.Context,
	params *iot.ListMitigationActionsInput,
	_ ...func(*iot.Options)) (*iot.ListMitigationActionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iot.ListMitigationActionsOutput), args.Error(1)
}

func (m *mockIoTClient) DeleteMitigationAction(ctx context.Context,
	params *iot.DeleteMitigationActionInput,
	_ ...func(*iot.Options)) (*iot.DeleteMitigationActionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iot.DeleteMitigationActionOutput), args.Error(1)
}

func (m *mockIoTClient) ListAuditSuppressions(ctx context.Context,
	params *iot.ListAuditSuppressionsInput,
	_ ...func(*iot.Options)) (*iot.ListAuditSuppressionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iot.ListAuditSuppressionsOutput), args.Error(1)
}

func (m *mockIoTClient) DeleteAuditSuppression(ctx context.Context,
	params *iot.DeleteAuditSuppressionInput,
	_ ...func(*iot.Options)) (*iot.DeleteAuditSuppressionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iot.DeleteAuditSuppressionOutput), args.Error(1)
}

func (m *mockIoTClient) ListCustomMetrics(ctx context.Context,
	params *iot.ListCustomMetricsInput,
	_ ...func(*iot.Options)) (*iot.ListCustomMetricsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iot.ListCustomMetricsOutput), args.Error(1)
}

func (m *mockIoTClient) DeleteCustomMetric(ctx context.Context,
	params *iot.DeleteCustomMetricInput,
	_ ...func(*iot.Options)) (*iot.DeleteCustomMetricOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iot.DeleteCustomMetricOutput), args.Error(1)
}

func (m *mockIoTClient) ListPolicies(ctx context.Context,
	params *iot.ListPoliciesInput,
	_ ...func(*iot.Options)) (*iot.ListPoliciesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iot.ListPoliciesOutput), args.Error(1)
}

func (m *mockIoTClient) ListPolicyVersions(ctx context.Context,
	params *iot.ListPolicyVersionsInput,
	_ ...func(*iot.Options)) (*iot.ListPolicyVersionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iot.ListPolicyVersionsOutput), args.Error(1)
}

func (m *mockIoTClient) DeletePolicyVersion(ctx context.Context,
	params *iot.DeletePolicyVersionInput,
	_ ...func(*iot.Options)) (*iot.DeletePolicyVersionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iot.DeletePolicyVersionOutput), args.Error(1)
}

var testIoTListerOpts = &nuke.ListerOpts{}
