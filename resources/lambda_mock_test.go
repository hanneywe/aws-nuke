package resources

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/lambda"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testLambdaListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

type mockLambdaClient struct {
	mock.Mock
}

func (m *mockLambdaClient) ListCodeSigningConfigs(ctx context.Context, params *lambda.ListCodeSigningConfigsInput,
	_ ...func(*lambda.Options)) (*lambda.ListCodeSigningConfigsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lambda.ListCodeSigningConfigsOutput), args.Error(1)
}

func (m *mockLambdaClient) DeleteCodeSigningConfig(ctx context.Context, params *lambda.DeleteCodeSigningConfigInput,
	_ ...func(*lambda.Options)) (*lambda.DeleteCodeSigningConfigOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lambda.DeleteCodeSigningConfigOutput), args.Error(1)
}

func (m *mockLambdaClient) ListFunctions(ctx context.Context, params *lambda.ListFunctionsInput,
	_ ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lambda.ListFunctionsOutput), args.Error(1)
}

func (m *mockLambdaClient) ListAliases(ctx context.Context, params *lambda.ListAliasesInput,
	_ ...func(*lambda.Options)) (*lambda.ListAliasesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lambda.ListAliasesOutput), args.Error(1)
}

func (m *mockLambdaClient) DeleteAlias(ctx context.Context, params *lambda.DeleteAliasInput,
	_ ...func(*lambda.Options)) (*lambda.DeleteAliasOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lambda.DeleteAliasOutput), args.Error(1)
}

func (m *mockLambdaClient) ListVersionsByFunction(ctx context.Context, params *lambda.ListVersionsByFunctionInput,
	_ ...func(*lambda.Options)) (*lambda.ListVersionsByFunctionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lambda.ListVersionsByFunctionOutput), args.Error(1)
}

func (m *mockLambdaClient) DeleteFunction(ctx context.Context, params *lambda.DeleteFunctionInput,
	_ ...func(*lambda.Options)) (*lambda.DeleteFunctionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lambda.DeleteFunctionOutput), args.Error(1)
}
