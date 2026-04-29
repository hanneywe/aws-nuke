package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/appsync"
)

// AppSyncClient is the interface for the AppSync SDK v2 client methods used by sub-resources.
type AppSyncClient interface {
	ListGraphqlApis(ctx context.Context, params *appsync.ListGraphqlApisInput,
		optFns ...func(*appsync.Options)) (*appsync.ListGraphqlApisOutput, error)
	GetApiCache(ctx context.Context, params *appsync.GetApiCacheInput,
		optFns ...func(*appsync.Options)) (*appsync.GetApiCacheOutput, error)
	DeleteApiCache(ctx context.Context, params *appsync.DeleteApiCacheInput,
		optFns ...func(*appsync.Options)) (*appsync.DeleteApiCacheOutput, error)
	ListApiKeys(ctx context.Context, params *appsync.ListApiKeysInput,
		optFns ...func(*appsync.Options)) (*appsync.ListApiKeysOutput, error)
	DeleteApiKey(ctx context.Context, params *appsync.DeleteApiKeyInput,
		optFns ...func(*appsync.Options)) (*appsync.DeleteApiKeyOutput, error)
	ListTypes(ctx context.Context, params *appsync.ListTypesInput,
		optFns ...func(*appsync.Options)) (*appsync.ListTypesOutput, error)
	DeleteType(ctx context.Context, params *appsync.DeleteTypeInput,
		optFns ...func(*appsync.Options)) (*appsync.DeleteTypeOutput, error)

	ListApis(ctx context.Context, params *appsync.ListApisInput,
		optFns ...func(*appsync.Options)) (*appsync.ListApisOutput, error)
	ListChannelNamespaces(ctx context.Context, params *appsync.ListChannelNamespacesInput,
		optFns ...func(*appsync.Options)) (*appsync.ListChannelNamespacesOutput, error)
	DeleteChannelNamespace(ctx context.Context, params *appsync.DeleteChannelNamespaceInput,
		optFns ...func(*appsync.Options)) (*appsync.DeleteChannelNamespaceOutput, error)
}
