package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

type CloudWatchV2Client interface {
	DescribeInsightRules(ctx context.Context, params *cloudwatch.DescribeInsightRulesInput,
		optFns ...func(*cloudwatch.Options)) (*cloudwatch.DescribeInsightRulesOutput, error)
	DeleteInsightRules(ctx context.Context, params *cloudwatch.DeleteInsightRulesInput,
		optFns ...func(*cloudwatch.Options)) (*cloudwatch.DeleteInsightRulesOutput, error)
}
