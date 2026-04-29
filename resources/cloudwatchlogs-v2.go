package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

// CloudWatchLogsV2Client is an interface for the AWS CloudWatch Logs SDK v2 client methods.
// This is separate from the existing CloudWatch Logs resources which use SDK v1.
type CloudWatchLogsV2Client interface {
	DescribeQueryDefinitions(ctx context.Context, params *cloudwatchlogs.DescribeQueryDefinitionsInput,
		optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeQueryDefinitionsOutput, error)
	DeleteQueryDefinition(ctx context.Context, params *cloudwatchlogs.DeleteQueryDefinitionInput,
		optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DeleteQueryDefinitionOutput, error)
	DescribeLogGroups(ctx context.Context, params *cloudwatchlogs.DescribeLogGroupsInput,
		optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error)
	DescribeLogStreams(ctx context.Context, params *cloudwatchlogs.DescribeLogStreamsInput,
		optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogStreamsOutput, error)
	DeleteLogStream(ctx context.Context, params *cloudwatchlogs.DeleteLogStreamInput,
		optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DeleteLogStreamOutput, error)
	DeleteRetentionPolicy(ctx context.Context, params *cloudwatchlogs.DeleteRetentionPolicyInput,
		optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DeleteRetentionPolicyOutput, error)
	PutBearerTokenAuthentication(ctx context.Context, params *cloudwatchlogs.PutBearerTokenAuthenticationInput,
		optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.PutBearerTokenAuthenticationOutput, error)
	ListIntegrations(ctx context.Context, params *cloudwatchlogs.ListIntegrationsInput,
		optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.ListIntegrationsOutput, error)
	DeleteIntegration(ctx context.Context, params *cloudwatchlogs.DeleteIntegrationInput,
		optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DeleteIntegrationOutput, error)
}
