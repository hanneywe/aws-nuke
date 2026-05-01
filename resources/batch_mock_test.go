package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/batch"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockBatchClient struct {
	mock.Mock
}

func (m *mockBatchClient) ListSchedulingPolicies(
	ctx context.Context, params *batch.ListSchedulingPoliciesInput,
	_ ...func(*batch.Options),
) (*batch.ListSchedulingPoliciesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*batch.ListSchedulingPoliciesOutput), args.Error(1)
}

func (m *mockBatchClient) DeleteSchedulingPolicy(
	ctx context.Context, params *batch.DeleteSchedulingPolicyInput,
	_ ...func(*batch.Options),
) (*batch.DeleteSchedulingPolicyOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*batch.DeleteSchedulingPolicyOutput), args.Error(1)
}

func (m *mockBatchClient) ListConsumableResources(
	ctx context.Context, params *batch.ListConsumableResourcesInput,
	_ ...func(*batch.Options),
) (*batch.ListConsumableResourcesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*batch.ListConsumableResourcesOutput), args.Error(1)
}

func (m *mockBatchClient) DeleteConsumableResource(
	ctx context.Context, params *batch.DeleteConsumableResourceInput,
	_ ...func(*batch.Options),
) (*batch.DeleteConsumableResourceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*batch.DeleteConsumableResourceOutput), args.Error(1)
}

func (m *mockBatchClient) DescribeJobDefinitions(
	ctx context.Context, params *batch.DescribeJobDefinitionsInput,
	_ ...func(*batch.Options),
) (*batch.DescribeJobDefinitionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*batch.DescribeJobDefinitionsOutput), args.Error(1)
}

func (m *mockBatchClient) DeregisterJobDefinition(
	ctx context.Context, params *batch.DeregisterJobDefinitionInput,
	_ ...func(*batch.Options),
) (*batch.DeregisterJobDefinitionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*batch.DeregisterJobDefinitionOutput), args.Error(1)
}

var testBatchListerOpts = &nuke.ListerOpts{}
