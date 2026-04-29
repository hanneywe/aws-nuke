package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/apigateway"
)

type APIGatewayV2Client interface {
	GetRestApis(ctx context.Context, params *apigateway.GetRestApisInput,
		optFns ...func(*apigateway.Options)) (*apigateway.GetRestApisOutput, error)
	GetGatewayResponses(ctx context.Context, params *apigateway.GetGatewayResponsesInput,
		optFns ...func(*apigateway.Options)) (*apigateway.GetGatewayResponsesOutput, error)
	DeleteGatewayResponse(ctx context.Context, params *apigateway.DeleteGatewayResponseInput,
		optFns ...func(*apigateway.Options)) (*apigateway.DeleteGatewayResponseOutput, error)
}
