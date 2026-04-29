package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/configservice"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockConfigServiceClient struct {
	mock.Mock
}

func (m *mockConfigServiceClient) DescribeAggregationAuthorizations(
	ctx context.Context, params *configservice.DescribeAggregationAuthorizationsInput,
	_ ...func(*configservice.Options),
) (*configservice.DescribeAggregationAuthorizationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*configservice.DescribeAggregationAuthorizationsOutput), args.Error(1)
}

func (m *mockConfigServiceClient) DeleteAggregationAuthorization(
	ctx context.Context, params *configservice.DeleteAggregationAuthorizationInput,
	_ ...func(*configservice.Options),
) (*configservice.DeleteAggregationAuthorizationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*configservice.DeleteAggregationAuthorizationOutput), args.Error(1)
}

func (m *mockConfigServiceClient) DescribeRetentionConfigurations(
	ctx context.Context, params *configservice.DescribeRetentionConfigurationsInput,
	_ ...func(*configservice.Options),
) (*configservice.DescribeRetentionConfigurationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*configservice.DescribeRetentionConfigurationsOutput), args.Error(1)
}

func (m *mockConfigServiceClient) DeleteRetentionConfiguration(
	ctx context.Context, params *configservice.DeleteRetentionConfigurationInput,
	_ ...func(*configservice.Options),
) (*configservice.DeleteRetentionConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*configservice.DeleteRetentionConfigurationOutput), args.Error(1)
}

var testConfigServiceListerOpts = &nuke.ListerOpts{}

func (m *mockConfigServiceClient) DescribeConfigurationAggregators(
	ctx context.Context, params *configservice.DescribeConfigurationAggregatorsInput,
	_ ...func(*configservice.Options),
) (*configservice.DescribeConfigurationAggregatorsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*configservice.DescribeConfigurationAggregatorsOutput), args.Error(1)
}

func (m *mockConfigServiceClient) DeleteConfigurationAggregator(
	ctx context.Context, params *configservice.DeleteConfigurationAggregatorInput,
	_ ...func(*configservice.Options),
) (*configservice.DeleteConfigurationAggregatorOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*configservice.DeleteConfigurationAggregatorOutput), args.Error(1)
}

func (m *mockConfigServiceClient) DescribeOrganizationConfigRules(
	ctx context.Context, params *configservice.DescribeOrganizationConfigRulesInput,
	_ ...func(*configservice.Options),
) (*configservice.DescribeOrganizationConfigRulesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*configservice.DescribeOrganizationConfigRulesOutput), args.Error(1)
}

func (m *mockConfigServiceClient) DeleteOrganizationConfigRule(
	ctx context.Context, params *configservice.DeleteOrganizationConfigRuleInput,
	_ ...func(*configservice.Options),
) (*configservice.DeleteOrganizationConfigRuleOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*configservice.DeleteOrganizationConfigRuleOutput), args.Error(1)
}

func (m *mockConfigServiceClient) ListStoredQueries(
	ctx context.Context, params *configservice.ListStoredQueriesInput,
	_ ...func(*configservice.Options),
) (*configservice.ListStoredQueriesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*configservice.ListStoredQueriesOutput), args.Error(1)
}

func (m *mockConfigServiceClient) DeleteStoredQuery(
	ctx context.Context, params *configservice.DeleteStoredQueryInput,
	_ ...func(*configservice.Options),
) (*configservice.DeleteStoredQueryOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*configservice.DeleteStoredQueryOutput), args.Error(1)
}

func (m *mockConfigServiceClient) ListTagsForResource(
	ctx context.Context, params *configservice.ListTagsForResourceInput,
	_ ...func(*configservice.Options),
) (*configservice.ListTagsForResourceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*configservice.ListTagsForResourceOutput), args.Error(1)
}
