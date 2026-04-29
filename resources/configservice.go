package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/configservice"
)

// ConfigServiceClient is the interface for the Config Service SDK v2 client methods.
type ConfigServiceClient interface {
	DescribeAggregationAuthorizations(ctx context.Context, params *configservice.DescribeAggregationAuthorizationsInput,
		optFns ...func(*configservice.Options)) (*configservice.DescribeAggregationAuthorizationsOutput, error)
	DeleteAggregationAuthorization(ctx context.Context, params *configservice.DeleteAggregationAuthorizationInput,
		optFns ...func(*configservice.Options)) (*configservice.DeleteAggregationAuthorizationOutput, error)
	DescribeRetentionConfigurations(ctx context.Context, params *configservice.DescribeRetentionConfigurationsInput,
		optFns ...func(*configservice.Options)) (*configservice.DescribeRetentionConfigurationsOutput, error)
	DeleteRetentionConfiguration(ctx context.Context, params *configservice.DeleteRetentionConfigurationInput,
		optFns ...func(*configservice.Options)) (*configservice.DeleteRetentionConfigurationOutput, error)
	DescribeConfigurationAggregators(ctx context.Context, params *configservice.DescribeConfigurationAggregatorsInput,
		optFns ...func(*configservice.Options)) (*configservice.DescribeConfigurationAggregatorsOutput, error)
	DeleteConfigurationAggregator(ctx context.Context, params *configservice.DeleteConfigurationAggregatorInput,
		optFns ...func(*configservice.Options)) (*configservice.DeleteConfigurationAggregatorOutput, error)
	DescribeOrganizationConfigRules(ctx context.Context, params *configservice.DescribeOrganizationConfigRulesInput,
		optFns ...func(*configservice.Options)) (*configservice.DescribeOrganizationConfigRulesOutput, error)
	DeleteOrganizationConfigRule(ctx context.Context, params *configservice.DeleteOrganizationConfigRuleInput,
		optFns ...func(*configservice.Options)) (*configservice.DeleteOrganizationConfigRuleOutput, error)
	ListStoredQueries(ctx context.Context, params *configservice.ListStoredQueriesInput,
		optFns ...func(*configservice.Options)) (*configservice.ListStoredQueriesOutput, error)
	DeleteStoredQuery(ctx context.Context, params *configservice.DeleteStoredQueryInput,
		optFns ...func(*configservice.Options)) (*configservice.DeleteStoredQueryOutput, error)
	ListTagsForResource(ctx context.Context, params *configservice.ListTagsForResourceInput,
		optFns ...func(*configservice.Options)) (*configservice.ListTagsForResourceOutput, error)
}
