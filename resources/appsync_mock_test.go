package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/appsync"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockAppSyncClient struct {
	mock.Mock
}

func (m *mockAppSyncClient) ListGraphqlApis(
	ctx context.Context, params *appsync.ListGraphqlApisInput,
	_ ...func(*appsync.Options),
) (*appsync.ListGraphqlApisOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*appsync.ListGraphqlApisOutput), args.Error(1)
}

func (m *mockAppSyncClient) GetApiCache(
	ctx context.Context, params *appsync.GetApiCacheInput,
	_ ...func(*appsync.Options),
) (*appsync.GetApiCacheOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*appsync.GetApiCacheOutput), args.Error(1)
}

func (m *mockAppSyncClient) DeleteApiCache(
	ctx context.Context, params *appsync.DeleteApiCacheInput,
	_ ...func(*appsync.Options),
) (*appsync.DeleteApiCacheOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*appsync.DeleteApiCacheOutput), args.Error(1)
}

func (m *mockAppSyncClient) ListApiKeys(
	ctx context.Context, params *appsync.ListApiKeysInput,
	_ ...func(*appsync.Options),
) (*appsync.ListApiKeysOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*appsync.ListApiKeysOutput), args.Error(1)
}

func (m *mockAppSyncClient) DeleteApiKey(
	ctx context.Context, params *appsync.DeleteApiKeyInput,
	_ ...func(*appsync.Options),
) (*appsync.DeleteApiKeyOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*appsync.DeleteApiKeyOutput), args.Error(1)
}

func (m *mockAppSyncClient) ListTypes(
	ctx context.Context, params *appsync.ListTypesInput,
	_ ...func(*appsync.Options),
) (*appsync.ListTypesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*appsync.ListTypesOutput), args.Error(1)
}

func (m *mockAppSyncClient) DeleteType(
	ctx context.Context, params *appsync.DeleteTypeInput,
	_ ...func(*appsync.Options),
) (*appsync.DeleteTypeOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*appsync.DeleteTypeOutput), args.Error(1)
}

func (m *mockAppSyncClient) ListApis(
	ctx context.Context, params *appsync.ListApisInput,
	_ ...func(*appsync.Options),
) (*appsync.ListApisOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*appsync.ListApisOutput), args.Error(1)
}

func (m *mockAppSyncClient) ListChannelNamespaces(
	ctx context.Context, params *appsync.ListChannelNamespacesInput,
	_ ...func(*appsync.Options),
) (*appsync.ListChannelNamespacesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*appsync.ListChannelNamespacesOutput), args.Error(1)
}

func (m *mockAppSyncClient) DeleteChannelNamespace(
	ctx context.Context, params *appsync.DeleteChannelNamespaceInput,
	_ ...func(*appsync.Options),
) (*appsync.DeleteChannelNamespaceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*appsync.DeleteChannelNamespaceOutput), args.Error(1)
}

var testAppSyncListerOpts = &nuke.ListerOpts{}
