package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/apigateway"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockApigatewayClient struct {
	mock.Mock
}

func (m *mockApigatewayClient) GetUsagePlans(
	ctx context.Context, params *apigateway.GetUsagePlansInput,
	_ ...func(*apigateway.Options),
) (*apigateway.GetUsagePlansOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apigateway.GetUsagePlansOutput), args.Error(1)
}

func (m *mockApigatewayClient) GetUsagePlanKeys(
	ctx context.Context, params *apigateway.GetUsagePlanKeysInput,
	_ ...func(*apigateway.Options),
) (*apigateway.GetUsagePlanKeysOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apigateway.GetUsagePlanKeysOutput), args.Error(1)
}

func (m *mockApigatewayClient) DeleteUsagePlanKey(
	ctx context.Context, params *apigateway.DeleteUsagePlanKeyInput,
	_ ...func(*apigateway.Options),
) (*apigateway.DeleteUsagePlanKeyOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apigateway.DeleteUsagePlanKeyOutput), args.Error(1)
}

var testApigatewayListerOpts = &nuke.ListerOpts{}
