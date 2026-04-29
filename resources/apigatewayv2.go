package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
)

// Apigatewayv2Client is the interface for the apigatewayv2 SDK client methods.
type Apigatewayv2Client interface {
	GetApis(ctx context.Context, params *apigatewayv2.GetApisInput,
		optFns ...func(*apigatewayv2.Options)) (*apigatewayv2.GetApisOutput, error)
	GetStages(ctx context.Context, params *apigatewayv2.GetStagesInput,
		optFns ...func(*apigatewayv2.Options)) (*apigatewayv2.GetStagesOutput, error)
	DeleteStage(ctx context.Context, params *apigatewayv2.DeleteStageInput,
		optFns ...func(*apigatewayv2.Options)) (*apigatewayv2.DeleteStageOutput, error)
}
