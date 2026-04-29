package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

type mockCloudWatchLogsV2Client struct {
	mock.Mock
}

func (m *mockCloudWatchLogsV2Client) DescribeQueryDefinitions(ctx context.Context, params *cloudwatchlogs.DescribeQueryDefinitionsInput,
	_ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeQueryDefinitionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudwatchlogs.DescribeQueryDefinitionsOutput), args.Error(1)
}

func (m *mockCloudWatchLogsV2Client) DeleteQueryDefinition(ctx context.Context, params *cloudwatchlogs.DeleteQueryDefinitionInput,
	_ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DeleteQueryDefinitionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudwatchlogs.DeleteQueryDefinitionOutput), args.Error(1)
}

func (m *mockCloudWatchLogsV2Client) DescribeLogGroups(ctx context.Context, params *cloudwatchlogs.DescribeLogGroupsInput,
	_ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudwatchlogs.DescribeLogGroupsOutput), args.Error(1)
}

func (m *mockCloudWatchLogsV2Client) DescribeLogStreams(ctx context.Context, params *cloudwatchlogs.DescribeLogStreamsInput,
	_ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogStreamsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudwatchlogs.DescribeLogStreamsOutput), args.Error(1)
}

func (m *mockCloudWatchLogsV2Client) DeleteLogStream(ctx context.Context, params *cloudwatchlogs.DeleteLogStreamInput,
	_ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DeleteLogStreamOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudwatchlogs.DeleteLogStreamOutput), args.Error(1)
}

func (m *mockCloudWatchLogsV2Client) DeleteRetentionPolicy(ctx context.Context, params *cloudwatchlogs.DeleteRetentionPolicyInput,
	_ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DeleteRetentionPolicyOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudwatchlogs.DeleteRetentionPolicyOutput), args.Error(1)
}

func (m *mockCloudWatchLogsV2Client) PutBearerTokenAuthentication(
	ctx context.Context, params *cloudwatchlogs.PutBearerTokenAuthenticationInput,
	_ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.PutBearerTokenAuthenticationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudwatchlogs.PutBearerTokenAuthenticationOutput), args.Error(1)
}

func (m *mockCloudWatchLogsV2Client) ListIntegrations(ctx context.Context, params *cloudwatchlogs.ListIntegrationsInput,
	_ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.ListIntegrationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudwatchlogs.ListIntegrationsOutput), args.Error(1)
}

func (m *mockCloudWatchLogsV2Client) DeleteIntegration(ctx context.Context, params *cloudwatchlogs.DeleteIntegrationInput,
	_ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DeleteIntegrationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudwatchlogs.DeleteIntegrationOutput), args.Error(1)
}
