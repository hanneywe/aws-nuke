package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockApigatewayv2Client struct {
	mock.Mock
}

func (m *mockApigatewayv2Client) GetApis(
	ctx context.Context, params *apigatewayv2.GetApisInput,
	_ ...func(*apigatewayv2.Options),
) (*apigatewayv2.GetApisOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apigatewayv2.GetApisOutput), args.Error(1)
}

func (m *mockApigatewayv2Client) GetStages(
	ctx context.Context, params *apigatewayv2.GetStagesInput,
	_ ...func(*apigatewayv2.Options),
) (*apigatewayv2.GetStagesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apigatewayv2.GetStagesOutput), args.Error(1)
}

func (m *mockApigatewayv2Client) DeleteStage(
	ctx context.Context, params *apigatewayv2.DeleteStageInput,
	_ ...func(*apigatewayv2.Options),
) (*apigatewayv2.DeleteStageOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apigatewayv2.DeleteStageOutput), args.Error(1)
}

var testApigatewayv2ListerOpts = &nuke.ListerOpts{}
