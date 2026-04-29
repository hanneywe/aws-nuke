package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/kendraranking"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockKendraRankingClient struct {
	mock.Mock
}

func (m *mockKendraRankingClient) ListRescoreExecutionPlans(ctx context.Context,
	params *kendraranking.ListRescoreExecutionPlansInput,
	_ ...func(*kendraranking.Options)) (*kendraranking.ListRescoreExecutionPlansOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*kendraranking.ListRescoreExecutionPlansOutput), args.Error(1)
}

func (m *mockKendraRankingClient) DeleteRescoreExecutionPlan(ctx context.Context,
	params *kendraranking.DeleteRescoreExecutionPlanInput,
	_ ...func(*kendraranking.Options)) (*kendraranking.DeleteRescoreExecutionPlanOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*kendraranking.DeleteRescoreExecutionPlanOutput), args.Error(1)
}

var testKendraRankingListerOpts = &nuke.ListerOpts{}
