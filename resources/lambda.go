package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type LambdaClient interface {
	ListCodeSigningConfigs(ctx context.Context, params *lambda.ListCodeSigningConfigsInput,
		optFns ...func(*lambda.Options)) (*lambda.ListCodeSigningConfigsOutput, error)
	DeleteCodeSigningConfig(ctx context.Context, params *lambda.DeleteCodeSigningConfigInput,
		optFns ...func(*lambda.Options)) (*lambda.DeleteCodeSigningConfigOutput, error)
	ListFunctions(ctx context.Context, params *lambda.ListFunctionsInput,
		optFns ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error)
	ListAliases(ctx context.Context, params *lambda.ListAliasesInput,
		optFns ...func(*lambda.Options)) (*lambda.ListAliasesOutput, error)
	DeleteAlias(ctx context.Context, params *lambda.DeleteAliasInput,
		optFns ...func(*lambda.Options)) (*lambda.DeleteAliasOutput, error)
	ListVersionsByFunction(ctx context.Context, params *lambda.ListVersionsByFunctionInput,
		optFns ...func(*lambda.Options)) (*lambda.ListVersionsByFunctionOutput, error)
	DeleteFunction(ctx context.Context, params *lambda.DeleteFunctionInput,
		optFns ...func(*lambda.Options)) (*lambda.DeleteFunctionOutput, error)
}
