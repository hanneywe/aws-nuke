package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

type mockCloudWatchV2Client struct {
	mock.Mock
}

func (m *mockCloudWatchV2Client) DescribeInsightRules(ctx context.Context, params *cloudwatch.DescribeInsightRulesInput,
	_ ...func(*cloudwatch.Options)) (*cloudwatch.DescribeInsightRulesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudwatch.DescribeInsightRulesOutput), args.Error(1)
}

func (m *mockCloudWatchV2Client) DeleteInsightRules(ctx context.Context, params *cloudwatch.DeleteInsightRulesInput,
	_ ...func(*cloudwatch.Options)) (*cloudwatch.DeleteInsightRulesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudwatch.DeleteInsightRulesOutput), args.Error(1)
}
