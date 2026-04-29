package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/apprunner"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockAppRunnerClient struct {
	mock.Mock
}

func (m *mockAppRunnerClient) ListAutoScalingConfigurations(
	ctx context.Context,
	params *apprunner.ListAutoScalingConfigurationsInput,
	_ ...func(*apprunner.Options),
) (*apprunner.ListAutoScalingConfigurationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apprunner.ListAutoScalingConfigurationsOutput), args.Error(1)
}

func (m *mockAppRunnerClient) DeleteAutoScalingConfiguration(
	ctx context.Context,
	params *apprunner.DeleteAutoScalingConfigurationInput,
	_ ...func(*apprunner.Options),
) (*apprunner.DeleteAutoScalingConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apprunner.DeleteAutoScalingConfigurationOutput), args.Error(1)
}

func (m *mockAppRunnerClient) ListObservabilityConfigurations(
	ctx context.Context,
	params *apprunner.ListObservabilityConfigurationsInput,
	_ ...func(*apprunner.Options),
) (*apprunner.ListObservabilityConfigurationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apprunner.ListObservabilityConfigurationsOutput), args.Error(1)
}

func (m *mockAppRunnerClient) DeleteObservabilityConfiguration(
	ctx context.Context,
	params *apprunner.DeleteObservabilityConfigurationInput,
	_ ...func(*apprunner.Options),
) (*apprunner.DeleteObservabilityConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apprunner.DeleteObservabilityConfigurationOutput), args.Error(1)
}

var testAppRunnerListerOpts = &nuke.ListerOpts{}
