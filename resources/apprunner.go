package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/apprunner"
)

// AppRunnerClient is the interface for the App Runner SDK v2 client methods.
type AppRunnerClient interface {
	ListAutoScalingConfigurations(ctx context.Context, params *apprunner.ListAutoScalingConfigurationsInput,
		optFns ...func(*apprunner.Options)) (*apprunner.ListAutoScalingConfigurationsOutput, error)
	DeleteAutoScalingConfiguration(ctx context.Context, params *apprunner.DeleteAutoScalingConfigurationInput,
		optFns ...func(*apprunner.Options)) (*apprunner.DeleteAutoScalingConfigurationOutput, error)
	ListObservabilityConfigurations(ctx context.Context, params *apprunner.ListObservabilityConfigurationsInput,
		optFns ...func(*apprunner.Options)) (*apprunner.ListObservabilityConfigurationsOutput, error)
	DeleteObservabilityConfiguration(ctx context.Context, params *apprunner.DeleteObservabilityConfigurationInput,
		optFns ...func(*apprunner.Options)) (*apprunner.DeleteObservabilityConfigurationOutput, error)
}
