package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/apigateway"
)

// ApigatewayClient is the interface for the apigateway SDK client methods.
type ApigatewayClient interface {
	GetUsagePlans(ctx context.Context, params *apigateway.GetUsagePlansInput,
		optFns ...func(*apigateway.Options)) (*apigateway.GetUsagePlansOutput, error)
	GetUsagePlanKeys(ctx context.Context, params *apigateway.GetUsagePlanKeysInput,
		optFns ...func(*apigateway.Options)) (*apigateway.GetUsagePlanKeysOutput, error)
	DeleteUsagePlanKey(ctx context.Context, params *apigateway.DeleteUsagePlanKeyInput,
		optFns ...func(*apigateway.Options)) (*apigateway.DeleteUsagePlanKeyOutput, error)
}
