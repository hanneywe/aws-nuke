package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/greengrass"
)

// GreengrassClient is the interface for the Greengrass SDK client methods.
type GreengrassClient interface {
	ListConnectorDefinitions(ctx context.Context, params *greengrass.ListConnectorDefinitionsInput,
		optFns ...func(*greengrass.Options)) (*greengrass.ListConnectorDefinitionsOutput, error)
	DeleteConnectorDefinition(ctx context.Context, params *greengrass.DeleteConnectorDefinitionInput,
		optFns ...func(*greengrass.Options)) (*greengrass.DeleteConnectorDefinitionOutput, error)
	ListCoreDefinitions(ctx context.Context, params *greengrass.ListCoreDefinitionsInput,
		optFns ...func(*greengrass.Options)) (*greengrass.ListCoreDefinitionsOutput, error)
	DeleteCoreDefinition(ctx context.Context, params *greengrass.DeleteCoreDefinitionInput,
		optFns ...func(*greengrass.Options)) (*greengrass.DeleteCoreDefinitionOutput, error)
	ListDeviceDefinitions(ctx context.Context, params *greengrass.ListDeviceDefinitionsInput,
		optFns ...func(*greengrass.Options)) (*greengrass.ListDeviceDefinitionsOutput, error)
	DeleteDeviceDefinition(ctx context.Context, params *greengrass.DeleteDeviceDefinitionInput,
		optFns ...func(*greengrass.Options)) (*greengrass.DeleteDeviceDefinitionOutput, error)
	ListFunctionDefinitions(ctx context.Context, params *greengrass.ListFunctionDefinitionsInput,
		optFns ...func(*greengrass.Options)) (*greengrass.ListFunctionDefinitionsOutput, error)
	DeleteFunctionDefinition(ctx context.Context, params *greengrass.DeleteFunctionDefinitionInput,
		optFns ...func(*greengrass.Options)) (*greengrass.DeleteFunctionDefinitionOutput, error)
	ListLoggerDefinitions(ctx context.Context, params *greengrass.ListLoggerDefinitionsInput,
		optFns ...func(*greengrass.Options)) (*greengrass.ListLoggerDefinitionsOutput, error)
	DeleteLoggerDefinition(ctx context.Context, params *greengrass.DeleteLoggerDefinitionInput,
		optFns ...func(*greengrass.Options)) (*greengrass.DeleteLoggerDefinitionOutput, error)
	ListSubscriptionDefinitions(ctx context.Context, params *greengrass.ListSubscriptionDefinitionsInput,
		optFns ...func(*greengrass.Options)) (*greengrass.ListSubscriptionDefinitionsOutput, error)
	DeleteSubscriptionDefinition(ctx context.Context, params *greengrass.DeleteSubscriptionDefinitionInput,
		optFns ...func(*greengrass.Options)) (*greengrass.DeleteSubscriptionDefinitionOutput, error)
	ListResourceDefinitions(ctx context.Context, params *greengrass.ListResourceDefinitionsInput,
		optFns ...func(*greengrass.Options)) (*greengrass.ListResourceDefinitionsOutput, error)
	DeleteResourceDefinition(ctx context.Context, params *greengrass.DeleteResourceDefinitionInput,
		optFns ...func(*greengrass.Options)) (*greengrass.DeleteResourceDefinitionOutput, error)
}
