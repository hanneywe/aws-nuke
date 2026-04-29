package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/apigateway"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockAPIGatewayV2Client struct {
	mock.Mock
}

func (m *mockAPIGatewayV2Client) GetRestApis(
	ctx context.Context, params *apigateway.GetRestApisInput,
	_ ...func(*apigateway.Options),
) (*apigateway.GetRestApisOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apigateway.GetRestApisOutput), args.Error(1)
}

func (m *mockAPIGatewayV2Client) GetGatewayResponses(
	ctx context.Context, params *apigateway.GetGatewayResponsesInput,
	_ ...func(*apigateway.Options),
) (*apigateway.GetGatewayResponsesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apigateway.GetGatewayResponsesOutput), args.Error(1)
}

func (m *mockAPIGatewayV2Client) DeleteGatewayResponse(
	ctx context.Context, params *apigateway.DeleteGatewayResponseInput,
	_ ...func(*apigateway.Options),
) (*apigateway.DeleteGatewayResponseOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apigateway.DeleteGatewayResponseOutput), args.Error(1)
}

var testAPIGatewayV2ListerOpts = &nuke.ListerOpts{}
