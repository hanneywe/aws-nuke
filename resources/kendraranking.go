package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/kendraranking"
)

// KendraRankingClient is the interface for the Kendra Ranking SDK client methods.
type KendraRankingClient interface {
	ListRescoreExecutionPlans(ctx context.Context, params *kendraranking.ListRescoreExecutionPlansInput,
		optFns ...func(*kendraranking.Options)) (*kendraranking.ListRescoreExecutionPlansOutput, error)
	DeleteRescoreExecutionPlan(ctx context.Context, params *kendraranking.DeleteRescoreExecutionPlanInput,
		optFns ...func(*kendraranking.Options)) (*kendraranking.DeleteRescoreExecutionPlanOutput, error)
}
